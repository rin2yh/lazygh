# lazygh

lazydocker風TUI for GitHub。gh CLI経由でGitHub APIにアクセスし、ターミナルからGitHub操作を快適に行えるツール。

```
┌─────────────────────────────────────────────────────────┐
│ lazygh                                                   │
├──────────────┬──────────────────────────────────────────┤
│ Repositories │                                           │
│              │  Main Content                             │
│  owner/repo1 │                                           │
│> owner/repo2 │  (PR一覧 / Issue一覧 / 詳細 / diff)       │
│              │                                           │
├──────────────┤                                           │
│ Items        │                                           │
│              │                                           │
│> PR #123 ... │                                           │
│  Issue #45   │                                           │
│              │                                           │
├──────────────┴──────────────────────────────────────────┤
│ [q]Quit  [tab]Panel  [j/k]Navigate  [enter]Select       │
└─────────────────────────────────────────────────────────┘
```

## 前提条件

- Go 1.21+
- [gh CLI](https://cli.github.com/) インストール済み・`gh auth login`で認証済み

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

## キーバインド

| Key | Action |
|-----|--------|
| `q` | Quit |
| `Tab` | Next panel |
| `Shift+Tab` | Prev panel |
| `j` / `↓` | Navigate down |
| `k` / `↑` | Navigate up |
| `Enter` | Select |

## ローカル実ghスモーク

実gh（モックなし）をローカルで軽く確認する最小手順:

```sh
gh auth status
go build -o lazygh ./cmd/lazygh
./lazygh
```

確認観点:

- 起動時にRepositoriesが表示される
- Repo選択でIssues/PRsが表示される
- `q` で終了できる
