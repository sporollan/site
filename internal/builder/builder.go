package builder

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/sporollan/site/internal/renderer"
	"github.com/sporollan/site/internal/site"
)

type Builder struct {
	site     *site.Site
	renderer *renderer.Renderer
	workers  int
}

func New(s *site.Site, r *renderer.Renderer, workers int) *Builder {
	if workers < 1 {
		workers = 1
	}
	return &Builder{
		site:     s,
		renderer: r,
		workers:  workers,
	}
}
func (b *Builder) Build() error {
	// copy static
	_ = os.RemoveAll(b.site.OutputDir)
	if err := copyDir(b.site.StaticDir, b.site.OutputDir); err != nil {
		return err
	}

	jobs := make(chan string)
	// start workers
	var wg sync.WaitGroup
	for i := 0; i < b.workers; i++ {
		wg.Go(func() {
			b.worker(jobs)
		})
	}

	// walk input directory
	err := filepath.WalkDir(b.site.InputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		jobs <- path
		return nil
	})

	close(jobs)
	wg.Wait()
	return err
}
