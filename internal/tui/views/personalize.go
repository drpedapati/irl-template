package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// ProfileField identifies which field is being edited
type ProfileField int

const (
	FieldName ProfileField = iota
	FieldTitle
	FieldInstitution
	FieldDepartment
	FieldEmail
	FieldInstructions
	FieldCount
)

// PersonalizeModel handles profile editing
type PersonalizeModel struct {
	width        int
	height       int
	inputs       []textinput.Model
	focusIndex   int
	saved        bool
	instructions string // Separate since it's multiline
}

// NewPersonalizeModel creates a new personalize view
func NewPersonalizeModel() PersonalizeModel {
	profile := config.GetProfile()

	inputs := make([]textinput.Model, FieldCount)

	// Name
	inputs[FieldName] = textinput.New()
	inputs[FieldName].Placeholder = "Your name"
	inputs[FieldName].SetValue(profile.Name)
	inputs[FieldName].Width = 40
	inputs[FieldName].Focus()

	// Title
	inputs[FieldTitle] = textinput.New()
	inputs[FieldTitle].Placeholder = "PhD Candidate, Professor, etc."
	inputs[FieldTitle].SetValue(profile.Title)
	inputs[FieldTitle].Width = 40

	// Institution
	inputs[FieldInstitution] = textinput.New()
	inputs[FieldInstitution].Placeholder = "University or organization"
	inputs[FieldInstitution].SetValue(profile.Institution)
	inputs[FieldInstitution].Width = 40

	// Department
	inputs[FieldDepartment] = textinput.New()
	inputs[FieldDepartment].Placeholder = "Department or lab"
	inputs[FieldDepartment].SetValue(profile.Department)
	inputs[FieldDepartment].Width = 40

	// Email
	inputs[FieldEmail] = textinput.New()
	inputs[FieldEmail].Placeholder = "email@example.com"
	inputs[FieldEmail].SetValue(profile.Email)
	inputs[FieldEmail].Width = 40

	// Instructions (shown as single line but stored as multiline)
	inputs[FieldInstructions] = textinput.New()
	inputs[FieldInstructions].Placeholder = "Common AI instructions for all projects"
	inputs[FieldInstructions].SetValue(profile.Instructions)
	inputs[FieldInstructions].Width = 50

	return PersonalizeModel{
		inputs:       inputs,
		focusIndex:   0,
		instructions: profile.Instructions,
	}
}

// SetSize sets the view dimensions
func (m *PersonalizeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// IsSaved returns true if profile was saved
func (m PersonalizeModel) IsSaved() bool {
	return m.saved
}

// Init returns the initial command
func (m PersonalizeModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m PersonalizeModel) Update(msg tea.Msg) (PersonalizeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "j":
			m.inputs[m.focusIndex].Blur()
			m.focusIndex++
			if m.focusIndex >= int(FieldCount) {
				m.focusIndex = 0
			}
			m.inputs[m.focusIndex].Focus()
			return m, textinput.Blink

		case "shift+tab", "up", "k":
			m.inputs[m.focusIndex].Blur()
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = int(FieldCount) - 1
			}
			m.inputs[m.focusIndex].Focus()
			return m, textinput.Blink

		case "enter":
			// Save profile
			profile := config.Profile{
				Name:         m.inputs[FieldName].Value(),
				Title:        m.inputs[FieldTitle].Value(),
				Institution:  m.inputs[FieldInstitution].Value(),
				Department:   m.inputs[FieldDepartment].Value(),
				Email:        m.inputs[FieldEmail].Value(),
				Instructions: m.inputs[FieldInstructions].Value(),
			}
			config.SetProfile(profile)
			m.saved = true
			return m, nil
		}
	}

	// Update current input
	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

// View renders the personalize view
func (m PersonalizeModel) View() string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().Foreground(theme.Muted).Width(14)
	focusedLabelStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Width(14)
	successStyle := lipgloss.NewStyle().Foreground(theme.Success).MarginLeft(2)

	if m.saved {
		b.WriteString("\n")
		b.WriteString(successStyle.Render("Profile saved"))
		b.WriteString("\n")
		return b.String()
	}

	b.WriteString("\n")

	fields := []string{"Name", "Title", "Institution", "Department", "Email", "Instructions"}

	for i, input := range m.inputs {
		style := labelStyle
		if i == m.focusIndex {
			style = focusedLabelStyle
		}

		b.WriteString("  ")
		b.WriteString(style.Render(fields[i]))
		b.WriteString(input.View())
		b.WriteString("\n")
	}

	return b.String()
}
