package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update templates from GitHub",
	Long:  "Fetch the latest templates from the IRL template repository",
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// Spinner model for the update operation
type updateModel struct {
	spinner  spinner.Model
	done     bool
	err      error
	list     []templates.Template
	quitting bool
}

type updateDoneMsg struct {
	err  error
	list []templates.Template
}

func (m updateModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, doUpdate())
}

func doUpdate() tea.Cmd {
	return func() tea.Msg {
		err := templates.Update()
		if err != nil {
			return updateDoneMsg{err: err}
		}
		list, _ := templates.ListTemplates()
		return updateDoneMsg{list: list}
	}
}

func (m updateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case updateDoneMsg:
		m.done = true
		m.err = msg.err
		m.list = msg.list
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m updateModel) View() string {
	if m.quitting {
		return ""
	}
	if m.done {
		return ""
	}
	return fmt.Sprintf("%s Grabbing the latest templates...\n", m.spinner.View())
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Check if we're in a TTY - if not, skip the spinner
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		// Not a TTY, use simple output
		fmt.Println(theme.Faint("Grabbing the latest templates..."))

		if err := templates.Update(); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		list, _ := templates.ListTemplates()
		printUpdateSuccess(list)
		return nil
	}

	// Create spinner with friendly MiniDot style
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(theme.Accent)

	m := updateModel{spinner: s}

	// Run the spinner program
	p := tea.NewProgram(m)

	// Add a small delay to ensure spinner is visible
	go func() {
		time.Sleep(100 * time.Millisecond)
	}()

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("spinner error: %w", err)
	}

	final := finalModel.(updateModel)

	if final.quitting {
		return nil
	}

	if final.err != nil {
		return fmt.Errorf("update failed: %w", final.err)
	}

	printUpdateSuccess(final.list)
	return nil
}

func printUpdateSuccess(list []templates.Template) {
	fmt.Printf("%s %s templates refreshed\n",
		theme.OK("All set!"),
		theme.Cmd(fmt.Sprintf("%d", len(list))))

	fmt.Println()
	for _, t := range list {
		fmt.Printf("  %s %s %s\n",
			theme.Cmd(theme.Dot),
			theme.B(t.Name),
			theme.Faint(t.Description))
	}
}
