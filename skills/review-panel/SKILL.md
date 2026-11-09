---
name: review-panel
description: Review a PR, branch or set of commits with independent subagent reviewers run in parallel — always an architect and a devil's advocate, plus a product owner, design/UX reviewer or QA engineer when the change calls for it — then vet, merge and rank their findings before reporting. Use this whenever the user asks to review a PR or changes "with subagents", "as an architect and a devil's advocate", asks for a second round of review after fixes, or asks whether a PR is ready to merge, even if they don't name this skill. Also use it to re-review only the newest commits after review fixes were applied.
---

# Review panel

Two reviewers with opposite instincts catch more than one generic pass: the
architect checks fit, correctness and consistency with the codebase and its
decision records; the devil's advocate tries to break it. You (the main agent)
are the third role — the editor who vets their findings before the user sees
them. Reviewers working from a plan will happily recommend "fixes" that violate
the plan's own constraints, so the vetting step is the part that protects the
user.

## 1. Gather the scope

Work out, from the request and the repo:

- **Target**: PR number(s), branch(es) and base, or specific commits. For
  stacked PRs, review each against the branch below it. For a re-review after
  fixes, name the newest commits explicitly and say what they were meant to fix.
- **Settled decisions**: product and design calls the user already made in this
  conversation, in the PR body, a plan file, or the repo's decision records
  (ADRs, DECISIONS.md and the like). List them so the reviewers don't
  relitigate them — relitigating settled calls is the most common way a review
  wastes the user's time.
- **Context pointers**: the plan file, the relevant ADR sections, related repos
  (e.g. an iOS client, index or rules files), and how to run checks (typecheck,
  lint, unit and integration tests, and any PATH setup the repo needs).

If the target is ambiguous, ask once; otherwise proceed.

## 2. Pick the panel

The architect and the devil's advocate always run. Add optional reviewers when
the change calls for them, or when the user names one ("add QA"):

| Reviewer | Add when the diff… | What it brings |
|---|---|---|
| **Product owner** | changes user-facing behaviour, defaults, privacy or copy | Product questions nobody answered, copy that over-claims, behaviour that was decided by accident |
| **Design / UX** | touches UI (views, components, styles, strings shown on screen) | Empty/loading/error states, truncation, design-system consistency, accessibility, match with designs if a link exists |
| **QA engineer** | adds endpoints, migrations, scripts, or anything with a manual rollout | Scenario list vs. existing tests, coverage gaps, a manual test script; runs smoke checks when it can |

Each extra reviewer costs time and adds noise for the vetting step, so add one
only when its row clearly applies, and say in one line which ones you picked
and why. Product owner and design/UX can't know the user's intent better than
the user: their output is **questions for the user**, never verdicts, so they
don't quietly make product decisions.

## 3. Launch the reviewers in one turn

Spawn the chosen reviewers as read-only subagents in parallel (general-purpose,
or Plan/Explore if they only need to read). Use `references/prompts.md` for the
prompt templates and fill in the scope from step 1. Every prompt must say:

- read-only: no edits, commits, pushes or GitHub comments;
- the settled decisions, marked "don't relitigate";
- what to focus on this round (newest commits, risky areas, merge readiness);
- integration tests that share an emulator or database must run one file at a
  time, and failures that don't reproduce alone are environmental;
- return a prioritized list (must-fix / should-fix / nit), each with file:line,
  a concrete failure scenario, and a fix.

While they run, tell the user in one line what each is checking.

## 4. Vet every finding before reporting

For each finding, decide: real, already handled, out of scope, or harmful.
Check the claim against the code when it's cheap to. In particular, reject or
flag fixes that:

- contradict a settled decision or a decision record (say which);
- leak information the design deliberately withholds (e.g. a pagination cursor
  that would carry a hidden item's id or timestamp);
- trade a rare edge case for a real regression.

When several reviewers report the same issue, merge it and say who found it —
that's a strong signal. When they disagree, say so and give your call.

## 5. Report

Lead with the verdict (ready / ready after X / not ready) and what was run
(typecheck, lint, test counts). Then:

- **Must-fix**, **Should-fix**, **Nits** — each item one or two sentences in
  plain language: what breaks, for whom, and the fix. Keep file:line where it
  helps; drop jargon the user hasn't used.
- **Questions for you** — product and design questions from the product owner
  or design/UX reviewer, phrased so the user can answer in a line.
- **Test plan** — QA's manual test script, if QA ran (also useful input for a
  review brief's "Try it" section).
- **Rejected suggestions** — any reviewer fix you're not recommending, and why.
- **Worth knowing** — non-defects that matter for rollout (deploy order,
  behaviour changes that go live on merge, flaky tests unrelated to the PR).

Then ask which fixes to apply, with your recommendation first. If a finding
raises a product question rather than a code question, ask that separately and
plainly — those are the user's calls.

## 6. Applying fixes (when asked)

Apply fixes on the branch where the code lives, bottom of the stack first, then
rebase the branches above. Re-run the checks, push, and summarise per PR what
changed. Re-reviews after a round of fixes should focus the reviewers on the new
commits plus overall coherence, not re-read everything from scratch.
