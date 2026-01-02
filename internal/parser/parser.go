package parser

import (
	"bytes"
	"path/filepath"
	"strings"
	"time"

	"github.com/sporollan/site/internal/site"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v2"
)

var markdownConverter = goldmark.New(
	goldmark.WithExtensions(),
)

// parseFrontMatter splits YAML front matter from markdown content
func parseFrontMatter(data []byte) (map[string]interface{}, []byte, error) {
	content := string(data)

	// Check if file starts with front matter
	if !strings.HasPrefix(content, "---") {
		return nil, data, nil
	}

	// Find the second "---"
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, data, nil
	}

	frontMatterStr := parts[1]
	markdownContent := []byte(strings.TrimSpace(parts[2]))

	metadata := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(frontMatterStr), &metadata)
	if err != nil {
		return nil, markdownContent, err
	}

	return metadata, markdownContent, nil
}

func Parse(path string, data []byte) (site.Page, error) {
	// Parse front matter
	metadata, markdownContent, err := parseFrontMatter(data)
	if err != nil {
		return site.Page{}, err
	}

	// Convert markdown to HTML
	var htmlBuf bytes.Buffer
	if err := markdownConverter.Convert(markdownContent, &htmlBuf); err != nil {
		return site.Page{}, err
	}

	// Extract title (from front matter or filename)
	baseName := filepath.Base(path)
	title := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	title = strings.Title(strings.ReplaceAll(title, "-", " ")) // "about-me" -> "About Me"

	if mdTitle, ok := metadata["title"].(string); ok && mdTitle != "" {
		title = mdTitle
	}

	// Extract template name
	templateName := "page.html" // default
	if tmpl, ok := metadata["template"].(string); ok && tmpl != "" {
		templateName = tmpl
	}

	// Extract date if present
	var pageDate time.Time
	if dateStr, ok := metadata["date"].(string); ok {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			pageDate = date
		}
	}

	// Extract other metadata
	tags := []string{}
	if tagsInterface, ok := metadata["tags"].([]interface{}); ok {
		for _, tag := range tagsInterface {
			if tagStr, ok := tag.(string); ok {
				tags = append(tags, tagStr)
			}
		}
	}

	// Extract draft status
	draft := false
	if draftVal, ok := metadata["draft"].(bool); ok {
		draft = draftVal
	}

	return site.Page{
		Path:         path,
		Title:        title,
		Body:         htmlBuf.String(),
		RawBody:      string(markdownContent),
		TemplateName: templateName,
		Date:         pageDate,
		Draft:        draft,
		Tags:         tags,
		Metadata:     metadata,
	}, nil
}
