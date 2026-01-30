package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/internal/tui/views"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/theme"
)

const (
	appWidth  = 72 // Fixed app width for Claude Code-like feel
	appHeight = 18 // Fixed app height (content area)
)

// Model is the main TUI model
type Model struct {
	version   string
	width     int
	height    int
	header    Header
	menu      Menu
	statusBar StatusBar
	view      ViewType
	quitting  bool

	// Sub-views
	templatesView views.TemplatesModel
	doctorView    views.DoctorModel
	initView      views.InitModel
	configView    views.ConfigModel

	// Loading state
	loading bool
	spinner spinner.Model
}

// New creates a new TUI model
func New(version string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	m := Model{
		version:       version,
		header:        NewHeader(version),
		menu:          NewMenu(),
		statusBar:     NewStatusBar(),
		view:          ViewMenu,
		templatesView: views.NewTemplatesModel(),
		doctorView:    views.NewDoctorModel(),
		initView:      views.NewInitModel(),
		configView:    views.NewConfigModel(),
		spinner:       s,
	}

	// Set fixed widths
	m.header.SetWidth(appWidth)
	m.menu.SetWidth(appWidth)
	m.statusBar.SetWidth(appWidth)
	m.templatesView.SetSize(appWidth, appHeight-4)
	m.doctorView.SetSize(appWidth, appHeight-4)
	m.initView.SetSize(appWidth, appHeight-4)
	m.configView.SetSize(appWidth, appHeight-4)

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		// Handle based on current view
		switch m.view {
		case ViewMenu:
			return m.updateMenu(msg)
		case ViewTemplates:
			return m.updateTemplates(msg)
		case ViewDoctor:
			return m.updateDoctor(msg)
		case ViewInit:
			return m.updateInit(msg)
		case ViewConfig:
			return m.updateConfig(msg)
		}

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	// Handle view-specific messages
	case views.TemplatesLoadedMsg:
		m.templatesView, _ = m.templatesView.Update(msg)
		m.loading = false

	case views.DoctorResultMsg:
		m.doctorView, _ = m.doctorView.Update(msg)
		m.loading = false

	case views.InitCompleteMsg:
		m.initView, _ = m.initView.Update(msg)

	case views.BackToMenuMsg:
		m.view = ViewMenu
		m.statusBar.SetKeys(DefaultMenuKeys())
	}

	return m, tea.Batch(cmds...)
}

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "up", "k":
		m.menu.Up()
	case "down", "j":
		m.menu.Down()
	case "enter", "right", "l":
		return m.selectView(m.menu.Select())
	case "n":
		if v, ok := m.menu.SelectByKey("n"); ok {
			return m.selectView(v)
		}
	case "t":
		if v, ok := m.menu.SelectByKey("t"); ok {
			return m.selectView(v)
		}
	case "d":
		if v, ok := m.menu.SelectByKey("d"); ok {
			return m.selectView(v)
		}
	case "c":
		if v, ok := m.menu.SelectByKey("c"); ok {
			return m.selectView(v)
		}
	}
	return m, nil
}

func (m Model) selectView(v ViewType) (tea.Model, tea.Cmd) {
	m.view = v
	var cmd tea.Cmd

	switch v {
	case ViewTemplates:
		m.statusBar.SetKeys(TemplateViewKeys())
		m.loading = true
		cmd = tea.Batch(m.spinner.Tick, m.templatesView.LoadTemplates())
	case ViewDoctor:
		m.statusBar.SetKeys(ViewKeys())
		m.loading = true
		cmd = tea.Batch(m.spinner.Tick, m.doctorView.RunChecks())
	case ViewInit:
		m.statusBar.SetKeys(InitViewKeys())
		m.initView = views.NewInitModel()
		m.initView.SetSize(appWidth, appHeight-4)
		cmd = m.initView.Init()
	case ViewConfig:
		m.statusBar.SetKeys(ViewKeys())
		cmd = m.configView.Load()
	}

	return m, cmd
}

func (m Model) updateTemplates(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc", "left", "h":
		m.view = ViewMenu
		m.statusBar.SetKeys(DefaultMenuKeys())
		return m, nil
	case "r":
		m.loading = true
		return m, tea.Batch(m.spinner.Tick, m.templatesView.RefreshTemplates())
	}

	var cmd tea.Cmd
	m.templatesView, cmd = m.templatesView.Update(msg)
	return m, cmd
}

