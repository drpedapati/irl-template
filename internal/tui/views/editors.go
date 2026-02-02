package views

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/editor"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// AppCategory groups applications by type
type AppCategory int

const (
	CategoryEditor AppCategory = iota
	CategoryIDE
	CategoryTool
)

func (c AppCategory) String() string {
	switch c {
	case CategoryEditor:
		return "Editor"
	case CategoryIDE:
		return "IDE"
	case CategoryTool:
		return "Tool"
	default:
		return ""
	}
}

// AppInfo represents an editor or utility application
type AppInfo struct {
	Name        string
	Description string
	Cmd         string      // Command line name
	AppName     string      // macOS .app name (if different)
	Key         string      // Hotkey for quick access (used in project actions)
	Category    AppCategory
	Installed   bool
	Favorite    bool   // User has marked this as a favorite
	URL         string // Website for installation
	IsTerminal  bool   // True for terminal-based editors (nvim, helix, fresh)
}

// GetInstalledEditors returns installed editors/IDEs for project actions.
// If favorites are set, only returns favorites. Otherwise returns all installed.
// This is the single source of truth for editor lists across the app.
func GetInstalledEditors() []AppInfo {
	apps := getAllApps()
	favorites := config.GetFavoriteEditors()
	hasFavorites := len(favorites) > 0

	var installed []AppInfo
	for _, app := range apps {
		// Only include editors and IDEs that are installed and have hotkeys
		if app.Installed && app.Key != "" && (app.Category == CategoryEditor || app.Category == CategoryIDE) {
			// If favorites are set, only include favorites
			if hasFavorites && !app.Favorite {
				continue
			}
			installed = append(installed, app)
		}
	}
	return installed
}

// GetInstalledTools returns installed tools for project actions (Finder, Terminal, etc.)
// Tools are always shown regardless of favorites setting (favorites only applies to editors/IDEs).
func GetInstalledTools() []AppInfo {
	apps := getAllApps()

	var installed []AppInfo
	for _, app := range apps {
		if app.Installed && app.Key != "" && app.Category == CategoryTool {
			installed = append(installed, app)
		}
	}
	return installed
}

// GetAllInstalledWithKeys returns all installed apps that have hotkeys assigned
func GetAllInstalledWithKeys() []AppInfo {
	apps := getAllApps()
	var installed []AppInfo
	for _, app := range apps {
		if app.Installed && app.Key != "" {
			installed = append(installed, app)
		}
	}
	return installed
}

// PlanEditorMode is the plan editor selection mode
type PlanEditorMode int

const (
	PlanEditorModeNone PlanEditorMode = iota
	PlanEditorModeTerminal
	PlanEditorModeGUI
)

// EditorsModel displays and manages editors/utilities
type EditorsModel struct {
	width    int
	height   int
	apps     []AppInfo
	cursor   int
	scroll   int
	loaded   bool
	message  string
	category AppCategory // Filter by category (-1 = all)

	// Plan editor configuration
	planEditorMode   PlanEditorMode
	planEditorCursor int
	currentPlanEditor string
	currentPlanEditorType string
}

// Calculate visible items based on available height
func (m EditorsModel) visibleItems() int {
	// Reserve lines for: tabs(2) + footer(3) + scroll(1) = 6
	available := m.height - 6
	if available < 5 {
		available = 5
	}
	if available > 20 {
		available = 20
	}
	return available
}

// EditorsLoadedMsg is sent when editor detection completes
type EditorsLoadedMsg struct {
	Apps []AppInfo
}

// NewEditorsModel creates a new editors view
func NewEditorsModel() EditorsModel {
	return EditorsModel{
		category: -1, // Show all by default
	}
}

