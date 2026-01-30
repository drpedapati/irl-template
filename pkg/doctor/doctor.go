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
	parts := []string{s.Platform, fmt.Sprintf("%d cores", s.Cores)}
	if s.Memory != "" {
		parts = append(parts, s.Memory)
	}
	if s.Disk != "" {
		parts = append(parts, s.Disk)
	}
	return strings.Join(parts, " · ")
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
		{Name: "Aider", Cmd: "aider", Install: "pip install aider-chat", Category: "AI Assistants"},
		{Name: "Copilot CLI", Cmd: "gh copilot", Install: "gh extension install github/gh-copilot", Category: "AI Assistants"},
		{Name: "Ollama", Cmd: "ollama", Install: "brew install ollama", Category: "AI Assistants"},
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
	case "gh copilot":
		return checkGhCopilot()
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

func checkGhCopilot() bool {
	out, err := exec.Command("gh", "extension", "list").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "copilot")
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
