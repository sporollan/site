package renderer

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sporollan/site/internal/site"
)

func setupTestTemplates(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Create templates without inheritance for testing
	templates := []struct {
		name    string
		content string
	}{
		{
			name: "page.html",
			content: `<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
	<h1>{{.Title}}</h1>
	{{.Body | safeHTML}}
</body>
</html>`,
		},
		{
			name: "post.html",
			content: `<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
	<article>
		<h1>{{.Title}}</h1>
		{{if not .Date.IsZero}}<time>{{.Date.Format "2006-01-02"}}</time>{{end}}
		{{.Body | safeHTML}}
	</article>
</body>
</html>`,
		},
	}

	for _, tmpl := range templates {
		path := filepath.Join(tmpDir, tmpl.name)
		if err := os.WriteFile(path, []byte(tmpl.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return tmpDir
}

func TestNew(t *testing.T) {
	t.Run("successful renderer creation", func(t *testing.T) {
		tmpDir := setupTestTemplates(t)

		r, err := New(tmpDir)
		if err != nil {
			t.Fatalf("New() error = %v, want nil", err)
		}

		if r.templates == nil {
			t.Error("templates should not be nil")
		}
	})

	t.Run("missing template directory", func(t *testing.T) {
		_, err := New("/non/existent/directory")
		if err == nil {
			t.Error("Expected error for missing directory")
		}
	})
}

func TestRender(t *testing.T) {
	tmpDir := setupTestTemplates(t)
	r, err := New(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	tests := []struct {
		name     string
		page     site.Page
		wantErr  bool
		contains []string
	}{
		{
			name: "basic page",
			page: site.Page{
				Title:        "Test Page",
				Body:         "<p>Test content</p>",
				TemplateName: "page.html",
				SiteName:     "Test Site",
			},
			wantErr:  false,
			contains: []string{"Test Page", "Test content", "<!DOCTYPE html>"},
		},
		{
			name: "post with date",
			page: site.Page{
				Title:        "Blog Post",
				Body:         "<p>Post content</p>",
				TemplateName: "post.html",
				Date:         time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				SiteName:     "Test Site",
			},
			wantErr:  false,
			contains: []string{"Blog Post", "2023-10-01", "Post content"},
		},
		{
			name: "default template when not specified",
			page: site.Page{
				Title:    "Default Page",
				Body:     "<p>Default content</p>",
				SiteName: "Test Site",
			},
			wantErr:  false,
			contains: []string{"Default Page"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Render(tt.page)

			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Check for required content
			for _, substr := range tt.contains {
				if !bytes.Contains(got, []byte(substr)) {
					t.Errorf("Output does not contain %q", substr)
				}
			}
		})
	}
}
