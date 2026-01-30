package views

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// Project represents a discovered IRL project
type Project struct {
	Name     string
	Path     string
	Modified time.Time
}

// ProjectsModel displays discovered IRL projects
type ProjectsModel struct {
	projects []Project
	cursor   int
	width    int
	height   int
	loaded   bool
	err      error
}

// ProjectsLoadedMsg is sent when projects are scanned
type ProjectsLoadedMsg struct {
	Projects []Project
	Err      error
}

// NewProjectsModel creates a new projects view
func NewProjectsModel() ProjectsModel {
	return ProjectsModel{}
}

// SetSize sets the view dimensions
func (m *ProjectsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// ScanProjects returns a command that scans for IRL projects
func (m *ProjectsModel) ScanProjects() tea.Cmd {
	return func() tea.Msg {
		baseDir := config.GetDefaultDirectory()
		if baseDir == "" {
			return ProjectsLoadedMsg{Err: nil, Projects: []Project{}}
		}

		projects, err := scanForProjects(baseDir)
		return ProjectsLoadedMsg{Projects: projects, Err: err}
	}
}

func scanForProjects(baseDir string) ([]Project, error) {
	var projects []Project

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check for main-plan.md
		planPath := filepath.Join(baseDir, entry.Name(), "main-plan.md")
		info, err := os.Stat(planPath)
		if err != nil {
			continue // Not an IRL project
		}

		projects = append(projects, Project{
			Name:     entry.Name(),
			Path:     filepath.Join(baseDir, entry.Name()),
			Modified: info.ModTime(),
		})
	}

	// Sort by modified time, most recent first
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Modified.After(projects[j].Modified)
	})

	return projects, nil
}

// Update handles messages
func (m ProjectsModel) Update(msg tea.Msg) (ProjectsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ProjectsLoadedMsg:
		m.loaded = true
		m.err = msg.Err
		m.projects = msg.Projects
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

// SelectedProject returns the currently selected project path
func (m ProjectsModel) SelectedProject() string {
	if m.cursor >= 0 && m.cursor < len(m.projects) {
		return m.projects[m.cursor].Path
	}
	return ""
}

// View renders the projects view
func (m ProjectsModel) View() string {
	var b strings.Builder

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
		b.WriteString(errStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
		return b.String()
	}

	if !m.loaded {
		return b.String()
	}

	if len(m.projects) == 0 {
		mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
		b.WriteString(mutedStyle.Render("No IRL projects found"))
		b.WriteString("\n\n")
		b.WriteString(mutedStyle.Render("Projects must contain a main-plan.md file"))
		b.WriteString("\n")
		return b.String()
	}

	itemStyle := lipgloss.NewStyle().MarginLeft(2)
	selectedStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true).
		MarginLeft(2)
	dateStyle := lipgloss.NewStyle().
		Foreground(theme.Muted).
		MarginLeft(4)

	cursorOn := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("●")
	cursorOff := "  "

	for i, p := range m.projects {
		cursor := cursorOff
		style := itemStyle

		if i == m.cursor {
			cursor = cursorOn
			style = selectedStyle
		}

		b.WriteString(cursor + " " + style.Render(p.Name))
		b.WriteString("\n")

		// Show relative time
		relTime := relativeTime(p.Modified)
		b.WriteString(dateStyle.Render("Modified " + relTime))
		b.WriteString("\n")
	}

	return b.String()
}

func relativeTime(t time.Time) string {
	diff := time.Since(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return strings.Replace("X minutes ago", "X", string(rune('0'+mins/10))+string(rune('0'+mins%10)), 1)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return t.Format("3:04 PM")
	} else if diff < 7*24*time.Hour {
		return t.Format("Mon 3:04 PM")
	}
	return t.Format("Jan 2")
}