func (m Model) updateDoctor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc", "left", "h":
		m.view = ViewMenu
		m.statusBar.SetKeys(DefaultMenuKeys())
		return m, nil
	}

	var cmd tea.Cmd
	m.doctorView, cmd = m.doctorView.Update(msg)
	return m, cmd
}

func (m Model) updateInit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "left", "h":
		if m.initView.CanGoBack() {
			var cmd tea.Cmd
			m.initView, cmd = m.initView.Update(msg)
			return m, cmd
		}
		m.view = ViewMenu
		m.statusBar.SetKeys(DefaultMenuKeys())
		return m, nil
	}

	var cmd tea.Cmd
	m.initView, cmd = m.initView.Update(msg)

	// Check if init completed and we should go back to menu
	if m.initView.Done() {
		m.view = ViewMenu
		m.statusBar.SetKeys(DefaultMenuKeys())
	}

	return m, cmd
}

func (m Model) updateConfig(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc", "left", "h":
		m.view = ViewMenu
		m.statusBar.SetKeys(DefaultMenuKeys())
		return m, nil
	}

	var cmd tea.Cmd
	m.configView, cmd = m.configView.Update(msg)
	return m, cmd
}

// View renders the TUI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Build inner content
	var inner strings.Builder

	// Add breathing room at top
	inner.WriteString("\n\n")

	// Header
	inner.WriteString(m.header.View())
	inner.WriteString("\n")

	// Subheader: folder path on left, datetime on right
	defaultDir := config.GetDefaultDirectory()
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	var pathText string
	if defaultDir == "" {
		pathText = "No default project path"
	} else {
		pathText = defaultDir
	}

	// Format: Mon Jan 30 2:45 PM MST
	now := time.Now()
	dateTime := now.Format("Mon Jan 2 3:04 PM MST")

	// Calculate padding between path and datetime
	padding := appWidth - lipgloss.Width(pathText) - lipgloss.Width(dateTime)
	if padding < 1 {
		padding = 1
	}

	inner.WriteString(mutedStyle.Render(pathText) + strings.Repeat(" ", padding) + mutedStyle.Render(dateTime))
	inner.WriteString("\n")

	inner.WriteString(Divider(appWidth))

	// Content area
	content := ""
	switch m.view {
	case ViewMenu:
		content = m.renderMenuView()
	case ViewTemplates:
		if m.loading {
			content = m.renderLoading("Loading templates...")
		} else {
			content = m.templatesView.View()
		}
	case ViewDoctor:
		if m.loading {
			content = m.renderLoading("Checking environment...")
		} else {
			content = m.doctorView.View()
		}
	case ViewInit:
		content = m.initView.View()
	case ViewConfig:
		content = m.configView.View()
	}

	// Truncate or pad content to fixed height
	contentLines := strings.Split(content, "\n")
	maxContentLines := appHeight - 4 // header, path, divider, status bar

	// Truncate if too long
	if len(contentLines) > maxContentLines {
		contentLines = contentLines[:maxContentLines]
	}

	// Pad if too short
	for len(contentLines) < maxContentLines {
		contentLines = append(contentLines, "")
	}

	inner.WriteString("\n")
	inner.WriteString(strings.Join(contentLines, "\n"))
	inner.WriteString("\n")
	inner.WriteString(m.statusBar.View())

	// No border - clean Claude Code style
	return inner.String()
}

func (m Model) renderMenuView() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.menu.View())
	b.WriteString("\n")

	// Tip at bottom
	tipStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	cmdStyle := lipgloss.NewStyle().Foreground(theme.Accent)

	b.WriteString(tipStyle.Render("  Or run: ") + cmdStyle.Render("irl init \"your project\""))
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderLoading(msg string) string {
	style := lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1)
	return style.Render(m.spinner.View() + " " + msg) + "\n"
}

// Run starts the TUI
func Run(version string) error {
	p := tea.NewProgram(
		New(version),
		// No alt screen - renders inline at current cursor position
	)

	_, err := p.Run()
	return err
}
