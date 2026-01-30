package views

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
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

// Implement list.Item interface
func (p Project) Title() string       { return p.Name }
func (p Project) Description() string { return smartDate(p.Modified) }
func (p Project) FilterValue() string { return p.Name }

// ProjectsModel displays discovered IRL projects
type ProjectsModel struct {
	projects    []Project
	filtered    []Project
	cursor      int
	scroll      int
	width       int
	height      int
	loaded      bool
	err         error
	filterInput textinput.Model
	filtering   bool
	sortBy      string // "date" or "name"
}

const projectsVisibleItems = 10

// ProjectsLoadedMsg is sent when projects are scanned
type ProjectsLoadedMsg struct {
	Projects []Project
	Err      error
}

// NewProjectsModel creates a new projects view
func NewProjectsModel() ProjectsModel {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Width = 30

	return ProjectsModel{
		filterInput: ti,
		sortBy:      "date", // Default sort by date
	}
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

		// Skip IRL internal folders - these are not projects
		name := entry.Name()
		if name == "01-plans" || name == "02-data" || name == "03-outputs" || name == "04-logs" {
			continue
		}

		projectDir := filepath.Join(baseDir, name)

		// Check for main-plan.md in multiple locations
		planPaths := []string{
			filepath.Join(projectDir, "main-plan.md"),             // root level
			filepath.Join(projectDir, "01-plans", "main-plan.md"), // IRL structure
		}

		var planInfo os.FileInfo
		for _, planPath := range planPaths {
			if info, err := os.Stat(planPath); err == nil {
				planInfo = info
				break
			}
		}

		if planInfo == nil {
			continue // Not an IRL project
		}

		projects = append(projects, Project{
			Name:     name,
			Path:     projectDir,
			Modified: planInfo.ModTime(),
		})
	}

	// Sort by modified time, most recent first (default)
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
		m.filtered = msg.Projects
		m.cursor = 0
		m.scroll = 0
		return m, nil

	case tea.KeyMsg:
		// If filtering, handle text input
		if m.filtering {
			switch msg.String() {
			case "esc":
				m.filtering = false
				m.filterInput.Blur()
				m.filterInput.SetValue("")
				m.applyFilter()
				return m, nil
			case "enter":
				m.filtering = false
				m.filterInput.Blur()
				return m, nil
			default:
				var cmd tea.Cmd
				m.filterInput, cmd = m.filterInput.Update(msg)
				m.applyFilter()
				return m, cmd
			}
		}

		switch msg.String() {
		case "/":
			// Start filtering
			m.filtering = true
			m.filterInput.Focus()
			return m, textinput.Blink
		case "s":
			// Toggle sort
			if m.sortBy == "date" {
				m.sortBy = "name"
			} else {
				m.sortBy = "date"
			}
			m.applySort()
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				if m.cursor >= m.scroll+projectsVisibleItems {
					m.scroll = m.cursor - projectsVisibleItems + 1
				}
			}
		}
	}

	return m, nil
}

func (m *ProjectsModel) applyFilter() {
	query := strings.ToLower(m.filterInput.Value())
	if query == "" {
		m.filtered = m.projects
	} else {
		m.filtered = []Project{}
		for _, p := range m.projects {
			if strings.Contains(strings.ToLower(p.Name), query) {
				m.filtered = append(m.filtered, p)
			}
		}
	}
	m.cursor = 0
	m.scroll = 0
	m.applySort()
}

func (m *ProjectsModel) applySort() {
	if m.sortBy == "name" {
		sort.Slice(m.filtered, func(i, j int) bool {
			return m.filtered[i].Name < m.filtered[j].Name
		})
	} else {
		sort.Slice(m.filtered, func(i, j int) bool {
			return m.filtered[i].Modified.After(m.filtered[j].Modified)
		})
	}
}

// SelectedProject returns the currently selected project path
func (m ProjectsModel) SelectedProject() string {
	if m.cursor >= 0 && m.cursor < len(m.filtered) {
		return m.filtered[m.cursor].Path
	}
	return ""
}

// View renders the projects view
func (m ProjectsModel) View() string {
	var b strings.Builder

	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	accentStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	normalStyle := lipgloss.NewStyle()
	dateStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
		b.WriteString(errStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
		return b.String()
	}

	if !m.loaded {
		return b.String()
	}

	// Filter bar
	b.WriteString("\n")
	if m.filtering {
		b.WriteString("  " + m.filterInput.View())
	} else {
		filterHint := mutedStyle.Render("/ filter")
		sortLabel := "date"
		if m.sortBy == "name" {
			sortLabel = "name"
		}
		sortHint := mutedStyle.Render("s sort:" + sortLabel)
		b.WriteString("  " + filterHint + "  " + sortHint)
		if m.filterInput.Value() != "" {
			b.WriteString("  " + accentStyle.Render("\""+m.filterInput.Value()+"\""))
		}
	}
	b.WriteString("\n\n")

	if len(m.filtered) == 0 {
		if len(m.projects) == 0 {
			b.WriteString(mutedStyle.Render("  No IRL projects found"))
			b.WriteString("\n\n")
			b.WriteString(mutedStyle.Render("  Projects need a main-plan.md file"))
		} else {
			b.WriteString(mutedStyle.Render("  No matches for \"" + m.filterInput.Value() + "\""))
		}
		b.WriteString("\n")
		return b.String()
	}

	// Calculate column widths
	nameColWidth := m.width - 20 // Leave room for date column
	if nameColWidth < 20 {
		nameColWidth = 20
	}
	if nameColWidth > 45 {
		nameColWidth = 45
	}

	cursorOn := accentStyle.Render("●")
	cursorOff := " "

	// Show visible projects with scrolling
	endIdx := m.scroll + projectsVisibleItems
	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}

	for i := m.scroll; i < endIdx; i++ {
		p := m.filtered[i]
		cursor := cursorOff
		nameStyle := normalStyle
		dateStyleLocal := dateStyle

		if i == m.cursor {
			cursor = cursorOn
			nameStyle = selectedStyle
			dateStyleLocal = selectedStyle
		}

		// Truncate long names
		displayName := p.Name
		if len(displayName) > nameColWidth {
			displayName = displayName[:nameColWidth-3] + "..."
		}

		// Pad name to align dates
		namePadded := displayName + strings.Repeat(" ", nameColWidth-len(displayName))

		// Smart date
		dateStr := smartDate(p.Modified)

		b.WriteString("  " + cursor + " " + nameStyle.Render(namePadded) + " " + dateStyleLocal.Render(dateStr))
		b.WriteString("\n")
	}

	// Scroll indicator
	if len(m.filtered) > projectsVisibleItems {
		showing := m.scroll + 1
		showingEnd := endIdx
		total := len(m.filtered)
		indicator := mutedStyle.Render("    " + itoa(showing) + "-" + itoa(showingEnd) + " of " + itoa(total))
		b.WriteString(indicator)
		b.WriteString("\n")
	}

	return b.String()
}

// smartDate returns a human-friendly date string
func smartDate(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	// Same day - show time
	if t.YearDay() == now.YearDay() && t.Year() == now.Year() {
		return t.Format("3:04 PM")
	}

	// Yesterday
	yesterday := now.AddDate(0, 0, -1)
	if t.YearDay() == yesterday.YearDay() && t.Year() == yesterday.Year() {
		return "Yesterday"
	}

	// Within last week - show day name
	if diff < 7*24*time.Hour {
		return t.Format("Monday")
	}

	// Same year - show month and day
	if t.Year() == now.Year() {
		return t.Format("Jan 2")
	}

	// Different year - show full date
	return t.Format("Jan 2, 2006")
}

