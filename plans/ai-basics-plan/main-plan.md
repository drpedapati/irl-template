# AI Basics for Physicians

## 🔧 First Time Setup — Run once when starting a new project
<!-- 👤 AUTHOR AREA: Define setup steps below -->

### Common Skill Library — Pre-installed tools available to all projects
<!-- Uncomment to use -->
Install Quarto Skill: https://github.com/posit-dev/skills/tree/main/quarto/authoring
Install https://github.com/henrybloomingdale/skills-monorepo/blob/main/humanize-text/SKILL.md skill

## ✅ Before Each Loop — Checklist to run before every iteration

- Script designed to be **idempotent**
- **Version control**: clean tree (`git status`), commit baseline, `.gitignore` enforced
- The **only permitted modification** is within the `## One-Time Instructions` section.
- New agents: review existing outputs before recreating; don’t overwrite unless explicitly instructed. Seek clarification if unsure.
---

## 🔁 Instruction Loop — Define the work for each iteration

<!-- 👤 AUTHOR AREA: Define each loop's work below -->

Plan for for revising 03-outputs/irl-explainer-site/ai-basics.html and related pages

1. Install and apply humanize text
2. generate a paired markdown shadow outline for this html. goal is for me to edit this outline and it can generate downstream html changes.
3. On each loop check outline and reasoning plan and execute changes

Revisions to outline:
1. Lede
  after "market shifts" add a new paragraph; current content is not great
  move to physicians who get most out of AI are those understand the usefulness as a process tool; strengths and weaknesses; understanding hybrid cognition where human excel and AI excel.

### One-Time Instructions — Tasks that should only execute once

<!-- 👤 AUTHOR AREA: Add one-time tasks below -->

### Formatting Guidelines — Rules for output style and structure

<!-- 👤 AUTHOR AREA: Add formatting rules below -->

---

## 📝 After Each Loop — Steps to complete after every iteration

> Uncomment your preferred options below.

- **Update activity log**
  - Use timestamps **only** when sequencing or causality matters
  - In `plans/main-plan-activity.md`, write 1–2 lines describing:
    - What you did
    - Timestamp
    - Git hash
  <!-- Optional CSV logging -->
  <!--
  - In `plans/main-plan-activity.csv`, add 1 row:
    - What you did
    - Timestamp
    - Git hash
  -->

- **Update plan log**: Update `plans/main-plan-log.csv`

- **Version control**: Commit intended changes only; verify no ignored or unintended files staged.

- **Give feedback to the AUTHOR** — concise and actionable:
  1. What was done, decisions needed, next steps
  2. Identify anything breaking idempotency, or obsolete/outdated instructions
  3. Identify critical reasoning errors

## 📚 Skill Library — Optional community skills to install per project
<!-- Uncomment to use -->

<!-- Install PPTX Posters -->
<!-- https://github.com/K-Dense-AI/claude-scientific-skills/tree/main/scientific-skills/pptx-posters -->

<!-- Install Scientific Writing Skill -->
<!-- https://github.com/K-Dense-AI/claude-scientific-skills/tree/main/scientific-skills/scientific-writing -->

<!-- Install BioRx Search -->
<!-- https://github.com/K-Dense-AI/claude-scientific-skills/tree/main/scientific-skills/biorxiv-database -->

<!-- Install PubMed Search -->
<!-- https://github.com/K-Dense-AI/claude-scientific-skills/tree/main/scientific-skills/pubmed-database -->

<!-- Install Flowcharts -->
<!-- https://github.com/lukilabs/beautiful-mermaid -->

<!-- Install PowerPoint -->
<!-- https://github.com/anthropics/skills/tree/main/skills/pptx -->

<!-- Install PDF -->
<!-- https://github.com/anthropics/skills/tree/main/skills/pdf -->
