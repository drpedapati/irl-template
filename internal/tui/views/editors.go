package views

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// AppCategory groups applications by type
type AppCategory int

const (
	CategoryEditor AppCategory = iota
	CategoryIDE
	CategoryUtility
)

// AppInfo represents an editor or utility application
type AppInfo struct {
	Name        string
	Description string
	Cmd         string      // Command line name
	AppName     string      // macOS .app name (if different)
	Key         string      // Hotkey
	Category    AppCategory
	Installed   bool
	URL         string // Website for installation
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

const editorsVisibleItems = 12

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
		// Code Editors
		{Name: "Cursor", Description: "AI-first code editor", Cmd: "cursor", AppName: "Cursor", Key: "u", Category: CategoryEditor, URL: "https://cursor.sh"},
		{Name: "VS Code", Description: "Popular extensible editor", Cmd: "code", Key: "v", Category: CategoryEditor, URL: "https://code.visualstudio.com"},
		{Name: "Zed", Description: "High-performance editor", Cmd: "zed", AppName: "Zed", Key: "z", Category: CategoryEditor, URL: "https://zed.dev"},
		{Name: "Sublime Text", Description: "Fast, lightweight editor", Cmd: "subl", AppName: "Sublime Text", Key: "s", Category: CategoryEditor, URL: "https://sublimetext.com"},
		{Name: "Neovim", Description: "Modern Vim-based editor", Cmd: "nvim", Key: "n", Category: CategoryEditor, URL: "https://neovim.io"},
		{Name: "Vim", Description: "Classic modal editor", Cmd: "vim", Key: "i", Category: CategoryEditor, URL: "https://vim.org"},
		{Name: "Emacs", Description: "Extensible text editor", Cmd: "emacs", Key: "e", Category: CategoryEditor, URL: "https://gnu.org/software/emacs"},
		{Name: "Helix", Description: "Post-modern modal editor", Cmd: "hx", Key: "h", Category: CategoryEditor, URL: "https://helix-editor.com"},

		// IDEs
		{Name: "Positron", Description: "Data science IDE", Cmd: "positron", AppName: "Positron", Key: "p", Category: CategoryIDE, URL: "https://github.com/posit-dev/positron"},
		{Name: "RStudio", Description: "R programming IDE", Cmd: "rstudio", AppName: "RStudio", Key: "r", Category: CategoryIDE, URL: "https://posit.co/products/open-source/rstudio"},
		{Name: "PyCharm", Description: "Python IDE", Cmd: "pycharm", AppName: "PyCharm", Key: "y", Category: CategoryIDE, URL: "https://jetbrains.com/pycharm"},
		{Name: "IntelliJ IDEA", Description: "Java/Kotlin IDE", Cmd: "idea", AppName: "IntelliJ IDEA", Key: "j", Category: CategoryIDE, URL: "https://jetbrains.com/idea"},
		{Name: "GoLand", Description: "Go IDE", Cmd: "goland", AppName: "GoLand", Key: "g", Category: CategoryIDE, URL: "https://jetbrains.com/go"},

		// Utilities
		{Name: "Terminal", Description: "System terminal", Cmd: "terminal", Key: "t", Category: CategoryUtility, URL: ""},
		{Name: "Finder", Description: "File manager (macOS)", Cmd: "open", Key: "f", Category: CategoryUtility, URL: ""},
		{Name: "iTerm2", Description: "Advanced terminal (macOS)", Cmd: "iterm", AppName: "iTerm", Key: "2", Category: CategoryUtility, URL: "https://iterm2.com"},
		{Name: "Warp", Description: "AI-powered terminal", Cmd: "warp", AppName: "Warp", Key: "w", Category: CategoryUtility, URL: "https://warp.dev"},
	}

	// Check installation status
	for i := range apps {
		apps[i].Installed = isAppInstalled(apps[i])
	}

	return apps
}

