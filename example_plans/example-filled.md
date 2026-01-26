# Example: Data Summary Report

## Objective

Generate a summary report analyzing the synthetic admissions dataset in `data/raw/admissions.csv`. The report should include:
- Descriptive statistics for key variables
- One visualization showing the distribution of admissions by department
- A brief narrative interpretation of the findings

## Constraints

- Do not modify this plan document
- Use only the data file `data/raw/admissions.csv`
- All outputs should be written to `outputs/` directory
- Use Quarto format for the executed analysis document (`executed.qmd`)
- Render outputs to both HTML and PDF formats

## Acceptance Criteria

- [ ] `executed.qmd` exists and executes successfully
- [ ] Rendered report exists in `outputs/report.html` and `outputs/report.pdf`
- [ ] Report includes descriptive statistics table
- [ ] Report includes at least one visualization
- [ ] Activity log entry created in `logs/activity_log.md`
- [ ] All changes committed with message: "Iteration 1: Generate baseline admissions report"

## Instructions for AI Executor

1. Read this plan and examine `data/raw/admissions.csv` to understand the data structure
2. Create `executed.qmd` with:
   - Code to load and examine the data
   - Code to compute descriptive statistics
   - Code to create a visualization (e.g., bar chart of admissions by department)
   - Narrative text interpreting the findings
3. Ensure the document uses proper Quarto YAML header for rendering
4. Execute the document to verify it runs without errors
5. Render to HTML and PDF formats in `outputs/`
6. Append an entry to `logs/activity_log.md`:
   ```
   ## Iteration 1 - [Date]
   - Created baseline report from admissions.csv
   - Generated descriptive statistics and department visualization
   - Rendered to HTML and PDF
   ```
7. Commit all changes with the message specified in Acceptance Criteria

## Post-Run Checklist

- [ ] Review `outputs/report.html` for correctness and completeness
- [ ] Check git diff to see what files were created/modified
- [ ] Verify `logs/activity_log.md` has new entry
- [ ] Test that `executed.qmd` can be rerun successfully
- [ ] Decision: [ ] Accept [ ] Revise Plan [ ] Rerun

## Notes and Decisions

- **Decision**: Using Quarto format for executable document to enable multi-format rendering
- **Rationale**: Quarto supports both code execution and narrative text, making it ideal for reproducible reports
- **Next steps**: If accepted, may refine visualization or add additional analyses in next iteration

---

**Iteration History:**
- Created: 2026-01-26
- Last modified: 2026-01-26
- Status: Draft
