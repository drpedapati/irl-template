Katie Stefani, DO<sup>1</sup>, Donald Gilbert, and Ernest Pedapati, MD<sup>1,2,3</sup>

<sup>1</sup>Division of Child and Adolescent Psychiatry, Cincinnati Children’s Hospital Medical Center, 3333 Burnet Ave., Cincinnati, OH, 45229-3039, USA. 

<sup>2</sup>Division of Neurology, Cincinnati Children’s Hospital Medical Center, 3333 Burnet Ave., Cincinnati, OH, 45229-3039, USA. 

<sup>3</sup>Department of Psychiatry and Behavioral Neuroscience, University of Cincinnati College of Medicine, Stetson Building Suite 3200,  260 Stetson Street, Cincinnati, OH 45267-0559, USA. 

Corresponding Author:

ernest.pedapati@cchmc.org

# Visual Abstract <img src="03-outputs/manuscript/media/media/image1.png" style="width:6.32661in;height:3.6224in" />

# Summary

Large language models (LLMs) are increasingly used in scientific work, yet chat-based usage patterns conflict with norms of reproducibility, provenance, and auditability. We describe the Idempotent Research Loop (IRL), a document-centric workflow architecture for human-AI collaboration in scientific analysis. IRL treats a single plan document as the control surface: a human editor specifies intent and constraints, and the AI executes tasks to produce versioned, structured research artifacts (e.g., executable analyses and rendered reports) rather than ephemeral conversational answers. IRL operationalizes six pillars: idempotency, document-centric control, separation of reasoning from execution, multimodal report outputs, provenance-by-default, and a managerial human-in-the-loop model. These are implemented through paired artifacts and pre/post execution hooks that bracket each iteration with version control and logging. We illustrate the loop with a case study based on a reproducible report-generation workflow (Quarto as one concrete instantiation) and show how practical idempotency can be assessed as convergence of repository diffs across reruns. IRL does not claim algorithmic novelty; it provides an operational architecture that makes AI-assisted analysis more reviewable, repeatable, and shareable under common scientific expectations for computational work.

Keywords: Reproducible research; provenance; literate programming; workflow patterns; human-AI collaboration; large language models

# Abstract

Conversational large language model (LLM) interfaces are attractive for exploratory analysis, but they often leave scientific work with weak provenance, ambiguous authority over results, and susceptibility to conversational drift. We propose the Idempotent Research Loop (IRL) as a workflow architecture that shifts AI-assisted analysis from chat transcripts to versioned, document-centric artifacts. In IRL, a single plan document functions as the control surface: the human editor specifies intent, constraints, and acceptance criteria; the AI reads the plan and repository state, executes tasks, and writes back structured outputs (executable analysis documents, rendered reports, figures, tables, and logs). IRL enforces a separation between reasoning/intent and execution via paired artifacts (e.g., a human-editable plan in Markdown and an executable analysis document that renders to HTML/PDF/Word). Each iteration is bracketed by provenance defaults (e.g., commits, change logs, decision records), enabling diffs to function as an audit interface. We define practical idempotency for AI-assisted research workflows as convergence: rerunning the loop with unchanged inputs yields no substantive changes to authoritative artifacts within declared tolerances. Through a case study, we illustrate how IRL produces reviewable, rerunnable outputs and makes human oversight explicit. IRL is not a tool or agent; it is a pattern that can be implemented with existing research infrastructure to support auditable AI assistance.

# 1. Introduction

Computational research increasingly depends on complex, multi-step analyses that must remain interpretable and reproducible beyond the original analyst. Reproducible research guidance emphasizes executable workflows, shared data and code, explicit environments, and transparent histories of decisions and changes \[1-3\]. In parallel, LLMs are now commonly used to accelerate drafting, coding, summarization, and exploratory data analysis. These systems can be useful, but they are not fully reliable and can produce confident inaccuracies, including fabricated citations and unsupported claims \[18-20\]. This tension is particularly acute when LLM usage is mediated primarily through chat: conversational history becomes an implicit state store, and what counts as the authoritative specification of the analysis is often unclear.

This manuscript describes a workflow architecture intended for scientific contexts where auditability and rerunnability matter: regulated environments, collaborative lab work, and analyses expected to survive peer review and post-publication scrutiny. Our core thesis is that many failure modes of chat-based AI usage in science can be mitigated not by proposing new model capabilities, but by changing the control surface and artifact model of human-AI collaboration.

