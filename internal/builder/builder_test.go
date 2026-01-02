package builder

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/sporollan/site/internal/renderer"
	"github.com/sporollan/site/internal/site"
)

func setupTestSite(t *testing.T) (*site.Site, *renderer.Renderer, string) {
	t.Helper()

	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Create directories
	contentDir := filepath.Join(tmpDir, "content")
	staticDir := filepath.Join(tmpDir, "static")
	templateDir := filepath.Join(tmpDir, "templates")
	outputDir := filepath.Join(tmpDir, "public")

	for _, dir := range []string{contentDir, staticDir, templateDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create content files
	contentFiles := []struct {
		path    string
		content string
	}{
		{
			path: filepath.Join(contentDir, "index.md"),
			content: `---
title: "Home"
template: "home.html"
---

# Welcome

Homepage content.`,
		},
		{
			path: filepath.Join(contentDir, "about.md"),
			content: `---
title: "About"
---

# About Me

About page content.`,
		},
		{
			path: filepath.Join(contentDir, "blog", "post1.md"),
			content: `---
title: "First Post"
date: 2023-10-01
tags: ["go", "test"]
---

# Post 1

Blog post content.`,
		},
		{
			path: filepath.Join(contentDir, "blog", "post2.md"),
			content: `---
title: "Second Post"
date: 2023-10-02
tags: ["ssg", "web"]
draft: true
---

# Post 2

Draft post.`,
		},
	}

	for _, file := range contentFiles {
		// Create subdirectory
		dir := filepath.Dir(file.path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(file.path, []byte(file.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create static file
	staticFile := filepath.Join(staticDir, "style.css")
	if err := os.WriteFile(staticFile, []byte("body { color: red; }"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create templates
	templates := []struct {
		name    string
		content string
	}{
		{
			name:    "home.html",
			content: `<!DOCTYPE html><html><body><h1>{{.Title}}</h1>{{.Body | safeHTML}}</body></html>`,
		},
		{
			name:    "page.html",
			content: `<!DOCTYPE html><html><body><h1>{{.Title}}</h1>{{.Body | safeHTML}}</body></html>`,
		},
		{
			name:    "post.html",
			content: `<!DOCTYPE html><html><body><h1>{{.Title}}</h1>{{if .Date}}<time>{{.Date.Format "2006-01-02"}}</time>{{end}}{{.Body | safeHTML}}</body></html>`,
		},
		{
			name:    "list.html",
			content: `<!DOCTYPE html><html><body><h1>{{.Title}}</h1>{{range .Pages}}<h2>{{.Title}}</h2>{{end}}</body></html>`,
		},
	}

	for _, tmpl := range templates {
		path := filepath.Join(templateDir, tmpl.name)
		if err := os.WriteFile(path, []byte(tmpl.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create site
	s := site.NewWithConfig(
		contentDir,
		outputDir,
		staticDir,
		templateDir,
		"Test Site",
		"https://example.com",
	)

	// Create renderer
	r, err := renderer.New(templateDir)
	if err != nil {
		t.Fatal(err)
	}

	return s, r, tmpDir
}

func TestBuilder_Build(t *testing.T) {
	s, r, tmpDir := setupTestSite(t)

	b := New(s, r, 4)

	t.Run("successful build", func(t *testing.T) {
		if err := b.Build(); err != nil {
			t.Fatalf("Build() error = %v, want nil", err)
		}

		// Check output files were created
		expectedFiles := []string{
			"index.html",
			"about/index.html",
			"blog/post1/index.html",
			"style.css",
		}

		for _, file := range expectedFiles {
			path := filepath.Join(s.OutputDir, file)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Expected file not created: %s", file)
			}
		}

		// Check draft file was not created
		draftPath := filepath.Join(s.OutputDir, "blog/post2/index.html")
		if _, err := os.Stat(draftPath); !os.IsNotExist(err) {
			t.Error("Draft file should not be created")
		}

		// Check site statistics
		if len(s.Pages) != 3 { // index, about, post1 (post2 is draft)
			t.Errorf("Pages count = %d, want 3", len(s.Pages))
		}

		// Check content of generated files
		indexPath := filepath.Join(s.OutputDir, "index.html")
		indexContent, err := os.ReadFile(indexPath)
		if err != nil {
			t.Fatal(err)
		}

		// Debug: Print what's in the file
		t.Logf("Index file content length: %d", len(indexContent))
		t.Logf(
			"Index file content (first 100 chars): %s",
			string(indexContent[:min(100, len(indexContent))]),
		)

		// Check index contains expected content
		if len(indexContent) == 0 {
			t.Error("Index file should not be empty")
		}

		// Check it contains expected HTML elements
		if !bytes.Contains(indexContent, []byte("<!DOCTYPE html>")) {
			t.Error("Index file should contain DOCTYPE")
		}

		// Check it contains expected content from markdown
		if !bytes.Contains(indexContent, []byte("Home")) {
			t.Error("Index file should contain 'Home'")
		}
	})

	t.Run("clean output directory", func(t *testing.T) {
		// Create a file in output directory
		testFile := filepath.Join(s.OutputDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		// Build again
		if err := b.Build(); err != nil {
			t.Fatal(err)
		}

		// Test file should be removed
		if _, err := os.Stat(testFile); !os.IsNotExist(err) {
			t.Error("Old files should be cleaned")
		}
	})

	t.Run("static files are copied", func(t *testing.T) {
		staticFile := filepath.Join(s.OutputDir, "style.css")
		content, err := os.ReadFile(staticFile)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != "body { color: red; }" {
			t.Errorf("Static file content = %s, want 'body { color: red; }'", string(content))
		}
	})

	t.Run("non-existent content directory", func(t *testing.T) {
		s2 := site.NewWithConfig(
			"/non/existent/content",
			filepath.Join(tmpDir, "public2"),
			s.StaticDir,
			s.TemplateDir,
			"Test Site",
			"",
		)

		b2 := New(s2, r, 4)
		if err := b2.Build(); err == nil {
			t.Error("Expected error for non-existent content directory")
		}
	})
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestBuilder_processMarkdownFile(t *testing.T) {
	s, r, _ := setupTestSite(t)
	b := New(s, r, 4)

	t.Run("process valid markdown", func(t *testing.T) {
		// Create a test markdown file
		testContent := `---
title: "Test Page"
date: 2023-01-01
---
# Test Content`

		tmpFile := filepath.Join(t.TempDir(), "test.md")
		if err := os.WriteFile(tmpFile, []byte(testContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Process it
		if err := b.processMarkdownFile(tmpFile); err != nil {
			t.Fatalf("processMarkdownFile() error = %v, want nil", err)
		}
	})

	t.Run("process draft file", func(t *testing.T) {
		draftContent := `---
title: "Draft"
draft: true
---
Draft content`

		tmpFile := filepath.Join(t.TempDir(), "draft.md")
		if err := os.WriteFile(tmpFile, []byte(draftContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Process should succeed but page should be marked as draft
		if err := b.processMarkdownFile(tmpFile); err != nil {
			t.Fatal(err)
		}
	})
}

func TestBuilder_copyStaticFiles(t *testing.T) {
	s, r, tmpDir := setupTestSite(t)
	b := New(s, r, 4)

	// Create a test static directory structure
	staticDir := filepath.Join(tmpDir, "test-static")
	if err := os.MkdirAll(filepath.Join(staticDir, "css"), 0755); err != nil {
		t.Fatal(err)
	}

	// Create files
	files := []struct {
		path    string
		content string
	}{
		{filepath.Join(staticDir, "css", "style.css"), "body { color: blue; }"},
		{filepath.Join(staticDir, "js", "app.js"), "console.log('test');"},
		{filepath.Join(staticDir, "image.jpg"), "fake image data"},
	}

	for _, file := range files {
		// Create directory
		dir := filepath.Dir(file.path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(file.path, []byte(file.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Set static directory
	s.StaticDir = staticDir

	t.Run("copy static files", func(t *testing.T) {
		if err := b.copyStaticFiles(); err != nil {
			t.Fatalf("copyStaticFiles() error = %v, want nil", err)
		}

		// Check files were copied
		for _, file := range files {
			relPath, err := filepath.Rel(staticDir, file.path)
			if err != nil {
				t.Fatal(err)
			}

			dstPath := filepath.Join(s.OutputDir, relPath)
			content, err := os.ReadFile(dstPath)
			if err != nil {
				t.Errorf("Static file not copied: %s", relPath)
				continue
			}

			if string(content) != file.content {
				t.Errorf("File content mismatch for %s", relPath)
			}
		}
	})

	t.Run("non-existent static directory", func(t *testing.T) {
		s.StaticDir = "/non/existent/static"

		// Should not error, just log
		if err := b.copyStaticFiles(); err != nil {
			t.Errorf("copyStaticFiles() error = %v, want nil for non-existent dir", err)
		}
	})
}
