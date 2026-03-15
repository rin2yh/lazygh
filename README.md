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

## タスク管理
GitHub Milestones で管理する。
詳細はマイルストーン内のIssueに記載する。

