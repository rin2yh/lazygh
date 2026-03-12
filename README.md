# lazygh

lazygit風の操作系で、現在のGitHubレポジトリに対するPRレビューを高速に行うTUI。
`gh` CLIを通してGitHub APIにアクセスする。

```
┌─────────────────────────────────────────────────────────┐
│ lazygh  repo: owner/repo                                │
├──────────────────────────────┬──────────────────────────┤
│ PRs (Open/Draft)             │ Detail (Overview / Diff) │
│> #123 Fix parser edge case   │                          │
│  #118 Improve CI cache       │  Title / Body            │
│  #116 Add auth retry         │  or                       │
│                              │  Unified diff            │
│                              │                          │
├──────────────────────────────┴──────────────────────────┤
│ [q]Quit [j/k]Move [h/l]Mode [enter]Reload [c]Comment   │
└─────────────────────────────────────────────────────────┘
```

## 前提条件

- Go 1.21+
- [gh CLI](https://cli.github.com/) インストール済み
- `gh auth login` 済み
- Git管理されたレポジトリ配下で実行

## インストール

```sh
go install github.com/rin2yh/lazygh/cmd/lazygh@latest
```

ソースからビルド:

```sh
git clone https://github.com/rin2yh/lazygh
cd lazygh
go build -o lazygh ./cmd/lazygh
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
| `h` | Overview表示に切替 |
| `l` | Diff表示に切替 |
| `Enter` | 選択PRの詳細を再取得 |
| `c` | PR全体コメント入力を開始 |

## v1スコープ

- 対象は現在のレポジトリのみ
- PR一覧は Open + Draft のみ
- 詳細は Overview と PR全体 Unified diff を表示
- 書き込み操作は PR全体コメント投稿のみ

## ローカル実ghスモーク

実gh（モックなし）での最小確認:

```sh
gh auth status
go build -o lazygh ./cmd/lazygh
./lazygh
```

確認観点:

- 起動時に現在レポジトリのPR一覧が表示される
- `h` / `l` で詳細モードを切り替えられる
- `q` で終了できる
