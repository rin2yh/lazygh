# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## コマンド

```sh
# ビルド
go build -o lazygh .

# テスト（全パッケージ）
go test ./...

# 特定パッケージのテスト
go test ./internal/core/...

# フォーマット
go fmt ./...

# 静的解析
go vet ./...
```

## アーキテクチャ

**lazygh** は lazydocker 風の Go TUI アプリ。`gh` CLI 経由で GitHub API にアクセスする。

```
lazygh (Go TUI)  →  gh CLI  →  GitHub API
```

認証・トークン管理はすべて `gh auth` に委譲。独自 OAuth 実装なし。

### レイヤー構成

| パッケージ | 責務 |
|-----------|------|
| `.` (`main`) | エントリーポイント、設定・App初期化 |
| `internal/app` | `App` 構造体・`bubbletea` ベースの TUIアダプタ・描画（GUI含む） |
| `internal/app/layout` | TUI レイアウト計算・フォーカス状態・ステータスライン描画 |
| `internal/pr/diff` | PR diff の選択状態管理・カラー描画 |
| `internal/gh` | `gh` CLI の `exec` ラッパー。`ClientInterface` でモック可能 |
| `internal/core` | UI状態遷移とドメインロジック |
| `internal/config` | テーマ・キーバインド設定の読み込み |
| `pkg/gui/viewport` | 汎用 BubbleTea viewport ラッパー（ドメイン非依存） |

### テスト戦略

`gh` コマンドをモックするため `exec.Cmd` を差し替える仕組みを使用。各テストファイルに fake process のエントリーポイント関数を定義し、`execCommand` 差し替え経由で fake コマンドを注入する。
