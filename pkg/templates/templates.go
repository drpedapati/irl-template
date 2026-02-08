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
	"irl-basic": {
		Name:        "irl-basic",
		Description: "IRL basic template",
		Content:     defaultIrlBasicTemplate,
	},
}

var defaultIrlBasicTemplate = `# IRL Basic Template

## 🔧 First Time Setup — Run once when starting a new project
<!-- 👤 AUTHOR AREA: Define setup steps below -->

### Common Skill Library — Pre-installed tools available to all projects
<!-- Uncomment to use -->
Install Quarto Skill: https://github.com/posit-dev/skills/tree/main/quarto/authoring
<!-- Install Word DOCX Skill: https://github.com/anthropics/skills/tree/main/skills/docx -->

## ✅ Before Each Loop — Checklist to run before every iteration

- Script designed to be **idempotent**
- **Version control**: clean tree (git status), commit baseline, .gitignore enforced
- The **only permitted modification** is within the ## One-Time Instructions section.

---

## 🔁 Instruction Loop — Define the work for each iteration

<!-- 👤 AUTHOR AREA: Define each loop's work below -->

### One-Time Instructions — Tasks that should only execute once

<!-- 👤 AUTHOR AREA: Add one-time tasks below -->

### Formatting Guidelines — Rules for output style and structure

<!-- 👤 AUTHOR AREA: Add formatting rules below -->

---

## 📝 After Each Loop — Steps to complete after every iteration

> Uncomment your preferred options below.

- **Update activity log**
  - Use timestamps **only** when sequencing or causality matters
  - In plans/main-plan-activity.md, write 1–2 lines describing:
    - What you did
    - Timestamp
    - Git hash
  <!-- Optional CSV logging -->
  <!--
  - In plans/main-plan-activity.csv, add 1 row:
    - What you did
    - Timestamp
    - Git hash
  -->

- **Update plan log**: Update plans/main-plan-log.csv

- **Version control**: Commit intended changes only; verify no ignored or unintended files staged.

- **Give feedback to the AUTHOR** — concise and actionable:
  1. What was done, decisions needed, next steps
  2. Identify anything breaking idempotency, or obsolete/outdated instructions
  3. Identify critical reasoning errors

## 📚 Skill Library — Optional community skills to install per project
<!-- Uncomment to use -->

<!-- Install PPTX Posters -->
<!-- https://github.com/K-Dense-AI/claude-scientific-skills/tree/main/scientific-skills/pptx-posters -->

<!-- Install Scientific Writing Skill -->
<!-- https://github.com/K-Dense-AI/claude-scientific-skills/tree/main/scientific-skills/scientific-writing -->

<!-- Install BioRx Search -->
<!-- https://github.com/K-Dense-AI/claude-scientific-skills/tree/main/scientific-skills/biorxiv-database -->

<!-- Install PubMed Search -->
<!-- https://github.com/K-Dense-AI/claude-scientific-skills/tree/main/scientific-skills/pubmed-database -->

<!-- Install Flowcharts -->
<!-- https://github.com/lukilabs/beautiful-mermaid -->

<!-- Install PowerPoint -->
<!-- https://github.com/anthropics/skills/tree/main/skills/pptx -->

<!-- Install PDF -->
<!-- https://github.com/anthropics/skills/tree/main/skills/pdf -->
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

	// Load content from cached files (since Content has json:"-")
	for i := range templates {
		cachePath := filepath.Join(cacheDir, templates[i].Name+".md")
		if content, err := os.ReadFile(cachePath); err == nil {
			templates[i].Content = string(content)
		}
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
		"irl-basic": "IRL basic template",
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

// CheckForNewTemplates returns the count of templates available on GitHub
// that are not in the local cache (without downloading them)
func CheckForNewTemplates() (int, error) {
	// Get cached template names
	cachedNames := make(map[string]bool)
	cacheDir, err := GetCacheDir()
	if err == nil {
		indexPath := filepath.Join(cacheDir, "index.json")
		if data, err := os.ReadFile(indexPath); err == nil {
			var cached []Template
			if json.Unmarshal(data, &cached) == nil {
				for _, t := range cached {
					cachedNames[t.Name] = true
				}
			}
		}
	}

	// Query GitHub for current template list (lightweight - just metadata)
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", GitHubRepo, TemplatesPath)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var files []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return 0, err
	}

	// Count new templates
	newCount := 0
	for _, f := range files {
		if !strings.HasSuffix(f.Name, ".md") {
			continue
		}
		name := strings.TrimSuffix(f.Name, ".md")
		if !cachedNames[name] {
			newCount++
		}
	}

	return newCount, nil
}
