package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sporollan/site/internal/renderer"
	"github.com/sporollan/site/internal/site"
)

func TestBuildCreatesHTML(t *testing.T) {
	input := t.TempDir()
	output := t.TempDir()
	static := t.TempDir()

	err := os.WriteFile(
		filepath.Join(input, "index.txt"),
		[]byte("hello"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	s := site.New(input, output, static)
	r, err := renderer.New("../../templates/page.html")
	if err != nil {
		t.Fatal(err)
	}
	b := New(s, r, 2)

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(output, "index.html")
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected output file %s to exist", out)
	}
}
