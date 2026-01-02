package site

import (
	"testing"
	"time"
)

func TestNewWithConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		output   string
		static   string
		template string
		siteName string
		baseURL  string
		want     *Site
	}{
		{
			name:     "basic configuration",
			input:    "content",
			output:   "public",
			static:   "static",
			template: "templates",
			siteName: "Test Site",
			baseURL:  "https://example.com",
			want: &Site{
				InputDir:    "content",
				OutputDir:   "public",
				StaticDir:   "static",
				TemplateDir: "templates",
				SiteName:    "Test Site",
				BaseURL:     "https://example.com",
				Collections: make(map[string][]*Page),
			},
		},
		{
			name:     "empty base URL",
			input:    "content",
			output:   "public",
			static:   "static",
			template: "templates",
			siteName: "Test Site",
			baseURL:  "",
			want: &Site{
				InputDir:    "content",
				OutputDir:   "public",
				StaticDir:   "static",
				TemplateDir: "templates",
				SiteName:    "Test Site",
				BaseURL:     "",
				Collections: make(map[string][]*Page),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWithConfig(
				tt.input,
				tt.output,
				tt.static,
				tt.template,
				tt.siteName,
				tt.baseURL,
			)

			if got.InputDir != tt.want.InputDir {
				t.Errorf("InputDir = %v, want %v", got.InputDir, tt.want.InputDir)
			}
			if got.OutputDir != tt.want.OutputDir {
				t.Errorf("OutputDir = %v, want %v", got.OutputDir, tt.want.OutputDir)
			}
			if got.StaticDir != tt.want.StaticDir {
				t.Errorf("StaticDir = %v, want %v", got.StaticDir, tt.want.StaticDir)
			}
			if got.TemplateDir != tt.want.TemplateDir {
				t.Errorf("TemplateDir = %v, want %v", got.TemplateDir, tt.want.TemplateDir)
			}
			if got.SiteName != tt.want.SiteName {
				t.Errorf("SiteName = %v, want %v", got.SiteName, tt.want.SiteName)
			}
			if got.BaseURL != tt.want.BaseURL {
				t.Errorf("BaseURL = %v, want %v", got.BaseURL, tt.want.BaseURL)
			}
			if got.Collections == nil {
				t.Error("Collections map not initialized")
			}
		})
	}
}

func TestPageMethods(t *testing.T) {
	now := time.Now()
	page := Page{
		Path:         "content/test.md",
		Title:        "Test Page",
		Description:  "A test description",
		Date:         now,
		Tags:         []string{"test", "go"},
		Draft:        false,
		TemplateName: "post.html",
		SiteName:     "Test Site",
		BaseURL:      "https://example.com",
	}

	t.Run("Page fields", func(t *testing.T) {
		if page.Title != "Test Page" {
			t.Errorf("Title = %v, want 'Test Page'", page.Title)
		}
		if page.Description != "A test description" {
			t.Errorf("Description = %v, want 'A test description'", page.Description)
		}
		if len(page.Tags) != 2 {
			t.Errorf("Tags length = %v, want 2", len(page.Tags))
		}
		if page.TemplateName != "post.html" {
			t.Errorf("TemplateName = %v, want 'post.html'", page.TemplateName)
		}
	})

	t.Run("Page date operations", func(t *testing.T) {
		if page.Date.IsZero() {
			t.Error("Date should not be zero")
		}
	})

	t.Run("Page draft status", func(t *testing.T) {
		if page.Draft {
			t.Error("Draft should be false")
		}
	})
}

func TestSiteCollections(t *testing.T) {
	site := NewWithConfig("content", "public", "static", "templates", "Test Site", "")

	page1 := &Page{Title: "Page 1"}
	page2 := &Page{Title: "Page 2"}

	site.Collections["pages"] = []*Page{page1, page2}
	site.Pages = []*Page{page1, page2}

	t.Run("Collection length", func(t *testing.T) {
		if len(site.Collections["pages"]) != 2 {
			t.Errorf("Collection length = %v, want 2", len(site.Collections["pages"]))
		}
	})

	t.Run("Site pages length", func(t *testing.T) {
		if len(site.Pages) != 2 {
			t.Errorf("Pages length = %v, want 2", len(site.Pages))
		}
	})

	t.Run("Add to collection", func(t *testing.T) {
		page3 := &Page{Title: "Page 3"}
		site.Collections["pages"] = append(site.Collections["pages"], page3)

		if len(site.Collections["pages"]) != 3 {
			t.Errorf("After add: Collection length = %v, want 3", len(site.Collections["pages"]))
		}
	})
}
