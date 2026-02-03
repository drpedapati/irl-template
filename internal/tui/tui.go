package tui

import (
	"fmt"
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
	"github.com/drpedapati/irl-template/pkg/editor"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// NewTemplatesAvailableMsg is sent when the update check completes
type NewTemplatesAvailableMsg struct {
	Count int
}

// SystemInfoLoadedMsg is sent when system info is loaded
type SystemInfoLoadedMsg struct {
	DiskInfo string
}

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
	templatesView   views.TemplatesModel
	projectsView    views.ProjectsModel
	folderView      views.FolderModel
	personalizeView views.PersonalizeModel
	doctorView      views.DoctorModel
	initView        views.InitModel
	configView      views.ConfigModel
	editorsView     views.EditorsModel
	helpView        views.HelpModel

	// Loading state
	loading bool
	spinner spinner.Model

	// Update check state
	newTemplateCount   int
	checkingForUpdates bool

	// Inline update progress (shown above footer)
	inlineUpdating       bool
	inlineUpdatePercent  float64
	inlineUpdateStatus   string
	inlineUpdateDone     bool
	inlineUpdateErr      error
	inlineUpdateCount    int

	// Cached values (avoid expensive calls in View())
	cachedDiskInfo   string
	cachedDefaultDir string

	// Twinkle animation for "What is IRL?" hint
	twinkleFrame int
}

// New creates a new TUI model
func New(version string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	m := Model{
		version:            version,
		header:             NewHeader(version),
		menu:               NewMenu(),
		statusBar:          NewStatusBar(),
		view:               ViewMenu,
		templatesView:      views.NewTemplatesModel(),
		projectsView:       views.NewProjectsModel(),
		folderView:         views.NewFolderModel(),
		personalizeView:    views.NewPersonalizeModel(),
		doctorView:         views.NewDoctorModel(),
		initView:           views.NewInitModel(),
		configView:         views.NewConfigModel(),
		helpView:           views.NewHelpModel(),
		spinner:            s,
		checkingForUpdates: true, // Will check on init
	}

	// Set fixed widths
	m.header.SetWidth(appWidth)
	m.menu.SetWidth(appWidth)
	m.statusBar.SetWidth(appWidth)
	m.templatesView.SetSize(appWidth, appHeight-7)
	m.projectsView.SetSize(appWidth, appHeight-7)
	m.folderView.SetSize(appWidth, appHeight-7)
	m.personalizeView.SetSize(appWidth, appHeight-7)
	m.doctorView.SetSize(appWidth, appHeight-7)
	m.initView.SetSize(appWidth, appHeight-7)
	m.configView.SetSize(appWidth, appHeight-7)
	m.helpView.SetSize(appWidth-2, appHeight-7) // -2 for left margin

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Check for new templates asynchronously on startup
	m.checkingForUpdates = true
	return tea.Batch(m.spinner.Tick, checkForNewTemplates(), loadSystemInfo(), m.animateTwinkle())
}

// loadSystemInfo loads system info in the background
func loadSystemInfo() tea.Cmd {
	return func() tea.Msg {
		sysInfo := doctor.GetSystemInfo()
		return SystemInfoLoadedMsg{DiskInfo: sysInfo.Disk}
	}
}

// checkForNewTemplates returns a command that checks for new templates
func checkForNewTemplates() tea.Cmd {
	return func() tea.Msg {
		count, _ := templates.CheckForNewTemplates()
		return NewTemplatesAvailableMsg{Count: count}
	}
}

// InlineUpdateTickMsg for progress animation
type InlineUpdateTickMsg struct{}

// clearInlineUpdateMsg hides the progress bar after completion
type clearInlineUpdateMsg struct{}

// TwinkleTickMsg for the "What is IRL?" animation
type TwinkleTickMsg struct{}

// InlineUpdateCompleteMsg when update finishes
type InlineUpdateCompleteMsg struct {
	Count int
	Err   error
}

func (m Model) animateInlineUpdate() tea.Cmd {
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return InlineUpdateTickMsg{}
	})
}

