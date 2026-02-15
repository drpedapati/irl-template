package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/drpedapati/irl-template/pkg/config"
)

// Create sets up a minimal IRL project - just the plan file
// Setup is defined in the plan and created by the AI on first run
func Create(projectPath string) error {
	// Create .gitignore
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// WritePlan writes the main-plan.md file in the plans/ folder
func WritePlan(projectPath, content string) error {
	// Create plans/ directory
	plansDir := filepath.Join(projectPath, "plans")
	if err := os.MkdirAll(plansDir, 0755); err != nil {
		return fmt.Errorf("failed to create plans directory: %w", err)
	}

	planPath := filepath.Join(plansDir, "main-plan.md")
	return os.WriteFile(planPath, []byte(content), 0644)
}

// GitInit initializes a git repository with initial commit
func GitInit(projectPath string) error {
	commands := [][]string{
		{"git", "init", "-q"},
		{"git", "add", "-A"},
		{"git", "commit", "-q", "-m", "Initial commit from IRL"},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = projectPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git %s failed: %w", args[1], err)
		}
	}

	return nil
}

// InjectProfile adds profile information (author, affiliation, AI instructions)
// as YAML front matter to plan content. Returns content unchanged if no profile is set.
func InjectProfile(content string) string {
	profile := config.GetProfile()

	// If no profile set, return content unchanged
	if profile.Name == "" && profile.Institution == "" && profile.Instructions == "" {
		return content
	}

	var header strings.Builder

	// Build author/affiliation block
	if profile.Name != "" || profile.Institution != "" {
		header.WriteString("---\n")
		if profile.Name != "" {
			header.WriteString("author: " + profile.Name)
			if profile.Title != "" {
				header.WriteString(", " + profile.Title)
			}
			header.WriteString("\n")
		}
		if profile.Institution != "" {
			header.WriteString("affiliation: " + profile.Institution)
			if profile.Department != "" {
				header.WriteString(", " + profile.Department)
			}
			header.WriteString("\n")
		}
		if profile.Email != "" {
			header.WriteString("email: " + profile.Email + "\n")
		}
		header.WriteString("---\n\n")
	}

	// Add AI instructions as a comment block if set
	if profile.Instructions != "" {
		header.WriteString("<!-- AI Instructions:\n")
		header.WriteString(profile.Instructions)
		header.WriteString("\n-->\n\n")
	}

	return header.String() + content
}

var gitignoreContent = `# OS
.DS_Store
Thumbs.db

# IDE
.vscode/
.idea/
*.swp

# Python
__pycache__/
*.pyc
.venv/
venv/

# R
.Rhistory
.RData

# Quarto
/.quarto/
*_files/

# Node
node_modules/

# Environment
.env
.env.local
`
