package views

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/editor"
	"github.com/drpedapati/irl-template/pkg/scaffold"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// AdoptStep represents the current step in the adopt wizard
type AdoptStep int

const (
	AdoptStepBrowse AdoptStep = iota
	AdoptStepTemplate
	AdoptStepAdopting
	AdoptStepDone
)

// AdoptModel is the folder adoption wizard
type AdoptModel struct {
	step          AdoptStep
	width         int
	height        int
	sourcePath    string // Selected source folder

	// Browse state (same pattern as init.go)
	browseDir     string
	browseFolders []browseEntry
	browseCursor  int
	browseScroll  int
	browseSortBy  string // "name-asc" or "date-desc"

	// Template state
	templates   []templates.Template
	templateIdx int

	// Result
	spinner     spinner.Model
	projectPath string
	err         error
	done        bool
	actionView  ProjectActionModel
}

const adoptBrowseVisibleItems = 8

// AdoptProjectMsg is sent when the adopt operation completes
type AdoptProjectMsg struct {
	Path string
	Err  error
}

// AdoptTemplatesLoadedMsg is sent when templates are loaded for adopt wizard
type AdoptTemplatesLoadedMsg struct {
	Templates []templates.Template
}

// NewAdoptModel creates a new adopt view
func NewAdoptModel() AdoptModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.Accent)

	home, _ := os.UserHomeDir()

	m := AdoptModel{
		step:      AdoptStepBrowse,
		browseDir: home,
		spinner:   s,
	}
	m.loadFolders()
	return m
}