// SetSize sets the view dimensions
func (m *EditorsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// DetectApps returns a command that detects available applications
func (m *EditorsModel) DetectApps() tea.Cmd {
	return func() tea.Msg {
		apps := getAllApps()
		return EditorsLoadedMsg{Apps: apps}
	}
}

// IsSelectingPlanEditor returns true if in plan editor selection mode
func (m EditorsModel) IsSelectingPlanEditor() bool {
	return m.planEditorMode != PlanEditorModeNone
}

// getAllApps returns all known applications with their install status
func getAllApps() []AppInfo {
	apps := []AppInfo{
		// Code Editors (ordered by popularity/preference)
		{Name: "Cursor", Description: "AI-powered code editor", Cmd: "cursor", AppName: "Cursor", Key: "c", Category: CategoryEditor, URL: "https://cursor.sh"},
		{Name: "VS Code", Description: "Microsoft's popular editor", Cmd: "code", AppName: "Visual Studio Code", Key: "v", Category: CategoryEditor, URL: "https://code.visualstudio.com"},
		{Name: "Zed", Description: "Fast, collaborative editor", Cmd: "zed", AppName: "Zed", Key: "z", Category: CategoryEditor, URL: "https://zed.dev"},
		{Name: "Sublime Text", Description: "Lightweight and fast", Cmd: "subl", AppName: "Sublime Text", Key: "s", Category: CategoryEditor, URL: "https://sublimetext.com"},
		{Name: "Neovim", Description: "Modern terminal editor", Cmd: "nvim", Key: "n", Category: CategoryEditor, URL: "https://neovim.io", IsTerminal: true},
		{Name: "Helix", Description: "Modal terminal editor", Cmd: "hx", Key: "x", Category: CategoryEditor, URL: "https://helix-editor.com", IsTerminal: true},
		{Name: "Fresh", Description: "Intuitive terminal editor", Cmd: "fresh", Key: "h", Category: CategoryEditor, URL: "https://getfresh.dev", IsTerminal: true},

		// IDEs
		{Name: "Positron", Description: "Data science IDE from Posit", Cmd: "positron", AppName: "Positron", Key: "p", Category: CategoryIDE, URL: "https://github.com/posit-dev/positron"},
		{Name: "RStudio", Description: "IDE for R programming", Cmd: "rstudio", AppName: "RStudio", Key: "r", Category: CategoryIDE, URL: "https://posit.co/products/open-source/rstudio"},
		{Name: "PyCharm", Description: "Python IDE from JetBrains", Cmd: "pycharm", AppName: "PyCharm", Key: "y", Category: CategoryIDE, URL: "https://jetbrains.com/pycharm"},

		// Tools
		{Name: "Finder", Description: "Open folder in Finder", Cmd: "open", Key: "f", Category: CategoryTool, URL: ""},
		{Name: "Terminal", Description: "macOS Terminal app", Cmd: "terminal", Key: "t", Category: CategoryTool, URL: ""},
		{Name: "iTerm2", Description: "Enhanced terminal for macOS", Cmd: "iterm", AppName: "iTerm", Key: "i", Category: CategoryTool, URL: "https://iterm2.com"},
		{Name: "Warp", Description: "Modern terminal with AI", Cmd: "warp", AppName: "Warp", Key: "w", Category: CategoryTool, URL: "https://warp.dev"},
	}

	// Check installation status and favorites
	favorites := config.GetFavoriteEditors()
	for i := range apps {
		apps[i].Installed = isAppInstalled(apps[i])
		apps[i].Favorite = isFavorite(apps[i].Cmd, favorites)
	}

	return apps
}

func isFavorite(cmd string, favorites []string) bool {
	for _, f := range favorites {
		if f == cmd {
			return true
		}
	}
	return false
}

func isAppInstalled(app AppInfo) bool {
	// Special cases - always available on macOS
	if app.Cmd == "terminal" || app.Cmd == "open" {
		return runtime.GOOS == "darwin"
	}

	// Check command line
	if _, err := exec.LookPath(app.Cmd); err == nil {
		return true
	}

	// Check macOS .app
	if runtime.GOOS == "darwin" && app.AppName != "" {
		paths := []string{
			filepath.Join("/Applications", app.AppName+".app"),
			filepath.Join(os.Getenv("HOME"), "Applications", app.AppName+".app"),
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return true
			}
		}
	}

	return false
}

