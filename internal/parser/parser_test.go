package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	data := []byte("hello")
	page, err := Parse("index.txt", data)
	if err != nil {
		t.Fatal(err)
	}

	if page.Path != "index.txt" {
		t.Fatalf("Path = %q", page.Path)
	}

	if page.Body != "hello" {
		t.Fatalf("Body = %q", page.Body)
	}
}
