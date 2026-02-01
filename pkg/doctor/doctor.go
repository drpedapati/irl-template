// Package doctor provides environment checking functionality.
package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Tool represents a tool to check
type Tool struct {
	Name     string
	Cmd      string
	Install  string
	Category string
}

// ToolResult represents the result of checking a tool
type ToolResult struct {
	Tool  Tool
	Found bool
}

// SystemInfo holds system information
type SystemInfo struct {
	Platform string
	Cores    int
	Memory   string
	Disk     string
}

// GetSystemInfo returns system information
func GetSystemInfo() SystemInfo {
	info := SystemInfo{
		Platform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Cores:    runtime.NumCPU(),
	}

	// Memory
	if runtime.GOOS == "darwin" {
		if out, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
			memBytes := strings.TrimSpace(string(out))
			if len(memBytes) > 9 {
				info.Memory = memBytes[:len(memBytes)-9] + " GB"
			}
		}
	} else if runtime.GOOS == "linux" {
		if out, err := exec.Command("free", "-g").Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) > 1 {
				fields := strings.Fields(lines[1])
				if len(fields) > 1 {
					info.Memory = fields[1] + " GB"
				}
			}
		}
	}

	// Disk
	if out, err := exec.Command("df", "-h", ".").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 4 {
				info.Disk = fields[3] + " free"
			}
		}
	}

	return info
}

// FormatSystemInfo returns a formatted string of system info
func (s SystemInfo) String() string {
	return s.Platform
}

// CoreTools returns the list of core tools to check
func CoreTools() []Tool {
	return []Tool{
		{Name: "Git", Cmd: "git", Install: "brew install git", Category: "Core Tools"},
		{Name: "Quarto", Cmd: "quarto", Install: "brew install --cask quarto", Category: "Core Tools"},
		{Name: "R", Cmd: "R", Install: "brew install r", Category: "Core Tools"},
		{Name: "Python", Cmd: "python3", Install: "brew install python", Category: "Core Tools"},
	}
}

// AITools returns the list of AI tools to check
func AITools() []Tool {
	return []Tool{
		{Name: "Claude Code", Cmd: "claude", Install: "npm i -g @anthropic-ai/claude-code", Category: "AI Assistants"},
		{Name: "Codex", Cmd: "codex", Install: "npm i -g @openai/codex", Category: "AI Assistants"},
		{Name: "Copilot", Cmd: "copilot", Install: "brew install --cask github-copilot", Category: "AI Assistants"},
	}
}

// IDETools returns the list of IDE tools to check
func IDETools() []Tool {
	return []Tool{
		{Name: "Positron", Cmd: "positron", Install: "brew install --cask positron", Category: "IDEs"},
		{Name: "VS Code", Cmd: "code", Install: "brew install --cask visual-studio-code", Category: "IDEs"},
		{Name: "Cursor", Cmd: "cursor", Install: "brew install --cask cursor", Category: "IDEs"},
		{Name: "RStudio", Cmd: "rstudio", Install: "brew install --cask rstudio", Category: "IDEs"},
	}
}

// SandboxTools returns the list of sandbox tools to check
func SandboxTools() []Tool {
	return []Tool{
		{Name: "Docker", Cmd: "docker", Install: "brew install --cask docker", Category: "Sandbox"},
	}
}

// AllTools returns all tools grouped by category
func AllTools() []Tool {
	var all []Tool
	all = append(all, CoreTools()...)
	all = append(all, AITools()...)
	all = append(all, IDETools()...)
	all = append(all, SandboxTools()...)
	return all
}

// CheckTool checks if a tool is available
func CheckTool(t Tool) bool {
	switch t.Cmd {
	case "positron":
		return checkApp("Positron")
	case "cursor":
		return checkApp("Cursor")
	case "rstudio":
		return checkApp("RStudio")
	default:
		return checkCmd(t.Cmd)
	}
}

// CheckAllTools checks all tools and returns results
func CheckAllTools() []ToolResult {
	tools := AllTools()
	results := make([]ToolResult, len(tools))
	for i, t := range tools {
		results[i] = ToolResult{
			Tool:  t,
			Found: CheckTool(t),
		}
	}
	return results
}

func checkCmd(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func checkApp(name string) bool {
	if runtime.GOOS != "darwin" {
		return checkCmd(strings.ToLower(name))
	}
	paths := []string{
		filepath.Join("/Applications", name+".app"),
		filepath.Join(os.Getenv("HOME"), "Applications", name+".app"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

// HasDocker returns true if Docker is available
func HasDocker() bool {
	return checkCmd("docker")
}

// PlanEditorResult represents the status of a plan editor
type PlanEditorResult struct {
	Name      string
	Command   string
	Available bool
	IsDefault bool
}

// PlanEditorTools returns the list of terminal editors to check
func PlanEditorTerminalTools() []Tool {
	return []Tool{
		{Name: "nano", Cmd: "nano", Category: "Plan Editors (Terminal)"},
		{Name: "vim", Cmd: "vim", Category: "Plan Editors (Terminal)"},
		{Name: "vi", Cmd: "vi", Category: "Plan Editors (Terminal)"},
	}
}

// PlanEditorGUITools returns the list of GUI editors to check
func PlanEditorGUITools() []Tool {
	return []Tool{
		{Name: "VS Code", Cmd: "code", Category: "Plan Editors (GUI)"},
		{Name: "Cursor", Cmd: "cursor", Category: "Plan Editors (GUI)"},
		{Name: "Zed", Cmd: "zed", Category: "Plan Editors (GUI)"},
	}
}

// CheckPlanEditors checks all plan editors and returns results
func CheckPlanEditors() (terminal []ToolResult, gui []ToolResult) {
	for _, t := range PlanEditorTerminalTools() {
		terminal = append(terminal, ToolResult{
			Tool:  t,
			Found: CheckTool(t),
		})
	}
	for _, t := range PlanEditorGUITools() {
		gui = append(gui, ToolResult{
			Tool:  t,
			Found: CheckTool(t),
		})
	}
	return terminal, gui
}