func (m EditorsModel) filteredApps() []AppInfo {
	if m.category < 0 {
		return m.apps
	}
	var filtered []AppInfo
	for _, app := range m.apps {
		if app.Category == m.category {
			filtered = append(filtered, app)
		}
	}
	return filtered
}

// Update handles messages
func (m EditorsModel) Update(msg tea.Msg) (EditorsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case EditorsLoadedMsg:
		m.apps = msg.Apps
		m.loaded = true
		m.cursor = 0
		m.scroll = 0
		m.currentPlanEditor = config.GetPlanEditor()
		m.currentPlanEditorType = config.GetPlanEditorType()
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// Handle plan editor selection mode
		if m.planEditorMode != PlanEditorModeNone {
			return m.updatePlanEditorSelection(msg)
		}

		filtered := m.filteredApps()
		visibleItems := m.visibleItems()

		switch key {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}
			m.message = ""
			return m, nil
		case "down":
			if m.cursor < len(filtered)-1 {
				m.cursor++
				if m.cursor >= m.scroll+visibleItems {
					m.scroll = m.cursor - visibleItems + 1
				}
			}
			m.message = ""
			return m, nil
		case "a":
			m.category = -1 // All
			m.cursor = 0
			m.scroll = 0
			m.message = ""
			return m, nil
		case "1":
			m.category = CategoryEditor
			m.cursor = 0
			m.scroll = 0
			m.message = ""
			return m, nil
		case "2":
			m.category = CategoryIDE
			m.cursor = 0
			m.scroll = 0
			m.message = ""
			return m, nil
		case "3":
			m.category = CategoryTool
			m.cursor = 0
			m.scroll = 0
			m.message = ""
			return m, nil
		case "t":
			// Select terminal editor for plans
			m.planEditorMode = PlanEditorModeTerminal
			m.planEditorCursor = 0
			m.message = ""
			return m, nil
		case "g":
			// Select GUI editor for plans
			m.planEditorMode = PlanEditorModeGUI
			m.planEditorCursor = 0
			m.message = ""
			return m, nil
		case "0":
			// Reset to auto-detect
			if err := config.SetPlanEditor("auto", ""); err != nil {
				m.message = "Failed to save preference"
			} else {
				m.currentPlanEditor = "auto"
				m.currentPlanEditorType = ""
				m.message = "Plan editor set to auto-detect"
			}
			return m, nil
		case "enter":
			// Test launch for installed apps, open website for uninstalled
			if m.cursor < len(filtered) {
				app := filtered[m.cursor]
				if app.Installed {
					m.launchApp(app)
				} else if app.URL != "" {
					if openURL(app.URL) {
						m.message = "Opening " + app.Name + " website..."
					} else {
						m.message = "Could not open browser"
					}
				}
			}
			return m, nil
		case "b":
			// Open website for any app
			if m.cursor < len(filtered) {
				app := filtered[m.cursor]
				if app.URL != "" {
					if openURL(app.URL) {
						m.message = "Opening " + app.Name + " website..."
					} else {
						m.message = "Could not open browser"
					}
				} else {
					m.message = "No website available"
				}
			}
			return m, nil
		case " ":
			// Toggle favorite for installed apps
			if m.cursor < len(filtered) {
				app := filtered[m.cursor]
				if app.Installed {
					if err := config.ToggleFavoriteEditor(app.Cmd); err != nil {
						m.message = "Failed to save preference"
					} else {
						// Refresh app list to update Favorite status
						m.apps = getAllApps()
						if app.Favorite {
							m.message = "Removed " + app.Name + " from favorites"
						} else {
							m.message = "Added " + app.Name + " to favorites"
						}
					}
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// updatePlanEditorSelection handles keys in plan editor selection mode
func (m EditorsModel) updatePlanEditorSelection(msg tea.KeyMsg) (EditorsModel, tea.Cmd) {
	key := msg.String()

	var editors []editor.Editor
	var editorType string
	if m.planEditorMode == PlanEditorModeTerminal {
		editors = editor.GetAvailableTerminal()
		editorType = "terminal"
	} else {
		editors = editor.GetAvailableGUI()
		editorType = "gui"
	}

	switch key {
	case "up":
		if m.planEditorCursor > 0 {
			m.planEditorCursor--
		}
		return m, nil
	case "down":
		if m.planEditorCursor < len(editors)-1 {
			m.planEditorCursor++
		}
		return m, nil
	case "enter":
		// Save selection
		if m.planEditorCursor < len(editors) {
			selected := editors[m.planEditorCursor]
			if err := config.SetPlanEditor(selected.Command, editorType); err != nil {
				m.message = "Failed to save preference"
			} else {
				m.currentPlanEditor = selected.Command
				m.currentPlanEditorType = editorType
				m.message = "Plan editor set to " + selected.Name
			}
		}
		m.planEditorMode = PlanEditorModeNone
		return m, nil
	case "esc":
		m.planEditorMode = PlanEditorModeNone
		return m, nil
	}

	return m, nil
}

func (m *EditorsModel) launchApp(app AppInfo) {
	var cmd *exec.Cmd

	switch app.Cmd {
	case "terminal":
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", "-a", "Terminal")
		}
	case "open":
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", ".")
		}
	default:
		if runtime.GOOS == "darwin" && app.AppName != "" {
			// Try to open the .app
			cmd = exec.Command("open", "-a", app.AppName)
		} else {
			cmd = exec.Command(app.Cmd)
		}
	}

	if cmd != nil {
		if err := cmd.Start(); err != nil {
			m.message = "Failed to launch " + app.Name
		} else {
			m.message = "Launched " + app.Name
		}
	}
}

// OpenProjectWith launches an app to open a project directory.
// Returns an error message or empty string on success.
func OpenProjectWith(app AppInfo, projectPath string) string {
	var cmd *exec.Cmd

	// Handle terminal-based editors specially - they need a terminal window
	if app.IsTerminal {
		if runtime.GOOS == "darwin" {
			// Open Terminal.app, cd to project, and run the editor
			script := `tell application "Terminal" to do script "cd '` + projectPath + `' && ` + app.Cmd + `"`
			cmd = exec.Command("osascript", "-e", script)
		} else {
			// Linux: open terminal with editor command
			cmd = exec.Command("x-terminal-emulator", "-e", "sh", "-c", "cd '"+projectPath+"' && "+app.Cmd)
		}
		if cmd != nil {
			if err := cmd.Start(); err != nil {
				return "Failed to launch " + app.Name
			}
		}
		return ""
	}

	switch app.Cmd {
	case "terminal":
		if runtime.GOOS == "darwin" {
			script := `tell application "Terminal" to do script "cd '` + projectPath + `'"`
			cmd = exec.Command("osascript", "-e", script)
		} else {
			cmd = exec.Command("x-terminal-emulator", "--working-directory", projectPath)
		}
	case "open":
		// Finder
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", projectPath)
		} else if runtime.GOOS == "linux" {
			cmd = exec.Command("xdg-open", projectPath)
		} else {
			cmd = exec.Command("explorer", projectPath)
		}
	default:
		if runtime.GOOS == "darwin" && app.AppName != "" {
			cmd = exec.Command("open", "-a", app.AppName, projectPath)
		} else {
			cmd = exec.Command(app.Cmd, projectPath)
		}
	}

	if cmd != nil {
		if err := cmd.Start(); err != nil {
			return "Failed to launch " + app.Name
		}
	}
	return ""
}

