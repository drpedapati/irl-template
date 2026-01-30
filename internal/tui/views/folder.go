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
	width        int
	height       int
	currentDir   string
	folders      []string
	cursor       int
	scroll       int
	saved        bool
}

const folderVisibleItems = 12

// NewFolderModel creates a new folder selection view
func NewFolderModel() FolderModel {
	// Start at current default or home
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
	m.cursor = 0
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
	sort.Strings(m.folders)
}

// Update handles messages
func (m FolderModel) Update(msg tea.Msg) (FolderModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.folders)-1 {
				m.cursor++
				if m.cursor >= m.scroll+folderVisibleItems {
					m.scroll = m.cursor - folderVisibleItems + 1
				}
			}
		case "right", "l":
			// Enter selected folder
			if len(m.folders) > 0 && m.cursor < len(m.folders) {
				m.currentDir = filepath.Join(m.currentDir, m.folders[m.cursor])
				m.loadFolders()
			}
		case "left", "h":
			// Go up one level
			parent := filepath.Dir(m.currentDir)
			if parent != m.currentDir {
				m.currentDir = parent
				m.loadFolders()
			}
		case "enter":
			// Select current directory as default
			config.SetDefaultDirectory(m.currentDir)
			m.saved = true
		}
	}
	return m, nil
}

// IsSaved returns true if folder was saved
func (m FolderModel) IsSaved() bool {
	return m.saved
}

// View renders the folder selection view
func (m FolderModel) View() string {
	var b strings.Builder

	pathStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).MarginLeft(2)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)

	// Current path
	b.WriteString(pathStyle.Render(m.currentDir))
	b.WriteString("\n\n")

	if m.saved {
		successStyle := lipgloss.NewStyle().Foreground(theme.Success).MarginLeft(2)
		b.WriteString(successStyle.Render("✓ Default folder saved"))
		b.WriteString("\n")
		return b.String()
	}

	cursorOn := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("●")
	cursorOff := "  "
	itemStyle := lipgloss.NewStyle()
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	if len(m.folders) == 0 {
		b.WriteString(hintStyle.Render("  (no subfolders)"))
		b.WriteString("\n")
	} else {
		// Show visible items with scrolling
		end := m.scroll + folderVisibleItems
		if end > len(m.folders) {
			end = len(m.folders)
		}

		for i := m.scroll; i < end; i++ {
			folder := m.folders[i]
			cursor := cursorOff
			style := itemStyle

			if i == m.cursor {
				cursor = cursorOn
				style = selectedStyle
			}

			b.WriteString("  " + cursor + " " + style.Render(folder+"/"))
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if len(m.folders) > folderVisibleItems {
			shown := lipgloss.NewStyle().Foreground(theme.Muted).Render(
				strings.Repeat(" ", 4) + "(" + string(rune('0'+m.scroll/10)) + string(rune('0'+m.scroll%10)) + "-" +
				string(rune('0'+(end-1)/10)) + string(rune('0'+(end-1)%10)) + " of " +
				string(rune('0'+len(m.folders)/10)) + string(rune('0'+len(m.folders)%10)) + ")",
			)
			b.WriteString(shown)
			b.WriteString("\n")
		}
	}

	return b.String()
}
