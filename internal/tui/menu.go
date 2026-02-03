package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// MenuItem represents a menu item
type MenuItem struct {
	Title        string
	Desc         string
	Key          string // Quick access key (n, t, d, c)
	ViewType     ViewType
	SeparatorAfter bool // Show a subtle separator line after this item
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
	ViewHelp        // Help/tutorial view
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
			{Title: "Templates", Desc: "Browse and manage templates", Key: "t", ViewType: ViewTemplates, SeparatorAfter: true},
			{Title: "Docs", Desc: "Open documentation in browser", Key: "o", ViewType: ViewDocs},
			{Title: "Editors", Desc: "See installed editors and tools", Key: "e", ViewType: ViewEditors},
			{Title: "Doctor", Desc: "Check environment setup", Key: "d", ViewType: ViewDoctor},
			{Title: "Help", Desc: "Learn how IRL works", Key: "?", ViewType: ViewHelp},
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
	// Format: "● [x] Title" - find the longest
	maxLabelWidth := 0
	for _, item := range m.items {
		// cursor(2) + [key](3) + space + title
		labelWidth := 2 + 3 + 1 + len(item.Title)
		if labelWidth > maxLabelWidth {
			maxLabelWidth = labelWidth
		}
	}

	// Find max description width
	maxDescWidth := 0
	for _, item := range m.items {
		if len(item.Desc) > maxDescWidth {
			maxDescWidth = len(item.Desc)
		}
	}

	// Column separator gap
	gap := 4

	// Total table width
	tableWidth := maxLabelWidth + gap + maxDescWidth

	// Calculate left margin to center the table
	leftMargin := (m.width - tableWidth) / 2
	if leftMargin < 0 {
		leftMargin = 0
	}
	marginStr := strings.Repeat(" ", leftMargin)

	// Top padding
	b.WriteString("\n")

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
		leftCol := cursor + key + " " + title

		// Pad to align descriptions
		leftWidth := lipgloss.Width(leftCol)
		padding := maxLabelWidth + gap - leftWidth
		if padding < 2 {
			padding = 2
		}

		// Right column: description
		desc := descStyle.Render(item.Desc)

		b.WriteString(marginStr + leftCol + strings.Repeat(" ", padding) + desc)
		b.WriteString("\n")

		// Add spacing between items (except after last)
		if i < len(m.items)-1 {
			b.WriteString("\n")
		}

		// Add subtle separator if flagged
		if item.SeparatorAfter && i < len(m.items)-1 {
			separatorStyle := lipgloss.NewStyle().Foreground(theme.Muted)
			separatorWidth := tableWidth - 4 // Slightly shorter than table
			separatorPadding := leftMargin + 2
			b.WriteString(strings.Repeat(" ", separatorPadding) + separatorStyle.Render(strings.Repeat("─", separatorWidth)))
			b.WriteString("\n\n")
		}
	}

	// Bottom padding
	b.WriteString("\n\n")

	return b.String()
}