We describe the Idempotent Research Loop (IRL): a document-centric, iterative execution loop in which a human editor acts as a managerial controller (analogous to a PI or lead analyst), while an AI assistant performs bounded execution steps that produce structured, versioned artifacts.

## 1.1 Design goals

- Reduce conversational drift by treating a stable plan document as the authoritative interface.

- Improve auditability by ensuring work products are versioned, diffable artifacts.

- Improve reproducibility by producing executable analysis documents and deterministic renders when environments are controlled.

- Make human oversight explicit through a managerial human-in-the-loop loop contract.

- Provide a practical definition of idempotency for AI-assisted scientific work.

# 2. Related Work

## 2.1 Literate programming and computational notebooks

Literate programming framed programs as narratives interleaving code and explanation \[7\]. Computational notebooks extend this idea for interactive analysis and communication \[8\]. However, notebooks can exhibit hidden state and execution-order ambiguity, complicating reproducibility and review \[9-10\]. IRL borrows the strengths of literate artifacts (narrative + executable code + outputs) while addressing control-surface weaknesses by separating plan from execution and by making idempotency checks and provenance defaults part of the workflow contract.

## 2.2 Workflow engines and pipeline standards

Workflow engines (e.g., Snakemake, Nextflow) and standards (e.g., CWL) provide explicit dependency graphs and portable execution descriptions \[11-13\]. These systems excel at deterministic execution once the pipeline is specified, but they are not optimized for the iterative human decision-making and narrative reporting common in scientific analysis. IRL is complementary: it is not a workflow engine, but a pattern that can incorporate workflow engines within an iteration as part of execution.

## 2.3 Provenance, versioning, and research object packaging

Version control and provenance models are widely recommended for computational research transparency \[2-4\]. W3C PROV provides a formal model for describing provenance \[6\]. FAIR principles emphasize making scientific objects findable and reusable \[5\]. IRL operationalizes these ideas pragmatically: each loop iteration produces versioned and logged artifacts, enabling diffs to serve as an audit interface and facilitating later packaging as a research object.

## 2.4 Human-AI collaboration and LLM reliability

Human-AI interaction research emphasizes that interactive systems should support oversight, error recovery, and appropriate expectation-setting \[15-16\]. LLMs remain fallible and can hallucinate, motivating workflows that encourage grounding and verification \[18-20\]. IRL adopts a managerial model in which the human specifies intent and acceptance criteria and retains responsibility for conclusions; the AI performs bounded tasks and updates artifacts under a plan-defined contract.

# 3. The Idempotent Research Loop

## 3.1 Formal idempotency and practical idempotency

In computer science, an operation is idempotent if applying it multiple times has the same effect as applying it once (e.g., f(f(x)) = f(x)). Scientific analysis workflows that involve stochastic components (including LLM generations) cannot generally guarantee byte-for-byte identical outputs across reruns. IRL therefore defines a pragmatic operational target.

Definition (Practical idempotency for IRL). Given a repository state S that includes the plan document, data, environment specification, and toolchain configuration, an IRL iteration is practically idempotent if repeated execution converges such that subsequent runs produce no substantive changes to the declared authoritative artifacts, within a declared tolerance model (e.g., permitting timestamps in logs, nondeterministic ordering in plots if documented).

In practice, no substantive changes is assessed via: repository diffs between successive runs; output checksums for declared deterministic renders when feasible; and explicit tolerance rules recorded in the plan (e.g., log timestamps may differ).

## 3.2 Core architecture and artifact model

IRL is defined by three ingredients: (1) a single authoritative plan document (plan.md) that encodes intent, constraints, and instructions, and is edited by the human controller; (2) an execution substrate (local or remote) that can run code, render reports, and modify files; and (3) a loop contract: pre-instruction hooks, execution steps, and post-instruction hooks applied each iteration.

We distinguish artifact roles.

Authoritative artifacts: plan.md (intent + constraints + acceptance criteria); executed analysis document (e.g., executed.qmd) that yields the reported outputs when run; and declared data inputs.

Derived artifacts: rendered reports (HTML/PDF/Word), figures, tables; activity and decision logs; diagrams of plan/structure.

## 3.3 Conceptual diagram

