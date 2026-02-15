# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

This is a template repository for the **Idempotent Research Loop (IRL)** - a document-centric workflow pattern for AI-assisted research. Instead of chat-based interaction, IRL uses a **plan document** (`main-plan.md`) as the control surface. The AI reads the plan, executes tasks, and produces versioned artifacts.

## CLI Agent Specification

The `irl` CLI is fully non-interactive when flags/args are provided. Every command below runs without prompts, making it safe for agent use.

### Discovery

```bash
# List all projects in workspace (JSON for parsing)
irl list --json
# Returns: {"projects":[{"name":"...","path":"...","modified":"..."},...]}
# Always returns exit 0 with {"projects":[]} on empty sets

# Read configuration
irl config --json
# Returns: {"default_directory":"...","profile":{...},"favorite_editors":[...]}

# Read profile
irl profile --json
# Returns: {"name":"...","title":"...","institution":"...","department":"...","email":"...","instructions":"..."}

# List templates
irl templates
irl templates show <name>   # Print raw template content to stdout
```

### Project Creation

```bash
# Non-interactive: provide purpose as arg, defaults to irl-basic template
irl init "ERP correlation analysis"
# → Creates YYMMDD-slug directory with plan, .gitignore, git init
# → Injects profile (author/affiliation) into plan front matter

# With specific template
irl init "my study" -t irl-basic

# With exact name (skip auto-naming)
irl init -n my-exact-name -t irl-basic

# In specific directory
irl init "purpose" -d ~/Research

# Adopt existing folder into workspace
irl adopt ~/Downloads/my-data
irl adopt ~/Downloads/my-data --rename -t irl-basic
```

**Important for agents:**
- When purpose is provided as an argument (non-interactive mode), the template defaults to `irl-basic` if `-t` is not specified. No prompts will appear.
- **Default directory resolution** (non-interactive): `--dir` flag > `config.default_directory` > current working directory. If no default is configured and no `--dir` is passed, the project is created in `cwd`.
- `irl list` and `irl adopt` require a configured default directory and will error if none is set. Run `irl config --dir ~/path` first.

### Configuration

```bash
# Set workspace directory
irl config --dir ~/Research

# Set preferred editor
irl config --editor cursor

# Set profile (merged with existing — only specified fields change)
irl profile --name "Jane Doe" --title "MD" --institution "UCSF"
irl profile --department "Neurology" --email "jane@ucsf.edu"
irl profile --instructions "Always cite sources in APA format"

# Clear profile
irl profile --clear
```

### Template Management

```bash
# List all templates (built-in + custom)
irl templates

# Show template content (pipe to file, parse, etc.)
irl templates show irl-basic

# Create custom template (stored in workspace/_templates/<name>/main-plan.md)
irl templates create my-template                    # Copy from irl-basic
irl templates create my-template --from irl-basic   # Explicit source

# Delete custom template
irl templates delete my-template
```

### Open Projects

```bash
irl open my-project                # Uses configured editor
irl open my-project --editor code  # Specific editor
```

### Environment

```bash
irl doctor    # Check tools and environment
irl update    # Refresh templates from GitHub
```

## Project Structure (created by `irl init`)

```
my-project/
├── .gitignore
├── plans/
│   └── main-plan.md    ← Control surface (with profile front matter)
```

Additional directories (`02-data/`, `03-outputs/`, `04-logs/`) are defined in the plan and created by the AI on first run.

## Profile Injection

When a profile is configured, `irl init` and `irl adopt` automatically prepend YAML front matter to the plan:

```yaml
---
author: Jane Doe, MD
affiliation: UCSF, Neurology
email: jane@ucsf.edu
---

<!-- AI Instructions:
Always cite sources in APA format
-->

# IRL Basic Template
...
```

## CLI Development

The `irl` CLI is written in Go:
```
cmd/           # Cobra commands (init, adopt, list, open, profile, config, templates, doctor, update)
pkg/naming/    # YYMMDD-slug generation
pkg/scaffold/  # Project structure creation + profile injection
pkg/templates/ # GitHub template fetching + caching
pkg/projects/  # Project scanning (shared between CLI and TUI)
pkg/config/    # Configuration + profile persistence (~/.irl/config.json)
pkg/editor/    # Editor detection and launching
pkg/doctor/    # Environment checks
pkg/theme/     # Terminal styling
internal/tui/  # Bubble Tea TUI (launched with bare `irl` command)
```

**Build commands (always use make):**
```bash
make build          # Build ./irl binary
make test           # Build and run quick tests
make clean          # Remove all build artifacts
```

**Never use `go build .`** - it creates `irl-template` (from module name) instead of `irl`.
A pre-commit hook prevents accidentally committing binaries.

## IRL Workflow Pattern

1. **Edit the plan** (`main-plan.md`)
2. **AI executes** tasks according to the plan
3. **Outputs rendered**, logs updated
4. **Review diffs**, accept or revise
5. **Commit** after each iteration

## Key Conventions

- **Do not edit the plan file unless explicitly permitted**
- **Make surgical git commits** before and after edits
- **Update logs** (`logs/activity.md`) after iterations
