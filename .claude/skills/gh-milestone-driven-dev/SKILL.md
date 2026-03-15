---
name: gh-milestone-driven-dev
description: Execute milestone-driven development with GitHub Issues and PR-adjacent workflows. Use when user needs to pick issues from a GitHub milestone, inspect local code, implement fixes or refactors, run project verification commands, commit the result, post the commit ID back to the issue, and close the issue.
---

# Gh Milestone Driven Dev

Execute milestone work as a repeatable issue lifecycle, not as isolated edits.

## Workflow

1. Identify the target milestone and list its open issues with title, priority, area, and size.
2. Pick the issue that is most urgent, riskiest, or explicitly described as broken before taking lower-value work.
3. Read the issue body and the affected local code before proposing structure changes.
4. Implement the smallest change that fully resolves the issue, preserving existing design constraints in the repo.
5. Run the repository's required verification commands before finishing.
6. Commit with a focused message.
7. Comment on the issue with the commit ID and then close the issue.

## Operating Rules

- Prefer milestone issues whose title or labels imply breakage, tech debt, or blocked progress.
- If the user's wording maps to a specific issue title, confirm that issue first instead of guessing broadly.
- Treat issue resolution as incomplete until code changes, verification, commit, and issue closure are all done.
- If the worktree is dirty, avoid reverting unrelated changes; stage only the files relevant to the issue.
- If a refactor suggests a follow-up issue, finish the current issue first and then surface the follow-up separately.

## Verification

- Read project-local instructions such as `AGENTS.md`, `CLAUDE.md`, or equivalent before editing.
- Run the repo's required format, lint/static analysis, test, and CI-adjacent commands before completion.
- If a command is blocked by sandbox or cache permissions, rerun it with the required escalation instead of skipping silently.
- If verification cannot run, state exactly what was blocked and why.

## Issue Closure

- Post a short issue comment that includes the commit ID.
- Close the issue only after the commit succeeds.
- Keep the issue comment factual and minimal.

## Output Style

- Give short progress updates while working.
- In the final response, report the resolved issue, the commit ID, and the verification commands that passed.