A conceptual IRL diagram contains a Human Controller node that edits plan.md; an AI Executor node that reads plan.md and repository state; and a Repository/Workspace node containing executed analysis, data, logs, and rendered outputs. Arrows form a loop: human edits plan -\> AI executes and updates artifacts -\> outputs rendered and logs/commits written -\> human reviews diffs and outputs -\> repeat.

## 3.4 One iteration lifecycle

A single IRL iteration proceeds as: Step 0 (Human edits plan): update objectives, constraints, acceptance criteria; record decisions and rationale when scope changes. Step 1 (Pre-instruction hooks): confirm clean working state; commit/snapshot current repository; ensure toolchain consistency (pinned environment, fixed seeds where relevant). Step 2 (AI execution): read plan.md; modify or create analysis artifacts; avoid modifying plan.md unless explicitly permitted. Step 3 (Post-instruction hooks): render executed analysis into declared formats; update logs; update diagrams if used; commit changes with a message linked to the iteration objective. Step 4 (Human review): review rendered report, diffs, and logs; accept, revise plan, or rerun.

<img src="03-outputs/manuscript/media/media/image2.png" style="width:6.5in;height:2.75694in" alt="A diagram of steps AI-generated content may be incorrect." />

# 4. Implementation Pattern

## 4.1 Invariant components

Across the provided instruction templates, the following components recur: (1) pre-instruction hooks (do not edit the instruction script; commit/snapshot before changes; enumerate files that may be modified); (2) paired artifacts (a reasoning/plan artifact and an executed analysis artifact); (3) post-instruction hooks (render outputs; update logs and diagrams; commit/snapshot after changes); (4) decision and activity logs; and (5) report rendering to durable formats.

## 4.2 Minimal directory skeleton

A minimal IRL repository can be organized as:

project/  
plan.md  
executed.qmd  
data/  
raw/  
derived/  
outputs/  
report.html  
report.pdf  
logs/  
activity_log.md  
decision_log.csv  
env/  
environment.yml (or renv.lock / requirements.txt)  
README.md

This structure is intentionally generic: Quarto and Git are examples, not requirements. The invariant requirements are a stable plan document, rerunnable executed artifacts, durable rendered outputs, and a versioned history with logs.

## 4.3 Loop contract pseudocode

repeat until accepted:  
human edits plan.md  
  
pre_hook:  
assert plan.md unchanged by AI  
snapshot repository state (commit/tag)  
  
ai_execute(plan.md):  
update executed analysis (executed.qmd, scripts, configs)  
run analysis + tests  
render outputs (HTML/PDF/Word)  
update logs (activity + decision)  
  
post_hook:  
snapshot repository state (commit/tag)  
record iteration summary in decision_log  
  
human reviews outputs + diffs

## 4.4 Practical idempotency checks

IRL treats rerun stability as an explicit check. A fixed-point check runs the same iteration twice without changing plan.md or inputs; the expected outcome is an empty or tolerance-limited diff. If diffs persist, the plan can be tightened (e.g., forbid rewriting sections unnecessarily), environments pinned, or nondeterministic work isolated into deterministic code steps.

## 4.5 Tool-specific instantiation: Quarto

Quarto is one concrete substrate for executed analysis documents because it supports executable documents and multi-format rendering (HTML/PDF/Word). IRL does not require Quarto; the essential requirement is a rerunnable document format that can be rendered into durable outputs.

# 5. Case Study: Iterative, auditable report generation from a synthetic admissions dataset

## 5.1 Goal and setup

We illustrate IRL using a report-generation workflow based on a synthetic admissions dataset. The demonstration uses a plan document to request synthetic data generation, descriptive summaries, plots, narrative interpretation, and report rendering.

Artifacts: plan.md (scope, constraints, acceptance criteria, post-run checklist); executed.qmd (data generation and analysis code); rendered outputs (HTML/PDF); logs (activity and decision); and version control commits bracketing each iteration.

## 5.2 Iteration 1: Create baseline report

Human action: draft plan.md specifying objectives (generate synthetic CSV; summarize admissions; produce plots and narrative; render report).

AI execution: generate executed.qmd, render to HTML/PDF, update logs, and commit.

Checkable outcomes: (i) a compiled report exists as a durable artifact; (ii) the analysis is rerunnable from executed.qmd under the specified environment.

## 5.3 Iteration 2: Plan edit to refine scope

Human action: edit plan.md to refine the cohort (e.g., an inpatient psychiatry unit) and request specific visualization changes.

