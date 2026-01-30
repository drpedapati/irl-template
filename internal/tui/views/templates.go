package views

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/glamour"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// TemplatesModel is the templates browser view
type TemplatesModel struct {
	templates   []TemplateItem
	filtered    []TemplateItem
	cursor      int
	scroll      int
	width       int
	height      int
	loaded      bool
	err         error
	previewing  bool
	previewName string // Name of template being previewed
	renderer    *glamour.TermRenderer
	filterInput textinput.Model
	sortBy      string // "name-asc", "name-desc"
	sourceMode  string // "all", "default", "custom"

	// Copy mode state
	copying       bool
	copyInput     textinput.Model
	copySourceIdx int

	// Edit mode state
	editing bool
	editors []Editor

	// Feedback message
	message    string
	messageErr bool
}

const templatesVisibleItems = 8

// NewTemplatesModel creates a new templates view
func NewTemplatesModel() TemplatesModel {
	// Create glamour renderer with dark style
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(68),
	)

	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Width = 40
	ti.Focus() // Auto-focus on creation

	return TemplatesModel{
		renderer:    r,
		filterInput: ti,
		sortBy:      "name-asc", // Default alphabetical
		sourceMode:  "all",      // Show all templates by default
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
		var items []TemplateItem

		// Load default templates from GitHub/cache
		list, err := templates.ListTemplates()
		if err == nil {
			for _, t := range list {
				items = append(items, TemplateItem{
					Name:        t.Name,
					Description: t.Description,
					Content:     t.Content,
					Source:      "default",
				})
			}
		}

		// Load custom templates from _templates folder
		customTemplates := loadCustomTemplates()
		items = append(items, customTemplates...)

		return TemplatesLoadedMsg{Templates: items, Err: err}
	}
}

// loadCustomTemplates scans the _templates folder for custom templates
func loadCustomTemplates() []TemplateItem {
	var items []TemplateItem

	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return items
	}

	templatesDir := filepath.Join(baseDir, "_templates")
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return items // Folder doesn't exist or can't be read
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip hidden folders
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Check for main-plan.md
		planPath := filepath.Join(templatesDir, entry.Name(), "main-plan.md")
		content, err := os.ReadFile(planPath)
		if err != nil {
			continue // No main-plan.md, skip
		}

		// Extract description from first heading or use folder name
		description := extractDescription(string(content))
		if description == "" {
			description = "Custom template"
		}

		items = append(items, TemplateItem{
			Name:        entry.Name(),
			Description: description,
			Content:     string(content),
			Source:      "custom",
		})
	}

	return items
}

// extractDescription tries to get a description from the template content
func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Look for first heading
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

// RefreshTemplates forces a refresh from GitHub
func (m *TemplatesModel) RefreshTemplates() tea.Cmd {
	return func() tea.Msg {
		var items []TemplateItem

		// Fetch fresh templates from GitHub
		list, err := templates.FetchTemplates()
		if err == nil {
			for _, t := range list {
				items = append(items, TemplateItem{
					Name:        t.Name,
					Description: t.Description,
					Content:     t.Content,
					Source:      "default",
				})
			}
		}

		// Also reload custom templates
		customTemplates := loadCustomTemplates()
		items = append(items, customTemplates...)

		return TemplatesLoadedMsg{Templates: items, Err: err}
	}
}

// IsPreviewing returns true if in preview mode
func (m TemplatesModel) IsPreviewing() bool {
	return m.previewing
}

// PreviewingName returns the name of the template being previewed (empty if not previewing)
func (m TemplatesModel) PreviewingName() string {
	if m.previewing {
		return m.previewName
	}
	return ""
}

// IsCopying returns true if in copy mode
func (m TemplatesModel) IsCopying() bool {
	return m.copying
}

// IsEditing returns true if in edit mode
func (m TemplatesModel) IsEditing() bool {
	return m.editing
}

// HasFilterText returns true if there's text in the filter input
func (m TemplatesModel) HasFilterText() bool {
	return m.filterInput.Value() != ""
}

// ClearFilter clears the filter input and resets the list
func (m *TemplatesModel) ClearFilter() {
	m.filterInput.SetValue("")
	m.applyFilter()
}

