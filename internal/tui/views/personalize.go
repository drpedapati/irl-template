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

type ProfileAction int

const (
	ActionClearProfile ProfileAction = iota
	ActionClearProjectDirectory
	ActionCount
)

// PersonalizeModel handles profile editing
type PersonalizeModel struct {
	width           int
	height          int
	inputs          []textinput.Model
	focusIndex      int
	done            bool
	doneMessage     string
	err             error
	confirming      bool
	pendingAction   ProfileAction
	confirmPrompt   string
	confirmIsDanger bool
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
		inputs:     inputs,
		focusIndex: 0,
	}
}

// SetSize sets the view dimensions
func (m *PersonalizeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// IsSaved returns true if profile was saved
func (m PersonalizeModel) IsSaved() bool {
	return m.done
}

// Init returns the initial command
func (m PersonalizeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PersonalizeModel) totalItems() int {
	return int(FieldCount) + int(ActionCount)
}

func (m PersonalizeModel) isFieldIndex(index int) bool {
	return index >= 0 && index < int(FieldCount)
}

func (m PersonalizeModel) isActionIndex(index int) bool {
	return index >= int(FieldCount) && index < m.totalItems()
}

func (m PersonalizeModel) focusedAction() ProfileAction {
	return ProfileAction(m.focusIndex - int(FieldCount))
}

func (m *PersonalizeModel) setFocus(index int) tea.Cmd {
	if m.isFieldIndex(m.focusIndex) {
		m.inputs[m.focusIndex].Blur()
	}

	m.focusIndex = index

	if m.isFieldIndex(m.focusIndex) {
		m.inputs[m.focusIndex].Focus()
		return textinput.Blink
	}

	return nil
}

func (m *PersonalizeModel) startConfirm(action ProfileAction, prompt string, danger bool) {
	m.confirming = true
	m.pendingAction = action
	m.confirmPrompt = prompt
	m.confirmIsDanger = danger
}

func (m *PersonalizeModel) performAction(action ProfileAction) {
	m.err = nil

	var err error
	var message string

	switch action {
	case ActionClearProfile:
		err = config.ClearProfile()
		message = "Profile cleared"
		if err == nil {
			for i := range m.inputs {
				m.inputs[i].SetValue("")
			}
		}
	case ActionClearProjectDirectory:
		err = config.ClearDefaultDirectory()
		message = "Project directory cleared"
	default:
		return
	}

	if err != nil {
		m.err = err
		return
	}

	m.done = true
	m.doneMessage = message
}

// Update handles messages
func (m PersonalizeModel) Update(msg tea.Msg) (PersonalizeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.confirming {
			switch msg.String() {
			case "y", "Y":
				m.confirming = false
				m.performAction(m.pendingAction)
				return m, nil
			case "n", "N", "esc":
				m.confirming = false
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "tab", "down":
			next := m.focusIndex + 1
			if next >= m.totalItems() {
				next = 0
			}
			return m, m.setFocus(next)

		case "shift+tab", "up":
			prev := m.focusIndex - 1
			if prev < 0 {
				prev = m.totalItems() - 1
			}
			return m, m.setFocus(prev)

		case "enter":
			if m.isActionIndex(m.focusIndex) {
				switch m.focusedAction() {
				case ActionClearProfile:
					m.startConfirm(ActionClearProfile, "Clear saved profile? (y/n)", true)
					return m, nil
				case ActionClearProjectDirectory:
					m.startConfirm(ActionClearProjectDirectory, "Clear default project directory? (y/n)", true)
					return m, nil
				}
				return m, nil
			}

			// Save profile
			profile := config.Profile{
				Name:         m.inputs[FieldName].Value(),
				Title:        m.inputs[FieldTitle].Value(),
				Institution:  m.inputs[FieldInstitution].Value(),
				Department:   m.inputs[FieldDepartment].Value(),
				Email:        m.inputs[FieldEmail].Value(),
				Instructions: m.inputs[FieldInstructions].Value(),
			}
			if err := config.SetProfile(profile); err != nil {
				m.err = err
				return m, nil
			}
			m.done = true
			m.doneMessage = "Profile saved"
			return m, nil
		}
	}

	// Update current input
	if !m.isFieldIndex(m.focusIndex) {
		return m, nil
	}

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
	warningStyle := lipgloss.NewStyle().Foreground(theme.Warning).MarginLeft(2)
	errorStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)

	if m.done {
		b.WriteString("\n")
		b.WriteString(successStyle.Render(m.doneMessage))
		b.WriteString("\n")
		return b.String()
	}

	b.WriteString("\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if m.confirming {
		style := warningStyle
		if m.confirmIsDanger {
			style = errorStyle
		}
		b.WriteString(style.Render(m.confirmPrompt))
		b.WriteString("\n\n")
	}

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

	// Actions menu
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("Actions"))
	b.WriteString("\n\n")

	actionNames := []string{
		"Clear profile",
		"Clear project directory",
	}
	actionDescs := []string{
		"Removes saved profile fields",
		"Unsets the default project directory (does not delete files)",
	}

	actionCursorStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	actionNormalStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	actionSelectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	for i := 0; i < int(ActionCount); i++ {
		cursor := "  "
		style := actionNormalStyle
		if m.focusIndex == int(FieldCount)+i {
			cursor = actionCursorStyle.Render("● ")
			style = actionSelectedStyle
		}

		b.WriteString("  " + cursor + style.Render(actionNames[i]))
		b.WriteString("\n")
		b.WriteString("     " + actionNormalStyle.Render(actionDescs[i]))
		b.WriteString("\n\n")
	}

	return b.String()
}
