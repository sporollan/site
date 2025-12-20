package builder

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "a/b/dst.txt")

	writeFile(t, src, "hello world")

	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	if got := readFile(t, dst); got != "hello world" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestCopyDir(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src")
	dst := filepath.Join(tmp, "dst")

	writeFile(t, filepath.Join(src, "a.txt"), "A")
	writeFile(t, filepath.Join(src, "nested/b.txt"), "B")

	if err := copyDir(src, dst); err != nil {
		t.Fatalf("copyDir failed: %v", err)
	}

	if got := readFile(t, filepath.Join(dst, "a.txt")); got != "A" {
		t.Fatalf("unexpected content: %q", got)
	}

	if got := readFile(t, filepath.Join(dst, "nested/b.txt")); got != "B" {
		t.Fatalf("unexpected content: %q", got)
	}
}
