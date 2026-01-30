package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// ConfigModel is the configuration view
type ConfigModel struct {
	width      int
	height     int
	defaultDir string
	editing    bool
	input      textinput.Model
	saved      bool
	err        error
}

// NewConfigModel creates a new config view
func NewConfigModel() ConfigModel {
	ti := textinput.New()
	ti.Placeholder = "Enter directory path..."
	ti.Width = 50

	return ConfigModel{
		input: ti,
	}
}

// SetSize sets the view dimensions
func (m *ConfigModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Load returns a command that loads the config
func (m *ConfigModel) Load() tea.Cmd {
	return func() tea.Msg {
		dir := config.GetDefaultDirectory()
		return ConfigLoadedMsg{DefaultDir: dir}
	}
}

// Update handles messages
func (m ConfigModel) Update(msg tea.Msg) (ConfigModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ConfigLoadedMsg:
		m.defaultDir = msg.DefaultDir
		m.err = msg.Err
		m.input.SetValue(m.defaultDir)
		return m, nil

	case tea.KeyMsg:
		if m.editing {
			return m.updateEditing(msg)
		}
		return m.updateViewing(msg)
	}

	return m, nil
}

func (m ConfigModel) updateViewing(msg tea.KeyMsg) (ConfigModel, tea.Cmd) {
	switch msg.String() {
	case "e":
		m.editing = true
		m.saved = false
		m.input.SetValue(m.defaultDir)
		m.input.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

func (m ConfigModel) updateEditing(msg tea.KeyMsg) (ConfigModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Save the config
		newDir := m.input.Value()
		if err := config.SetDefaultDirectory(newDir); err != nil {
			m.err = err
		} else {
			m.defaultDir = newDir
			m.saved = true
		}
		m.editing = false
		return m, nil
	case "esc":
		m.editing = false
		m.input.SetValue(m.defaultDir)
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the config view
func (m ConfigModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Primary).
		MarginLeft(2).
		MarginTop(1)

	b.WriteString(titleStyle.Render("Configuration"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	valueStyle := lipgloss.NewStyle().Foreground(theme.Accent).MarginLeft(4)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)

	b.WriteString(labelStyle.Render("Default project directory"))
	b.WriteString("\n")

	if m.editing {
		b.WriteString("  " + m.input.View())
		b.WriteString("\n\n")
		b.WriteString(hintStyle.Render("Enter to save, Esc to cancel"))
	} else {
		if m.defaultDir == "" {
			b.WriteString(valueStyle.Render("(not set - uses current directory)"))
		} else {
			b.WriteString(valueStyle.Render(m.defaultDir))
		}
		b.WriteString("\n\n")

		if m.saved {
			checkStyle := lipgloss.NewStyle().Foreground(theme.Success).MarginLeft(2)
			b.WriteString(checkStyle.Render("✓ Saved"))
			b.WriteString("\n\n")
		}

		if m.err != nil {
			errStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
			b.WriteString(errStyle.Render("Error: " + m.err.Error()))
			b.WriteString("\n\n")
		}

		keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
		b.WriteString(hintStyle.Render("Press ") + keyStyle.Render("e") + hintStyle.Render(" to edit"))
	}

	b.WriteString("\n")

	return b.String()
}
