package views

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
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
	step         InitStep
	width        int
	height       int
	baseDir      string
	purposeInput textinput.Model
	purpose      string
	projectName  string
	templates    []templates.Template
	templateIdx  int
	spinner      spinner.Model
	projectPath  string
	err          error
	done         bool

	// Browse state
	browseDir     string
	browseFolders []string
	browseCursor  int
	browseScroll  int
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
	if baseDir == "" {
		home, _ := os.UserHomeDir()
		baseDir = filepath.Join(home, "Documents", "irl_projects")
	}

	return InitModel{
		step:         StepDirectory,
		baseDir:      baseDir,
		purposeInput: ti,
		spinner:      s,
	}
}

// SetSize sets the view dimensions
func (m *InitModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init initializes the model
func (m InitModel) Init() tea.Cmd {
	return nil
}

// CanGoBack returns true if we can go back a step
func (m InitModel) CanGoBack() bool {
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
			if msg.String() == "enter" || msg.String() == "esc" || msg.String() == "left" {
				m.done = true
				return m, nil
			}
		}

	case spinner.TickMsg:
		if m.step == StepCreating {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case templatesLoadedMsg:
		m.templates = msg.templates
		m.step = StepTemplate

	case projectCreatedMsg:
		m.step = StepDone
		m.projectPath = msg.path
		m.err = msg.err
	}

	return m, tea.Batch(cmds...)
}

func (m InitModel) updateDirectory(msg tea.KeyMsg) (InitModel, tea.Cmd) {
	switch msg.String() {
	case "enter", "right":
		// Confirm current directory and proceed
		m.step = StepPurpose
		m.purposeInput.Focus()
		return m, textinput.Blink
	case "down", "j", "b":
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
	case "esc", "left":
		// Let parent handle going back to menu
		return m, nil
	}
	return m, nil
}

func (m *InitModel) loadFolders() {
	m.browseFolders = []string{}
	m.browseCursor = 0
	m.browseScroll = 0

	entries, err := os.ReadDir(m.browseDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			m.browseFolders = append(m.browseFolders, e.Name())
		}
	}
	sort.Strings(m.browseFolders)
}

func (m InitModel) updateBrowse(msg tea.KeyMsg) (InitModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.browseCursor > 0 {
			m.browseCursor--
			if m.browseCursor < m.browseScroll {
				m.browseScroll = m.browseCursor
			}
		}
	case "down", "j":
		if m.browseCursor < len(m.browseFolders)-1 {
			m.browseCursor++
			if m.browseCursor >= m.browseScroll+browseVisibleItems {
				m.browseScroll = m.browseCursor - browseVisibleItems + 1
			}
		}
	case "right", "l":
		// Enter selected folder
		if len(m.browseFolders) > 0 && m.browseCursor < len(m.browseFolders) {
			m.browseDir = filepath.Join(m.browseDir, m.browseFolders[m.browseCursor])
			m.loadFolders()
		}
	case "left", "h":
		// Go up one level
		parent := filepath.Dir(m.browseDir)
		if parent != m.browseDir {
			m.browseDir = parent
			m.loadFolders()
		}
	case "enter":
		// Select current directory
		m.baseDir = m.browseDir
		config.SetDefaultDirectory(m.baseDir)
		m.step = StepDirectory
		return m, nil
	case "esc":
		// Cancel browse, go back to directory confirm
		m.step = StepDirectory
		return m, nil
	}
	return m, nil
}

func (m InitModel) updatePurpose(msg tea.KeyMsg) (InitModel, tea.Cmd) {
	switch msg.String() {
	case "enter", "right":
		m.purpose = m.purposeInput.Value()
		if m.purpose == "" {
			return m, nil
		}
		m.projectName = naming.GenerateName(m.purpose)
		// Load templates
		return m, tea.Batch(m.loadTemplates(), m.spinner.Tick)
	case "esc", "left":
		// Go back to directory step
		m.step = StepDirectory
		return m, nil
	}

	var cmd tea.Cmd
	m.purposeInput, cmd = m.purposeInput.Update(msg)
	return m, cmd
}

