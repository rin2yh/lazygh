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

## 現在の実装範囲

- 起動時に現在レポジトリを解決
- PR一覧（Open状態）を表示
- PR概要（タイトル・ステータス・アサイン）を表示
- PR差分（Unified diff）を表示
- Diffの変更ファイルツリー（変更種別・`+/-` 行数）を表示
- Diff行ごとのコメントと複数行範囲コメントを pending review に追加
- レビュー送信前の概要入力と下部レビュー下書きドロワー

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
