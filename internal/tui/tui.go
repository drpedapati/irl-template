package tui

import (
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/internal/tui/views"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/doctor"
	"github.com/drpedapati/irl-template/pkg/theme"
)

const (
	appWidth  = 72 // Fixed app width for Claude Code-like feel
	appHeight = 24 // Fixed app height (taller for markdown preview)
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
	updateView    views.UpdateModel

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
		updateView:    views.NewUpdateModel(),
		spinner:       s,
	}

	// Set fixed widths
	m.header.SetWidth(appWidth)
	m.menu.SetWidth(appWidth)
	m.statusBar.SetWidth(appWidth)
	m.templatesView.SetSize(appWidth, appHeight-7)
	m.doctorView.SetSize(appWidth, appHeight-7)
	m.initView.SetSize(appWidth, appHeight-7)
	m.configView.SetSize(appWidth, appHeight-7)

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
		case ViewUpdate:
			return m.updateUpdateView(msg)
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

	case views.UpdateProgressMsg:
		m.updateView, _ = m.updateView.Update(msg)

	case views.UpdateCompleteMsg:
		m.updateView, _ = m.updateView.Update(msg)
		// Also refresh templates view so it has latest
		m.templatesView = views.NewTemplatesModel()
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
	case "o":
		if v, ok := m.menu.SelectByKey("o"); ok {
			return m.selectView(v)
		}
	case "u":
		if v, ok := m.menu.SelectByKey("u"); ok {
			return m.selectView(v)
		}
	}
	return m, nil
}

func (m Model) selectView(v ViewType) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch v {
	case ViewDocs:
		// Open docs in browser, stay on menu
		openBrowser("https://docs.irloop.org")
		return m, nil
	case ViewUpdate:
		// Show update view with progress bar
		m.view = v
		m.updateView = views.NewUpdateModel()
		m.updateView.SetSize(appWidth, appHeight-7)
		return m, m.updateView.StartUpdate()
	case ViewTemplates:
		m.view = v
		m.statusBar.SetKeys(TemplateViewKeys())
		m.loading = true
		cmd = tea.Batch(m.spinner.Tick, m.templatesView.LoadTemplates())
	case ViewDoctor:
		m.view = v
		m.statusBar.SetKeys(ViewKeys())
		m.loading = true
		cmd = tea.Batch(m.spinner.Tick, m.doctorView.RunChecks())
	case ViewInit:
		m.view = v
		m.statusBar.SetKeys(InitViewKeys())
		m.initView = views.NewInitModel()
		m.initView.SetSize(appWidth, appHeight-7)
		cmd = m.initView.Init()
	case ViewConfig:
		m.view = v
		m.statusBar.SetKeys(ViewKeys())
		cmd = m.configView.Load()
	default:
		m.view = v
	}

	return m, cmd
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}

func (m Model) updateTemplates(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc", "left", "h":
		// If previewing, let the view handle it to exit preview
		if m.templatesView.IsPreviewing() {
			var cmd tea.Cmd
			m.templatesView, cmd = m.templatesView.Update(msg)
			return m, cmd
		}
		// Otherwise go back to menu
		m.view = ViewMenu
		m.statusBar.SetKeys(DefaultMenuKeys())
		return m, nil
	case "r":
		if !m.templatesView.IsPreviewing() {
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, m.templatesView.RefreshTemplates())
		}
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