// Update handles messages
func (m TemplatesModel) Update(msg tea.Msg) (TemplatesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case TemplatesLoadedMsg:
		m.loaded = true
		m.err = msg.Err
		m.templates = msg.Templates
		m.filtered = msg.Templates
		m.applySort()
		return m, textinput.Blink

	case tea.KeyMsg:
		// Clear any message on keypress
		m.message = ""
		m.messageErr = false

		if m.copying {
			return m.updateCopying(msg)
		}
		if m.editing {
			return m.updateEditing(msg)
		}
		if m.previewing {
			return m.updatePreview(msg)
		}
		return m.updateList(msg)
	}

	return m, nil
}

func (m TemplatesModel) updateList(msg tea.KeyMsg) (TemplatesModel, tea.Cmd) {
	key := msg.String()

	switch key {
	case "a":
		// Show all templates
		m.sourceMode = "all"
		m.applyFilter()
		return m, nil
	case "d":
		// Show only default (GitHub) templates
		m.sourceMode = "default"
		m.applyFilter()
		return m, nil
	case "c":
		// Show only custom templates
		m.sourceMode = "custom"
		m.applyFilter()
		return m, nil
	case "t":
		// Copy as custom template
		if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
			m.startCopyMode()
			return m, textinput.Blink
		}
		return m, nil
	case "e":
		// Edit custom template
		if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
			t := m.filtered[m.cursor]
			if t.Source == "custom" {
				m.startEditMode()
			} else {
				m.message = "Press t to create a custom copy first"
				m.messageErr = false
			}
		}
		return m, nil
	case "s":
		// Toggle sort: A-Z → Z-A → A-Z
		if m.sortBy == "name-asc" {
			m.sortBy = "name-desc"
		} else {
			m.sortBy = "name-asc"
		}
		m.applySort()
		return m, nil
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.scroll {
				m.scroll = m.cursor
			}
		}
		return m, nil
	case "down", "j":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
			if m.cursor >= m.scroll+templatesVisibleItems {
				m.scroll = m.cursor - templatesVisibleItems + 1
			}
		}
		return m, nil
	case "right", "enter", "l":
		// Enter preview mode
		if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
			m.previewing = true
			m.previewName = m.filtered[m.cursor].Name
			m.scroll = 0
		}
		return m, nil
	default:
		// Pass other keys to filter input
		var cmd tea.Cmd
		m.filterInput, cmd = m.filterInput.Update(msg)
		m.applyFilter()
		return m, cmd
	}
}

func (m *TemplatesModel) startCopyMode() {
	t := m.filtered[m.cursor]
	m.copying = true
	m.copySourceIdx = m.cursor

	// Create input with suggested name
	ti := textinput.New()
	ti.Placeholder = "template-name"
	ti.Width = 40
	ti.SetValue(t.Name + "-custom")
	ti.Focus()
	ti.CursorEnd()
	m.copyInput = ti
}

func (m *TemplatesModel) startEditMode() {
	m.editing = true
	m.editors = detectEditors()
}

func (m TemplatesModel) updateCopying(msg tea.KeyMsg) (TemplatesModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.copying = false
		m.filterInput.Focus()
		return m, textinput.Blink
	case "enter":
		name := m.copyInput.Value()
		if name == "" {
			return m, nil
		}
		// Create the custom template
		err := m.createCustomTemplate(name)
		m.copying = false
		m.filterInput.Focus()
		if err != nil {
			m.message = "Error: " + err.Error()
			m.messageErr = true
		} else {
			m.message = "✓ Created: _templates/" + name
			m.messageErr = false
			// Reload templates to include the new one
			return m, m.LoadTemplates()
		}
		return m, textinput.Blink
	default:
		var cmd tea.Cmd
		m.copyInput, cmd = m.copyInput.Update(msg)
		return m, cmd
	}
}

func (m *TemplatesModel) createCustomTemplate(name string) error {
	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return os.ErrNotExist
	}

	// Sanitize name - replace spaces with dashes, lowercase
	name = strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	// Create _templates folder if needed
	templatesDir := filepath.Join(baseDir, "_templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return err
	}

	// Create template folder
	templateDir := filepath.Join(templatesDir, name)
	if _, err := os.Stat(templateDir); !os.IsNotExist(err) {
		return os.ErrExist // Already exists
	}
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return err
	}

	// Copy content
	source := m.filtered[m.copySourceIdx]
	planPath := filepath.Join(templateDir, "main-plan.md")
	return os.WriteFile(planPath, []byte(source.Content), 0644)
}

