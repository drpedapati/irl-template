# IRL: Idempotent Research Loop

A CLI tool and workflow for reproducible AI-assisted research.

## Install

```bash
brew tap drpedapati/tap
brew install irl
```

## Quick Start

```bash
# Set your projects directory (one time)
irl config --dir ~/Research

# Set your profile (one time)
irl profile --name "Jane Doe" --institution "UCSF" --title "MD"

# Create a project
irl init "ERP correlation analysis"
# → Creates ~/Research/260129-erp-correlation-analysis

# Check your environment
irl doctor
```

## What is IRL?

IRL is a document-centric workflow for working with AI assistants. Instead of chat conversations that drift and get lost, you maintain a **plan document** that serves as both your instructions and your record of intent.

```
my-project/
├── plans/
│   └── main-plan.md    ← Your control surface
├── 02-data/
│   ├── raw/
│   └── derived/
├── 03-outputs/
└── 04-logs/
```

**The workflow:**
1. Edit `main-plan.md` with your objectives
2. AI reads the plan and executes tasks
3. Review outputs and diffs
4. Commit, revise, repeat

Everything is versioned. Diffs become your audit trail.

## Commands

### Projects

| Command | Description |
|---------|-------------|
| `irl init "purpose"` | Create project with auto-naming (YYMMDD-slug) |
| `irl init` | Interactive mode with directory browser |
| `irl init -t template` | Use specific template |
| `irl init -n name` | Use exact project name |
| `irl init -d ~/path` | Override workspace directory |
| `irl adopt ~/folder` | Copy existing folder into workspace |
| `irl adopt ~/folder --rename` | Adopt with YYMMDD prefix |
| `irl list` | List all projects (table) |
| `irl list --json` | List projects as JSON |
| `irl open my-project` | Open project in preferred editor |
| `irl open my-project --editor code` | Open in specific editor |

### Templates

| Command | Description |
|---------|-------------|
| `irl templates` | List all templates (built-in + custom) |
| `irl templates show <name>` | Print template content |
| `irl templates create <name>` | Create custom template from irl-basic |
| `irl templates create <name> --from X` | Create from another template |
| `irl templates delete <name>` | Delete a custom template |
| `irl update` | Refresh built-in templates from GitHub |

### Configuration

| Command | Description |
|---------|-------------|
| `irl config` | View current configuration |
| `irl config --json` | Configuration as JSON |
| `irl config --dir ~/path` | Set default workspace directory |
| `irl config --editor cursor` | Set preferred editor |
| `irl profile` | View current profile |
| `irl profile --json` | Profile as JSON |
| `irl profile --name "..." --institution "..."` | Set profile fields |
| `irl profile --clear` | Clear all profile fields |
| `irl doctor` | Check environment and tools |

### TUI (Terminal UI)

Run `irl` with no arguments to launch the interactive terminal UI, which provides all the above capabilities plus a project browser, editor configuration, and visual template management.

## Agent Usage

The CLI is designed for both humans and AI agents. Every TUI capability has a non-interactive CLI equivalent. See `CLAUDE.md` for agent-specific patterns, or use `--json` flags for machine-readable output:

```bash
irl list --json          # {"projects":[...]} — always exit 0
irl config --json        # Full config object
irl profile --json       # Profile fields
irl templates show X     # Raw template content to stdout
irl init "purpose"       # Create project (no prompts when args provided)
```

**Default directory resolution** (non-interactive `init`): `--dir` flag > configured default > current working directory. `list` and `adopt` require a configured default directory.

## Why IRL?

| Problem | IRL Solution |
|---------|--------------|
| Chat histories get lost | Plan document is your source of truth |
| Hard to reproduce results | Everything versioned in git |
| Unclear what changed | Diffs show exactly what happened |
| Who decided what? | Human owns plan; AI executes tasks |

## Associated Manuscript

**"The Idempotent Research Loop (IRL): A Document-Centric Framework for AI-Assisted Scientific Analysis"**

The manuscript provides theoretical foundation and case studies.

## License

MIT
