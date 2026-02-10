package views

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// FolderModel handles folder selection
type FolderModel struct {
	width       int
	height      int
	currentDir  string
	folders     []browseEntry // All folders
	filtered    []browseEntry // Filtered folders
	sortBy      string        // "name-asc" or "date-desc"
	cursor      int           // 0 = "Use this folder", 1+ = filtered subfolders
	scroll      int
	saved       bool
	wantsBack   bool // True when user presses back while on "Use this folder"
	filterInput textinput.Model
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

	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Width = 30
	ti.Focus()

	m := FolderModel{
		currentDir:  startDir,
		cursor:      0, // Start on "Use this folder"
		filterInput: ti,
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
	m.folders = []browseEntry{}
	m.scroll = 0

	entries, err := os.ReadDir(m.currentDir)
	if err != nil {
		m.filtered = m.folders
		return
	}

	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			modTime := time.Time{}
			if info, err := e.Info(); err == nil {
				modTime = info.ModTime()
			}
			m.folders = append(m.folders, browseEntry{Name: e.Name(), ModTime: modTime})
		}
	}

	m.applyFilter()
}

func (m *FolderModel) applyFilter() {
	query := strings.ToLower(m.filterInput.Value())
	if query == "" {
		m.filtered = m.folders
	} else {
		m.filtered = []browseEntry{}
		for _, f := range m.folders {
			if strings.Contains(strings.ToLower(f.Name), query) {
				m.filtered = append(m.filtered, f)
			}
		}
	}
	m.applyFolderSort()
	// Reset cursor to "Use this folder" when filter changes
	m.cursor = 0
	m.scroll = 0
}

func (m *FolderModel) applyFolderSort() {
	switch m.sortBy {
	case "date-desc":
		sort.Slice(m.filtered, func(i, j int) bool {
			return m.filtered[i].ModTime.After(m.filtered[j].ModTime)
		})
	default: // "name-asc"
		sort.Slice(m.filtered, func(i, j int) bool {
			return strings.ToLower(m.filtered[i].Name) < strings.ToLower(m.filtered[j].Name)
		})
	}
}

// totalItems returns the total number of selectable items (1 for "use this" + filtered subfolders or 1 for "go up")
func (m FolderModel) totalItems() int {
	if len(m.filtered) == 0 {
		// When no subfolders, show "Go up one level" option (unless at root)
		parent := filepath.Dir(m.currentDir)
		if parent != m.currentDir {
			return 2 // "Use this folder" + "Go up one level"
		}
		return 1 // At root with no subfolders
	}
	return 1 + len(m.filtered)
}

// HasFilterText returns true if there's text in the filter input
func (m FolderModel) HasFilterText() bool {
	return m.filterInput.Value() != ""
}

// ClearFilter clears the filter input
func (m *FolderModel) ClearFilter() {
	m.filterInput.SetValue("")
	m.applyFilter()
}

