# lazygh コンセプト

## Why lazygh

lazygitのようなキーバインドで、現在レポジトリのPRレビューを最短導線で行うためのTUI。

### 既存ツールとの差別化

| 観点 | gh CLI | lazygit | lazygh |
|------|--------|---------|--------|
| 操作性 | コマンド記憶が必要 | 単キー中心 | 単キー中心（GitHub PR向け） |
| 対象 | GitHub全般 | Gitローカル操作 | 現在レポのPRレビュー |
| 認証 | `gh auth` | 不要 | `gh auth` に委譲 |
| 差分閲覧 | コマンド出力 | 高速 | PR Unified diffを中央表示 |

## アーキテクチャ概要

```
lazygh (Go TUI)
    ↕ exec
gh CLI
    ↕
GitHub API
```

### 各レイヤーの責務

- **lazygh (Go TUI)**: 画面描画、キー入力、状態遷移
- **gh CLI**: APIアクセスと認証
- **GitHub API**: PR/コメント/差分データ

## 設計方針

### 現在レポジトリ固定
起動時に現在レポジトリを解決し、そのレポのPRだけを扱う。
複数レポ切替はv2以降。

### PRレビュー特化
v1はPRレビューに必要な最小機能に絞る。
Issue導線はv1では提供しない。

### 操作はlazygit寄せ
`h/j/k/l` を中心に、選択と表示切替を高速に行える導線を優先する。

### レイアウト

```
┌─────────────────────────────────────────────────────────┐
│ lazygh  repo: owner/repo                                │
├──────────────────────────────┬──────────────────────────┤
│ PRs (Open/Draft)             │ Detail (Overview / Diff) │
│> #123 Fix parser edge case   │                          │
│  #118 Improve CI cache       │  PR本文 or Unified diff  │
│  #116 Add auth retry         │                          │
├──────────────────────────────┴──────────────────────────┤
│ [q]Quit [j/k]Move [h/l]Panel [enter]Reload [c]Comment  │
└─────────────────────────────────────────────────────────┘
```

- 左: PR一覧
- 中央: PR詳細（Overview / Diff）
- 下部: ステータス/キーガイド