AI execution: update only necessary sections of executed.qmd, re-render outputs, append to logs, and commit.

Auditability: reviewers can inspect diffs between iteration commits and see exactly which code, figures, and text changed, along with log entries that justify changes.

## 5.4 Iteration 3: Practical idempotency check

Human action: rerun the loop without changing the plan or inputs.

Expected result: authoritative artifacts (executed.qmd and rendered outputs) show no substantive diffs; logs may record a rerun confirmation depending on policy. If differences occur (e.g., random seeds), they are diagnosed and corrected by pinning randomness and documenting tolerances.

# 6. Discussion

## 6.1 Strengths (mechanism-based)

Auditability via diffs. By making work products files in a versioned repository, IRL enables review through standard tooling (diffs, commit history), operationalizing reproducibility guidance \[2-4\].

Reduced hidden state. IRL relocates state from conversational context to explicit artifacts (plan, executed analysis, environment specifications, logs), addressing known reproducibility challenges in interactive environments \[9-10\].

Separation of intent from implementation. Paired artifacts reduce the need to steer by editing code directly: humans adjust a short plan and the executed document is regenerated and rendered under contract.

Clear accountability model. IRL is managerial: the human owns intent and acceptance; the AI performs bounded tasks. This aligns with Human-AI interaction guidelines emphasizing oversight and recovery when systems are wrong \[15-16\].

## 6.2 Limitations and non-ideal contexts

LLM stochasticity and model drift. IRL does not eliminate nondeterminism; it contains it by forcing changes into reviewable diffs and encouraging deterministic downstream execution and environment pinning where feasible.

External dependencies. Workflows relying on APIs, web content, or changing databases may not be stable across reruns. IRL mitigates this by anchoring claims to stable identifiers (DOI, PMID, accession IDs) and recording retrieval times and inputs.

Overhead. The loop imposes structure (plan upkeep, logging, commits). For small tasks or informal exploration, IRL may be unnecessarily heavyweight.

Not a substitute for rigorous statistics or domain expertise. IRL organizes work; it does not validate scientific correctness. The human controller remains responsible for methodological decisions and interpretations.

## 6.3 Failure modes and mitigations

Persistent diffs on rerun: pin seeds; isolate nondeterministic steps; define tolerated differences explicitly.

Citation hallucinations: require stable identifiers and verification steps; enforce this in the plan for literature synthesis workflows.

Scope creep in executed document: enforce plan constraints such as only updating indicated sections.

Unreviewable changes: require short commit messages linked to plan changes and decision-log entries summarizing rationale.

## 6.4 Relationship to agent loops

IRL is not an autonomous agent architecture. It can incorporate agent-like tooling, but it enforces a single authoritative plan document, explicit artifact updates, versioning and logging, and human acceptance as the default endpoint. This differs from approaches that interleave reasoning and actions in extended conversational contexts \[21\], which can be effective for some tasks but complicate provenance unless paired with explicit artifact management.

# 7. Implications for Scientific Practice

## 7.1 Lab adoption and onboarding

IRL lends itself to lab-scale standardization because the plan file functions like a living protocol for an analysis. New contributors can read the plan, run the executed document, and inspect diffs/logs to understand what changed and why.

## 7.2 Peer review and post-publication audit

Because IRL produces durable reports and versioned histories, it can support transparency aligned with reproducibility norms \[1-3\]. Repositories can be packaged and shared alongside manuscripts, improving reviewers' ability to audit claims.

## 7.3 Alignment with open science and FAIR principles

IRL encourages structured, identifiable, and versioned research objects, complementing FAIR goals for making outputs findable and reusable \[5\].

# 8. Conclusion

The Idempotent Research Loop (IRL) is a document-centric workflow pattern for AI-assisted scientific analysis. It reframes LLM usage from conversational interaction to artifact-based execution: a human-controlled plan document guides AI execution that produces rerunnable analysis documents, durable report renders where feasible, and provenance-by-default histories. IRL does not claim algorithmic novelty or model autonomy. Its contribution is an operational architecture that makes AI assistance more compatible with scientific expectations for reproducibility, auditability, and accountable collaboration.

# 9. Reproducibility and Availability

A reproducibility package for IRL should include: (i) a minimal template repository (plan + executed document + logs); (ii) at least one public case study (synthetic or openly licensed data); (iii) environment specifications (e.g., Conda, renv, Docker); and (iv) a short guide describing how to rerun and verify practical idempotency.