func (m Model) updateUpdateView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc", "left", "h", "enter":
		// Allow going back when done
		if m.updateView.IsDone() {
			m.view = ViewMenu
			m.statusBar.SetKeys(DefaultMenuKeys())
			return m, nil
		}
	}
	return m, nil
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

	// View title (when in submenu) - italic style under header
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	viewTitleStyle := lipgloss.NewStyle().Foreground(theme.Muted).Italic(true)
	now := time.Now()
	dateTime := now.Format("Mon Jan 2 3:04 PM MST")

	viewTitle := "Main Menu"
	switch m.view {
	case ViewTemplates:
		viewTitle = "Templates"
	case ViewDoctor:
		viewTitle = "Environment"
	case ViewInit:
		viewTitle = "New Project"
	case ViewConfig:
		viewTitle = "Configuration"
	case ViewUpdate:
		viewTitle = "Update"
	}

	// Show view title on left, datetime on right
	titleText := viewTitleStyle.Render(viewTitle)
	padding := appWidth - lipgloss.Width(viewTitle) - lipgloss.Width(dateTime)
	if padding < 1 {
		padding = 1
	}
	inner.WriteString(titleText + strings.Repeat(" ", padding) + mutedStyle.Render(dateTime))
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
	case ViewUpdate:
		content = m.updateView.View()
	}

	// Truncate or pad content to fixed height (top-justified)
	contentLines := strings.Split(content, "\n")
	maxContentLines := appHeight - 7 // header, subheader, top divider, hint, bottom divider, path, footer

	// Truncate if too long
	if len(contentLines) > maxContentLines {
		contentLines = contentLines[:maxContentLines]
	}

	// Pad at bottom if too short (top-justified)
	for len(contentLines) < maxContentLines {
		contentLines = append(contentLines, "")
	}

	inner.WriteString("\n")
	inner.WriteString(strings.Join(contentLines, "\n"))

	// Context hint (centered, above footer)
	hint := m.getContextHint()
	if hint != "" {
		hintPadding := (appWidth - lipgloss.Width(hint)) / 2
		if hintPadding < 0 {
			hintPadding = 0
		}
		inner.WriteString("\n")
		inner.WriteString(strings.Repeat(" ", hintPadding) + hint)
	}

	// Footer divider
	inner.WriteString("\n")
	inner.WriteString(Divider(appWidth))
	inner.WriteString("\n")

	// Folder path with disk space (centered)
	defaultDir := config.GetDefaultDirectory()
	var pathText string
	if defaultDir == "" {
		pathText = "No default project path"
	} else {
		sysInfo := doctor.GetSystemInfo()
		if sysInfo.Disk != "" {
			pathText = defaultDir + " (" + sysInfo.Disk + ")"
		} else {
			pathText = defaultDir
		}
	}
	// Truncate long paths from the left
	maxPathLen := appWidth - 4
	if len(pathText) > maxPathLen {
		pathText = "..." + pathText[len(pathText)-maxPathLen+3:]
	}
	// Center the path
	pathPadding := (appWidth - lipgloss.Width(pathText)) / 2
	if pathPadding < 0 {
		pathPadding = 0
	}
	inner.WriteString(strings.Repeat(" ", pathPadding) + mutedStyle.Render(pathText))
	inner.WriteString("\n")

	// Footer: centered command help
	inner.WriteString(m.statusBar.View())
	inner.WriteString("\n")

	// No border - clean Claude Code style
	return inner.String()
}

func (m Model) renderMenuView() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.menu.View())
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderLoading(msg string) string {
	style := lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1)
	return style.Render(m.spinner.View() + " " + msg) + "\n"
}

// getContextHint returns the appropriate hint for current view state
func (m Model) getContextHint() string {
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	switch m.view {
	case ViewMenu:
		return mutedStyle.Render("Or run: ") + keyStyle.Render("irl init \"your project\"")
	case ViewTemplates:
		if m.templatesView.IsPreviewing() {
			return keyStyle.Render("↑↓") + mutedStyle.Render(" scroll  ") + keyStyle.Render("←") + mutedStyle.Render(" back")
		}
		return keyStyle.Render("→") + mutedStyle.Render(" preview  ") + keyStyle.Render("r") + mutedStyle.Render(" refresh")
	case ViewInit:
		return keyStyle.Render("↑↓") + mutedStyle.Render(" navigate  ") + keyStyle.Render("→") + mutedStyle.Render(" select")
	case ViewDoctor:
		return ""
	case ViewConfig:
		return keyStyle.Render("e") + mutedStyle.Render(" edit")
	}
	return ""
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
