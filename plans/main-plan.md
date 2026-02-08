# IRL Basic Template

🔧 First Time Setup — Run once when starting a new project
<!-- 👤 AUTHOR AREA: Define setup steps below -->


### Common Skill Library — Pre-installed tools available to all projects
<!-- Uncomment to use -->
Install Quarto Skill: https://github.com/posit-dev/skills/tree/main/quarto/authoring
Install Explainer Skill: https://github.com/henrybloomingdale/explainer-site-skill
<!-- Install Word DOCX Skill: https://github.com/anthropics/skills/tree/main/skills/docx -->

## ✅ Before Each Loop — Checklist to run before every iteration

- Script designed to be **idempotent**
- **Version control**: clean tree (`git status`), commit baseline, `.gitignore` enforced
- The **only permitted modification** is within the `## One-Time Instructions` and `AI Feedback` section **ONLY** on request.

---

## 🔁 Instruction Loop — Define the work for each iteration

<!-- 👤 AUTHOR AREA: Define each loop's work below -->

1. Install skills 
2. learn about IRL
3. use the explainer skill to draft a site explaining IRL

locally serve the website so I can test it.
add makefile in website dir for easy management
AI Feedback: Add ZeroTier address to access site

4. The senior reviewer has reviewed your site and has provided extensive critical feedback. Before making any edits, please review all of the files in '/Users/ernie/Dropbox/Henry Projects/irl-feedback'

it is hard to believe that you made all of the revisions in just one pass. Look at where you've been and what you've done in terms of responding to feedback. Go back to the feedback docs, read some more, and then continue to apply. Iteration will get better and better

review the typography md feedback in particular and apply
review sample website in explainer skill and apply similar styling
there is an elegance in the original that is not here
for example the title area seems so basic
where are the subtle shading int he background etc
I want some better white space management on my page, especially the header. Add author Ernest Pedapati, MD

remove this: Artifact: 03-outputs/irl-explainer-site/index.html
revisions:
-probably explain what idempotent means early, audience is non-technical
-add intuitive "toy" (see skill) for idempotent for non-technical users
i want the title area to look like this: http://10.241.64.217:3000/ with centered, subtitle below title and I need a new main title that is better for non-technical users and subtitle can be idempotent research loop (technical name)

new first paragraphs (and remove table of contents)
I'm obsessed with a specific kind of failure: you do something smart with an AI assistant, it works, you feel productive... and then the whole thing evaporates into a chat transcript nobody will ever read again. A week later you can't answer basic questions. What inputs did we use? What assumptions did we make? If the result is wrong, where did it go wrong?
Think about the difference between cooking from memory and cooking from a recipe. From memory, you might make something great, but you can't teach it to someone else, you can't reliably make it again, and if it goes wrong you're just guessing at what changed. A recipe externalizes the knowledge. You can trace problems, share it with others, and build on it deliberately.
IRL is a recipe for AI-assisted work. You write a plan file with your inputs, assumptions, and what the output should be. The AI follows the plan and produces files you can actually inspect: a report, a dataset, a diagram. You start small and layer complexity one step at a time. "Idempotent" just means: same recipe, same ingredients, same dish. You can rerun any step, add to it, and know the foundation hasn't shifted underneath you.
The design plays to what each side is good at. Humans decide what matters, judge quality, notice when something is off. AI processes data, generates drafts, follows instructions without getting tired. But AI has no memory of why it did something last week, and humans lose track of assumptions fast. The plan file sits in the middle. It gives you a place to think clearly about what you want and gives the AI stable instructions that don't depend on anyone's memory. The files that come out become the shared record that holds both sides honest.

New Task: 
I want to experiment with design a bit. Make a separate page where I can build and isolate some toys. I want to create a toy that generlaly explains the concept of idempotent. Show contrast with a non-idempotent.
revision: see example. I love these beautiful visual animations. what you have done is very concrete, and not elegant and taking advantage of the technology.
i think both sides need some nice animations. Maybe a subtle touch of color. Really think about how to communicate idempotency for non technical users. maybe takeup less vertical space.
2. next toy: explain paragraph "IRL is a recipe for AI-assisted work" elegant, intutive animation

### One-Time Instructions — Tasks that should only execute once

<!-- 👤 AUTHOR AREA: Add one-time tasks below -->
1. again explainer skill has been updated and should be redownloaded, header update


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

## AI Feedback
<!-- AI Edits AREA -->


## 📚 Skill Library — Optional community skills to install per project
<!-- Uncomment to use -->

<!-- Install Posters -->
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
