---
name: gh-milestone-driven-dev
description: Execute milestone-driven development with GitHub Issues and PR-adjacent workflows. Use when user needs to pick issues from a GitHub milestone, inspect local code, implement fixes or refactors, run project verification commands, commit the result, create a PR linked to the issue, and post the PR URL back to the issue.
---

# Gh Milestone Driven Dev

Execute milestone work as a repeatable issue lifecycle, not as isolated edits.

## GitHub CLI Commands

`gh milestone list` は存在しない。マイルストーン操作には以下のコマンドのみ使用すること。

```sh
# マイルストーン一覧（ローカルプロキシ環境では動作しない。代わりに gh issue list を使う）
# gh api repos/OWNER/REPO/milestones  ← 使用不可

# Issue一覧（マイルストーン指定なしで全件取得し、tech-debt ラベル等で絞る）
gh issue list --label "tech-debt" --json number,title,labels,body

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

## Workflow

1. `gh issue list --milestone "<name>"` でマイルストーンのオープン Issue を取得し、title・priority・area・size を確認する。
2. 最も緊急・リスクが高い、または壊れていると明示された Issue を選ぶ。
3. Issue 本文と関連コードを読んでから変更を始める。
4. Issue を解決する最小限の変更を実装する。
5. リポジトリの必須検証コマンドを実行する。
6. `/simplify` スキルを2回実行してコードをシンプルにする。
7. 集中したメッセージでコミットし、ブランチにプッシュする。
8. PR を作成し、本文に `Closes #<number>` で Issue を紐づける。
9. Issue に PR の URL をコメントする。Issue は自分でクローズしない（PR マージ時に自動クローズされる）。

## Operating Rules

- Prefer milestone issues whose title or labels imply breakage, tech debt, or blocked progress.
- If the user's wording maps to a specific issue title, confirm that issue first instead of guessing broadly.
- Treat issue resolution as incomplete until code changes, verification, commit, PR creation, and issue comment are all done.
- If the worktree is dirty, avoid reverting unrelated changes; stage only the files relevant to the issue.
- If a refactor suggests a follow-up issue, finish the current issue first and then surface the follow-up separately.

## Verification

- Read project-local instructions such as `AGENTS.md`, `CLAUDE.md`, or equivalent before editing.
- Run the repo's required format, lint/static analysis, test, and CI-adjacent commands before completion.
- If a command is blocked by sandbox or cache permissions, rerun it with the required escalation instead of skipping silently.
- If verification cannot run, state exactly what was blocked and why.

## PR Creation and Issue Linking

- PR 本文に `Closes #<number>` を含めて Issue を紐づける。
- PR マージ時に Issue が自動クローズされるため、自分で Issue をクローズしない。
- PR 作成後、Issue に PR の URL を短いコメントで通知する。
- Keep the issue comment factual and minimal.

### PR 作成コマンド

git remote がローカルプロキシを指しているため、`gh pr create` は `GH_HOST=github.com` を付与して実行すること。

```sh
GH_HOST=github.com gh pr create \
  --title "<title>" \
  --base main \
  --head "<branch-name>" \
  --body "$(cat <<'EOF'
## Summary
<bullet points>

## Test plan
<checklist>

Closes #<number>

https://claude.ai/code/session_0127NTVSSgRcBVbzR9vWHJYR
EOF
)"
```

**禁止コマンド:**
- `gh pr create` (GH_HOST なし) — ローカルプロキシを解決できず失敗する

## Output Style

- Give short progress updates while working.
- In the final response, report the resolved issue, the commit ID, the PR URL, and the verification commands that passed.
