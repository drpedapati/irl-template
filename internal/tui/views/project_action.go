package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// ProjectActionModel displays a project with available actions
type ProjectActionModel struct {
	projectPath string
	projectName string
	editors     []AppInfo // Uses unified AppInfo from editors.go
	tools       []AppInfo // Finder, Terminal, etc.
	message     string    // Feedback message after opening
	done        bool
	isNew       bool // True if this is a newly created project
}

// NewProjectActionModel creates a new project action view
func NewProjectActionModel(projectPath string, isNew bool) ProjectActionModel {
	// Extract project name from path
	parts := strings.Split(projectPath, "/")
	name := projectPath
	if len(parts) > 0 {
		name = parts[len(parts)-1]
	}

	return ProjectActionModel{
		projectPath: projectPath,
		projectName: name,
		editors:     GetInstalledEditors(), // Uses unified source
		tools:       GetInstalledTools(),   // Uses unified source
		isNew:       isNew,
	}
}

// ProjectPath returns the project path
func (m ProjectActionModel) ProjectPath() string {
	return m.projectPath
}

// IsDone returns true when user wants to exit
func (m ProjectActionModel) IsDone() bool {
	return m.done
}

// Update handles messages
func (m ProjectActionModel) Update(msg tea.Msg) (ProjectActionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		// Check for editor hotkeys
		for _, editor := range m.editors {
			if key == editor.Key {
				if errMsg := OpenProjectWith(editor, m.projectPath); errMsg != "" {
					m.message = errMsg
				} else {
					m.message = "Opened in " + editor.Name
				}
				return m, nil
			}
		}

		// Check for tool hotkeys
		for _, tool := range m.tools {
			if key == tool.Key {
				if errMsg := OpenProjectWith(tool, m.projectPath); errMsg != "" {
					m.message = errMsg
				} else {
					m.message = "Opened in " + tool.Name
				}
				return m, nil
			}
		}

		// Handle other keys
		switch key {
		case "enter", "esc", "left":
			m.done = true
			return m, nil
		}
	}

	return m, nil
}

// View renders the project action view
func (m ProjectActionModel) View() string {
	var b strings.Builder

	checkStyle := lipgloss.NewStyle().Foreground(theme.Success)
	pathStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	nameStyle := lipgloss.NewStyle().Foreground(theme.Primary)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	headerStyle := lipgloss.NewStyle().Foreground(theme.Muted).Bold(true)
	messageStyle := lipgloss.NewStyle().Foreground(theme.Success).MarginLeft(2)

	// Header
	b.WriteString("\n")
	if m.isNew {
		b.WriteString("  " + checkStyle.Render("✓") + " Project created successfully")
	} else {
		b.WriteString("  " + headerStyle.Render("Project"))
	}
	b.WriteString("\n\n")

	// Project path
	b.WriteString("  " + pathStyle.Render(m.projectPath))
	b.WriteString("\n\n")

	// Feedback message (if any)
	if m.message != "" {
		if strings.HasPrefix(m.message, "Failed") {
			errorStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
			b.WriteString(errorStyle.Render("✗ " + m.message))
		} else {
			b.WriteString(messageStyle.Render("✓ " + m.message))
		}
		b.WriteString("\n\n")
	}

	// Editors section
	if len(m.editors) > 0 {
		b.WriteString("  " + headerStyle.Render("Editors & IDEs:"))
		b.WriteString("\n\n")

		// Show available editors in a grid (3 columns)
		colWidth := 16
		col := 0
		for _, editor := range m.editors {
			if col == 0 {
				b.WriteString("  ")
			}

			item := keyStyle.Render(editor.Key) + " " + nameStyle.Render(editor.Name)

			// Pad to column width
			padding := colWidth - len(editor.Key) - 1 - len(editor.Name)
			if padding > 0 {
				item += strings.Repeat(" ", padding)
			}

			b.WriteString(item)
			col++
			if col >= 3 {
				b.WriteString("\n")
				col = 0
			}
		}
		if col != 0 {
			b.WriteString("\n")
		}
	}

	// Tools section
	if len(m.tools) > 0 {
		b.WriteString("\n")
		b.WriteString("  " + headerStyle.Render("Tools:"))
		b.WriteString("\n\n")

		colWidth := 16
		col := 0
		for _, tool := range m.tools {
			if col == 0 {
				b.WriteString("  ")
			}

			item := keyStyle.Render(tool.Key) + " " + nameStyle.Render(tool.Name)

			padding := colWidth - len(tool.Key) - 1 - len(tool.Name)
			if padding > 0 {
				item += strings.Repeat(" ", padding)
			}

			b.WriteString(item)
			col++
			if col >= 3 {
				b.WriteString("\n")
				col = 0
			}
		}
		if col != 0 {
			b.WriteString("\n")
		}
	}

	// Footer hint
	b.WriteString("\n")
	b.WriteString("  " + hintStyle.Render("Press keys to open, Enter to finish"))

	return b.String()
}
