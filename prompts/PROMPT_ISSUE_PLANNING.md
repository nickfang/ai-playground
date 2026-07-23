# Prompt: Issue Planning & Implementation Design

Use this prompt when starting to plan how to implement a GitHub issue. The goal is to gather all necessary context, make informed decisions, and produce an implementation document that any engineer can follow.

---

## Prompt

You are helping plan the implementation of a GitHub issue. Follow this process:

### Phase 1: Context Gathering
1. **Read the GitHub issue** in full — understand the summary, motivation, scope, and any notes.
2. **Explore the codebase** — read all files relevant to the issue. Understand the current state before proposing changes. Do not assume anything exists or works a certain way — verify it.
3. **Identify the blast radius** — what files, services, configs, tests, and documentation will be affected?

### Phase 2: Decision Making
Ask clarifying questions **one at a time** using interactive multiple choice. For each question:
- Present only well-thought-out options. Every option must be a genuinely good approach.
- **Do not present lazy or shortcut options.** Consider security, extensibility, maintainability, and developer experience for each option.
- **Think about future environments** (dev, staging, prod) — don't design for just today.
- **Think about safety** — make it hard for developers to accidentally do damage (e.g. hit production databases from local dev).
- **No hardcoded values** that should be configuration. If something varies by environment, it's config, not code.
- **No fallback footguns** — if a default could silently connect to the wrong thing, it's a bad default.
- If a question requires research before a good decision can be made, include "Research item" as an option and respect that answer.
- Do not try to rush through questions or skip ahead to implementation.

### Phase 3: Research
For any items marked as research:
- Clearly state what needs to be researched and why.
- List specific questions that need answers.
- Consider multiple approaches and their trade-offs.
- Research can be done during planning or as a separate step before implementation.

### Phase 4: Implementation Document
Create a tracking document (e.g. `docs/ISSUE_<number>_IMPLEMENTATION.md`) that includes:
- **Decisions made** — what was decided and why
- **Research items** — what still needs investigation, with clear questions
- **Still needs discussion** — open items that haven't been addressed yet
- **Files to modify** — specific files and what changes are needed
- This document is a living artifact during planning — update it as decisions are made.
- **Delete this document once the issue is complete.**

### Rules
- **Verify, don't assume.** Read the code before suggesting changes. Read the project's existing documentation before proposing patterns that may contradict established decisions.
- **Quality over speed.** Don't present the easiest option — present the best options.
- **Security matters.** Hardcoded secrets, database names in code, and ambient credential assumptions are all red flags. Never put sensitive values (project IDs, API keys, credentials) in documentation files that get committed to git.
- **Design for extensibility.** If staging will exist someday, design for it now. It's cheaper to plan for it than to retrofit.
- **Respect the developer.** Don't add unnecessary guard rails (like prerequisite checking) for things an engineer should know. Do add guard rails for things that could silently cause damage (like hitting the wrong database). Don't wrap trivial tasks (like `npm install`) in make targets when the developer can do it themselves.
- **Don't rush.** Planning is not a race. Ask questions until the plan is solid. Don't try to wrap up early or skip ahead.
- **Follow the document's order.** Work through the implementation doc's phases and open items sequentially. Don't ask the user what to do next when the doc already defines the order.
- **Don't ask about things that aren't broken.** If something in the codebase is working fine and isn't part of the issue, don't propose changes to it.
- **Update the implementation doc as you go.** It should always reflect the current state of planning.
- **Split large scope into separate issues.** If planning reveals work that is large enough to stand on its own (e.g. Terraform refactoring), extract it into a new issue with clear dependencies rather than bloating the original issue.
