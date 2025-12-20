package main

import (
	"log"

	"github.com/sporollan/site/internal/builder"
	"github.com/sporollan/site/internal/renderer"
	"github.com/sporollan/site/internal/site"
)

func run() error {
	s := site.New("content", "public", "static")
	r, err := renderer.New("templates/page.html")
	if err != nil {
		log.Fatal(err)
	}

	b := builder.New(s, r, 4)

	return b.Build()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
