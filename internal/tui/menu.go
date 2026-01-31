package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// MenuItem represents a menu item
type MenuItem struct {
	Title    string
	Desc     string
	Key      string // Quick access key (n, t, d, c)
	ViewType ViewType
}

// ViewType identifies which view to show
type ViewType int

const (
	ViewMenu ViewType = iota
	ViewInit
	ViewProjects    // List IRL projects
	ViewFolder      // Set default folder
	ViewTemplates
	ViewEditors     // Editors & utilities manager
	ViewDoctor
	ViewConfig
	ViewPersonalize // Academic profile settings
	ViewDocs        // Opens browser to www.irloop.org
	ViewUpdate      // Updates templates from GitHub
)

// Menu represents the main menu
type Menu struct {
	items    []MenuItem
	cursor   int
	width    int
	selected ViewType
}

// NewMenu creates a new menu with default items
func NewMenu() Menu {
	return Menu{
		items: []MenuItem{
			{Title: "New project", Desc: "Create a new IRL project", Key: "n", ViewType: ViewInit},
			{Title: "Projects", Desc: "Browse existing IRL projects", Key: "p", ViewType: ViewProjects},
			{Title: "Folder", Desc: "Set default project folder", Key: "f", ViewType: ViewFolder},
			{Title: "Templates", Desc: "Browse available templates", Key: "t", ViewType: ViewTemplates},
			{Title: "Editors", Desc: "See installed editors and tools", Key: "e", ViewType: ViewEditors},
			{Title: "Doctor", Desc: "Check environment setup", Key: "d", ViewType: ViewDoctor},
			{Title: "Docs", Desc: "Open documentation in browser", Key: "o", ViewType: ViewDocs},
		},
		cursor:   0,
		selected: ViewMenu,
	}
}

// SetWidth sets the menu width
func (m *Menu) SetWidth(width int) {
	m.width = width
}

// Up moves cursor up
func (m *Menu) Up() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// Down moves cursor down
func (m *Menu) Down() {
	if m.cursor < len(m.items)-1 {
		m.cursor++
	}
}

// Select returns the selected item's view type
func (m *Menu) Select() ViewType {
	if m.cursor >= 0 && m.cursor < len(m.items) {
		return m.items[m.cursor].ViewType
	}
	return ViewMenu
}

// SelectByKey selects an item by its quick key
func (m *Menu) SelectByKey(key string) (ViewType, bool) {
	for i, item := range m.items {
		if item.Key == key {
			m.cursor = i
			return item.ViewType, true
		}
	}
	return ViewMenu, false
}

// Cursor returns current cursor position
func (m *Menu) Cursor() int {
	return m.cursor
}

// View renders the menu
func (m Menu) View() string {
	var b strings.Builder

	keyStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	selectedStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	normalStyle := lipgloss.NewStyle()

	descStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	cursorStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	// Calculate max label width for two-column alignment
	// Format: "  ● [x] Title" - find the longest
	maxLabelWidth := 0
	for _, item := range m.items {
		// cursor(2) + space + [key](3) + space + title
		labelWidth := 2 + 2 + 3 + 1 + len(item.Title)
		if labelWidth > maxLabelWidth {
			maxLabelWidth = labelWidth
		}
	}

	// Column separator position
	colSep := maxLabelWidth + 4

	// Top padding (uniform with bottom)
	b.WriteString("\n\n")

	for i, item := range m.items {
		cursor := "  "
		titleStyle := normalStyle

		if i == m.cursor {
			cursor = cursorStyle.Render("● ")
			titleStyle = selectedStyle
		}

		// Left column: cursor [key] Title
		key := keyStyle.Render("[" + item.Key + "]")
		title := titleStyle.Render(item.Title)
		leftCol := "  " + cursor + key + " " + title

		// Pad to align descriptions
		leftWidth := lipgloss.Width(leftCol)
		padding := colSep - leftWidth
		if padding < 2 {
			padding = 2
		}

		// Right column: description
		desc := descStyle.Render(item.Desc)

		b.WriteString(leftCol + strings.Repeat(" ", padding) + desc)
		b.WriteString("\n")

		// Add spacing between items (except after last)
		if i < len(m.items)-1 {
			b.WriteString("\n")
		}
	}

	// Bottom padding (uniform with top)
	b.WriteString("\n\n")

	return b.String()
}