In this submission, the demonstration transcript and example instruction templates are provided as supplementary materials (Supplementary Files S1-S3).

# References

1\. Peng, R.D. (2011). Reproducible research in computational science. Science 334, 1226-1227. doi:10.1126/science.1213847.

2\. Sandve, G.K., Nekrutenko, A., Taylor, J., and Hovig, E. (2013). Ten simple rules for reproducible computational research. PLoS Comput. Biol. 9, e1003285. doi:10.1371/journal.pcbi.1003285.

3\. Wilson, G., Aruliah, D.A., Brown, C.T., et al. (2014). Best practices for scientific computing. PLoS Biol. 12, e1001745. doi:10.1371/journal.pbio.1001745.

4\. Ram, K. (2013). Git can facilitate greater reproducibility and increased transparency in science. Source Code Biol. Med. 8, 7. doi:10.1186/1751-0473-8-7.

5\. Wilkinson, M.D., Dumontier, M., Aalbersberg, I.J., et al. (2016). The FAIR Guiding Principles for scientific data management and stewardship. Sci. Data 3, 160018. doi:10.1038/sdata.2016.18.

6\. Moreau, L., Missier, P., and the W3C Provenance Working Group (2013). PROV-DM: The PROV Data Model (W3C Recommendation).

7\. Knuth, D.E. (1984). Literate programming. Comput. J. 27, 97-111. doi:10.1093/comjnl/27.2.97.

8\. Kluyver, T., Ragan-Kelley, B., Perez, F., et al. (2016). Jupyter Notebooks - a publishing format for reproducible computational workflows. In Proceedings of ELPUB 2016. doi:10.3233/978-1-61499-649-1-87.

9\. Pimentel, J.F., Murta, L., Braganholo, V., and Freire, J. (2021). Understanding and improving the quality and reproducibility of Jupyter notebooks. Empir. Softw. Eng. 26, 65. doi:10.1007/s10664-021-09961-9.

10\. Rule, A., Birmingham, A., Zuniga, C., et al. (2019). Ten simple rules for writing and sharing computational analyses in Jupyter Notebooks. PLoS Comput. Biol. 15, e1007007. doi:10.1371/journal.pcbi.1007007.

11\. Koster, J., and Rahmann, S. (2012). Snakemake - a scalable bioinformatics workflow engine. Bioinformatics 28, 2520-2522. doi:10.1093/bioinformatics/bts480.

12\. Di Tommaso, P., Chatzou, M., Floden, E.W., et al. (2017). Nextflow enables reproducible computational workflows. Nat. Biotechnol. 35, 316-319. doi:10.1038/nbt.3820.

13\. Amstutz, P., Crusoe, M.R., Tijanic, N., et al. (2016). Common Workflow Language, v1.0. Specification. doi:10.6084/m9.figshare.3115156.v2.

14\. Quarto Development Team. Quarto Documentation. https://quarto.org/docs/ (accessed 2026-01-23).

15\. Amershi, S., Cakmak, M., Knox, W.B., and Kulesza, T. (2014). Power to the People: The Role of Humans in Interactive Machine Learning. AI Magazine 35(4). doi:10.1609/aimag.v35i4.2513.

16\. Amershi, S., Weld, D., Vorvoreanu, M., et al. (2019). Guidelines for Human-AI Interaction. In Proceedings of CHI 2019. doi:10.1145/3290605.3300233.

17\. Reynolds, L., and McDonell, K. (2021). Prompt programming for large language models: Beyond the few-shot paradigm. arXiv:2102.07350. doi:10.48550/arXiv.2102.07350.

18\. OpenAI (2023). GPT-4 Technical Report. arXiv:2303.08774.

19\. Huang, L., Yu, W., Ma, W., et al. (2025). A survey on hallucination in large language models: principles, taxonomy, challenges, and open questions. ACM Computing Surveys. doi:10.1145/3703155.

20\. Bender, E.M., Gebru, T., McMillan-Major, A., and Shmitchell, S. (2021). On the dangers of stochastic parrots: Can language models be too big? In FAccT 2021. doi:10.1145/3442188.3445922.

21\. Yao, S., Zhao, J., Yu, D., et al. (2023). ReAct: Synergizing reasoning and acting in language models. arXiv:2210.03629.
