package main

import (
	"log"
	"os"
	"strings"

	"github.com/sporollan/site/internal/builder"
	"github.com/sporollan/site/internal/renderer"
	"github.com/sporollan/site/internal/site"
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func run() error {
	// Load configuration from environment variables
	inputDir := getEnv("SITE_INPUT_DIR", "content")
	outputDir := getEnv("SITE_OUTPUT_DIR", "public")
	staticDir := getEnv("SITE_STATIC_DIR", "static")
	templateDir := getEnv("SITE_TEMPLATE_DIR", "templates")
	siteName := getEnv("SITE_NAME", "Santiago Porollan")
	baseURL := getEnv("SITE_BASE_URL", "http://localhost:8080")
	
	// Ensure base URL has proper protocol and no trailing slash
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}
	baseURL = strings.TrimSuffix(baseURL, "/")
	
	// Detect if we're in production (GitHub Actions sets GITHUB_ACTIONS=true)
	isProduction := os.Getenv("GITHUB_ACTIONS") == "true"
	
	if isProduction {
		log.Printf("Production build for: %s", baseURL)
	} else {
		log.Printf("Development build for: %s", baseURL)
	}
	
	// Initialize site
	s := site.NewWithConfig(
		inputDir,
		outputDir,
		staticDir,
		templateDir,
		siteName,
		baseURL,
	)
	
	// Create renderer
	r, err := renderer.New(s.TemplateDir)
	if err != nil {
		return err
	}
	
	// Create builder and build site
	b := builder.New(s, r, 4)
	return b.Build()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}