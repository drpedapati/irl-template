package views

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/editor"
	"github.com/drpedapati/irl-template/pkg/naming"
	"github.com/drpedapati/irl-template/pkg/scaffold"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// InitStep represents the current step in the wizard
type InitStep int

const (
	StepDirectory InitStep = iota
	StepBrowse
	StepPurpose
	StepTemplate
	StepCreating
	StepDone
)

// InitModel is the project creation wizard
type InitModel struct {
	step           InitStep
	width          int
	height         int
	baseDir        string
	purposeInput   textinput.Model
	purpose        string
	projectName    string
	templates      []templates.Template
	templateIdx    int
	spinner        spinner.Model
	projectPath    string
	err            error
	done           bool
	skippedDirStep bool // True if we skipped directory selection (default was set)

	// Directory selection state
	directoryCursor int // 0 = use current, 1 = browse

	// Browse state
	browseDir     string
	browseFolders []browseEntry
	browseCursor  int
	browseScroll  int
	browseSortBy  string // "name-asc" or "date-desc"

	// Project action view for success screen
	actionView ProjectActionModel
}

// browseEntry holds a folder name and its modification time
type browseEntry struct {
	Name    string
	ModTime time.Time
}

const browseVisibleItems = 8

// NewInitModel creates a new init view
func NewInitModel() InitModel {
	ti := textinput.New()
	ti.Placeholder = "Enter your project purpose..."
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.Accent)

	// Get default directory
	baseDir := config.GetDefaultDirectory()
	hasDefaultDir := baseDir != ""

	if !hasDefaultDir {
		home, _ := os.UserHomeDir()
		baseDir = filepath.Join(home, "Documents", "irl_projects")
	}

	// Skip folder selection if default is already set
	startStep := StepDirectory
	if hasDefaultDir {
		startStep = StepPurpose
		ti.Focus()
	}

	return InitModel{
		step:           startStep,
		baseDir:        baseDir,
		purposeInput:   ti,
		spinner:        s,
		skippedDirStep: hasDefaultDir,
	}
}

