package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// TemplatesModel is the templates browser view
type TemplatesModel struct {
	templates []TemplateItem
	cursor    int
	width     int
	height    int
	loaded    bool
	err       error
}

// NewTemplatesModel creates a new templates view
func NewTemplatesModel() TemplatesModel {
	return TemplatesModel{}
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
			}
		}
		return TemplatesLoadedMsg{Templates: items}
	}
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
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.templates)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

// View renders the templates view
func (m TemplatesModel) View() string {
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
	b.WriteString(hintStyle.Render("Press r to refresh from GitHub"))
	b.WriteString("\n")

	return b.String()
}
