package parser

import (
	"path/filepath"
	"strings"

	"github.com/sporollan/site/internal/site"
)

func Parse(path string, data []byte) (site.Page, error) {
	name := filepath.Base(path)
	title := strings.TrimSuffix(name, filepath.Ext(name))

	return site.Page{
		Path:  path,
		Title: title,
		Body:  string(data),
	}, nil
}