func (m Model) animateTwinkle() tea.Cmd {
	return tea.Tick(time.Millisecond*400, func(t time.Time) tea.Msg {
		return TwinkleTickMsg{}
	})
}

func (m Model) fetchTemplatesInline() tea.Cmd {
	return func() tea.Msg {
		list, err := templates.FetchTemplates()
		if err != nil {
			return InlineUpdateCompleteMsg{Err: err}
		}
		return InlineUpdateCompleteMsg{Count: len(list)}
	}
}

// getMenuKeys returns the appropriate menu keys based on current state
func (m Model) getMenuKeys() []KeyBinding {
	return DefaultMenuKeysWithUpdate(m.newTemplateCount, m.checkingForUpdates)
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
		case ViewProjects:
			return m.updateProjects(msg)
		case ViewTemplates:
			return m.updateTemplates(msg)
		case ViewEditors:
			return m.updateEditors(msg)
		case ViewDoctor:
			return m.updateDoctor(msg)
		case ViewInit:
			return m.updateInit(msg)
		case ViewConfig:
			return m.updateConfig(msg)
		case ViewFolder:
			return m.updateFolder(msg)
		case ViewPersonalize:
			return m.updatePersonalize(msg)
		case ViewHelp:
			return m.updateHelp(msg)
		}

	case spinner.TickMsg:
		if m.loading || m.checkingForUpdates {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case NewTemplatesAvailableMsg:
		m.checkingForUpdates = false
		m.newTemplateCount = msg.Count
		// Update status bar with new count
		if m.view == ViewMenu {
			m.statusBar.SetKeys(m.getMenuKeys())
		}

	case SystemInfoLoadedMsg:
		m.cachedDiskInfo = msg.DiskInfo
		m.cachedDefaultDir = config.GetDefaultDirectory()

	case InlineUpdateTickMsg:
		if m.inlineUpdating && !m.inlineUpdateDone {
			// Animate progress up to 90% while waiting
			if m.inlineUpdatePercent < 0.9 {
				m.inlineUpdatePercent += 0.03
				if m.inlineUpdatePercent > 0.3 && m.inlineUpdatePercent < 0.35 {
					m.inlineUpdateStatus = "Fetching..."
				} else if m.inlineUpdatePercent > 0.6 && m.inlineUpdatePercent < 0.65 {
					m.inlineUpdateStatus = "Caching..."
				}
			}
			return m, m.animateInlineUpdate()
		}

	case TwinkleTickMsg:
		// Advance twinkle animation frame (only when on menu)
		if m.view == ViewMenu {
			m.twinkleFrame = (m.twinkleFrame + 1) % 12
		}
		return m, m.animateTwinkle()

	case InlineUpdateCompleteMsg:
		m.inlineUpdateDone = true
		m.inlineUpdatePercent = 1.0
		m.inlineUpdateCount = msg.Count
		m.inlineUpdateErr = msg.Err
		if msg.Err != nil {
			m.inlineUpdateStatus = "Failed"
		} else {
			m.inlineUpdateStatus = fmt.Sprintf("%d templates", msg.Count)
		}
		// Reset update count and refresh templates
		m.newTemplateCount = 0
		m.templatesView = views.NewTemplatesModel()
		m.statusBar.SetKeys(m.getMenuKeys())
		// Clear inline update after a delay
		return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
			return clearInlineUpdateMsg{}
		})

	case clearInlineUpdateMsg:
		m.inlineUpdating = false
		m.inlineUpdateDone = false
		m.inlineUpdatePercent = 0

	// Handle view-specific messages
	case views.TemplatesLoadedMsg:
		m.templatesView, _ = m.templatesView.Update(msg)
		m.loading = false

	case views.ProjectsLoadedMsg:
		m.projectsView, _ = m.projectsView.Update(msg)
		m.loading = false

	case views.DoctorResultMsg:
		m.doctorView, _ = m.doctorView.Update(msg)
		m.loading = false

	case views.EditorsLoadedMsg:
		m.editorsView, _ = m.editorsView.Update(msg)
		m.loading = false

	case views.InitCompleteMsg:
		m.initView, _ = m.initView.Update(msg)

	case views.InitTemplatesLoadedMsg:
		m.initView, _ = m.initView.Update(msg)

	case views.InitProjectCreatedMsg:
		m.initView, _ = m.initView.Update(msg)

	case views.BackToMenuMsg:
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())

	case views.HelpLoadedMsg:
		m.helpView, _ = m.helpView.Update(msg)

	case editor.EditorFinishedMsg, editor.EditorOpenedMsg:
		// Pass editor messages to the appropriate view
		if m.view == ViewProjects {
			var cmd tea.Cmd
			m.projectsView, cmd = m.projectsView.Update(msg)
			return m, cmd
		}
		if m.view == ViewInit {
			var cmd tea.Cmd
			m.initView, cmd = m.initView.Update(msg)
			return m, cmd
		}
		if m.view == ViewTemplates {
			var cmd tea.Cmd
			m.templatesView, cmd = m.templatesView.Update(msg)
			return m, cmd
		}

	}

	return m, tea.Batch(cmds...)
}

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "up":
		m.menu.Up()
	case "down":
		m.menu.Down()
	case "enter", "right":
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
	case "e":
		if v, ok := m.menu.SelectByKey("e"); ok {
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
		// Start inline update (don't switch views)
		if !m.inlineUpdating {
			m.inlineUpdating = true
			m.inlineUpdatePercent = 0
			m.inlineUpdateStatus = "Connecting..."
			m.inlineUpdateDone = false
			m.inlineUpdateErr = nil
			return m, tea.Batch(m.animateInlineUpdate(), m.fetchTemplatesInline())
		}
	case "p":
		if v, ok := m.menu.SelectByKey("p"); ok {
			return m.selectView(v)
		}
	case "f":
		if v, ok := m.menu.SelectByKey("f"); ok {
			return m.selectView(v)
		}
	case "?":
		if v, ok := m.menu.SelectByKey("?"); ok {
			return m.selectView(v)
		}
	case "i":
		return m.selectView(ViewPersonalize)
	}
	return m, nil
}

