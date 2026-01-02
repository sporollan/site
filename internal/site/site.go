package site

import "time"

type Page struct {
	Path         string
	Permalink    string
	Title        string
	Body         string
	RawBody      string
	TemplateName string
	Date         time.Time
	Draft        bool
	Tags         []string
	Categories   []string
	Summary      string
	Description  string
	Metadata     map[string]interface{}

	// For lists
	Pages []*Page

	// Site context
	SiteName string
	BaseURL  string
	Language string
}

type Site struct {
	InputDir    string
	OutputDir   string
	StaticDir   string
	TemplateDir string
	SiteName    string
	BaseURL     string
	Pages       []*Page
	Collections map[string][]*Page // "posts", "pages", etc.
}

func NewWithConfig(input, output, static, templateDir, siteName, baseURL string) *Site {
	return &Site{
		InputDir:    input,
		OutputDir:   output,
		StaticDir:   static,
		TemplateDir: templateDir,
		SiteName:    siteName,
		BaseURL:     baseURL,
		Collections: make(map[string][]*Page),
	}
}
