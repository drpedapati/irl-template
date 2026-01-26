# IRL Template: A Practical Framework for AI-Assisted Research

This repository is a template implementation of the **Idempotent Research Loop (IRL)**, a workflow pattern designed to help researchers use AI coding assistants (like terminal-based agents) in ways that are reproducible, auditable, and aligned with scientific best practices.

## What is IRL?

IRL is a document-centric approach to working with AI assistants. Instead of relying on chat conversations that can drift and be hard to track, IRL uses a **plan document** (`plans/main-plan.md`) as your control surface. You edit the plan to specify what you want done; the AI reads it, executes tasks, and produces structured outputs that are versioned and reviewable.

Think of it like this: rather than having a conversation where context can get lost, you maintain a living document that serves as both your instructions and your record of intent. The AI acts as an executor that follows your plan and produces artifacts you can review, rerun, and share.

## Why Does This Matter?

When using AI assistants for research, several challenges emerge:

- **Conversational drift**: Chat histories become hard to navigate, and it's unclear what the "authoritative" version of your analysis is
- **Reproducibility**: Without clear artifacts and version control, it's difficult to rerun analyses or understand what changed
- **Auditability**: Reviewers (or future you) need to see what decisions were made and why
- **Accountability**: It should be clear what the human decided versus what the AI generated

IRL addresses these by shifting from chat-based interaction to artifact-based workflows. Everything becomes a file in a versioned repository, making it easier to track changes, rerun analyses, and maintain clear provenance.

## Associated Manuscript

This template repository accompanies the manuscript:

**"The Idempotent Research Loop (IRL): A Document-Centric Framework for AI-Assisted Scientific Analysis"**

The full manuscript is available at:
`/Users/ernie/Dropbox/DEEPPROJECTS/Idempotent Research Loop/IRL_Patterns_Download_Package/IRL_Patterns_Manuscript_Only.pdf`

The manuscript provides the theoretical foundation, design rationale, and detailed case studies. This repository provides a practical starting point you can adapt for your own work.

## Core Concepts

### 1. The Plan Document (`plans/main-plan.md`)

This is your control surface. You edit it to specify:
- What you want to accomplish
- Constraints and requirements
- Acceptance criteria
- Decisions and rationale

The AI reads this document but doesn't modify it (unless explicitly permitted). This keeps your intent stable and reviewable.

### 2. Paired Artifacts

IRL separates reasoning from execution:
- **Plan document** (`plans/main-plan.md`): Human-editable intent and constraints
- **Executed analysis** (`executed.qmd` or similar): The actual code and analysis that produces results

This separation means you can adjust your plan without diving into code, and the executed document can be regenerated and rerun.

### 3. Versioned Artifacts

Everything is tracked in version control (Git). Each iteration produces:
- Executable analysis documents
- Rendered reports (HTML, PDF, Word)
- Logs of activity and decisions
- Commits that bracket each iteration

This makes diffs your audit interface: you can see exactly what changed and when.

### 4. Practical Idempotency

In computer science, "idempotent" means doing something multiple times has the same effect as doing it once. For AI-assisted workflows, perfect idempotency isn't always possible (LLMs can be stochastic), but IRL aims for **practical idempotency**: rerunning the same plan with the same inputs should produce substantially the same results.

If it doesn't, that's a signal to investigate: pin random seeds, tighten constraints, or document acceptable tolerances.

### 5. Managerial Human-in-the-Loop

IRL uses a managerial model: you (the human) specify intent and retain responsibility for conclusions. The AI performs bounded execution tasks under your plan. This aligns with best practices for human-AI collaboration: clear oversight, error recovery, and appropriate expectation-setting.

## Repository Structure

A minimal IRL repository looks like this:

```
project/
├── plans/
│   └── main-plan.md           # Your control surface: intent, constraints, acceptance criteria
├── executed.qmd               # Executable analysis document (or .py, .R, etc.)
├── data/
│   ├── raw/                   # Original data files
│   └── derived/               # Processed data
├── outputs/
│   ├── report.html            # Rendered outputs
│   ├── report.pdf
│   └── figures/               # Generated figures
├── logs/
│   ├── activity_log.md        # What happened in each iteration
│   └── decision_log.csv       # Decisions and rationale
├── env/
│   └── environment.yml        # Environment specification (Conda, renv, etc.)
└── README.md                  # This file
```

This structure is intentionally generic. You can adapt it to your tools and needs. The essential requirements are:
- A stable plan document
- Rerunnable executed artifacts
- Durable rendered outputs
- Versioned history with logs

## How to Use This Template

1. **Clone or fork this repository** to start your own IRL project

2. **Read the associated manuscript** to understand the full framework and design rationale

3. **Customize the structure** for your domain and tools:
   - Replace `executed.qmd` with your preferred executable format (Python scripts, R markdown, Jupyter notebooks, etc.)
   - Adjust the directory structure as needed
   - Set up your environment specification

4. **Start with `plans/main-plan.md`**: Write your first plan specifying what you want to accomplish, any constraints, and how you'll know when it's done. See `templates/irl-basic-template.md` for a template.

5. **Run iterations**: Each iteration follows this pattern:
   - Edit `plans/main-plan.md` with your objectives
   - AI reads the plan and executes tasks
   - Outputs are rendered and logs are updated
   - Review the results and diffs
   - Accept, revise the plan, or rerun

6. **Maintain provenance**: Commit after each iteration with clear messages linked to your plan objectives

## Key Benefits

- **Auditability**: Everything is versioned and diffable. Reviewers can see exactly what changed and why.
- **Reproducibility**: Executable documents can be rerun in controlled environments.
- **Reduced hidden state**: Intent, code, and outputs are explicit artifacts, not buried in chat history.
- **Clear accountability**: The human owns intent and acceptance; the AI performs bounded tasks.
- **Lab-scale standardization**: The plan file functions like a living protocol that new contributors can read and follow.

## Limitations and Considerations

- **LLM stochasticity**: IRL doesn't eliminate nondeterminism, but it contains it by forcing changes into reviewable diffs
- **External dependencies**: Workflows relying on changing APIs or databases may not be stable across reruns
- **Overhead**: The loop imposes structure (plan upkeep, logging, commits). For very small tasks, this may be unnecessarily heavyweight
- **Not a substitute for expertise**: IRL organizes work; it doesn't validate scientific correctness. You remain responsible for methodological decisions

## Getting Help

- Review the associated manuscript for detailed explanations and case studies
- Check the `logs/` directory for examples of activity and decision tracking
- Adapt the structure to your needs—IRL is a pattern, not a rigid prescription

## Contributing

This is a template repository. Feel free to adapt it for your own work. If you develop improvements or variations that might benefit others, consider sharing them.

## License

[Specify your license here]

---

**Note**: IRL is not a tool or agent—it's a workflow pattern that can be implemented with existing research infrastructure. The goal is to make AI assistance more compatible with scientific expectations for reproducibility, auditability, and accountable collaboration.
