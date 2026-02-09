# Enhancements From irloop.org Docs (Draft Text + Placement Notes)

Source docs reviewed:
- `https://www.irloop.org/` (explainer / tutorial landing page)
- `https://irloop.org/cli/doctor` (doctor command and install guidance)

This file is intentionally written as "drop-in" blocks you can paste into the manuscript, plus a suggested placement for each.

---

## A. Plain-Language Orientation (Suggested: end of Abstract or start of Introduction)

**Rationale:** The manuscript is already strong on computational reproducibility and provenance, but the public docs add a clearer *reader-oriented* explanation: what IRL is for, why it is worth learning, and what "idempotent" means in non-technical language.

### Draft paragraph: "AI work you can retrace"

Many practitioners now use AI assistants to draft text, write code, and summarize literature. The challenge is that chat-first workflows often leave little durable structure: useful work evaporates into conversational transcripts, and weeks later it is difficult to reconstruct inputs, assumptions, or exactly what changed. IRL reframes AI assistance around reviewable files. A single plan document functions as the recipe: it encodes inputs, assumptions, acceptance criteria, and the outputs to be produced. The assistant executes against this plan and writes artifacts that can be inspected, shared, and rerun.

### Draft paragraph: "Idempotent" in plain language

The term idempotent sounds technical, but the core idea is simple. If you follow the same recipe with the same ingredients, you should get the same dish. In IRL, rerunning a step should be safe: when inputs and instructions have not changed, authoritative outputs should not change either (or should change only within explicitly declared tolerances, such as timestamps in logs). This makes iteration trustworthy and makes deviations easy to locate and explain.

### Draft paragraph: "Expertise, not model features"

IRL is not a claim about model intelligence or algorithmic novelty. It is a workflow that makes human judgment explicit. The human specifies what matters and evaluates outputs; the assistant performs bounded execution steps. This division of labor reduces overreliance on conversational memory and makes verification a first-class part of the loop, similar to supervising a capable human assistant whose work must be reviewed.

---

## B. A Short "Getting Started" Sidebar (Suggested: end of Section 3 or start of Methods/Case Study)

**Rationale:** The docs provide a concrete operational sequence that is easy for new users: plan file, repeatable command, outputs, checkpoint, repeat.

### Draft boxed sidebar text: "The IRL rhythm"

An IRL iteration can be summarized as: (1) edit the plan; (2) execute the plan; (3) review the outputs; (4) save a checkpoint; (5) repeat. The plan replaces conversational history as the source of truth. Outputs replace explanations as evidence. Checkpoints (e.g., version control commits plus activity logs) provide a durable audit trail across iterations.

---

## C. CLI Environment Checking Appendix (Suggested: new Appendix section)

**Rationale:** The manuscript discusses pre/post execution hooks and tooling assumptions, but the docs include a concrete command (`irl doctor`) and "what to do if something is missing" guidance with install commands.

### Draft appendix: "Environment checks"

IRL includes an optional environment check that inventories tool availability and provides recommendations. The command `irl doctor` reports system information and groups dependencies into categories (Core Tools, AI Assistants, IDEs, Sandbox). When a dependency is missing, users can install it using standard system installers (e.g., Git, Quarto, Python/R). Importantly, most tools are optional: the minimal requirement is version control plus at least one execution environment and one assistant interface.

### Recommended additions to Appendix (bullet form)

- Core tools:
  - Git (required for versioning/audit)
  - One analysis runtime (R and/or Python)
  - One publishing tool (e.g., Quarto) if rendering reports
- Optional:
  - Terminal-based assistants (Claude Code / Codex / etc.)
  - Docker for sandboxed, reproducible execution
- Recommendation text:
  - If something is missing, proceed with what you have; install incrementally.

---

## D. Where This Helps The Manuscript (Quick Mapping)

- Section 1 (Introduction): add the "AI work you can retrace" framing to clarify the user problem.
- Section 3 (IRL): add the plain-language idempotency paragraph as a complement to the formal/practical definition.
- Case study / Implementation: add the "IRL rhythm" sidebar to make the procedure scannable.
- Appendix: add the `irl doctor` environment-check description and dependency grouping.

