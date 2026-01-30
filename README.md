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
├── main-plan.md    ← Your control surface
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

| Command | Description |
|---------|-------------|
| `irl init "purpose"` | Create project with auto-naming |
| `irl init` | Interactive mode |
| `irl init -t template` | Use specific template |
| `irl init -d ~/path` | Override directory |
| `irl config` | View configuration |
| `irl config --dir ~/path` | Set default directory |
| `irl templates` | List available templates |
| `irl update` | Refresh templates from GitHub |
| `irl doctor` | Check environment and tools |

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
