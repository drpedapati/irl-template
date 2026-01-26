# Basic IRL Plan Template

## Objective

[Describe what you want to accomplish in this iteration. Be specific about the goal.]

Example: "Generate a summary report of the data in `data/raw/sample.csv`, including descriptive statistics and one visualization."

## Constraints

[Specify any limitations, requirements, or boundaries for this work.]

- Do not modify this plan document unless explicitly requested
- Use only data from `data/raw/` directory
- All outputs should be written to `outputs/` directory
- Code should be documented and executable

## Acceptance Criteria

[Define how you'll know the iteration is complete and successful.]

- [ ] Executed analysis document (`executed.qmd` or equivalent) exists and runs without errors
- [ ] Rendered report exists in `outputs/report.html` (and optionally `outputs/report.pdf`)
- [ ] Activity log entry created in `logs/activity_log.md`
- [ ] All changes committed to version control with a clear message
- [ ] No substantive changes to this plan document (unless explicitly requested)

## Instructions for AI Executor

[Provide specific instructions for what the AI should do. This is where you describe the tasks.]

1. Read the current repository state and this plan document
2. Create or update the executed analysis document (`executed.qmd`, `analysis.py`, etc.)
3. Implement the analysis specified in the Objective section
4. Ensure the analysis is executable and well-documented
5. Render outputs to the formats specified in Acceptance Criteria
6. Update `logs/activity_log.md` with a brief entry describing what was done
7. Commit all changes with a message that references this plan's objective

## Post-Run Checklist

[Items to verify after the iteration completes.]

- [ ] Review rendered outputs for correctness
- [ ] Check diffs to see what changed
- [ ] Verify logs are updated appropriately
- [ ] Confirm executable document runs successfully
- [ ] Decide: Accept, Revise Plan, or Rerun

## Notes and Decisions

[Use this space to record rationale, decisions, or context that might be useful for future iterations or reviewers.]

- Decision: [Record any important decisions made]
- Rationale: [Explain why]
- Next steps: [What might come next]

---

**Iteration History:**
- Created: [Date]
- Last modified: [Date]
- Status: [Draft / In Progress / Complete]