// Update handles messages
func (m FolderModel) Update(msg tea.Msg) (FolderModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				// Adjust scroll if needed
				if m.cursor > 0 && m.cursor-1 < m.scroll {
					m.scroll = m.cursor - 1
				}
			}
			return m, nil
		case "down":
			if m.cursor < m.totalItems()-1 {
				m.cursor++
				// Adjust scroll if needed
				if m.cursor > 0 && m.cursor-1 >= m.scroll+folderVisibleItems {
					m.scroll = m.cursor - folderVisibleItems
				}
			}
			return m, nil
		case "right", "enter":
			if m.cursor == 0 {
				// Select current directory
				config.SetDefaultDirectory(m.currentDir)
				m.saved = true
			} else if len(m.filtered) == 0 {
				// "Go up one level" option (cursor == 1 with no subfolders)
				parent := filepath.Dir(m.currentDir)
				if parent != m.currentDir {
					m.currentDir = parent
					m.filterInput.SetValue("")
					m.loadFolders()
					if len(m.filtered) > 0 {
						m.cursor = 1
					} else {
						m.cursor = 0
					}
				}
			} else {
				// Enter the selected subfolder
				folderIdx := m.cursor - 1
				if folderIdx < len(m.filtered) {
					m.currentDir = filepath.Join(m.currentDir, m.filtered[folderIdx].Name)
					m.filterInput.SetValue("") // Clear filter when entering folder
					m.loadFolders()
					// Stay in folder list area
					if len(m.filtered) > 0 {
						m.cursor = 1
					} else {
						m.cursor = 0
					}
				}
			}
			return m, nil
		case "left":
			if m.cursor == 0 {
				// On "Use this folder" - signal to go back to menu
				m.wantsBack = true
			} else {
				// In folder list area or on "Go up" - go up one directory level
				parent := filepath.Dir(m.currentDir)
				if parent != m.currentDir {
					m.currentDir = parent
					m.filterInput.SetValue("") // Clear filter when going up
					m.loadFolders()
					if len(m.filtered) > 0 {
						m.cursor = 1
					} else {
						m.cursor = 0
					}
				}
			}
			return m, nil
		case "s":
			// Toggle sort
			if m.sortBy == "date-desc" {
				m.sortBy = "name-asc"
			} else {
				m.sortBy = "date-desc"
			}
			m.applyFolderSort()
			if m.cursor > 0 {
				m.cursor = 1
			}
			m.scroll = 0
			return m, nil
		case "esc":
			// Two-stage: clear filter first, then signal back
			if m.HasFilterText() {
				m.ClearFilter()
				return m, nil
			}
			m.wantsBack = true
			return m, nil
		}

		// Pass other keys to filter input
		var cmd tea.Cmd
		m.filterInput, cmd = m.filterInput.Update(msg)
		m.applyFilter()
		return m, cmd
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

	pathStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	successStyle := lipgloss.NewStyle().Foreground(theme.Success)

	cursorOn := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("●")
	cursorOff := " "

	// Success state
	if m.saved {
		b.WriteString("\n")
		b.WriteString("  " + successStyle.Render("✓ Default folder set to:"))
		b.WriteString("\n\n")
		b.WriteString("  " + pathStyle.Render(m.currentDir))
		b.WriteString("\n\n")
		b.WriteString("  " + hintStyle.Render("New projects will be created here."))
		b.WriteString("\n")
		return b.String()
	}

	// Centered title
	title := "Choose where to save new projects"
	if m.width > 0 {
		padding := (m.width - len(title)) / 2
		if padding > 0 {
			title = strings.Repeat(" ", padding) + title
		}
	}
	b.WriteString("\n")
	b.WriteString(hintStyle.Render(title))
	b.WriteString("\n\n")

	// Current path with folder icon
	b.WriteString("  📁 " + pathStyle.Render(m.currentDir))
	b.WriteString("\n\n")

	// Filter and "Use this folder" on same line
	filterView := m.filterInput.View()
	filterCount := ""
	if len(m.folders) > 0 {
		sortLabel := "A-Z"
		if m.sortBy == "date-desc" {
			sortLabel = "newest"
		}
		filterCount = hintStyle.Render(" (" + itoa(len(m.filtered)) + "/" + itoa(len(m.folders)) + ")  s:" + sortLabel)
	}

	// "Use this folder" button
	cursor := cursorOff
	style := normalStyle
	if m.cursor == 0 {
		cursor = cursorOn
		style = selectedStyle
	}
	useFolder := cursor + " " + style.Render("✓ Use this folder")

	b.WriteString("  " + filterView + filterCount + "       " + useFolder)
	b.WriteString("\n\n")

	// Subfolders section
	if len(m.filtered) > 0 {
		// Show visible subfolders with scrolling
		endIdx := m.scroll + folderVisibleItems
		if endIdx > len(m.filtered) {
			endIdx = len(m.filtered)
		}

		for i := m.scroll; i < endIdx; i++ {
			entry := m.filtered[i]
			cursor = cursorOff
			style = normalStyle

			if m.cursor == i+1 { // +1 because cursor 0 is "Use this folder"
				cursor = cursorOn
				style = selectedStyle
			}

			b.WriteString("  " + cursor + " " + style.Render(entry.Name+"/"))
			b.WriteString("\n")
		}

		// Scroll indicator
		if len(m.filtered) > folderVisibleItems {
			showing := m.scroll + 1
			showingEnd := endIdx
			b.WriteString(hintStyle.Render("    " + itoa(showing) + "-" + itoa(showingEnd) + " of " + itoa(len(m.filtered))))
			b.WriteString("\n")
		}
	} else if len(m.folders) > 0 {
		b.WriteString(hintStyle.Render("  No matches for \"" + m.filterInput.Value() + "\""))
		b.WriteString("\n")
	} else {
		// No subfolders - show "Go up" option if not at root
		parent := filepath.Dir(m.currentDir)
		if parent != m.currentDir {
			cursor = cursorOff
			style = normalStyle
			if m.cursor == 1 {
				cursor = cursorOn
				style = selectedStyle
			}
			b.WriteString("  " + cursor + " " + style.Render("← Go up one level"))
			b.WriteString("\n")
		} else {
			b.WriteString(hintStyle.Render("  No subfolders in this location"))
			b.WriteString("\n")
		}
	}

	return b.String()
}