func openURL(url string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}
	if cmd != nil {
		return cmd.Start() == nil
	}
	return false
}

// View renders the editors view
func (m EditorsModel) View() string {
	// Show plan editor selection modal if active
	if m.planEditorMode != PlanEditorModeNone {
		return m.renderPlanEditorSelection()
	}

	var b strings.Builder

	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	nameStyle := lipgloss.NewStyle().Foreground(theme.Primary)
	descStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	installedStyle := lipgloss.NewStyle().Foreground(theme.Success)
	notInstalledStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	cursorStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	messageStyle := lipgloss.NewStyle().Foreground(theme.Success)
	errorStyle := lipgloss.NewStyle().Foreground(theme.Error)
	activeTabStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	favoriteStyle := lipgloss.NewStyle().Foreground(theme.Warning)

	if !m.loaded {
		b.WriteString("\n  Scanning...\n")
		return b.String()
	}

	// Count stats
	filtered := m.filteredApps()
	installedCount := 0
	favoriteCount := 0
	for _, app := range filtered {
		if app.Installed {
			installedCount++
		}
		if app.Favorite {
			favoriteCount++
		}
	}

	// Header line: category tabs + stats on same line
	b.WriteString("\n  ")
	categories := []struct {
		key  string
		name string
		cat  AppCategory
	}{
		{"a", "All", -1},
		{"1", "Editors", CategoryEditor},
		{"2", "IDEs", CategoryIDE},
		{"3", "Tools", CategoryTool},
	}

	for i, cat := range categories {
		if i > 0 {
			b.WriteString(" ")
		}
		isActive := int(cat.cat) == int(m.category) || (cat.cat < 0 && m.category < 0)
		if isActive {
			b.WriteString(activeTabStyle.Render("[" + cat.name + "]"))
		} else {
			b.WriteString(keyStyle.Render(cat.key) + mutedStyle.Render(cat.name))
		}
	}

	// Stats inline
	stats := "  " + itoa(installedCount) + "/" + itoa(len(filtered))
	if favoriteCount > 0 {
		stats += " " + favoriteStyle.Render(itoa(favoriteCount)+"★")
	}
	b.WriteString(mutedStyle.Render(stats))
	b.WriteString("\n\n")

	// App list
	if len(filtered) == 0 {
		b.WriteString("  " + descStyle.Render("No applications in this category"))
		b.WriteString("\n")
	} else {
		visibleItems := m.visibleItems()
		endIdx := m.scroll + visibleItems
		if endIdx > len(filtered) {
			endIdx = len(filtered)
		}

		for i := m.scroll; i < endIdx; i++ {
			app := filtered[i]

			// Cursor
			cursor := "  "
			style := nameStyle
			if i == m.cursor {
				cursor = cursorStyle.Render("> ")
				style = selectedStyle
			}

			// Status indicator
			var status string
			if app.Installed {
				if app.Favorite {
					status = favoriteStyle.Render("★")
				} else {
					status = installedStyle.Render("●")
				}
			} else {
				status = notInstalledStyle.Render("○")
			}

			// Build line
			line := cursor + status + " " + style.Render(app.Name)

			// Add description inline
			desc := " " + app.Description
			maxDescLen := m.width - lipgloss.Width(line) - 6
			if maxDescLen > 10 && len(desc) > maxDescLen {
				desc = desc[:maxDescLen-3] + "..."
			}
			if maxDescLen > 10 {
				line += descStyle.Render(desc)
			}

			b.WriteString("  " + line + "\n")
		}

		// Scroll indicator (inline with list)
		if len(filtered) > visibleItems {
			b.WriteString("  " + mutedStyle.Render("  "+itoa(m.scroll+1)+"-"+itoa(endIdx)+" of "+itoa(len(filtered))) + "\n")
		}
	}

	// Fixed footer area
	b.WriteString("\n")

	// Status message (fixed position, single line)
	if m.message != "" {
		if strings.HasPrefix(m.message, "Failed") || strings.HasPrefix(m.message, "Could not") {
			b.WriteString("  " + errorStyle.Render("✗ "+m.message) + "\n")
		} else {
			b.WriteString("  " + messageStyle.Render("✓ "+m.message) + "\n")
		}
	} else {
		// Empty line to maintain consistent height
		b.WriteString("\n")
	}

	// Plan Editor section
	b.WriteString("\n")
	headerStyle := lipgloss.NewStyle().Foreground(theme.Muted).Bold(true)
	b.WriteString("  " + headerStyle.Render("Plan Editor"))
	b.WriteString("\n")

	// Show current editor
	currentName := "auto-detect"
	if m.currentPlanEditor != "" && m.currentPlanEditor != "auto" {
		currentName = m.currentPlanEditor
		if m.currentPlanEditorType != "" {
			currentName += " (" + m.currentPlanEditorType + ")"
		}
	}
	b.WriteString("  " + mutedStyle.Render("Current: ") + nameStyle.Render(currentName))
	b.WriteString("\n")
	b.WriteString("  " + keyStyle.Render("t") + mutedStyle.Render(" terminal  ") + keyStyle.Render("g") + mutedStyle.Render(" gui  ") + keyStyle.Render("0") + mutedStyle.Render(" auto"))
	b.WriteString("\n")

	// Single footer with all hints
	b.WriteString("\n")
	var hints []string
	hints = append(hints, keyStyle.Render("↑↓")+"nav")

	if m.cursor < len(filtered) {
		app := filtered[m.cursor]
		if app.Installed {
			hints = append(hints, keyStyle.Render("Enter")+"test")
			if app.Favorite {
				hints = append(hints, keyStyle.Render("Space")+"−★")
			} else {
				hints = append(hints, keyStyle.Render("Space")+"+★")
			}
		} else if app.URL != "" {
			hints = append(hints, keyStyle.Render("Enter")+"get")
		}
		if app.URL != "" {
			hints = append(hints, keyStyle.Render("b")+"web")
		}
	}

	b.WriteString("  " + mutedStyle.Render(strings.Join(hints, " ")))

	return b.String()
}

