package views

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/editor"
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
	sortBy      string    // "date-desc", "date-asc", "name-asc", "name-desc"
	editors     []AppInfo // Uses unified AppInfo from editors.go
	openMsg     string    // Message shown after opening
	warningMsg  string    // Warning message (e.g., editor not found)

	// Project detail view
	viewing    bool
	actionView ProjectActionModel

	// Editor launch state
	launchingEditor bool

	// Delete confirmation
	confirmDelete bool
	deleteTarget  string // Path of project to delete
}

const projectsVisibleItems = 10

// ProjectsLoadedMsg is sent when projects are scanned
type ProjectsLoadedMsg struct {
	Projects []Project
	Editors  []AppInfo // Uses unified AppInfo from editors.go
	Err      error
}

// NewProjectsModel creates a new projects view
func NewProjectsModel() ProjectsModel {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Width = 40
	ti.Focus() // Auto-focus on creation

	return ProjectsModel{
		filterInput: ti,
		sortBy:      "date-desc", // Default sort by date, newest first
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
		editors := GetInstalledEditors() // Uses unified source
		return ProjectsLoadedMsg{Projects: projects, Editors: editors, Err: err}
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
			filepath.Join(projectDir, "plans", "main-plan.md"),    // new standard
			filepath.Join(projectDir, "main-plan.md"),             // legacy root level
			filepath.Join(projectDir, "01-plans", "main-plan.md"), // legacy IRL structure
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

// IsViewing returns true when viewing a project detail
func (m ProjectsModel) IsViewing() bool {
	return m.viewing
}

// IsConfirmingDelete returns true when confirming a delete
func (m ProjectsModel) IsConfirmingDelete() bool {
	return m.confirmDelete
}

// trashProject moves the delete target to trash using the trash CLI
func (m *ProjectsModel) trashProject() error {
	if m.deleteTarget == "" {
		return fmt.Errorf("no project selected")
	}

	// Check if trash command is available
	if _, err := exec.LookPath("trash"); err != nil {
		return fmt.Errorf("'trash' command not found (brew install trash)")
	}

	// Move to trash
	cmd := exec.Command("trash", m.deleteTarget)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to trash project: %w", err)
	}

	return nil
}

// backupProject copies the selected project to a _backups directory with timestamp
func (m *ProjectsModel) backupProject() error {
	projectPath := m.SelectedProject()
	if projectPath == "" {
		return fmt.Errorf("no project selected")
	}

	// Get base directory and project name
	baseDir := config.GetDefaultDirectory()
	projectName := filepath.Base(projectPath)

	// Create _backups directory in base dir
	backupsDir := filepath.Join(baseDir, "_backups")
	if err := os.MkdirAll(backupsDir, 0755); err != nil {
		return fmt.Errorf("failed to create backups directory: %w", err)
	}

	// Create backup name with timestamp: projectname_20260202_091500
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s_%s", projectName, timestamp)
	backupPath := filepath.Join(backupsDir, backupName)

	// Copy directory using cp -r
	cmd := exec.Command("cp", "-r", projectPath, backupPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy project: %w", err)
	}

	return nil
}

// Update handles messages
func (m ProjectsModel) Update(msg tea.Msg) (ProjectsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ProjectsLoadedMsg:
		m.loaded = true
		m.err = msg.Err
		m.projects = msg.Projects
		m.filtered = msg.Projects
		m.editors = msg.Editors
		m.cursor = 0
		m.scroll = 0
		return m, textinput.Blink

	case editor.EditorFinishedMsg:
		// Terminal editor closed, TUI resumes - clear screen to remove artifacts
		m.launchingEditor = false
		if msg.Err != nil {
			m.warningMsg = "Editor error: " + msg.Err.Error()
		} else {
			m.openMsg = "Plan saved"
		}
		return m, tea.ClearScreen

	case editor.EditorOpenedMsg:
		// GUI editor launched
		m.launchingEditor = false
		if msg.Err != nil {
			m.warningMsg = "Failed to open editor"
		} else {
			m.openMsg = "Opened in editor"
		}
		return m, nil

	case tea.KeyMsg:
		// If viewing a project, delegate to actionView
		if m.viewing {
			var cmd tea.Cmd
			m.actionView, cmd = m.actionView.Update(msg)
			if m.actionView.IsDone() {
				m.viewing = false
			}
			return m, cmd
		}

		// Handle delete confirmation
		if m.confirmDelete {
			switch msg.String() {
			case "y", "Y":
				// Confirm delete - use trash command
				if err := m.trashProject(); err != nil {
					m.warningMsg = err.Error()
				} else {
					m.openMsg = "Moved to trash"
					// Rescan projects
					return m, m.ScanProjects()
				}
				m.confirmDelete = false
				m.deleteTarget = ""
			case "n", "N", "esc":
				// Cancel delete
				m.confirmDelete = false
				m.deleteTarget = ""
			}
			return m, nil
		}

		key := msg.String()

		// Navigation keys always work
		switch key {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}
			return m, nil
		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				if m.cursor >= m.scroll+projectsVisibleItems {
					m.scroll = m.cursor - projectsVisibleItems + 1
				}
			}
			return m, nil
		case "right", "enter":
			// Enter project detail view
			if m.SelectedProject() != "" {
				m.viewing = true
				m.actionView = NewProjectActionModel(m.SelectedProject(), false)
			}
			return m, nil
		}

		// If filter has text, pass most keys to filter (not shortcuts)
		if m.filterInput.Value() != "" {
			var cmd tea.Cmd
			m.filterInput, cmd = m.filterInput.Update(msg)
			m.applyFilter()
			return m, cmd
		}

		// Shortcuts only work when filter is empty
		// Check for editor shortcuts first (only when there's a selection)
		if m.SelectedProject() != "" {
			for _, e := range m.editors {
				if key == e.Key {
					m.openInEditor(e)
					return m, nil
				}
			}
		}

		switch key {
		case "e":
			// Edit plan file
			if m.SelectedProject() != "" {
				return m.editPlanFile()
			}
			return m, nil
		case "b":
			// Backup project
			if m.SelectedProject() != "" {
				if err := m.backupProject(); err != nil {
					m.warningMsg = err.Error()
					m.openMsg = ""
				} else {
					m.openMsg = "Backup created"
					m.warningMsg = ""
				}
			}
			return m, nil
		case "x":
			// Delete project (with confirmation)
			if m.SelectedProject() != "" {
				m.confirmDelete = true
				m.deleteTarget = m.SelectedProject()
				m.warningMsg = ""
				m.openMsg = ""
			}
			return m, nil
		case "s":
			// Cycle through sort options: date newest → date oldest → name A-Z → name Z-A
			switch m.sortBy {
			case "date-desc":
				m.sortBy = "date-asc"
			case "date-asc":
				m.sortBy = "name-asc"
			case "name-asc":
				m.sortBy = "name-desc"
			default:
				m.sortBy = "date-desc"
			}
			m.applySort()
			return m, nil
		}

		// Pass other keys to filter input
		var cmd tea.Cmd
		m.filterInput, cmd = m.filterInput.Update(msg)
		m.applyFilter()
		return m, cmd
	}

	return m, nil
}