func (m TemplatesModel) updateEditing(msg tea.KeyMsg) (TemplatesModel, tea.Cmd) {
	key := msg.String()

	switch key {
	case "esc":
		m.editing = false
		return m, nil
	default:
		// Check if it's an editor shortcut
		for _, e := range m.editors {
			if key == e.Key {
				m.openInEditor(e)
				m.editing = false
				m.message = "Opened in " + e.Name
				return m, nil
			}
		}
	}
	return m, nil
}

func (m *TemplatesModel) openInEditor(e Editor) {
	if m.cursor >= len(m.filtered) {
		return
	}

	t := m.filtered[m.cursor]
	if t.Source != "custom" {
		return
	}

	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return
	}

	filePath := filepath.Join(baseDir, "_templates", t.Name, "main-plan.md")

	var cmd *exec.Cmd
	switch e.Cmd {
	case "terminal":
		dir := filepath.Dir(filePath)
		if runtime.GOOS == "darwin" {
			script := `tell application "Terminal" to do script "cd '` + dir + `'"`
			cmd = exec.Command("osascript", "-e", script)
		} else {
			cmd = exec.Command("x-terminal-emulator", "--working-directory", dir)
		}
	case "positron":
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", "-a", "Positron", filePath)
		} else {
			cmd = exec.Command("positron", filePath)
		}
	case "cursor":
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", "-a", "Cursor", filePath)
		} else {
			cmd = exec.Command("cursor", filePath)
		}
	case "rstudio":
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", "-a", "RStudio", filePath)
		} else {
			cmd = exec.Command("rstudio", filePath)
		}
	default:
		cmd = exec.Command(e.Cmd, filePath)
	}

	if cmd != nil {
		cmd.Start()
	}
}

func (m TemplatesModel) updatePreview(msg tea.KeyMsg) (TemplatesModel, tea.Cmd) {
	switch msg.String() {
	case "left", "esc", "h":
		// Exit preview mode
		m.previewing = false
		m.previewName = ""
		m.scroll = 0
	case "up", "k":
		if m.scroll > 0 {
			m.scroll--
		}
	case "down", "j":
		m.scroll++
	case "g":
		// Open template on GitHub
		m.openOnGitHub()
	}
	return m, nil
}

func (m *TemplatesModel) openOnGitHub() {
	if m.cursor >= len(m.filtered) {
		return
	}
	name := m.filtered[m.cursor].Name
	url := "https://github.com/" + templates.GitHubRepo + "/blob/main/" + templates.TemplatesPath + "/" + name + ".md"

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}

func (m *TemplatesModel) applyFilter() {
	query := strings.ToLower(m.filterInput.Value())
	m.filtered = []TemplateItem{}

	for _, t := range m.templates {
		// Apply source filter
		if m.sourceMode == "default" && t.Source != "default" {
			continue
		}
		if m.sourceMode == "custom" && t.Source != "custom" {
			continue
		}

		// Apply text filter
		if query != "" {
			if !strings.Contains(strings.ToLower(t.Name), query) &&
				!strings.Contains(strings.ToLower(t.Description), query) {
				continue
			}
		}

		m.filtered = append(m.filtered, t)
	}

	m.cursor = 0
	m.scroll = 0
	m.applySort()
}

func (m *TemplatesModel) applySort() {
	switch m.sortBy {
	case "name-asc":
		sort.Slice(m.filtered, func(i, j int) bool {
			return strings.ToLower(m.filtered[i].Name) < strings.ToLower(m.filtered[j].Name)
		})
	case "name-desc":
		sort.Slice(m.filtered, func(i, j int) bool {
			return strings.ToLower(m.filtered[i].Name) > strings.ToLower(m.filtered[j].Name)
		})
	}
}

// View renders the templates view
func (m TemplatesModel) View() string {
	if m.copying {
		return m.viewCopying()
	}
	if m.editing {
		return m.viewEditing()
	}
	if m.previewing {
		return m.viewPreview()
	}
	return m.viewList()
}

func (m TemplatesModel) viewCopying() string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().MarginLeft(2)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	accentStyle := lipgloss.NewStyle().Foreground(theme.Accent)

	source := m.filtered[m.copySourceIdx]

	b.WriteString("\n")
	b.WriteString(labelStyle.Render("Copy template as custom"))
	b.WriteString("\n\n")

	b.WriteString(hintStyle.Render("Source: ") + accentStyle.Render(source.Name))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Template name:"))
	b.WriteString("\n")
	b.WriteString("  " + m.copyInput.View())
	b.WriteString("\n\n")

	b.WriteString(hintStyle.Render("Enter to create, Esc to cancel"))
	b.WriteString("\n")

	return b.String()
}