func isAppInstalled(app AppInfo) bool {
	// Special cases
	if app.Cmd == "terminal" || app.Cmd == "open" {
		return true // Always available
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
		key := msg.String()

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
			if m.cursor < len(filtered)-1 {
				m.cursor++
				if m.cursor >= m.scroll+editorsVisibleItems {
					m.scroll = m.cursor - editorsVisibleItems + 1
				}
			}
			return m, nil
		case "a":
			m.category = -1 // All
			m.cursor = 0
			m.scroll = 0
			return m, nil
		case "1":
			m.category = CategoryEditor
			m.cursor = 0
			m.scroll = 0
			return m, nil
		case "2":
			m.category = CategoryIDE
			m.cursor = 0
			m.scroll = 0
			return m, nil
		case "3":
			m.category = CategoryUtility
			m.cursor = 0
			m.scroll = 0
			return m, nil
		case "enter", "right":
			// Open URL or launch app
			if m.cursor < len(filtered) {
				app := filtered[m.cursor]
				if app.Installed {
					m.message = "Launched " + app.Name
					// Could launch the app here if desired
				} else if app.URL != "" {
					// Open URL in browser
					openURL(app.URL)
					m.message = "Opening " + app.Name + " website..."
				}
			}
			return m, nil
		}
	}

	return m, nil
}

func openURL(url string) {
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
		cmd.Start()
	}
}

// View renders the editors view
func (m EditorsModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().Foreground(theme.Muted).Bold(true)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	nameStyle := lipgloss.NewStyle().Foreground(theme.Primary).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	installedStyle := lipgloss.NewStyle().Foreground(theme.Success)
	notInstalledStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	cursorStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	messageStyle := lipgloss.NewStyle().Foreground(theme.Success)
	filterStyle := lipgloss.NewStyle().Foreground(theme.Accent)

	if !m.loaded {
		b.WriteString("\n  Detecting applications...\n")
		return b.String()
	}

	b.WriteString("\n")

	// Message if any
	if m.message != "" {
		b.WriteString("  " + messageStyle.Render("✓ "+m.message))
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
		{"3", "Utilities", CategoryUtility},
	}

	for i, cat := range categories {
		if i > 0 {
			b.WriteString("  ")
		}
		label := keyStyle.Render(cat.key) + " " + cat.name
		if int(cat.cat) == int(m.category) || (cat.cat < 0 && m.category < 0) {
			label = filterStyle.Render("[" + cat.key + " " + cat.name + "]")
		}
		b.WriteString(label)
	}
	b.WriteString("\n\n")

	// Count installed
	filtered := m.filteredApps()
	installedCount := 0
	for _, app := range filtered {
		if app.Installed {
			installedCount++
		}
	}
	b.WriteString("  " + headerStyle.Render(itoa(installedCount)+"/"+itoa(len(filtered))+" installed"))
	b.WriteString("\n\n")

	// App list
	if len(filtered) == 0 {
		b.WriteString("  " + descStyle.Render("No applications in this category"))
		b.WriteString("\n")
	} else {
		endIdx := m.scroll + editorsVisibleItems
		if endIdx > len(filtered) {
			endIdx = len(filtered)
		}

		for i := m.scroll; i < endIdx; i++ {
			app := filtered[i]

			// Cursor
			cursor := "  "
			style := nameStyle
			if i == m.cursor {
				cursor = cursorStyle.Render("● ")
				style = selectedStyle
			}

			// Status indicator
			status := notInstalledStyle.Render("○")
			if app.Installed {
				status = installedStyle.Render("●")
			}

			// Key hint
			keyHint := keyStyle.Render(app.Key)

			b.WriteString("  " + cursor + status + " " + keyHint + " " + style.Render(app.Name))
			b.WriteString("\n")
			b.WriteString("       " + descStyle.Render(app.Description))
			b.WriteString("\n")
		}

		// Scroll indicator
		if len(filtered) > editorsVisibleItems {
			b.WriteString("\n")
			b.WriteString("  " + descStyle.Render(itoa(m.scroll+1)+"-"+itoa(endIdx)+" of "+itoa(len(filtered))))
			b.WriteString("\n")
		}
	}

	// Footer hint
	b.WriteString("\n")
	hint := keyStyle.Render("↑↓") + descStyle.Render(" navigate  ")
	hint += keyStyle.Render("Enter") + descStyle.Render(" open/install")
	b.WriteString("  " + hint)

	return b.String()
}
