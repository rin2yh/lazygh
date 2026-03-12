# lazygh

lazygit風の操作系で、現在のGitHubレポジトリに対するPRレビューを高速に行うTUI。
`gh` CLIを通してGitHub APIにアクセスする。

## 前提条件

- Go 1.25+
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

### 全体

| Key | Action |
|-----|--------|
| `q` / `Ctrl+C` | 終了 |
| `d` | Diff表示へ切替（選択PRのdiffを取得） |
| `o` | Overview表示へ切替 |
| `Enter` | 現在フォーカス中PRの内容を再取得 |
| `Esc` | フォーカスをPR一覧へ戻す |
| `Tab` | Diff時にフォーカスを `PRs → Files → Diff` で循環 |

### パネル別（Diff時）

| Panel | Key | Action |
|-------|-----|--------|
| `PRs` | `j` / `k` / `↑` / `↓` | PR選択を移動 |
| `PRs` | `l` | Overviewに戻して内容を再取得 |
| `Files` | `j` / `k` / `↑` / `↓` | 変更ファイル選択を移動 |
| `Files` | `l` | Diff本文パネルへフォーカス |
| `Diff` | `j` / `k` / `↑` / `↓` | 1行スクロール |
| `Diff` | `Space` / `b` | 1ページ下/上へスクロール |
| `Diff` | `g` / `G` | 先頭/末尾へ移動 |
| `Diff` | `h` | Filesパネルへフォーカス |

## 現在の実装範囲

- 起動時に現在レポジトリを解決
- PR一覧（Open状態）を表示
- PR概要（タイトル・ステータス・アサイン）を表示
- PR差分（Unified diff）を表示
- Diffの変更ファイルツリー（変更種別・`+/-` 行数）を表示

ロードマップ上のPhase1未達項目（コメント投稿 / `c` 操作）は `docs/roadmap.md` の「Phase1 未達項目（別タスク管理）」に記載。

## ローカル実ghスモーク

実gh（モックなし）での最小確認:

```sh
gh auth status
go build -o lazygh .
./lazygh
```

確認観点:

- 起動時に現在レポジトリのPR一覧が表示される
- `o/d` で Overview / Diff を切り替えできる
- `Enter` で現在モードの詳細を再取得できる
- `q` で終了できる