// SetSize sets the view dimensions
func (m *InitModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init initializes the model
func (m InitModel) Init() tea.Cmd {
	if m.step == StepPurpose {
		return textinput.Blink
	}
	return nil
}

// CanGoBack returns true if we can go back a step within the wizard
func (m InitModel) CanGoBack() bool {
	if m.step == StepPurpose && m.skippedDirStep {
		return false // Go back to menu, not within wizard
	}
	return m.step == StepBrowse || m.step == StepPurpose || m.step == StepTemplate
}

// Done returns true if the wizard is complete
func (m InitModel) Done() bool {
	return m.done
}

// Update handles messages
func (m InitModel) Update(msg tea.Msg) (InitModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case editor.EditorFinishedMsg, editor.EditorOpenedMsg:
		// Pass editor messages to action view when in StepDone
		if m.step == StepDone {
			var cmd tea.Cmd
			m.actionView, cmd = m.actionView.Update(msg)
			return m, cmd
		}

	case tea.KeyMsg:
		switch m.step {
		case StepDirectory:
			return m.updateDirectory(msg)
		case StepBrowse:
			return m.updateBrowse(msg)
		case StepPurpose:
			return m.updatePurpose(msg)
		case StepTemplate:
			return m.updateTemplate(msg)
		case StepDone:
			return m.updateDone(msg)
		}

	case spinner.TickMsg:
		if m.step == StepCreating {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case InitTemplatesLoadedMsg:
		m.templates = msg.Templates
		m.step = StepTemplate

	case InitProjectCreatedMsg:
		m.step = StepDone
		m.projectPath = msg.Path
		m.err = msg.Err
		if msg.Err == nil {
			m.actionView = NewProjectActionModel(msg.Path, true)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m InitModel) updateDirectory(msg tea.KeyMsg) (InitModel, tea.Cmd) {
	switch msg.String() {
	case "up":
		if m.directoryCursor > 0 {
			m.directoryCursor--
		}
	case "down":
		if m.directoryCursor < 1 {
			m.directoryCursor++
		}
	case "enter", "right":
		if m.directoryCursor == 0 {
			// Confirm current directory and proceed
			m.step = StepPurpose
			m.purposeInput.Focus()
			return m, textinput.Blink
		} else {
			// Enter browse mode
			m.browseDir = m.baseDir
			if info, err := os.Stat(m.browseDir); err != nil || !info.IsDir() {
				m.browseDir = filepath.Dir(m.browseDir)
			}
			home, _ := os.UserHomeDir()
			if info, err := os.Stat(m.browseDir); err != nil || !info.IsDir() {
				m.browseDir = home
			}
			m.loadFolders()
			m.step = StepBrowse
			return m, nil
		}
	case "esc", "left":
		// Let parent handle going back to menu
		return m, nil
	}
	return m, nil
}

func (m *InitModel) loadFolders() {
	m.browseFolders = []browseEntry{}
	m.browseCursor = 0
	m.browseScroll = 0

	entries, err := os.ReadDir(m.browseDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			modTime := time.Time{}
			if info, err := e.Info(); err == nil {
				modTime = info.ModTime()
			}
			m.browseFolders = append(m.browseFolders, browseEntry{Name: e.Name(), ModTime: modTime})
		}
	}
	m.applyBrowseSort()
}

func (m *InitModel) applyBrowseSort() {
	switch m.browseSortBy {
	case "date-desc":
		sort.Slice(m.browseFolders, func(i, j int) bool {
			return m.browseFolders[i].ModTime.After(m.browseFolders[j].ModTime)
		})
	default: // "name-asc"
		sort.Slice(m.browseFolders, func(i, j int) bool {
			return strings.ToLower(m.browseFolders[i].Name) < strings.ToLower(m.browseFolders[j].Name)
		})
	}
}

func (m InitModel) updateBrowse(msg tea.KeyMsg) (InitModel, tea.Cmd) {
	switch msg.String() {
	case "up":
		if m.browseCursor > 0 {
			m.browseCursor--
			if m.browseCursor < m.browseScroll {
				m.browseScroll = m.browseCursor
			}
		}
	case "down":
		if m.browseCursor < len(m.browseFolders)-1 {
			m.browseCursor++
			if m.browseCursor >= m.browseScroll+browseVisibleItems {
				m.browseScroll = m.browseCursor - browseVisibleItems + 1
			}
		}
	case "s":
		// Toggle sort
		if m.browseSortBy == "date-desc" {
			m.browseSortBy = "name-asc"
		} else {
			m.browseSortBy = "date-desc"
		}
		m.applyBrowseSort()
		m.browseCursor = 0
		m.browseScroll = 0
	case "right":
		// Enter selected folder
		if len(m.browseFolders) > 0 && m.browseCursor < len(m.browseFolders) {
			m.browseDir = filepath.Join(m.browseDir, m.browseFolders[m.browseCursor].Name)
			m.loadFolders()
		}
	case "left":
		// Go up one level
		parent := filepath.Dir(m.browseDir)
		if parent != m.browseDir {
			m.browseDir = parent
			m.loadFolders()
		}
	case "enter":
		// Select current directory for THIS project only (don't change default)
		m.baseDir = m.browseDir
		m.directoryCursor = 0 // Reset cursor to selected directory
		m.step = StepDirectory
		return m, nil
	case "esc":
		// Cancel browse, go back to directory confirm
		m.directoryCursor = 0 // Reset cursor
		m.step = StepDirectory
		return m, nil
	}
	return m, nil
}

func (m InitModel) updatePurpose(msg tea.KeyMsg) (InitModel, tea.Cmd) {
	key := msg.String()

	// Handle navigation keys first
	switch key {
	case "esc", "left":
		if m.skippedDirStep {
			// Let parent handle going back to menu
			return m, nil
		}
		// Go back to directory step
		m.step = StepDirectory
		return m, nil
	}

	// Update text input first so it captures the latest keystroke
	var cmd tea.Cmd
	m.purposeInput, cmd = m.purposeInput.Update(msg)

	// Then check for confirmation keys
	switch key {
	case "enter":
		m.purpose = m.purposeInput.Value()
		if m.purpose == "" {
			return m, cmd // Still return the textinput command
		}
		m.projectName = naming.GenerateName(m.purpose)
		// Load templates
		return m, tea.Batch(m.loadTemplates(), m.spinner.Tick)
	}

	return m, cmd
}

func (m InitModel) updateTemplate(msg tea.KeyMsg) (InitModel, tea.Cmd) {
	switch msg.String() {
	case "up":
		if m.templateIdx > 0 {
			m.templateIdx--
		}
	case "down":
		if m.templateIdx < len(m.templates) {
			m.templateIdx++ // +1 for "None" option
		}
	case "enter", "right":
		m.step = StepCreating
		return m, tea.Batch(m.createProject(), m.spinner.Tick)
	case "esc", "left":
		m.step = StepPurpose
		m.purposeInput.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

func (m InitModel) updateDone(msg tea.KeyMsg) (InitModel, tea.Cmd) {
	var cmd tea.Cmd
	m.actionView, cmd = m.actionView.Update(msg)

	// Check if user wants to exit
	if m.actionView.IsDone() {
		m.done = true
	}

	return m, cmd
}

// InitTemplatesLoadedMsg is sent when templates are loaded for init wizard
type InitTemplatesLoadedMsg struct {
	Templates []templates.Template
}

// InitProjectCreatedMsg is sent when project creation completes
type InitProjectCreatedMsg struct {
	Path string
	Err  error
}

func (m InitModel) loadTemplates() tea.Cmd {
	return func() tea.Msg {
		list, _ := templates.ListTemplates()
		return InitTemplatesLoadedMsg{Templates: list}
	}
}

// injectProfile adds profile information to the plan content
func injectProfile(content string) string {
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

func (m InitModel) createProject() tea.Cmd {
	return func() tea.Msg {
		projectPath := filepath.Join(m.baseDir, m.projectName)

		// Check if exists
		if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
			return InitProjectCreatedMsg{Err: fmt.Errorf("'%s' already exists", projectPath)}
		}

		// Create base directory if needed
		if err := os.MkdirAll(m.baseDir, 0755); err != nil {
			return InitProjectCreatedMsg{Err: err}
		}

		// Create project directory
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return InitProjectCreatedMsg{Err: err}
		}

		// Create scaffold
		if err := scaffold.Create(projectPath); err != nil {
			return InitProjectCreatedMsg{Err: err}
		}

		// Apply template
		var planContent string
		if m.templateIdx > 0 && m.templateIdx <= len(m.templates) {
			tmpl := m.templates[m.templateIdx-1]
			planContent = tmpl.Content
		} else {
			planContent = "# IRL Plan\n\n[Edit this file to define your research plan]\n"
		}

		// Inject profile information
		planContent = injectProfile(planContent)

		if err := scaffold.WritePlan(projectPath, planContent); err != nil {
			return InitProjectCreatedMsg{Err: err}
		}

		// Git init
		scaffold.GitInit(projectPath) // Ignore errors

		return InitProjectCreatedMsg{Path: projectPath}
	}
}

// View renders the init view
func (m InitModel) View() string {
	var b strings.Builder

	switch m.step {
	case StepDirectory:
		b.WriteString(m.viewDirectory())
	case StepBrowse:
		b.WriteString(m.viewBrowse())
	case StepPurpose:
		b.WriteString(m.viewPurpose())
	case StepTemplate:
		b.WriteString(m.viewTemplate())
	case StepCreating:
		b.WriteString(m.viewCreating())
	case StepDone:
		b.WriteString(m.viewDone())
	}

	return b.String()
}

func (m InitModel) viewDirectory() string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().MarginLeft(2)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	pathStyle := lipgloss.NewStyle().Foreground(theme.Accent).MarginLeft(2)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	b.WriteString(labelStyle.Render("Project location"))
	b.WriteString("\n\n")

	// Option 1: Use current directory
	cursor0 := "  "
	style0 := pathStyle
	if m.directoryCursor == 0 {
		cursor0 = keyStyle.Render("●") + " "
		style0 = selectedStyle
	} else {
		cursor0 = "  "
	}
	b.WriteString("  " + cursor0 + style0.Render(m.baseDir))
	b.WriteString("\n")

	// Option 2: Browse
	cursor1 := "  "
	style1 := hintStyle
	if m.directoryCursor == 1 {
		cursor1 = keyStyle.Render("●") + " "
		style1 = selectedStyle
	} else {
		cursor1 = "  "
	}
	b.WriteString("  " + cursor1 + style1.Render("Browse..."))
	b.WriteString("\n\n")

	// Hints
	b.WriteString(hintStyle.Render(keyStyle.Render("↑↓") + " navigate  " + keyStyle.Render("→") + " select"))
	b.WriteString("\n")

	return b.String()
}

func (m InitModel) viewBrowse() string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().MarginLeft(2)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	pathStyle := lipgloss.NewStyle().Foreground(theme.Accent).MarginLeft(2)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	b.WriteString(labelStyle.Render("Browse folders"))
	b.WriteString("\n")
	b.WriteString(pathStyle.Render(m.browseDir))
	b.WriteString("\n\n")

	cursorOn := keyStyle.Render("●")
	cursorOff := "  "
	itemStyle := lipgloss.NewStyle()
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	if len(m.browseFolders) == 0 {
		b.WriteString(hintStyle.Render("  (empty)"))
		b.WriteString("\n")
	} else {
		// Show visible items with scrolling
		end := m.browseScroll + browseVisibleItems
		if end > len(m.browseFolders) {
			end = len(m.browseFolders)
		}

		for i := m.browseScroll; i < end; i++ {
			entry := m.browseFolders[i]
			cursor := cursorOff
			style := itemStyle

			if i == m.browseCursor {
				cursor = cursorOn
				style = selectedStyle
			}

			b.WriteString("  " + cursor + " " + style.Render(entry.Name+"/"))
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if len(m.browseFolders) > browseVisibleItems {
			shown := fmt.Sprintf("%d-%d of %d", m.browseScroll+1, end, len(m.browseFolders))
			b.WriteString(hintStyle.Render("  " + shown))
			b.WriteString("\n")
		}
	}

	// Sort indicator
	sortLabel := "A-Z"
	if m.browseSortBy == "date-desc" {
		sortLabel = "newest"
	}

	b.WriteString("\n")
	b.WriteString(hintStyle.Render(keyStyle.Render("←") + " up  " + keyStyle.Render("→") + " open  " + keyStyle.Render("Enter") + " select  " + keyStyle.Render("s") + ":" + sortLabel))
	b.WriteString("\n")

	return b.String()
}

func (m InitModel) viewPurpose() string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().MarginLeft(2)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	pathStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)

	// Show selected directory
	b.WriteString(pathStyle.Render("📁 " + m.baseDir))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("What's this project for?"))
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("Brief description (e.g., 'ERP correlation analysis')"))
	b.WriteString("\n\n")
	b.WriteString("  " + m.purposeInput.View())
	b.WriteString("\n\n")

	// Show preview of generated folder name
	if m.purposeInput.Value() != "" {
		preview := naming.GenerateName(m.purposeInput.Value())
		labelStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
		nameStyle := lipgloss.NewStyle().Foreground(theme.Accent)
		b.WriteString(labelStyle.Render("Folder: ") + nameStyle.Render(preview))
		b.WriteString("\n\n")
		// Hint for how to proceed
		keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
		b.WriteString(labelStyle.Render(keyStyle.Render("Enter") + " or " + keyStyle.Render("→") + " to continue"))
		b.WriteString("\n")
	}

	return b.String()
}