func (m Model) selectView(v ViewType) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch v {
	case ViewDocs:
		// Open docs in browser, stay on menu
		openBrowser("https://www.irloop.org")
		return m, nil
	case ViewProjects:
		m.view = v
		m.statusBar.SetKeys(ProjectsFilterKeys()) // Start in filter mode
		m.loading = true
		cmd = tea.Batch(m.spinner.Tick, m.projectsView.ScanProjects())
	case ViewTemplates:
		m.view = v
		m.statusBar.SetKeys(TemplateViewKeys())
		m.loading = true
		cmd = tea.Batch(m.spinner.Tick, m.templatesView.LoadTemplates())
	case ViewEditors:
		m.view = v
		m.statusBar.SetKeys(ViewKeys())
		m.editorsView = views.NewEditorsModel()
		m.editorsView.SetSize(appWidth, appHeight-7)
		m.loading = true
		cmd = tea.Batch(m.spinner.Tick, m.editorsView.DetectApps())
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
	case ViewFolder:
		m.view = v
		m.statusBar.SetKeys(ViewKeys())
		m.folderView = views.NewFolderModel()
		m.folderView.SetSize(appWidth, appHeight-7)
	case ViewPersonalize:
		m.view = v
		m.statusBar.SetKeys(ViewKeys())
		m.personalizeView = views.NewPersonalizeModel()
		m.personalizeView.SetSize(appWidth, appHeight-7)
		cmd = m.personalizeView.Init()
	case ViewHelp:
		m.view = v
		m.statusBar.SetKeys(ViewKeys())
		m.helpView = views.NewHelpModel()
		m.helpView.SetSize(appWidth-2, appHeight-7)
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

// openInEditor opens a directory in the default editor
func openInEditor(path string) {
	var cmd *exec.Cmd
	// Try VS Code first, then fall back to opening the folder
	if _, err := exec.LookPath("code"); err == nil {
		cmd = exec.Command("code", path)
	} else if _, err := exec.LookPath("cursor"); err == nil {
		cmd = exec.Command("cursor", path)
	} else {
		// Fall back to opening the folder
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", path)
		case "linux":
			cmd = exec.Command("xdg-open", path)
		case "windows":
			cmd = exec.Command("explorer", path)
		}
	}
	if cmd != nil {
		cmd.Start()
	}
}

