package site

import "testing"

func TestNewSite(t *testing.T) {
	s := New("content", "public", "static")

	if s.InputDir != "content" {
		t.Fatalf("InputDir = %q, want %q", s.InputDir, "content")
	}

	if s.OutputDir != "public" {
		t.Fatalf("OutputDir = %q, want %q", s.OutputDir, "public")
	}

	if s.StaticDir != "static" {
		t.Fatalf("StaticDir = %q, want %q", s.StaticDir, "static")
	}
}