func (m *ProjectsModel) openInEditor(e AppInfo) {
	path := m.SelectedProject()
	if path == "" {
		return
	}

	if errMsg := OpenProjectWith(e, path); errMsg != "" {
		m.openMsg = errMsg
	} else {
		m.openMsg = "Opened in " + e.Name
	}
}

// editPlanFile opens the main-plan.md in the preferred editor
func (m ProjectsModel) editPlanFile() (ProjectsModel, tea.Cmd) {
	projectPath := m.SelectedProject()
	if projectPath == "" {
		return m, nil
	}

	// Get preferred editor
	ed, found := editor.GetPreferred()
	if !found {
		m.warningMsg = "No editor found. Configure in Editors view."
		return m, nil
	}

	// Get plan file path
	planPath := editor.GetPlanPath(projectPath)

	// Check if plan file exists
	if _, err := os.Stat(planPath); os.IsNotExist(err) {
		m.warningMsg = "main-plan.md not found"
		return m, nil
	}

	// Launch editor
	m.launchingEditor = true
	m.warningMsg = ""
	m.openMsg = ""

	if ed.Type == editor.EditorTypeTerminal {
		// Terminal editor suspends TUI
		return m, editor.Open(ed, planPath)
	}

	// GUI editor runs in background
	m.openMsg = "Opening in " + ed.Name + "..."
	return m, editor.Open(ed, planPath)
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
	switch m.sortBy {
	case "date-desc":
		sort.Slice(m.filtered, func(i, j int) bool {
			return m.filtered[i].Modified.After(m.filtered[j].Modified)
		})
	case "date-asc":
		sort.Slice(m.filtered, func(i, j int) bool {
			return m.filtered[i].Modified.Before(m.filtered[j].Modified)
		})
	case "name-asc":
		sort.Slice(m.filtered, func(i, j int) bool {
			return strings.ToLower(m.filtered[i].Name) < strings.ToLower(m.filtered[j].Name)
		})
	case "name-desc":
		sort.Slice(m.filtered, func(i, j int) bool {
			return strings.ToLower(m.filtered[i].Name) > strings.ToLower(m.filtered[j].Name)
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

// HasFilterText returns true if there's text in the filter input
func (m ProjectsModel) HasFilterText() bool {
	return m.filterInput.Value() != ""
}

// ClearFilter clears the filter input and resets the list
func (m *ProjectsModel) ClearFilter() {
	m.filterInput.SetValue("")
	m.applyFilter()
}

// GetEditorHints returns hints for available editors
func (m ProjectsModel) GetEditorHints() string {
	if len(m.editors) == 0 {
		return ""
	}

	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	var hints []string
	for _, e := range m.editors {
		hints = append(hints, e.Key+" "+e.Name)
	}
	return mutedStyle.Render(strings.Join(hints, "  "))
}

// View renders the projects view
func (m ProjectsModel) View() string {
	// Show project action view if viewing
	if m.viewing {
		return m.actionView.View()
	}

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

	// Filter input (always visible, always focused)
	b.WriteString("\n")
	b.WriteString("  " + m.filterInput.View())

	// Sort indicator
	var sortLabel string
	switch m.sortBy {
	case "date-desc":
		sortLabel = "newest"
	case "date-asc":
		sortLabel = "oldest"
	case "name-asc":
		sortLabel = "A-Z"
	case "name-desc":
		sortLabel = "Z-A"
	}
	b.WriteString("  " + mutedStyle.Render("s:"+sortLabel))
	b.WriteString("\n\n")

	// Delete confirmation
	if m.confirmDelete {
		warningStyle := lipgloss.NewStyle().Foreground(theme.Warning).Bold(true)
		projectName := filepath.Base(m.deleteTarget)
		b.WriteString("  " + warningStyle.Render("Delete \""+projectName+"\"? (y/n)"))
		b.WriteString("\n\n")
	} else if m.warningMsg != "" {
		// Warning message (inline, above list)
		warningStyle := lipgloss.NewStyle().Foreground(theme.Warning)
		b.WriteString("  " + warningStyle.Render("⚠ "+m.warningMsg))
		b.WriteString("\n\n")
	} else if m.openMsg != "" {
		// Open message (temporary feedback)
		successStyle := lipgloss.NewStyle().Foreground(theme.Success)
		b.WriteString("  " + successStyle.Render("✓ "+m.openMsg))
		b.WriteString("\n\n")
	}

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

	cursorOn := accentStyle.Render(">")
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
