package views

import (
	"strings"

	"github.com/charmbracelet/glamour"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// TemplatesModel is the templates browser view
type TemplatesModel struct {
	templates  []TemplateItem
	cursor     int
	width      int
	height     int
	loaded     bool
	err        error
	previewing bool
	scroll     int
	renderer   *glamour.TermRenderer
}

// NewTemplatesModel creates a new templates view
func NewTemplatesModel() TemplatesModel {
	// Create glamour renderer with dark style
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(68),
	)
	return TemplatesModel{
		renderer: r,
	}
}

// SetSize sets the view dimensions
func (m *TemplatesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// LoadTemplates returns a command that loads templates
func (m *TemplatesModel) LoadTemplates() tea.Cmd {
	return func() tea.Msg {
		list, err := templates.ListTemplates()
		if err != nil {
			return TemplatesLoadedMsg{Err: err}
		}

		items := make([]TemplateItem, len(list))
		for i, t := range list {
			items[i] = TemplateItem{
				Name:        t.Name,
				Description: t.Description,
				Content:     t.Content,
			}
		}
		return TemplatesLoadedMsg{Templates: items}
	}
}

// RefreshTemplates forces a refresh from GitHub
func (m *TemplatesModel) RefreshTemplates() tea.Cmd {
	return func() tea.Msg {
		list, err := templates.FetchTemplates()
		if err != nil {
			return TemplatesLoadedMsg{Err: err}
		}

		items := make([]TemplateItem, len(list))
		for i, t := range list {
			items[i] = TemplateItem{
				Name:        t.Name,
				Description: t.Description,
				Content:     t.Content,
			}
		}
		return TemplatesLoadedMsg{Templates: items}
	}
}

// IsPreviewing returns true if in preview mode
func (m TemplatesModel) IsPreviewing() bool {
	return m.previewing
}

// Update handles messages
func (m TemplatesModel) Update(msg tea.Msg) (TemplatesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case TemplatesLoadedMsg:
		m.loaded = true
		m.err = msg.Err
		m.templates = msg.Templates
		return m, nil

	case tea.KeyMsg:
		if m.previewing {
			return m.updatePreview(msg)
		}
		return m.updateList(msg)
	}

	return m, nil
}

func (m TemplatesModel) updateList(msg tea.KeyMsg) (TemplatesModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.templates)-1 {
			m.cursor++
		}
	case "right", "enter", "l":
		// Enter preview mode
		if len(m.templates) > 0 {
			m.previewing = true
			m.scroll = 0
		}
	}
	return m, nil
}

func (m TemplatesModel) updatePreview(msg tea.KeyMsg) (TemplatesModel, tea.Cmd) {
	switch msg.String() {
	case "left", "esc", "h":
		// Exit preview mode
		m.previewing = false
		m.scroll = 0
	case "up", "k":
		if m.scroll > 0 {
			m.scroll--
		}
	case "down", "j":
		m.scroll++
	}
	return m, nil
}

// View renders the templates view
func (m TemplatesModel) View() string {
	if m.previewing {
		return m.viewPreview()
	}
	return m.viewList()
}

func (m TemplatesModel) viewList() string {
	var b strings.Builder

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
		b.WriteString(errStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
		return b.String()
	}

	if !m.loaded {
		return b.String()
	}

	if len(m.templates) == 0 {
		mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
		b.WriteString(mutedStyle.Render("No templates found"))
		b.WriteString("\n")
		return b.String()
	}

	itemStyle := lipgloss.NewStyle().MarginLeft(2)
	selectedStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true).
		MarginLeft(2)
	descStyle := lipgloss.NewStyle().
		Foreground(theme.Muted).
		MarginLeft(4)

	cursorOn := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("●")
	cursorOff := "  "

	for i, t := range m.templates {
		cursor := cursorOff
		style := itemStyle

		if i == m.cursor {
			cursor = cursorOn
			style = selectedStyle
		}

		b.WriteString(cursor + " " + style.Render(t.Name))
		b.WriteString("\n")
		b.WriteString(descStyle.Render(t.Description))
		b.WriteString("\n")
	}

	// Hint at bottom
	b.WriteString("\n")
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	b.WriteString(hintStyle.Render(keyStyle.Render("→") + " preview  " + keyStyle.Render("r") + " refresh"))
	b.WriteString("\n")

	return b.String()
}

func (m TemplatesModel) viewPreview() string {
	var b strings.Builder

	if m.cursor >= len(m.templates) {
		return b.String()
	}

	t := m.templates[m.cursor]

	// Render markdown with glamour
	var rendered string
	if m.renderer != nil && t.Content != "" {
		out, err := m.renderer.Render(t.Content)
		if err == nil {
			rendered = out
		} else {
			rendered = t.Content
		}
	} else {
		rendered = t.Content
	}

	// Split into lines and handle scrolling
	lines := strings.Split(rendered, "\n")

	// Calculate visible area (leave room for hint)
	visibleLines := m.height - 2
	if visibleLines < 1 {
		visibleLines = 10
	}

	// Clamp scroll
	maxScroll := len(lines) - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}

	// Get visible portion
	end := m.scroll + visibleLines
	if end > len(lines) {
		end = len(lines)
	}

	for i := m.scroll; i < end; i++ {
		b.WriteString(lines[i])
		b.WriteString("\n")
	}

	// Scroll indicator
	if len(lines) > visibleLines {
		hintStyle := lipgloss.NewStyle().Foreground(theme.Muted)
		keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
		scrollInfo := hintStyle.Render("  " + keyStyle.Render("↑↓") + " scroll  " + keyStyle.Render("←") + " back")
		b.WriteString(scrollInfo)
	}

	return b.String()
}
