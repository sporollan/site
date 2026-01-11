package builder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sporollan/site/internal/parser"
	"github.com/sporollan/site/internal/renderer"
	"github.com/sporollan/site/internal/site"
)

type Builder struct {
	site     *site.Site
	renderer *renderer.Renderer
	workers  int
}

func New(s *site.Site, r *renderer.Renderer, workers int) *Builder {
	return &Builder{
		site:     s,
		renderer: r,
		workers:  workers,
	}
}

func (b *Builder) Build() error {
	// Clean output directory
	if err := b.cleanOutputDir(); err != nil {
		return fmt.Errorf("failed to clean output: %w", err)
	}

	// Reset site data
	b.site.Pages = []*site.Page{}
	b.site.Collections = make(map[string][]*site.Page)

	// Process all content files
	if err := b.processContent(); err != nil {
		return err
	}

	// Copy static files
	if err := b.copyStaticFiles(); err != nil {
		return err
	}

	// Generate index/archive pages
	if err := b.generateIndexPages(); err != nil {
		return err
	}

	// Update home page with recent posts
	if err := b.generateHomePage(); err != nil {
		return err
	}

	log.Printf("Build complete! Generated %d pages", len(b.site.Pages))
	return nil
}

func (b *Builder) cleanOutputDir() error {
	// Skip if output dir doesn't exist
	if _, err := os.Stat(b.site.OutputDir); os.IsNotExist(err) {
		return os.MkdirAll(b.site.OutputDir, 0755)
	}

	log.Printf("Cleaning output directory: %s", b.site.OutputDir)

	entries, err := os.ReadDir(b.site.OutputDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// Preserve .git directory if it exists
		if entry.Name() == ".git" {
			continue
		}

		path := filepath.Join(b.site.OutputDir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}

	return nil
}

func (b *Builder) processContent() error {
	return filepath.Walk(b.site.InputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Only process markdown files
		if strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			return b.processMarkdownFile(path)
		}

		return nil
	})
}

func (b *Builder) processMarkdownFile(path string) error {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	// Parse markdown
	page, err := parser.Parse(path, data)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	// Skip drafts
	if page.Draft {
		log.Printf("Skipping draft: %s", page.Title)
		return nil
	}

	// Calculate output path
	relPath, err := filepath.Rel(b.site.InputDir, path)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	// Create clean URL structure
	baseName := strings.TrimSuffix(relPath, ".md")

	var outputDir, outputPath, permalink string
	if filepath.Base(baseName) == "index" {
		// Root index
		outputDir = b.site.OutputDir
		outputPath = filepath.Join(outputDir, "index.html")
		permalink = "/"
	} else {
		// Regular page
		outputDir = filepath.Join(b.site.OutputDir, baseName)
		outputPath = filepath.Join(outputDir, "index.html")
		permalink = "/" + baseName + "/"
	}

	// Set page metadata
	page.Permalink = permalink
	page.SiteName = b.site.SiteName
	page.BaseURL = b.site.BaseURL

	// Determine template based on directory
	dir := filepath.Dir(relPath)
	if page.TemplateName == "" {
		switch dir {
		case "blog":
			page.TemplateName = "post.html"
		default:
			page.TemplateName = "page.html"
		}
	}

	// Render page
	html, err := b.renderer.Render(page)
	if err != nil {
		return fmt.Errorf("failed to render %s: %w", path, err)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	// Write HTML file
	if err := os.WriteFile(outputPath, html, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}

	// Store the page
	b.site.Pages = append(b.site.Pages, &page)

	// Add to collections based on directory
	if dir == "." {
		dir = "pages"
	}
	b.site.Collections[dir] = append(b.site.Collections[dir], &page)

	log.Printf("Generated: %s", outputPath)
	return nil
}

func (b *Builder) generateIndexPages() error {
	// Generate blog index page if we have blog posts
	if posts, exists := b.site.Collections["blog"]; exists && len(posts) > 0 {
		// Sort posts by date (newest first)
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Date.After(posts[j].Date)
		})

		// Generate summaries for posts (first 150 chars of raw markdown)
		for _, post := range posts {
			if len(post.RawBody) > 150 {
				post.Summary = strings.TrimSpace(post.RawBody[:150]) + "..."
			} else {
				post.Summary = strings.TrimSpace(post.RawBody)
			}
		}

		// Create blog index page
		blogIndex := site.Page{
			Title:        "Blog",
			Body:         "", // Not used for list template
			TemplateName: "list.html",
			SiteName:     b.site.SiteName,
			BaseURL:      b.site.BaseURL,
			Permalink:    "/blog/",
			Pages:        posts, // Pass posts to the template
		}

		// Render and write blog index
		html, err := b.renderer.Render(blogIndex)
		if err != nil {
			return fmt.Errorf("failed to render blog index: %w", err)
		}

		blogDir := filepath.Join(b.site.OutputDir, "blog")
		if err := os.MkdirAll(blogDir, 0755); err != nil {
			return fmt.Errorf("failed to create blog directory: %w", err)
		}

		blogIndexPath := filepath.Join(blogDir, "index.html")
		if err := os.WriteFile(blogIndexPath, html, 0644); err != nil {
			return fmt.Errorf("failed to write blog index: %w", err)
		}

		log.Printf("Generated blog index with %d posts: %s", len(posts), blogIndexPath)
	}

	return nil
}

func (b *Builder) copyStaticFiles() error {
	// Check if static directory exists
	if _, err := os.Stat(b.site.StaticDir); os.IsNotExist(err) {
		log.Printf("Note: Static directory %s does not exist", b.site.StaticDir)
		return nil
	}

	log.Printf("Copying static files from %s", b.site.StaticDir)

	return filepath.Walk(b.site.StaticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(b.site.StaticDir, path)
		if err != nil {
			return err
		}

		// Destination path
		destPath := filepath.Join(b.site.OutputDir, relPath)

		// Create destination directory
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return err
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return err
		}

		log.Printf("Copied static file: %s", destPath)
		return nil
	})
}

func (b *Builder) generateHomePage() error {
	// Find the home page (permalink "/")
	var homePage *site.Page
	for _, page := range b.site.Pages {
		if page.Permalink == "/" {
			homePage = page
			break
		}
	}

	if homePage == nil {
		// No home page found, nothing to do
		return nil
	}

	// Get recent blog posts (max 5)
	if posts, exists := b.site.Collections["blog"]; exists && len(posts) > 0 {
		// Posts are already sorted newest first in generateIndexPages()
		recentCount := 5
		if len(posts) < recentCount {
			recentCount = len(posts)
		}
		recentPosts := posts[:recentCount]

		// Initialize Metadata map if nil
		if homePage.Metadata == nil {
			homePage.Metadata = make(map[string]interface{})
		}
		homePage.Metadata["RecentPosts"] = recentPosts

		// Re-render the home page with updated metadata
		html, err := b.renderer.Render(*homePage)
		if err != nil {
			return fmt.Errorf("failed to re-render home page: %w", err)
		}

		// Write updated home page
		outputPath := filepath.Join(b.site.OutputDir, "index.html")
		if err := os.WriteFile(outputPath, html, 0644); err != nil {
			return fmt.Errorf("failed to write updated home page: %w", err)
		}

		log.Printf("Updated home page with %d recent posts", len(recentPosts))
	}

	return nil
}
