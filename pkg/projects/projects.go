package projects

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/drpedapati/irl-template/pkg/config"
)

// Project represents a discovered IRL project in the workspace
type Project struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Modified time.Time `json:"modified"`
}

// Scan discovers IRL projects in the configured workspace directory.
// A folder is considered a project if it contains main-plan.md in one of:
//   - plans/main-plan.md (current standard)
//   - main-plan.md (legacy root level)
//   - 01-plans/main-plan.md (legacy IRL structure)
func Scan() ([]Project, error) {
	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return []Project{}, nil
	}
	return ScanDir(baseDir)
}

// ScanDir discovers IRL projects in the given directory.
func ScanDir(baseDir string) ([]Project, error) {
	var projects []Project

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip IRL internal folders
		name := entry.Name()
		if name == "01-plans" || name == "02-data" || name == "03-outputs" || name == "04-logs" {
			continue
		}

		projectDir := filepath.Join(baseDir, name)

		// Check for main-plan.md in multiple locations
		planPaths := []string{
			filepath.Join(projectDir, "plans", "main-plan.md"),
			filepath.Join(projectDir, "main-plan.md"),
			filepath.Join(projectDir, "01-plans", "main-plan.md"),
		}

		var planInfo os.FileInfo
		for _, planPath := range planPaths {
			if info, err := os.Stat(planPath); err == nil {
				planInfo = info
				break
			}
		}

		if planInfo == nil {
			continue
		}

		projects = append(projects, Project{
			Name:     name,
			Path:     projectDir,
			Modified: planInfo.ModTime(),
		})
	}

	// Sort by modified time, most recent first
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Modified.After(projects[j].Modified)
	})

	return projects, nil
}
