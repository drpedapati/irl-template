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
// If favorites are set, only returns favorite tools. Otherwise returns all installed.
func GetInstalledTools() []AppInfo {
	apps := getAllApps()
	favorites := config.GetFavoriteEditors()
	hasFavorites := len(favorites) > 0

	var installed []AppInfo
	for _, app := range apps {
		if app.Installed && app.Key != "" && app.Category == CategoryTool {
			// If favorites are set, only include favorites
			if hasFavorites && !app.Favorite {
				continue
			}
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
}

// Calculate visible items based on available height
// Each item takes 1 line, with some header/footer space
func (m EditorsModel) visibleItems() int {
	// Reserve lines for: header(3) + filter tabs(2) + count(2) + footer(3) + margins(2) = 12
	available := m.height - 12
	if available < 5 {
		available = 5
	}
	if available > 15 {
		available = 15
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

// getAllApps returns all known applications with their install status
func getAllApps() []AppInfo {
	apps := []AppInfo{
		// Code Editors (ordered by popularity/preference)
		{Name: "Cursor", Description: "AI-powered code editor", Cmd: "cursor", AppName: "Cursor", Key: "c", Category: CategoryEditor, URL: "https://cursor.sh"},
		{Name: "VS Code", Description: "Microsoft's popular editor", Cmd: "code", Key: "v", Category: CategoryEditor, URL: "https://code.visualstudio.com"},
		{Name: "Zed", Description: "Fast, collaborative editor", Cmd: "zed", AppName: "Zed", Key: "z", Category: CategoryEditor, URL: "https://zed.dev"},
		{Name: "Sublime Text", Description: "Lightweight and fast", Cmd: "subl", AppName: "Sublime Text", Key: "s", Category: CategoryEditor, URL: "https://sublimetext.com"},
		{Name: "Neovim", Description: "Modern terminal editor", Cmd: "nvim", Key: "n", Category: CategoryEditor, URL: "https://neovim.io"},
		{Name: "Helix", Description: "Modal terminal editor", Cmd: "hx", Key: "x", Category: CategoryEditor, URL: "https://helix-editor.com"},
		{Name: "Fresh", Description: "Intuitive terminal editor", Cmd: "fresh", Key: "h", Category: CategoryEditor, URL: "https://getfresh.dev"},

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
		return m, nil

	case tea.KeyMsg:
		filtered := m.filteredApps()
		visibleItems := m.visibleItems()
		key := msg.String()

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
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(theme.Primary).Bold(true)
	headerStyle := lipgloss.NewStyle().Foreground(theme.Muted)
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

	if !m.loaded {
		b.WriteString("\n  Scanning for installed applications...\n")
		return b.String()
	}

	// Title and description
	b.WriteString("\n")
	b.WriteString("  " + titleStyle.Render("Editors & Tools"))
	b.WriteString("\n")
	b.WriteString("  " + headerStyle.Render("See what's installed and test launching applications"))
	b.WriteString("\n\n")

	// Message if any
	if m.message != "" {
		if strings.HasPrefix(m.message, "Failed") || strings.HasPrefix(m.message, "Could not") {
			b.WriteString("  " + errorStyle.Render("✗ "+m.message))
		} else {
			b.WriteString("  " + messageStyle.Render("✓ "+m.message))
		}
		b.WriteString("\n\n")
	}

	// Category filter tabs
	b.WriteString("  ")
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
			b.WriteString("  ")
		}
		isActive := int(cat.cat) == int(m.category) || (cat.cat < 0 && m.category < 0)
		if isActive {
			b.WriteString(activeTabStyle.Render("[" + cat.key + "] " + cat.name))
		} else {
			b.WriteString(keyStyle.Render(cat.key) + " " + descStyle.Render(cat.name))
		}
	}
	b.WriteString("\n\n")

	// Count installed and favorites
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
	stats := itoa(installedCount) + " of " + itoa(len(filtered)) + " installed"
	if favoriteCount > 0 {
		stats += ", " + itoa(favoriteCount) + " favorites"
	}
	b.WriteString("  " + headerStyle.Render(stats))
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

		favoriteStyle := lipgloss.NewStyle().Foreground(theme.Warning)

		for i := m.scroll; i < endIdx; i++ {
			app := filtered[i]

			// Cursor
			cursor := "  "
			style := nameStyle
			if i == m.cursor {
				cursor = cursorStyle.Render("> ")
				style = selectedStyle
			}

			// Status indicator (installed + favorite)
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

			// Build line: cursor + status + name + description
			line := cursor + status + " " + style.Render(app.Name)

			// Add description inline, truncated if needed
			desc := " - " + app.Description
			maxDescLen := m.width - lipgloss.Width(line) - 6
			if maxDescLen > 10 && len(desc) > maxDescLen {
				desc = desc[:maxDescLen-3] + "..."
			}
			if maxDescLen > 10 {
				line += descStyle.Render(desc)
			}

			b.WriteString("  " + line + "\n")
		}

		// Scroll indicator
		if len(filtered) > visibleItems {
			b.WriteString("\n")
			b.WriteString("  " + descStyle.Render("↑↓ "+itoa(m.scroll+1)+"-"+itoa(endIdx)+" of "+itoa(len(filtered))))
			b.WriteString("\n")
		}
	}

	// Footer hints
	b.WriteString("\n")
	var hints []string
	hints = append(hints, keyStyle.Render("↑↓")+" navigate")

	if m.cursor < len(filtered) {
		app := filtered[m.cursor]
		if app.Installed {
			hints = append(hints, keyStyle.Render("Enter")+" test")
			if app.Favorite {
				hints = append(hints, keyStyle.Render("Space")+" unfavorite")
			} else {
				hints = append(hints, keyStyle.Render("Space")+" favorite")
			}
		} else if app.URL != "" {
			hints = append(hints, keyStyle.Render("Enter")+" get it")
		}
		if app.URL != "" {
			hints = append(hints, keyStyle.Render("b")+" website")
		}
	}

	b.WriteString("  " + descStyle.Render(strings.Join(hints, "  ")))

	// Explanation of favorites
	if favoriteCount > 0 {
		b.WriteString("\n")
		b.WriteString("  " + descStyle.Render("★ = favorite (only favorites shown in project actions)"))
	}

	return b.String()
}
