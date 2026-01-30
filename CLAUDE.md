# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

This is a template repository for the **Idempotent Research Loop (IRL)** - a document-centric workflow pattern for AI-assisted research. Instead of chat-based interaction, IRL uses a **plan document** (`main-plan.md`) as the control surface. The AI reads the plan, executes tasks, and produces versioned artifacts.

## Creating IRL Projects

### Via CLI (for users)
```bash
# Install (via Homebrew, coming soon)
brew install drpedapati/tap/irl

# Create project - auto-names as YYMMDD-slug
irl init "ERP correlation analysis"    # → 260129-erp-correlation-analysis/
irl init "APA poster" -t meeting-abstract
irl init                               # Interactive mode
```

### Via Agent (when user asks)
When a user asks to create an IRL project, run:
```bash
./irl init "PURPOSE" -t TEMPLATE
```
Or build and use: `go build -o irl . && ./irl init "PURPOSE"`

### Via Makefile (legacy)
```bash
make irl                    # Interactive
make irl my-project         # Named project
```

## Directory Structure

- `main-plan.md` - The control surface
- `01-plans/templates/` - Reusable templates (repo source)
- `02-data/` - Data files (raw/, derived/)
- `03-outputs/` - Rendered outputs
- `04-logs/` - Activity logs

## CLI Development

The `irl` CLI is written in Go:
```
cmd/           # Cobra commands
pkg/naming/    # YYMMDD-slug generation
pkg/scaffold/  # Project structure creation
pkg/templates/ # GitHub template fetching + caching
```

Build: `go build -o irl .`
Test: `./irl init "test project" && rm -rf 260129-test-project`

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
