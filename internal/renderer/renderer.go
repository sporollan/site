package renderer

import (
	"bytes"
	"text/template"

	"github.com/sporollan/site/internal/site"
)

type Renderer struct {
	tmpl *template.Template
}

func New(templatePath string) (*Renderer, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &Renderer{tmpl: tmpl}, nil
}

func (r *Renderer) Render(p site.Page) ([]byte, error) {
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