func (m InitModel) updateTemplate(msg tea.KeyMsg) (InitModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.templateIdx > 0 {
			m.templateIdx--
		}
	case "down", "j":
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

type templatesLoadedMsg struct {
	templates []templates.Template
}

type projectCreatedMsg struct {
	path string
	err  error
}

func (m InitModel) loadTemplates() tea.Cmd {
	return func() tea.Msg {
		list, _ := templates.ListTemplates()
		return templatesLoadedMsg{templates: list}
	}
}

func (m InitModel) createProject() tea.Cmd {
	return func() tea.Msg {
		projectPath := filepath.Join(m.baseDir, m.projectName)

		// Check if exists
		if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
			return projectCreatedMsg{err: fmt.Errorf("'%s' already exists", projectPath)}
		}

		// Create base directory if needed
		if err := os.MkdirAll(m.baseDir, 0755); err != nil {
			return projectCreatedMsg{err: err}
		}

		// Create project directory
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return projectCreatedMsg{err: err}
		}

		// Create scaffold
		if err := scaffold.Create(projectPath); err != nil {
			return projectCreatedMsg{err: err}
		}

		// Apply template
		var planContent string
		if m.templateIdx > 0 && m.templateIdx <= len(m.templates) {
			tmpl := m.templates[m.templateIdx-1]
			planContent = tmpl.Content
		} else {
			planContent = "# IRL Plan\n\n[Edit this file to define your research plan]\n"
		}

		if err := scaffold.WritePlan(projectPath, planContent); err != nil {
			return projectCreatedMsg{err: err}
		}

		// Git init
		scaffold.GitInit(projectPath) // Ignore errors

		return projectCreatedMsg{path: projectPath}
	}
}

// View renders the init view
func (m InitModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Primary).
		MarginLeft(2).
		MarginTop(1)

	b.WriteString(titleStyle.Render("New Project"))
	b.WriteString("\n\n")

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

	b.WriteString(labelStyle.Render("Project location"))
	b.WriteString("\n\n")

	// Show current directory with cursor
	b.WriteString("  " + keyStyle.Render("●") + " " + pathStyle.Render(m.baseDir))
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("    Browse..."))
	b.WriteString("\n\n")

	// Hints
	b.WriteString(hintStyle.Render(keyStyle.Render("→") + " continue  " + keyStyle.Render("↓") + " browse"))
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
			folder := m.browseFolders[i]
			cursor := cursorOff
			style := itemStyle

			if i == m.browseCursor {
				cursor = cursorOn
				style = selectedStyle
			}

			b.WriteString("  " + cursor + " " + style.Render(folder+"/"))
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if len(m.browseFolders) > browseVisibleItems {
			shown := fmt.Sprintf("%d-%d of %d", m.browseScroll+1, end, len(m.browseFolders))
			b.WriteString(hintStyle.Render("  " + shown))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(hintStyle.Render(keyStyle.Render("←") + " up  " + keyStyle.Render("→") + " open  " + keyStyle.Render("Enter") + " select"))
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

	// Show preview of generated name
	if m.purposeInput.Value() != "" {
		preview := naming.GenerateName(m.purposeInput.Value())
		previewStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
		nameStyle := lipgloss.NewStyle().Foreground(theme.Accent)
		b.WriteString(previewStyle.Render("→ ") + nameStyle.Render(preview))
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
	var b strings.Builder

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
		b.WriteString(errStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
		hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
		b.WriteString(hintStyle.Render("Press ← to go back"))
		return b.String()
	}

	checkStyle := lipgloss.NewStyle().Foreground(theme.Success)
	pathStyle := lipgloss.NewStyle().Foreground(theme.Accent)
	cmdStyle := lipgloss.NewStyle().Foreground(theme.Accent).MarginLeft(2)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)

	b.WriteString("  " + checkStyle.Render("✓") + " Created " + pathStyle.Render(m.projectPath))
	b.WriteString("\n\n")

	b.WriteString(hintStyle.Render("Next:"))
	b.WriteString("\n")
	b.WriteString(cmdStyle.Render("cd " + m.projectPath))
	b.WriteString("\n")
	b.WriteString(cmdStyle.Render("cursor .") + hintStyle.Render(" or ") + cmdStyle.Render("code ."))
	b.WriteString("\n\n")

	b.WriteString(hintStyle.Render("Press ← to go back"))
	b.WriteString("\n")

	return b.String()
}
