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

revision: I don't know the right way of saying it, but I want to basically emphasize that this is worth learning because you're learning the thing that builds all the other things.

New Task: 
I want to experiment with design a bit. Make a separate page where I can build and isolate some toys. I want to create a toy that generlaly explains the concept of idempotent. Show contrast with a non-idempotent.
revision: see example. I love these beautiful visual animations. what you have done is very concrete, and not elegant and taking advantage of the technology.
i think both sides need some nice animations. Maybe a subtle touch of color. Really think about how to communicate idempotency for non technical users. maybe takeup less vertical space.
2. next toy: explain paragraph "IRL is a recipe for AI-assisted work" elegant, intutive animation.
    revision: can we demonstrate the principal of the human as reasoning and planning and the AI being the execution layer. but the thing that is missing is the loop part.
    revision: Not bad. I think you can make more effective use of space and make the figure more dense. words and labels are too small. I wonder if it can play on a loop a few times so the user doesn't have to keep clicking. for this they can just see how it works. 
    revision: better. Still words plan etc. are too small. I worry about the proportions. I don't think non-technical users understand what artifacts are. needs better caption. I gave you access to the chrome skill so you can now directly observe and iteriate! use it.
    Revision, can you make the arrows more prominent? They're kind of dotted and hard to see.
    The more I think about it, instead of files, maybe it's better to show that it's a report that's being generated or some kind of single output that's carefully being crafted. I don't know the right way of doing it. 
    I really like this next revision. What I'm thinking is very subtly, if the human and AI had whatever they're doing, like instructions pop in, fade out, pop in, fade out. It would kind of show the activity over time of how that report is being created directly from the actual instructions.
    It's too fast. I also want the instructions in the ample space around the objects. Maybe accumulating overtime. nice fade in etc. 

    Let's make a new toy. In this case, I'd like to have a simplified version of the main plan, just the critical details to communicate the main idea in an animation to the reader. It would be nice to have real output. To kind of show the process. Of course, simplify.
    Revision, base it on the actual IRL Basic Template. 
    Revision: we're almost there. There's a couple things I want to change. One is that for the actual template, there should be the before each loop, the instruction loop, and then after each loop. That's important to explain the process. I also think it would be better if you showed it side by side instead of vertically so we can see the main plan being executed and then the outputs in real time. I think also doing two or three loops to show what can happen over time, including a revision part, would be good. Kind of like you're seeing me do right now. I think the git diff is not essential. It'd be hard for non-technical users to understand that. I do like the checkpoint save part though.
    Revision really close, but when I'm looking at the actual plan itself, it has to pass muster if data scientists were to read this. So maybe make sure that it's actually realistic. It just says produce O3 outputs report. HTML, it doesn't make a lot of sense. I think you can add a little bit more detail there. Given that the loop column is longer than the main plan column, you probably have a little bit of room in the main plan.

    New Toy: Anatomy of a main plan. Make it interactive like clicking on things and can give more details after a nice layout.
    Revision: I want to see a more visual look of the plan as opposed to this accordion thing. I don't think it works very well for this. I think we want to see the plan on one side and then on the right side we have the captions and the little labels people could click on and see the description.

    New Toy: how to execute them main-plan with a terminal based ai. "Review main-plan.md, check for any revisions, and execute" as the command you run over and over.

    New Toy: i'm thinking a layout of a VS Code screen would be helpful. Showing the file list, the editor, the preview, and the terminal could really lay out exactly how the tools work. Again, visual that somebody could click on the different things, or even animated walking through each of the little areas. I think, in terms of the editor being on the left-hand split, and then the preview of the HTML report in the right-hand split would be nice. Go ahead and try a draft.
    Revision. This is a minor revision, but if you're going to show the editor, I think right now the bottom doesn't have a fixed height, so it just looks kind of squished. It doesn't have the proportions of a real editor. So I would just go ahead and make it have the right proportions the terminal should have some minimal height.

    New Toy: following up on our last revision, the thing that builds all the other things, I think it would be nice to have a very elegant animation, stepping through different potential uses of IRL and showing how the main plan essentially stays the same, but there's different instructions or custom things. I want to introduce things like document skills, PubMed skills, analysis skills. And so maybe our three basic use cases are building a web page, doing a data analysis, maybe doing something like a document, a literature review, something like that.

    For all the toys please update the text and captions to reflect the changes.
    These are not comptuer programers so refer to the real irl template where we aren't using words like git, use words like version control.
    Add reset button to each toy

