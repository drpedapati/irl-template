package templates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	GitHubRepo    = "drpedapati/irl-template"
	TemplatesPath = "01-plans/templates"
	CacheDir      = ".irl"
)

// Template represents an IRL project template
type Template struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"-"`
}

// Embedded fallback templates (minimal set)
var EmbeddedTemplates = map[string]Template{
	"basic": {
		Name:        "basic",
		Description: "General purpose IRL template",
		Content:     defaultBasicTemplate,
	},
}

var defaultBasicTemplate = `# IRL Plan

## Overall Strategy
- [Describe your approach]

## Instruction Sets

### Task 1: [Description]
- [Instructions]

## Post-Instruction Hooks
- Update logs/activity.md
- Make a surgical git commit
`

// GetCacheDir returns the template cache directory
func GetCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, CacheDir, "templates"), nil
}

// ListTemplates returns available templates (cached or embedded)
func ListTemplates() ([]Template, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return embeddedList(), nil
	}

	// Check if cache exists and is recent (< 24h)
	if templates, ok := loadCachedTemplates(cacheDir); ok {
		return templates, nil
	}

	// Try to fetch from GitHub
	if templates, err := FetchTemplates(); err == nil {
		return templates, nil
	}

	// Fall back to embedded
	return embeddedList(), nil
}

func embeddedList() []Template {
	var list []Template
	for _, t := range EmbeddedTemplates {
		list = append(list, t)
	}
	return list
}

func loadCachedTemplates(cacheDir string) ([]Template, bool) {
	indexPath := filepath.Join(cacheDir, "index.json")
	info, err := os.Stat(indexPath)
	if err != nil {
		return nil, false
	}

	// Check if cache is fresh (< 24 hours)
	if time.Since(info.ModTime()) > 24*time.Hour {
		return nil, false
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, false
	}

	var templates []Template
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, false
	}

	return templates, true
}

// FetchTemplates downloads latest templates from GitHub
func FetchTemplates() ([]Template, error) {
	// Get template list from GitHub API
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", GitHubRepo, TemplatesPath)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var files []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"download_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	cacheDir, _ := GetCacheDir()
	os.MkdirAll(cacheDir, 0755)

	var templates []Template
	for _, f := range files {
		if !strings.HasSuffix(f.Name, ".md") {
			continue
		}

		name := strings.TrimSuffix(f.Name, ".md")
		desc := getDescription(name)

		// Download content
		content, err := downloadFile(f.DownloadURL)
		if err != nil {
			continue
		}

		// Cache it
		cachePath := filepath.Join(cacheDir, f.Name)
		os.WriteFile(cachePath, []byte(content), 0644)

		templates = append(templates, Template{
			Name:        name,
			Description: desc,
			Content:     content,
		})
	}

	// Save index
	indexData, _ := json.Marshal(templates)
	indexPath := filepath.Join(cacheDir, "index.json")
	os.WriteFile(indexPath, indexData, 0644)

	return templates, nil
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getDescription(name string) string {
	descriptions := map[string]string{
		"irl-basic-template":  "General purpose IRL template",
		"scientific-abstract": "For journal article abstracts",
		"meeting-abstract":    "For conference/meeting abstracts",
		"admin-planning":      "Plans, reports, memos, proposals",
		"manuscript":          "Academic manuscript writing",
		"project-tracker":     "Task list and schedule management",
		"coding-project":      "Software development with docs",
	}
	if d, ok := descriptions[name]; ok {
		return d
	}
	return "Template"
}

// GetTemplate returns a specific template by name
func GetTemplate(name string) (Template, error) {
	// Check cache first
	cacheDir, _ := GetCacheDir()
	cachePath := filepath.Join(cacheDir, name+".md")

	if content, err := os.ReadFile(cachePath); err == nil {
		return Template{
			Name:        name,
			Description: getDescription(name),
			Content:     string(content),
		}, nil
	}

	// Check embedded
	if t, ok := EmbeddedTemplates[name]; ok {
		return t, nil
	}

	// Try to fetch
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/main/%s/%s.md",
		GitHubRepo, TemplatesPath, name)

	content, err := downloadFile(url)
	if err != nil {
		return Template{}, fmt.Errorf("template '%s' not found", name)
	}

	return Template{
		Name:        name,
		Description: getDescription(name),
		Content:     content,
	}, nil
}

// Update forces a refresh of templates from GitHub
func Update() error {
	_, err := FetchTemplates()
	return err
}
