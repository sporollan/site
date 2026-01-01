package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/sporollan/site/internal/site"
)

type Renderer struct {
	templates *template.Template
}

func New(templateDir string) (*Renderer, error) {
	// Parse all templates with a common name
	pattern := filepath.Join(templateDir, "*.html")

	// Create a template with functions first
	tmpl := template.New("").Funcs(template.FuncMap{
		"now": func() time.Time { return time.Now() },
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"first": func(n int, pages []*site.Page) []*site.Page {
			if n > len(pages) {
				n = len(pages)
			}
			return pages[:n]
		},
	})

	// Parse all template files
	tmpl, err := tmpl.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Debug: list all templates
	fmt.Printf("Loaded %d templates:\n", len(tmpl.Templates()))
	for _, t := range tmpl.Templates() {
		fmt.Printf("  - %s\n", t.Name())
	}

	return &Renderer{templates: tmpl}, nil
}

func (r *Renderer) Render(p site.Page) ([]byte, error) {
	var buf bytes.Buffer

	// Determine which template to use
	tmplName := p.TemplateName
	if tmplName == "" {
		tmplName = "page.html"
	}

	fmt.Printf("Rendering with template: %s\n", tmplName)

	// Look up the template
	tmpl := r.templates.Lookup(tmplName)
	if tmpl == nil {
		// Fallback to base.html if available
		tmpl = r.templates.Lookup("base.html")
		if tmpl == nil {
			return nil, fmt.Errorf("template %s not found and no base.html available", tmplName)
		}
		fmt.Printf("Falling back to base.html\n")
	}

	// Execute the template
	if err := tmpl.Execute(&buf, p); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}
