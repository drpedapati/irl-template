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
	ViewTemplates
	ViewDoctor
	ViewConfig
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
			{Title: "Templates", Desc: "Browse available templates", Key: "t", ViewType: ViewTemplates},
			{Title: "Doctor", Desc: "Check environment setup", Key: "d", ViewType: ViewDoctor},
			{Title: "Config", Desc: "View and edit settings", Key: "c", ViewType: ViewConfig},
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

	for i, item := range m.items {
		cursor := "  "
		titleStyle := normalStyle

		if i == m.cursor {
			cursor = cursorStyle.Render("● ")
			titleStyle = selectedStyle
		}

		// Format: cursor [key] Title - description
		key := keyStyle.Render("[" + item.Key + "]")
		title := titleStyle.Render(item.Title)
		desc := descStyle.Render(item.Desc)

		b.WriteString("  " + cursor + key + " " + title)
		b.WriteString("  " + desc)
		b.WriteString("\n")

		// Add spacing between items for breathing room
		if i < len(m.items)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
