package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Create sets up a minimal IRL project - just the plan file
// Structure is defined in the plan and created by the AI on first run
func Create(projectPath string) error {
	// Create .gitignore
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// WritePlan writes the main-plan.md file at the project root
func WritePlan(projectPath, content string) error {
	planPath := filepath.Join(projectPath, "main-plan.md")
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