// SetSize sets the view dimensions
func (m *AdoptModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init returns the initial command
func (m AdoptModel) Init() tea.Cmd {
	return nil
}

// CanGoBack returns true if we can go back within the wizard
func (m AdoptModel) CanGoBack() bool {
	return m.step == AdoptStepTemplate
}

// Done returns true if the wizard is complete
func (m AdoptModel) Done() bool {
	return m.done
}

// Update handles messages
func (m AdoptModel) Update(msg tea.Msg) (AdoptModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case editor.EditorFinishedMsg, editor.EditorOpenedMsg:
		if m.step == AdoptStepDone {
			var cmd tea.Cmd
			m.actionView, cmd = m.actionView.Update(msg)
			return m, cmd
		}

	case tea.KeyMsg:
		switch m.step {
		case AdoptStepBrowse:
			return m.updateBrowse(msg)
		case AdoptStepTemplate:
			return m.updateTemplate(msg)
		case AdoptStepDone:
			return m.updateDone(msg)
		}

	case spinner.TickMsg:
		if m.step == AdoptStepAdopting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case AdoptTemplatesLoadedMsg:
		m.templates = msg.Templates
		m.step = AdoptStepTemplate

	case AdoptProjectMsg:
		m.step = AdoptStepDone
		m.projectPath = msg.Path
		m.err = msg.Err
		if msg.Err == nil {
			m.actionView = NewProjectActionModel(msg.Path, true)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *AdoptModel) loadFolders() {
	m.browseFolders = []browseEntry{}
	m.browseCursor = 0
	m.browseScroll = 0

	entries, err := os.ReadDir(m.browseDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			modTime := time.Time{}
			if info, err := e.Info(); err == nil {
				modTime = info.ModTime()
			}
			m.browseFolders = append(m.browseFolders, browseEntry{Name: e.Name(), ModTime: modTime})
		}
	}
	m.applyBrowseSort()
}

func (m *AdoptModel) applyBrowseSort() {
	switch m.browseSortBy {
	case "date-desc":
		sort.Slice(m.browseFolders, func(i, j int) bool {
			return m.browseFolders[i].ModTime.After(m.browseFolders[j].ModTime)
		})
	default: // "name-asc"
		sort.Slice(m.browseFolders, func(i, j int) bool {
			return strings.ToLower(m.browseFolders[i].Name) < strings.ToLower(m.browseFolders[j].Name)
		})
	}
}

func (m AdoptModel) updateBrowse(msg tea.KeyMsg) (AdoptModel, tea.Cmd) {
	switch msg.String() {
	case "up":
		if m.browseCursor > 0 {
			m.browseCursor--
			if m.browseCursor < m.browseScroll {
				m.browseScroll = m.browseCursor
			}
		}
	case "down":
		if m.browseCursor < len(m.browseFolders)-1 {
			m.browseCursor++
			if m.browseCursor >= m.browseScroll+adoptBrowseVisibleItems {
				m.browseScroll = m.browseCursor - adoptBrowseVisibleItems + 1
			}
		}
	case "s":
		// Toggle sort
		if m.browseSortBy == "date-desc" {
			m.browseSortBy = "name-asc"
		} else {
			m.browseSortBy = "date-desc"
		}
		m.applyBrowseSort()
		m.browseCursor = 0
		m.browseScroll = 0
	case "right":
		// Enter selected folder
		if len(m.browseFolders) > 0 && m.browseCursor < len(m.browseFolders) {
			m.browseDir = filepath.Join(m.browseDir, m.browseFolders[m.browseCursor].Name)
			m.loadFolders()
		}
	case "left":
		// Go up one level
		parent := filepath.Dir(m.browseDir)
		if parent != m.browseDir {
			m.browseDir = parent
			m.loadFolders()
		}
	case "enter":
		// Select current folder as source
		if len(m.browseFolders) > 0 && m.browseCursor < len(m.browseFolders) {
			m.sourcePath = filepath.Join(m.browseDir, m.browseFolders[m.browseCursor].Name)
		} else {
			// If empty folder list, select current directory
			m.sourcePath = m.browseDir
		}
		// Load templates for next step
		m.step = AdoptStepTemplate
		return m, tea.Batch(m.loadTemplates(), m.spinner.Tick)
	case "esc":
		// Let parent handle going back to menu
		return m, nil
	}
	return m, nil
}

func (m AdoptModel) updateTemplate(msg tea.KeyMsg) (AdoptModel, tea.Cmd) {
	switch msg.String() {
	case "up":
		if m.templateIdx > 0 {
			m.templateIdx--
		}
	case "down":
		if m.templateIdx < len(m.templates) {
			m.templateIdx++ // +1 for "None" option
		}
	case "enter", "right":
		m.step = AdoptStepAdopting
		return m, tea.Batch(m.adoptProject(), m.spinner.Tick)
	case "esc", "left":
		// Go back to browse
		m.step = AdoptStepBrowse
		return m, nil
	}
	return m, nil
}

func (m AdoptModel) updateDone(msg tea.KeyMsg) (AdoptModel, tea.Cmd) {
	var cmd tea.Cmd
	m.actionView, cmd = m.actionView.Update(msg)

	if m.actionView.IsDone() {
		m.done = true
	}

	return m, cmd
}

func (m AdoptModel) loadTemplates() tea.Cmd {
	return func() tea.Msg {
		list, _ := templates.ListTemplates()
		return AdoptTemplatesLoadedMsg{Templates: list}
	}
}

func (m AdoptModel) adoptProject() tea.Cmd {
	return func() tea.Msg {
		baseDir := config.GetDefaultDirectory()
		if baseDir == "" {
			return AdoptProjectMsg{Err: fmt.Errorf("no default directory configured")}
		}

		folderName := filepath.Base(m.sourcePath)
		destPath := filepath.Join(baseDir, folderName)

		// Check destination doesn't exist
		if _, err := os.Stat(destPath); !os.IsNotExist(err) {
			return AdoptProjectMsg{Err: fmt.Errorf("'%s' already exists in workspace", folderName)}
		}

		// Ensure workspace exists
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return AdoptProjectMsg{Err: err}
		}

		// Copy folder
		cpCmd := exec.Command("cp", "-r", m.sourcePath, destPath)
		if err := cpCmd.Run(); err != nil {
			return AdoptProjectMsg{Err: fmt.Errorf("copy failed: %w", err)}
		}

		// Add .gitignore if not present
		gitignorePath := filepath.Join(destPath, ".gitignore")
		if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
			scaffold.Create(destPath)
		}

		// Add plans/main-plan.md if no plan file exists
		hasPlan := fileExists(filepath.Join(destPath, "plans", "main-plan.md")) ||
			fileExists(filepath.Join(destPath, "main-plan.md")) ||
			fileExists(filepath.Join(destPath, "01-plans", "main-plan.md"))

		if !hasPlan {
			var planContent string
			if m.templateIdx > 0 && m.templateIdx <= len(m.templates) {
				planContent = m.templates[m.templateIdx-1].Content
			} else {
				planContent = templates.EmbeddedTemplates["irl-basic"].Content
			}
			planContent = scaffold.InjectProfile(planContent)
			scaffold.WritePlan(destPath, planContent)
		}

		// Git init if not already a repo
		gitDir := filepath.Join(destPath, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			scaffold.GitInit(destPath)
		}

		return AdoptProjectMsg{Path: destPath}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// View renders the adopt view
func (m AdoptModel) View() string {
	var b strings.Builder

	switch m.step {
	case AdoptStepBrowse:
		b.WriteString(m.viewBrowse())
	case AdoptStepTemplate:
		if len(m.templates) == 0 && m.step == AdoptStepTemplate {
			// Still loading
			b.WriteString("  " + m.spinner.View() + " Loading templates...")
		} else {
			b.WriteString(m.viewTemplate())
		}
	case AdoptStepAdopting:
		b.WriteString("  " + m.spinner.View() + " Adopting folder...")
	case AdoptStepDone:
		b.WriteString(m.viewDone())
	}

	return b.String()
}

func (m AdoptModel) viewBrowse() string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().MarginLeft(2)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	pathStyle := lipgloss.NewStyle().Foreground(theme.Accent).MarginLeft(2)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	b.WriteString(labelStyle.Render("Select folder to adopt"))
	b.WriteString("\n")
	b.WriteString(pathStyle.Render(m.browseDir))
	b.WriteString("\n\n")

	cursorOn := keyStyle.Render("●")
	cursorOff := "  "
	itemStyle := lipgloss.NewStyle()
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	if len(m.browseFolders) == 0 {
		b.WriteString(hintStyle.Render("  (empty)"))
		b.WriteString("\n")
	} else {
		end := m.browseScroll + adoptBrowseVisibleItems
		if end > len(m.browseFolders) {
			end = len(m.browseFolders)
		}

		for i := m.browseScroll; i < end; i++ {
			entry := m.browseFolders[i]
			cursor := cursorOff
			style := itemStyle

			if i == m.browseCursor {
				cursor = cursorOn
				style = selectedStyle
			}

			b.WriteString("  " + cursor + " " + style.Render(entry.Name+"/"))
			b.WriteString("\n")
		}

		// Scroll indicator
		if len(m.browseFolders) > adoptBrowseVisibleItems {
			shown := fmt.Sprintf("%d-%d of %d", m.browseScroll+1, end, len(m.browseFolders))
			b.WriteString(hintStyle.Render("  " + shown))
			b.WriteString("\n")
		}
	}

	// Sort indicator
	sortLabel := "A-Z"
	if m.browseSortBy == "date-desc" {
		sortLabel = "newest"
	}

	b.WriteString("\n")
	b.WriteString(hintStyle.Render(keyStyle.Render("←") + " up  " + keyStyle.Render("→") + " open  " + keyStyle.Render("Enter") + " adopt  " + keyStyle.Render("s") + ":" + sortLabel))
	b.WriteString("\n")

	return b.String()
}