func (m InitModel) viewTemplate() string {
	var b strings.Builder

	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)

	// Show project info
	pathStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	nameStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).MarginLeft(2)

	b.WriteString(pathStyle.Render("📁 " + m.baseDir))
	b.WriteString("\n")
	b.WriteString(nameStyle.Render(m.projectName))
	b.WriteString("\n\n")

	b.WriteString(hintStyle.Render("Pick a template"))
	b.WriteString("\n\n")

	cursorOn := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("●")
	cursorOff := "  "
	itemStyle := lipgloss.NewStyle().MarginLeft(2)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).MarginLeft(2)

	// "None" option
	cursor := cursorOff
	style := itemStyle
	if m.templateIdx == 0 {
		cursor = cursorOn
		style = selectedStyle
	}
	b.WriteString(cursor + " " + style.Render("None (empty plan)"))
	b.WriteString("\n")

	// Template options
	for i, t := range m.templates {
		cursor = cursorOff
		style = itemStyle
		if m.templateIdx == i+1 {
			cursor = cursorOn
			style = selectedStyle
		}
		b.WriteString(cursor + " " + style.Render(t.Name))
		b.WriteString("\n")
	}

	return b.String()
}

func (m InitModel) viewCreating() string {
	var b strings.Builder
	b.WriteString("  " + m.spinner.View() + " Creating project...")
	b.WriteString("\n")
	return b.String()
}

func (m InitModel) viewDone() string {
	if m.err != nil {
		var b strings.Builder
		errStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
		b.WriteString(errStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
		hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
		b.WriteString(hintStyle.Render("Press ← to go back"))
		return b.String()
	}

	return m.actionView.View()
}
