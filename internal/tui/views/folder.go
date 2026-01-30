package views

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// FolderModel handles folder selection
type FolderModel struct {
	width      int
	height     int
	currentDir string
	folders    []string
	cursor     int // 0 = "Use this folder", 1+ = subfolders
	scroll     int
	saved      bool
	wantsBack  bool // True when user presses back while on "Use this folder"
}

const folderVisibleItems = 8

// NewFolderModel creates a new folder selection view
func NewFolderModel() FolderModel {
	// Start at current default or home/Documents
	startDir := config.GetDefaultDirectory()
	if startDir == "" {
		home, _ := os.UserHomeDir()
		startDir = filepath.Join(home, "Documents")
	}

	// Make sure directory exists, fall back to home
	if info, err := os.Stat(startDir); err != nil || !info.IsDir() {
		home, _ := os.UserHomeDir()
		startDir = home
	}

	m := FolderModel{
		currentDir: startDir,
		cursor:     0, // Start on "Use this folder"
	}
	m.loadFolders()
	return m
}

// SetSize sets the view dimensions
func (m *FolderModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *FolderModel) loadFolders() {
	m.folders = []string{}
	m.scroll = 0

	entries, err := os.ReadDir(m.currentDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			m.folders = append(m.folders, e.Name())
		}
	}
	// Case-insensitive alphabetical sort
	sort.Slice(m.folders, func(i, j int) bool {
		return strings.ToLower(m.folders[i]) < strings.ToLower(m.folders[j])
	})
}

// totalItems returns the total number of selectable items (1 for "use this" + subfolders)
func (m FolderModel) totalItems() int {
	return 1 + len(m.folders)
}

// Update handles messages
func (m FolderModel) Update(msg tea.Msg) (FolderModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < m.totalItems()-1 {
				m.cursor++
			}
		case "right", "l", "enter":
			if m.cursor == 0 {
				// Select current directory
				config.SetDefaultDirectory(m.currentDir)
				m.saved = true
			} else {
				// Enter the selected subfolder
				folderIdx := m.cursor - 1
				if folderIdx < len(m.folders) {
					m.currentDir = filepath.Join(m.currentDir, m.folders[folderIdx])
					m.loadFolders()
					// Stay in folder list area (cursor 1 = first folder)
					// but clamp to valid range if new folder has fewer items
					if len(m.folders) > 0 {
						m.cursor = 1 // First folder in new list
					} else {
						m.cursor = 0 // Only "Use this folder" available
					}
				}
			}
		case "left", "h":
			if m.cursor == 0 {
				// On "Use this folder" - signal to go back to menu
				m.wantsBack = true
			} else {
				// In folder list area - go up one directory level
				parent := filepath.Dir(m.currentDir)
				if parent != m.currentDir {
					m.currentDir = parent
					m.loadFolders()
					// Stay in folder list area
					if len(m.folders) > 0 {
						m.cursor = 1 // First folder in parent
					} else {
						m.cursor = 0 // Only "Use this folder" available
					}
				}
			}
		}
	}
	return m, nil
}

// IsSaved returns true if folder was saved
func (m FolderModel) IsSaved() bool {
	return m.saved
}

// WantsBack returns true if user pressed back while on "Use this folder"
func (m FolderModel) WantsBack() bool {
	return m.wantsBack
}

// View renders the folder selection view
func (m FolderModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).MarginLeft(2)
	pathStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	successStyle := lipgloss.NewStyle().Foreground(theme.Success).MarginLeft(2)

	cursorOn := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("●")
	cursorOff := " "

	// Success state
	if m.saved {
		b.WriteString("\n")
		b.WriteString(successStyle.Render("✓ Default folder set to:"))
		b.WriteString("\n\n")
		b.WriteString("  " + pathStyle.Render(m.currentDir))
		b.WriteString("\n\n")
		b.WriteString(hintStyle.Render("New projects will be created here."))
		b.WriteString("\n")
		return b.String()
	}

	// Title
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("Choose where to save new projects"))
	b.WriteString("\n\n")

	// Current location breadcrumb
	b.WriteString(hintStyle.Render("Current location:"))
	b.WriteString("\n")
	b.WriteString("  " + pathStyle.Render(m.currentDir))
	b.WriteString("\n\n")

	// Option 0: Use this folder
	cursor := cursorOff
	style := normalStyle
	if m.cursor == 0 {
		cursor = cursorOn
		style = selectedStyle
	}
	b.WriteString("  " + cursor + " " + style.Render("✓ Use this folder"))
	b.WriteString("\n\n")

	// Subfolders section
	if len(m.folders) > 0 {
		b.WriteString(hintStyle.Render("Or navigate to a subfolder:"))
		b.WriteString("\n")

		// Show visible subfolders with scrolling
		startIdx := 0
		if m.cursor > folderVisibleItems {
			startIdx = m.cursor - folderVisibleItems
		}
		endIdx := startIdx + folderVisibleItems
		if endIdx > len(m.folders) {
			endIdx = len(m.folders)
		}

		for i := startIdx; i < endIdx; i++ {
			folder := m.folders[i]
			cursor = cursorOff
			style = normalStyle

			if m.cursor == i+1 { // +1 because cursor 0 is "Use this folder"
				cursor = cursorOn
				style = selectedStyle
			}

			b.WriteString("  " + cursor + " " + style.Render(folder+"/"))
			b.WriteString("\n")
		}

		// Scroll indicator
		if len(m.folders) > folderVisibleItems {
			b.WriteString(hintStyle.Render("    (" + itoa(len(m.folders)) + " folders)"))
			b.WriteString("\n")
		}
	} else {
		b.WriteString(hintStyle.Render("No subfolders in this location"))
		b.WriteString("\n")
	}

	return b.String()
}

// Simple int to string without fmt
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
