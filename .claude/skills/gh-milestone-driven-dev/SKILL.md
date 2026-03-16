---
name: gh-milestone-driven-dev
description: Execute milestone-driven development with GitHub Issues and PR-adjacent workflows. Use when user needs to pick issues from a GitHub milestone, inspect local code, implement fixes or refactors, run project verification commands, commit the result, post the commit ID back to the issue, and close the issue.
---

# Gh Milestone Driven Dev

Execute milestone work as a repeatable issue lifecycle, not as isolated edits.

## GitHub CLI Commands

`gh milestone list` は存在しない。マイルストーン操作には以下のコマンドのみ使用すること。

```sh
# マイルストーン一覧
gh api repos/{owner}/{repo}/milestones

# マイルストーン名でIssue一覧（ラベル・本文付き）
gh issue list --milestone "<milestone-name>" --json number,title,labels,body

# Issue にコメント
gh issue comment <number> --body "<message>"

# Issue をクローズ
gh issue close <number>

# Issue を作成（マイルストーン・ラベル付き）
gh issue create --title "<title>" --body "<body>" --milestone "<name>" --label "<label1>,<label2>"

# Issue にラベルを追加
gh issue edit <number> --add-label "<label1>,<label2>"
```

**禁止コマンド:**
- `gh milestone list` — 存在しない、使用禁止
- `gh api repos/{owner}/{repo}/milestones` の `{owner}/{repo}` リテラル展開 — 必ず実際のオーナー/リポジトリ名に置換すること

## Workflow

1. `gh issue list --milestone "<name>"` でマイルストーンのオープン Issue を取得し、title・priority・area・size を確認する。
2. 最も緊急・リスクが高い、または壊れていると明示された Issue を選ぶ。
3. Issue 本文と関連コードを読んでから変更を始める。
4. Issue を解決する最小限の変更を実装する。
5. リポジトリの必須検証コマンドを実行する。
6. 集中したメッセージでコミットする。
7. Issue にコミット ID をコメントしてからクローズする。

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