func (m AdoptModel) viewTemplate() string {
	var b strings.Builder

	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	pathStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	nameStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).MarginLeft(2)

	b.WriteString(pathStyle.Render("Adopting:"))
	b.WriteString("\n")
	b.WriteString(nameStyle.Render(filepath.Base(m.sourcePath)))
	b.WriteString("\n\n")

	b.WriteString(hintStyle.Render("Pick a template for the plan file"))
	b.WriteString("\n\n")

	cursorOn := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("●")
	cursorOff := "  "
	itemStyle := lipgloss.NewStyle().MarginLeft(2)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).MarginLeft(2)

	// "Basic" option (default)
	cursor := cursorOff
	style := itemStyle
	if m.templateIdx == 0 {
		cursor = cursorOn
		style = selectedStyle
	}
	b.WriteString(cursor + " " + style.Render("irl-basic (default)"))
	b.WriteString("\n")

	// Template options
	for i, t := range m.templates {
		cursor = cursorOff
		style = itemStyle
		if m.templateIdx == i+1 {
			cursor = cursorOn
			style = selectedStyle
		}
		b.WriteString(cursor + " " + style.Render(t.Name))
		b.WriteString("\n")
	}

	return b.String()
}

func (m AdoptModel) viewDone() string {
	var b strings.Builder

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
		b.WriteString(errStyle.Render("✗ " + m.err.Error()))
		b.WriteString("\n\n")
		hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
		b.WriteString(hintStyle.Render("Press Esc to go back"))
		return b.String()
	}

	return m.actionView.View()
}
