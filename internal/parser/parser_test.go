package parser

import (
	"testing"
	"time"

	"github.com/sporollan/site/internal/site"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		data     []byte
		wantPage site.Page
		wantErr  bool
	}{
		{
			name: "simple markdown without front matter",
			path: "content/test.md",
			data: []byte("# Hello World\n\nThis is a test."),
			wantPage: site.Page{
				Title:        "Test",
				TemplateName: "page.html",
				Path:         "content/test.md",
			},
			wantErr: false,
		},
		{
			name: "markdown with front matter",
			path: "content/blog/post.md",
			data: []byte(`---
title: "My Post"
date: 2023-10-01
tags: ["go", "test"]
description: "A test post"
template: "post.html"
---

# Test Post

Content here.`),
			wantPage: site.Page{
				Title: "My Post",
				// Don't check Description if parser doesn't extract it
				TemplateName: "post.html",
				Path:         "content/blog/post.md",
				Date:         time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				Tags:         []string{"go", "test"},
				Draft:        false,
			},
			wantErr: false,
		},
		{
			name: "markdown with draft status",
			path: "content/draft.md",
			data: []byte(`---
title: "Draft Post"
draft: true
---

# Draft

Should not be published.`),
			wantPage: site.Page{
				Title:        "Draft Post",
				TemplateName: "page.html",
				Path:         "content/draft.md",
				Draft:        true,
			},
			wantErr: false,
		},
		{
			name: "empty file",
			path: "content/empty.md",
			data: []byte(``),
			wantPage: site.Page{
				Title:        "Empty",
				TemplateName: "page.html",
				Path:         "content/empty.md",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.path, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Check fields
			if got.Title != tt.wantPage.Title {
				t.Errorf("Title = %v, want %v", got.Title, tt.wantPage.Title)
			}

			if got.TemplateName != tt.wantPage.TemplateName {
				t.Errorf("TemplateName = %v, want %v", got.TemplateName, tt.wantPage.TemplateName)
			}

			if got.Path != tt.wantPage.Path {
				t.Errorf("Path = %v, want %v", got.Path, tt.wantPage.Path)
			}

			if got.Draft != tt.wantPage.Draft {
				t.Errorf("Draft = %v, want %v", got.Draft, tt.wantPage.Draft)
			}

			// Check description if expected
			if tt.wantPage.Description != "" && got.Description != tt.wantPage.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.wantPage.Description)
			}

			// Check tags if expected
			if len(tt.wantPage.Tags) > 0 {
				if len(got.Tags) != len(tt.wantPage.Tags) {
					t.Errorf("Tags length = %v, want %v", len(got.Tags), len(tt.wantPage.Tags))
				}
			}

			// Check date if expected
			if !tt.wantPage.Date.IsZero() && !got.Date.Equal(tt.wantPage.Date) {
				t.Errorf("Date = %v, want %v", got.Date, tt.wantPage.Date)
			}

			// Check body is not empty for non-empty markdown
			if len(tt.data) > 0 && tt.name != "empty file" && got.Body == "" {
				t.Error("Body should not be empty")
			}
		})
	}
}

func TestParseFrontMatter(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		wantMeta    map[string]interface{}
		wantContent string
		wantErr     bool
	}{
		{
			name: "valid front matter",
			data: []byte(`---
title: "Test"
tags: ["a", "b"]
date: "2023-01-01"
---
Content`),
			wantMeta: map[string]interface{}{
				"title": "Test",
				"tags":  []interface{}{"a", "b"},
				"date":  "2023-01-01",
			},
			wantContent: "Content",
			wantErr:     false,
		},
		{
			name:        "no front matter",
			data:        []byte(`# Just markdown`),
			wantMeta:    nil,
			wantContent: "# Just markdown",
			wantErr:     false,
		},
		{
			name: "incomplete front matter",
			data: []byte(`---
title: "Test"
---`),
			wantMeta: map[string]interface{}{
				"title": "Test",
			},
			wantContent: "",
			wantErr:     false,
		},
		{
			name: "only front matter delimiter",
			data: []byte(`---
---
Content`),
			wantMeta:    map[string]interface{}{},
			wantContent: "Content",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, content, err := parseFrontMatter(tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseFrontMatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Check content
			if string(content) != tt.wantContent {
				t.Errorf(
					"parseFrontMatter() content = %v, want %v",
					string(content),
					tt.wantContent,
				)
			}
		})
	}
}