// renderPlanEditorSelection renders the plan editor selection modal
func (m EditorsModel) renderPlanEditorSelection() string {
	var b strings.Builder

	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	headerStyle := lipgloss.NewStyle().Foreground(theme.Primary).Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(theme.Primary)
	checkStyle := lipgloss.NewStyle().Foreground(theme.Success)

	var title string
	var editors []editor.Editor
	if m.planEditorMode == PlanEditorModeTerminal {
		title = "Select Terminal Editor"
		editors = editor.GetAvailableTerminal()
	} else {
		title = "Select GUI Editor"
		editors = editor.GetAvailableGUI()
	}

	b.WriteString("\n")
	b.WriteString("  " + headerStyle.Render(title))
	b.WriteString("\n\n")

	if len(editors) == 0 {
		b.WriteString("  " + mutedStyle.Render("No editors available"))
		b.WriteString("\n")
	} else {
		for i, ed := range editors {
			cursor := "  "
			style := normalStyle
			if i == m.planEditorCursor {
				cursor = keyStyle.Render("●") + " "
				style = selectedStyle
			}

			// Show check if this is the current editor
			suffix := ""
			if ed.Command == m.currentPlanEditor {
				suffix = " " + checkStyle.Render("✓")
			}

			b.WriteString("  " + cursor + style.Render(ed.Name) + suffix)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString("  " + keyStyle.Render("↑↓") + mutedStyle.Render(" navigate  ") + keyStyle.Render("Enter") + mutedStyle.Render(" select  ") + keyStyle.Render("Esc") + mutedStyle.Render(" cancel"))

	return b.String()
}