func (m Model) updateProjects(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If viewing project details or confirming delete, let the view handle all keys
	if m.projectsView.IsViewing() || m.projectsView.IsConfirmingDelete() {
		var cmd tea.Cmd
		m.projectsView, cmd = m.projectsView.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "left":
		// Go back to menu (only in filter mode, or from selection mode explicitly)
		if m.projectsView.IsFilterMode() {
			m.view = ViewMenu
			m.statusBar.SetKeys(m.getMenuKeys())
			return m, nil
		}
		// In selection mode, left goes back to filter mode first
		// (handled by projects view)
	}

	// Pass all other keys to the projects view
	var cmd tea.Cmd
	m.projectsView, cmd = m.projectsView.Update(msg)

	// Update status bar based on mode change
	if m.projectsView.IsFilterMode() {
		m.statusBar.SetKeys(ProjectsFilterKeys())
	} else {
		m.statusBar.SetKeys(ProjectsViewKeys())
	}

	return m, cmd
}

func (m Model) updateTemplates(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If in copy, edit, delete, or post-copy edit mode, let the view handle all keys
	if m.templatesView.IsCopying() || m.templatesView.IsEditing() || m.templatesView.IsDeleting() || m.templatesView.IsPostCopyEdit() {
		var cmd tea.Cmd
		m.templatesView, cmd = m.templatesView.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		// If previewing, let the view handle it to exit preview
		if m.templatesView.IsPreviewing() {
			var cmd tea.Cmd
			m.templatesView, cmd = m.templatesView.Update(msg)
			return m, cmd
		}
		// Two-stage escape: clear filter first, then go back
		if m.templatesView.HasFilterText() {
			m.templatesView.ClearFilter()
			return m, nil
		}
		// Otherwise go back to menu
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
		return m, nil
	case "left":
		// If previewing, let the view handle it to exit preview
		if m.templatesView.IsPreviewing() {
			var cmd tea.Cmd
			m.templatesView, cmd = m.templatesView.Update(msg)
			return m, cmd
		}
		// Go back to menu
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
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

func (m Model) updateEditors(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If in plan editor selection mode, let the view handle all keys
	if m.editorsView.IsSelectingPlanEditor() {
		var cmd tea.Cmd
		m.editorsView, cmd = m.editorsView.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc", "left":
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
		return m, nil
	}

	var cmd tea.Cmd
	m.editorsView, cmd = m.editorsView.Update(msg)
	return m, cmd
}

func (m Model) updateDoctor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc", "left":
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
		return m, nil
	}

	var cmd tea.Cmd
	m.doctorView, cmd = m.doctorView.Update(msg)
	return m, cmd
}

func (m Model) updateInit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "left":
		if m.initView.CanGoBack() {
			var cmd tea.Cmd
			m.initView, cmd = m.initView.Update(msg)
			return m, cmd
		}
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
		return m, nil
	}

	var cmd tea.Cmd
	m.initView, cmd = m.initView.Update(msg)

	// Check if init completed and we should go back to menu
	if m.initView.Done() {
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
	}

	return m, cmd
}

func (m Model) updateConfig(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc", "left":
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
		// Refresh cached directory in case it changed
		m.cachedDefaultDir = config.GetDefaultDirectory()
		return m, nil
	}

	var cmd tea.Cmd
	m.configView, cmd = m.configView.Update(msg)
	return m, cmd
}

func (m Model) updateFolder(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
		return m, nil
	}

	var cmd tea.Cmd
	m.folderView, cmd = m.folderView.Update(msg)

	// If saved or wants back, go back to menu
	if m.folderView.IsSaved() {
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
		// Refresh cached directory
		m.cachedDefaultDir = config.GetDefaultDirectory()
	} else if m.folderView.WantsBack() {
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
	}

	return m, cmd
}