Let's draft a new complete doc that will include elements from both index.html and our toy page
I do not want to integrate the toys yet
instead make a list of the toys
think about the most reasonable outline and order
then make the page with an outline only with some brief descriptions. I will review this first.

revision: Don't you think for the combined document that people really need a hook and learning how to use IRL to build all the things? It's not the technique, it's that people want to use AI, but here is a way to use it, learn it, not tied to a particular commericial product.

revision: should emphasize early that since this isn't about the "features" of the AI or the software or coding, this is about your expertise, judgement, and reasoning and the goal is to review the outputs. No different than if you were to have a human assistant. The key loop is writing out reasoning and critical review of the outputs. Could you integrate this into the outline but also draft in the toys a new toy that can capture this well. People get concerned about AI making mistakes, but the key like human assistant make mistakes, the key is runderstanding your reasoning and reviewing the outputs, then looping back and asking for clarification, running verification, restating your reasoning, and lopping over an dover till you trust and like the output.

Revision: Toy. Almost perfect but like the other toys i like to have a right and left side and an animation of this review process loop. beautiful and elegant.

Next: Good work on the toy. Let's get back to refining the combined outlines.

I think the Getting started can be two things - you can start with a text plan file and have your terminal AI run it. You can also introduce the IRL go app.

New Task: Now make the final combined webpage based on the outline. Copy the outline into a new file first and then copy over sections to avoid any attrition or regression. Then go back and refine.

Revision: 
    - Prose tone and flow across all 13 sections
    - Check need more/less content
    - integrate the actual toys from toys.html

New Task:
Review /Users/ernie/Documents/GitHub/docfleet 
I want to move the final documentation page here
Please write out a plan in markdown on integration plan

Next for deployment:
hostname - www.irloop.org
no toys page needed 
hel2 is in ssh config (that's the ip)
just fyi hel2 has kamal proxy on 8080 and cloudflare tunnel terminates there
itis an amd64 server
please check status of kamal, containers, and www.irloop.org on hel2 so we can understand the landscape before pushing it. I agree with turning off ssl. 
i think currently there is a kamal container for the irloop (using mintify docs) that I want to retire and replace with this site.

reponse:
do an op service token check
pre connect test
I removed the existing kamal app the old docs container is removed so www.irloop.org is now free. 

read this to fix the 1password issue
https://github.com/drpedapati/fmrplean/issues/62

please turn docfleet repo private and make sure no secrets are commited

revision: clean up text no em dashes or clichees

copy and insert tastefully into first few paragraphs /Users/ernie/Downloads/irl_framework_recipe.png

new task: make a new toy to explain to non-tech users what markdown is. Make the component slim so it doesn't take alot of vertical space. use the left right pattern again. it's to support his paragraph: "IRL is different. It's a plain-text pattern, not a product. The plan file you write is just a text document. Any AI assistant can read it. The outputs are files on your computer that you own. Nothing is stored inside a proprietary platform. If you switch AI providers tomorrow, the pattern still works exactly the same way."

Revision: Modify the new plain text #9 to briefly comment that the funny characters in the plain text are markdown, still plain text but having some formatting, say this in a professional way that will make sense to non-tech user.

New Task: Found a critical point: there is no lead-in prose or friendly invitation to what each toy is for. It just plops the toy in the middle of the text. Please systematically address this for each toy.

New Task: I want a new page (you can copy the current draft.html as a template). replace content with an example task of a literature review using the pubmed skill and docx skill. I want progressive teaching starting with the basic template. Build it for the new user both to using an AI and terminal interfaces. TO do this possible, you need a new parameterized toy that is the IDE framework (vs code mockup that you built) and you can use it with animation through out the tutorial to walk people through. If you need to test real outputs you can do that.
revision: base it on the real basic template and add the pubmed and docx skill. learn about them so you can make a complete tutorial. We run the AI by opening a terminal AI and running the command "review main-plan.md for any revisions and execute" over and over.

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

- **Version control**: Commit intended changes only + main-plan.md; verify no ignored or unintended files staged.

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
