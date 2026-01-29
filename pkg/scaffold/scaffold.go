package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Directories defines the IRL project structure
var Directories = []string{
	"01-plans",
	"02-data/raw",
	"02-data/derived",
	"03-outputs/figures",
	"04-logs",
}

// Create sets up a new IRL project directory structure
func Create(projectPath string) error {
	for _, dir := range Directories {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	// Create .gitkeep files in empty directories
	gitkeeps := []string{
		"02-data/raw/.gitkeep",
		"02-data/derived/.gitkeep",
		"03-outputs/.gitkeep",
		"03-outputs/figures/.gitkeep",
		"04-logs/.gitkeep",
	}

	for _, gk := range gitkeeps {
		path := filepath.Join(projectPath, gk)
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", gk, err)
		}
	}

	// Create default activity log
	logPath := filepath.Join(projectPath, "04-logs", "activity_log.md")
	logContent := "# Activity Log\n\n## Iteration History\n\n"
	if err := os.WriteFile(logPath, []byte(logContent), 0644); err != nil {
		return fmt.Errorf("failed to create activity log: %w", err)
	}

	// Create .gitignore
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// WritePlan writes the main-plan.md file
func WritePlan(projectPath, content string) error {
	planPath := filepath.Join(projectPath, "01-plans", "main-plan.md")
	return os.WriteFile(planPath, []byte(content), 0644)
}

// GitInit initializes a git repository with initial commit
func GitInit(projectPath string) error {
	commands := [][]string{
		{"git", "init", "-q"},
		{"git", "add", "-A"},
		{"git", "commit", "-q", "-m", "Initial commit from IRL template"},
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

# Environment
.env
.env.local
`