func (m Model) updatePersonalize(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
		return m, nil
	}

	var cmd tea.Cmd
	m.personalizeView, cmd = m.personalizeView.Update(msg)

	// If saved, go back to menu
	if m.personalizeView.IsSaved() {
		m.view = ViewMenu
		m.statusBar.SetKeys(m.getMenuKeys())
	}

	return m, cmd
}

func (m Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		// If at help menu, go back to main menu
		if m.helpView.IsAtMenu() {
			m.view = ViewMenu
			m.statusBar.SetKeys(m.getMenuKeys())
			return m, nil
		}
		// Otherwise let helpView handle (goes back to help menu from slides)
		var cmd tea.Cmd
		m.helpView, cmd = m.helpView.Update(msg)
		return m, cmd
	case "left", "h":
		// If at help menu, go back to main menu
		if m.helpView.IsAtMenu() {
			m.view = ViewMenu
			m.statusBar.SetKeys(m.getMenuKeys())
			return m, nil
		}
		// Otherwise let helpView handle navigation
		var cmd tea.Cmd
		m.helpView, cmd = m.helpView.Update(msg)
		return m, cmd
	}

	// Pass all other keys to the help view
	var cmd tea.Cmd
	m.helpView, cmd = m.helpView.Update(msg)
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

	// View title (when in submenu) - italic style under header
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	viewTitleStyle := lipgloss.NewStyle().Foreground(theme.Muted).Italic(true)
	now := time.Now()
	dateTime := now.Format("Mon Jan 2 3:04 PM MST")

	viewTitle := "Main Menu"
	switch m.view {
	case ViewProjects:
		viewTitle = "Projects"
	case ViewFolder:
		viewTitle = "Default Folder"
	case ViewTemplates:
		if name := m.templatesView.PreviewingName(); name != "" {
			viewTitle = "Template: " + name
		} else {
			viewTitle = "Templates"
		}
	case ViewEditors:
		viewTitle = "Editors & Utilities"
	case ViewDoctor:
		viewTitle = "Environment"
	case ViewInit:
		viewTitle = "New Project"
	case ViewConfig:
		viewTitle = "Configuration"
	case ViewPersonalize:
		viewTitle = "Profile"
	case ViewHelp:
		viewTitle = "Help"
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

	// Animated "What is IRL?" hint centered above menu (only on main menu)
	if m.view == ViewMenu {
		inner.WriteString(m.renderTwinkleHint())
	}

	// Content area
	content := ""
	switch m.view {
	case ViewMenu:
		content = m.renderMenuView()
	case ViewProjects:
		if m.loading {
			content = m.renderLoading("Scanning for projects...")
		} else {
			content = m.projectsView.View()
		}
	case ViewTemplates:
		if m.loading {
			content = m.renderLoading("Loading templates...")
		} else {
			content = m.templatesView.View()
		}
	case ViewEditors:
		if m.loading {
			content = m.renderLoading("Detecting applications...")
		} else {
			content = m.editorsView.View()
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
	case ViewFolder:
		content = m.folderView.View()
	case ViewPersonalize:
		content = m.personalizeView.View()
	case ViewHelp:
		content = m.helpView.View()
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

	// Inline update progress bar (above footer)
	if m.inlineUpdating {
		inner.WriteString("\n")
		progressBar := m.renderInlineProgress()
		progressPadding := (appWidth - lipgloss.Width(progressBar)) / 2
		if progressPadding < 0 {
			progressPadding = 0
		}
		inner.WriteString(strings.Repeat(" ", progressPadding) + progressBar)
	}

	// Footer divider
	inner.WriteString("\n")
	inner.WriteString(Divider(appWidth))
	inner.WriteString("\n")

	// Folder path with disk space (centered) and [f] indicator
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	defaultDir := m.cachedDefaultDir
	var pathText string
	if defaultDir == "" {
		pathText = "No default project path"
	} else {
		if m.cachedDiskInfo != "" {
			pathText = defaultDir + " (" + m.cachedDiskInfo + ")"
		} else {
			pathText = defaultDir
		}
	}
	// Truncate long paths from the left (account for [f] prefix)
	maxPathLen := appWidth - 8
	if len(pathText) > maxPathLen {
		pathText = "..." + pathText[len(pathText)-maxPathLen+3:]
	}
	// Build the full path line with [f] indicator
	pathLine := keyStyle.Render("[f]") + " " + mutedStyle.Render(pathText)
	// Center the path
	pathPadding := (appWidth - lipgloss.Width(pathLine)) / 2
	if pathPadding < 0 {
		pathPadding = 0
	}
	inner.WriteString(strings.Repeat(" ", pathPadding) + pathLine)
	inner.WriteString("\n")

	// Footer: centered command help
	inner.WriteString(m.statusBar.View())
	inner.WriteString("\n")

	// Add small left margin for visual breathing room
	marginStyle := lipgloss.NewStyle().PaddingLeft(2)
	return marginStyle.Render(inner.String())
}

func (m Model) renderMenuView() string {
	return m.menu.View()
}

func (m Model) renderLoading(msg string) string {
	style := lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1)
	return style.Render(m.spinner.View() + " " + msg) + "\n"
}

func (m Model) renderInlineProgress() string {
	// Subtle compact progress bar
	barWidth := 20
	filled := int(m.inlineUpdatePercent * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	barStyle := lipgloss.NewStyle().Foreground(theme.Accent)
	emptyStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	statusStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	var icon string
	if m.inlineUpdateDone {
		if m.inlineUpdateErr != nil {
			icon = lipgloss.NewStyle().Foreground(theme.Error).Render("✗")
		} else {
			icon = lipgloss.NewStyle().Foreground(theme.Success).Render("✓")
		}
	} else {
		icon = "⟳"
	}

	bar := barStyle.Render(strings.Repeat("━", filled)) + emptyStyle.Render(strings.Repeat("─", barWidth-filled))
	return icon + " " + bar + " " + statusStyle.Render(m.inlineUpdateStatus)
}

// renderTwinkleHint renders the animated "What is IRL?" hint
func (m Model) renderTwinkleHint() string {
	// Sparkle characters that cycle through
	sparkles := []string{"✦", "✧", "⋆", "✶", "✦", "✧", "⋆", "✶", "✦", "✧", "⋆", "✶"}

	// Colors that shift subtly - warm palette
	colors := []lipgloss.Color{
		lipgloss.Color("#E07A5F"), // coral
		lipgloss.Color("#D4A373"), // warm tan
		lipgloss.Color("#E9C46A"), // golden
		lipgloss.Color("#F4A261"), // sandy
		lipgloss.Color("#E07A5F"), // coral
		lipgloss.Color("#C9ADA7"), // dusty rose
		lipgloss.Color("#E07A5F"), // coral
		lipgloss.Color("#D4A373"), // warm tan
		lipgloss.Color("#E9C46A"), // golden
		lipgloss.Color("#F4A261"), // sandy
		lipgloss.Color("#E07A5F"), // coral
		lipgloss.Color("#C9ADA7"), // dusty rose
	}

	// Get current frame sparkles and colors
	leftIdx := m.twinkleFrame
	rightIdx := (m.twinkleFrame + 6) % 12 // Offset for asymmetry

	leftSparkle := sparkles[leftIdx]
	rightSparkle := sparkles[rightIdx]
	leftColor := colors[leftIdx]
	rightColor := colors[rightIdx]

	// Styles
	leftStyle := lipgloss.NewStyle().Foreground(leftColor)
	rightStyle := lipgloss.NewStyle().Foreground(rightColor)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	textStyle := lipgloss.NewStyle().Foreground(theme.Primary).Bold(true)

	// Build the hint: ✦ [?] Learn IRL Fast! ✧
	hint := leftStyle.Render(leftSparkle) + "  " +
		keyStyle.Render("[?]") + " " +
		textStyle.Render("Learn IRL Fast!") + "  " +
		rightStyle.Render(rightSparkle)

	// Center it with even spacing above and below
	hintWidth := lipgloss.Width(hint)
	padding := (appWidth - hintWidth) / 2
	if padding < 0 {
		padding = 0
	}

	// One blank line above, hint, then menu adds its own spacing below
	return "\n\n" + strings.Repeat(" ", padding) + hint
}

// getContextHint returns the appropriate hint for current view state
func (m Model) getContextHint() string {
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	switch m.view {
	case ViewMenu:
		return ""
	case ViewProjects:
		if m.projectsView.IsFilterMode() {
			return keyStyle.Render("↓") + mutedStyle.Render(" to select projects")
		}
		// Selection mode: show action hints
		baseHints := keyStyle.Render("e") + mutedStyle.Render(" edit  ")
		editorHints := m.projectsView.GetEditorHints()
		if editorHints != "" {
			return baseHints + editorHints
		}
		return baseHints
	case ViewTemplates:
		if m.templatesView.IsCopying() || m.templatesView.IsEditing() || m.templatesView.IsDeleting() {
			return "" // Modal has its own hints
		}
		if m.templatesView.IsPreviewing() {
			return keyStyle.Render("↑↓") + mutedStyle.Render(" scroll  ") + keyStyle.Render("g") + mutedStyle.Render(" GitHub  ") + keyStyle.Render("←") + mutedStyle.Render(" back")
		}
		// Context-sensitive: show edit/del only for custom templates
		if m.templatesView.SelectedIsCustom() {
			return keyStyle.Render("t") + mutedStyle.Render(" new  ") + keyStyle.Render("e") + mutedStyle.Render(" edit  ") + keyStyle.Render("x") + mutedStyle.Render(" del  ") + mutedStyle.Render("│ ") + keyStyle.Render("a") + mutedStyle.Render("ll ") + keyStyle.Render("d") + mutedStyle.Render("efault ") + keyStyle.Render("c") + mutedStyle.Render("ustom  ") + keyStyle.Render("r") + mutedStyle.Render(" refresh")
		}
		return keyStyle.Render("t") + mutedStyle.Render(" new  ") + keyStyle.Render("→") + mutedStyle.Render(" preview  ") + mutedStyle.Render("│ ") + keyStyle.Render("a") + mutedStyle.Render("ll ") + keyStyle.Render("d") + mutedStyle.Render("efault ") + keyStyle.Render("c") + mutedStyle.Render("ustom  ") + keyStyle.Render("r") + mutedStyle.Render(" refresh")
	case ViewInit:
		return keyStyle.Render("↑↓") + mutedStyle.Render(" navigate  ") + keyStyle.Render("Enter") + mutedStyle.Render(" select")
	case ViewEditors:
		return keyStyle.Render("↑↓") + mutedStyle.Render(" navigate  ") + keyStyle.Render("1-3") + mutedStyle.Render(" filter  ") + keyStyle.Render("Enter") + mutedStyle.Render(" open/install")
	case ViewDoctor:
		return ""
	case ViewConfig:
		return keyStyle.Render("e") + mutedStyle.Render(" edit")
	case ViewFolder:
		return keyStyle.Render("←") + mutedStyle.Render(" up a level  ") + keyStyle.Render("↑↓") + mutedStyle.Render(" navigate  ") + keyStyle.Render("Enter") + mutedStyle.Render(" select")
	case ViewPersonalize:
		return keyStyle.Render("Tab") + mutedStyle.Render(" next item  ") + keyStyle.Render("Enter") + mutedStyle.Render(" save/select  ") + keyStyle.Render("y/n") + mutedStyle.Render(" confirm")
	case ViewHelp:
		// Help view has its own footer with navigation hints
		return ""
	}
	return ""
}

// Run starts the TUI
func Run(version string) error {
	// Clear screen and move cursor to top
	fmt.Print("\033[2J\033[H")

	p := tea.NewProgram(
		New(version),
		// No alt screen - renders inline at current cursor position
	)

	_, err := p.Run()
	return err
}
