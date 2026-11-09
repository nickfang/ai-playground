---
name: review-brief
description: Write a one-screen reviewer's brief for a PR (or a stack of PRs) so a busy human can review it in minutes — what changes for users, the decisions baked in and who made them, blast radius, what's irreversible, how it was verified, test names as a plain-English spec, and the two or three spots actually worth reading — and publish it as an easy-to-scan page. Use this whenever the user says they can't keep up with reviewing PRs, asks what they should look at in a PR, asks for a summary or explanation of a PR before merging, or asks how to review quickly, even if they don't say "brief".
---

# Review brief

The person reading this is short on time and can't read every diff. Agents are
good at checking code against intent and bad at noticing when the intent is
wrong, so the brief's job is to put the *decisions* and the *risky spots* in
front of the human, not to narrate the diff. A change-by-change walkthrough is
the diff in prose — don't write one.

## 1. Gather

For the target PR(s):

- `gh pr view` (title, body, base, checks) and the diff stat; the diff itself
  only as needed to answer the questions below.
- Decisions: the PR body, any plan file, decision-record addenda (ADRs,
  DECISIONS.md) the PR adds or edits, and the conversation. For each decision,
  note who made it — the user, someone they quoted (e.g. the CEO), or the agent.
  Agent-made decisions are the ones the human most needs to see.
- Test names added or changed (`it(...)`, `test(...)`, `func test…`) — titles
  only.
- Verification: CI status, what was run locally, emulator or manual runs, and
  anything not verified (e.g. "iOS not compiled").
- If `review-panel` ran in this conversation, its verdict and open findings.

## 2. Write the brief

Use this structure; keep the whole thing to roughly one screen. Plain
language, no jargon the user hasn't used.

1. **Verdict** — one line: ready / ready after X / not ready, and why.
2. **What changes for users** — 2–4 bullets of behaviour, not code.
3. **Decisions baked in** — a table: decision · who decided · where it's
   recorded. Mark agent-made decisions so they stand out.
4. **Blast radius** — what breaks or leaks if this is wrong, and who notices.
5. **Irreversible or hard to undo** — data migrations, deletes, things users
   have already seen, deploy-order constraints. Say how to roll back, if at all.
6. **How it was verified** — checks and counts; call out gaps honestly.
7. **Read these** — the 2–3 spots worth the human's eyes, as file:line links to
   the PR, each with one line on why. Prefer risk tier: privacy, data
   migrations and irreversible scripts get code links; read-only endpoints and
   copy get the decision-record addendum instead.
8. **Spec in plain English** — the test titles, grouped, so missing cases jump
   out.
9. **Try it** — a short checklist to exercise the change on a dev/staging
   environment, if there is one.
10. **Open findings** — from review-panel, if any, with the user's options.

## 3. Publish

Present the brief as a page the user can scan:

- If an artifact/page publishing tool is available, follow its design guidance
  (load any design skill it requires first), build a single HTML page with the
  sections above — risk tags as coloured labels, the decisions table, links to
  the PR's files — and publish it privately. Give the user the link.
- Otherwise write the HTML to the scratch directory and tell the user the path,
  and paste the verdict and "Read these" sections into the reply.

Keep the reply itself short: the verdict, the link, and the one question (if
any) the user needs to answer.

## For a stack of PRs

One page with a section per PR in merge order, plus a top strip showing the
order, what each depends on, and any deploy-order constraint (e.g. "A must be
live before running B's backfill").