func (m TemplatesModel) viewEditing() string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().MarginLeft(2)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	accentStyle := lipgloss.NewStyle().Foreground(theme.Accent)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	t := m.filtered[m.cursor]

	b.WriteString("\n")
	b.WriteString(labelStyle.Render("Edit custom template"))
	b.WriteString("\n\n")

	b.WriteString(hintStyle.Render("Template: ") + accentStyle.Render(t.Name))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Open in:"))
	b.WriteString("\n\n")

	if len(m.editors) == 0 {
		b.WriteString(hintStyle.Render("No editors detected"))
		b.WriteString("\n")
	} else {
		for _, e := range m.editors {
			b.WriteString("  " + keyStyle.Render("["+e.Key+"]") + " " + e.Name)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(hintStyle.Render("Esc to cancel"))
	b.WriteString("\n")

	return b.String()
}

func (m TemplatesModel) viewList() string {
	var b strings.Builder

	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	accentStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	normalStyle := lipgloss.NewStyle()
	descStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
		b.WriteString(errStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
		return b.String()
	}

	if !m.loaded {
		return b.String()
	}

	// Filter input (always visible, always focused)
	b.WriteString("\n")
	b.WriteString("  " + m.filterInput.View())

	// Source filter indicator
	var sourceLabel string
	switch m.sourceMode {
	case "all":
		sourceLabel = "All"
	case "default":
		sourceLabel = "Default"
	case "custom":
		sourceLabel = "Custom"
	}
	b.WriteString("  " + mutedStyle.Render("view:"+sourceLabel))

	// Sort indicator
	var sortLabel string
	if m.sortBy == "name-asc" {
		sortLabel = "A-Z"
	} else {
		sortLabel = "Z-A"
	}
	b.WriteString("  " + mutedStyle.Render("s:"+sortLabel))
	b.WriteString("\n\n")

	// Show any message (success or hint)
	if m.message != "" {
		msgStyle := lipgloss.NewStyle().Foreground(theme.Success).MarginLeft(2)
		if m.messageErr {
			msgStyle = lipgloss.NewStyle().Foreground(theme.Error).MarginLeft(2)
		}
		b.WriteString(msgStyle.Render(m.message))
		b.WriteString("\n\n")
	}

	if len(m.filtered) == 0 {
		if len(m.templates) == 0 {
			b.WriteString(mutedStyle.Render("  No templates found"))
		} else {
			b.WriteString(mutedStyle.Render("  No matches for \"" + m.filterInput.Value() + "\""))
		}
		b.WriteString("\n")
		return b.String()
	}

	// Calculate column widths for two-column layout
	nameColWidth := 22 // Fixed width for name column
	descColWidth := m.width - nameColWidth - 8

	cursorOn := accentStyle.Render(">")
	cursorOff := " "

	// Show visible templates with scrolling
	endIdx := m.scroll + templatesVisibleItems
	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}

	for i := m.scroll; i < endIdx; i++ {
		t := m.filtered[i]
		cursor := cursorOff
		nameStyle := normalStyle
		descStyleLocal := descStyle

		if i == m.cursor {
			cursor = cursorOn
			nameStyle = selectedStyle
			descStyleLocal = mutedStyle // Keep description muted even when selected
		}

		// Truncate long names
		displayName := t.Name
		if len(displayName) > nameColWidth {
			displayName = displayName[:nameColWidth-3] + "..."
		}

		// Pad name to align descriptions
		namePadded := displayName + strings.Repeat(" ", nameColWidth-len(displayName))

		// Truncate long descriptions
		displayDesc := t.Description
		if len(displayDesc) > descColWidth {
			displayDesc = displayDesc[:descColWidth-3] + "..."
		}

		b.WriteString("  " + cursor + " " + nameStyle.Render(namePadded) + " " + descStyleLocal.Render(displayDesc))
		b.WriteString("\n")
	}

	// Scroll indicator
	if len(m.filtered) > templatesVisibleItems {
		showing := m.scroll + 1
		showingEnd := endIdx
		total := len(m.filtered)
		indicator := mutedStyle.Render("    " + itoa(showing) + "-" + itoa(showingEnd) + " of " + itoa(total))
		b.WriteString(indicator)
		b.WriteString("\n")
	}

	return b.String()
}

func (m TemplatesModel) viewPreview() string {
	var b strings.Builder

	if m.cursor >= len(m.filtered) {
		return b.String()
	}

	t := m.filtered[m.cursor]

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

	return b.String()
}
