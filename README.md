# lazygh

lazygit風の操作系で、現在のGitHubレポジトリに対するPRレビューを高速に行うTUI。
`gh` CLIを通してGitHub APIにアクセスする。

```
┌─────────────────────────────────────────────────────────┐
│ lazygh  repo: owner/repo                                │
├──────────────────────────────┬──────────────────────────┤
│ PRs (Open/Draft)             │ Detail                   │
│> PR #123 Fix parser edge case│                          │
│  PR #118 Improve CI cache    │  Title / Body            │
│  PR #116 Add auth retry      │                          │
│                              │                          │
├──────────────────────────────┴──────────────────────────┤
│ [q]Quit [j/k]Move [enter]Reload detail                 │
└─────────────────────────────────────────────────────────┘
```

## 前提条件

- Go 1.21+
- [gh CLI](https://cli.github.com/) インストール済み
- `gh auth login` 済み
- Git管理されたレポジトリ配下で実行

## インストール

```sh
go install github.com/rin2yh/lazygh@latest
```

ソースからビルド:

```sh
git clone https://github.com/rin2yh/lazygh
cd lazygh
go build -o lazygh .
```

## 実行

```sh
lazygh
```

## 操作

| Key | Action |
|-----|--------|
| `q` | Quit |
| `j` / `↓` | PRを下へ移動 |
| `k` / `↑` | PRを上へ移動 |
| `Enter` | 選択PRの詳細を再取得 |

## 現在の実装範囲

- 起動時に現在レポジトリを解決
- PR一覧（Open状態）を表示
- PR詳細（タイトル・本文）を表示

ロードマップ上のPhase1未達項目（Diff表示 / コメント投稿 / `h/l/c` 操作）は `docs/roadmap.md` の「Phase1 未達項目（別タスク管理）」に記載。

## ローカル実ghスモーク

実gh（モックなし）での最小確認:

```sh
gh auth status
go build -o lazygh .
./lazygh
```

確認観点:

- 起動時に現在レポジトリのPR一覧が表示される
- `Enter` で選択PRの詳細を再取得できる
- `q` で終了できる
