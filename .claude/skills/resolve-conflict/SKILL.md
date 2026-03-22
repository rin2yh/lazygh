---
name: resolve-conflict
description: Resolve git merge conflicts on a feature branch. Use when the user reports conflicts, a PR shows conflicts, or a push is rejected because the branch diverged from main. Detects conflicting files, understands both sides, resolves each file, verifies, and pushes.
---

# Resolve Conflict

Resolve merge conflicts by merging the base branch into the feature branch (not rebasing), then push.

## Why Merge Instead of Rebase

Force push (`git push --force`) is blocked by branch protection rules on this repository.
Rebase rewrites history and requires force push — use `git merge origin/main` instead.

## Workflow

1. **状況把握**: `git status` と `git log --oneline -5` でブランチ状態を確認する。
2. **最新 main を取得**: `git fetch origin main` で最新状態にする。
3. **merge を試みる**: `git merge origin/main`。コンフリクトがなければ完了。
4. **コンフリクトファイルを特定**: `git status` でコンフリクトしているファイルを列挙する。
5. **両側の変更を把握**: 各コンフリクトファイルについて、HEAD 側（自分のブランチ）と origin/main 側（ベース）が「何を」「なぜ」変更したかを理解してから解消する。
6. **コンフリクトを解消**: コンフリクトマーカー（`<<<<<<<`, `=======`, `>>>>>>>`）を取り除き、両方の意図を正しく統合した内容にする。
7. **検証**: `go fmt ./...` → `go vet ./...` → `go test ./...` をすべて通す。
8. **merge コミット**: `git add <resolved-files>` → `git commit`（メッセージは自動生成される）。
9. **push**: `git push -u origin <branch-name>`。

## Resolving Each Conflict

コンフリクトマーカーの読み方:

```
<<<<<<< HEAD          ← 自分のブランチの内容（維持したい変更）
...our changes...
=======
...their changes...
>>>>>>> origin/main   ← main の内容（取り込む必要がある変更）
```

解消方針:
- **両方必要**: 両側の変更を手動でマージする（最も一般的）
- **自分側を採用**: main の変更が不要な場合（例: main に同じ変更が既にある、またはファイルを削除した）
- **main 側を採用**: 自分の変更が不要な場合

**このリポジトリで多いパターン**:
- 自分のブランチでファイルを分割・移動した場合: main に追加されたコードを分割先の新ファイルに反映する
- main で定数やヘルパーが追加された場合: 自分のブランチのインライン実装をその定数・ヘルパーに差し替える

## Verification

```sh
go fmt ./...
go vet ./...
go test ./...
```

すべて通過してから commit すること。

## Git Commands

```sh
# main の最新を取得
git fetch origin main

# merge（force push 不要なので rebase より安全）
git merge origin/main

# コンフリクト確認
git status

# 解消後にステージ
git add <file>

# merge コミット（--no-edit でデフォルトメッセージを使う）
git commit --no-edit

# push
git push -u origin <branch-name>
```

**禁止コマンド:**
- `git push --force` / `git push --force-with-lease` — ブランチ保護ルールにより拒否される
- `git rebase origin/main` — 履歴が書き換わり force push が必要になる

## Output Style

- コンフリクトしていたファイルと解消内容を簡潔に報告する。
- 最後に push が成功したことと、通過した検証コマンドを報告する。
