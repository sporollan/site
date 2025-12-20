package builder

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sporollan/site/internal/parser"
)

func (b *Builder) worker(jobs <-chan string) {
	for path := range jobs {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		page, err := parser.Parse(path, data)
		if err != nil {
			continue
		}

		html, err := b.renderer.Render(page)
		if err != nil {
			continue
		}

		rel, err := filepath.Rel(b.site.InputDir, path)
		if err != nil {
			continue
		}

		outPath := filepath.Join(
			b.site.OutputDir,
			strings.TrimSuffix(rel, filepath.Ext(rel))+".html",
		)

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			log.Printf("mkdir failed for %s: %v", outPath, err)
			continue
		}

		if err := os.WriteFile(outPath, html, 0644); err != nil {
			log.Printf("write failed for %s: %v", outPath, err)
		}

	}
}
