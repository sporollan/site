package renderer

import (
	"strings"
	"testing"

	"github.com/sporollan/site/internal/site"
)

func TestRender(t *testing.T) {
	page := site.Page{
		Body: "hello",
	}

	r, err := New("../../templates/page.html")

	out, err := r.Render(page)
	if err != nil {
		t.Fatal(err)
	}

	html := string(out)

	if !strings.Contains(html, "<body>") {
		t.Fatal("missing <body> tag")
	}

	if !strings.Contains(html, "hello") {
		t.Fatal("missing body content")
	}
}

func TestRenderTemplate(t *testing.T) {
	r, err := New("../../templates/page.html")
	if err != nil {
		t.Fatal(err)
	}

	page := site.Page{
		Title: "Hello",
		Body:  "World",
	}

	out, err := r.Render(page)
	if err != nil {
		t.Fatal(err)
	}

	html := string(out)

	if !strings.Contains(html, "<title>Hello</title>") {
		t.Fatal("title not rendered")
	}

	if !strings.Contains(html, "World") {
		t.Fatal("body not rendered")
	}
}
