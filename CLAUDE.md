# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## コマンド

```sh
# ビルド
go build -o lazygh .

# テスト（全パッケージ）
go test ./...

# 特定パッケージのテスト
go test ./internal/app/...

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
| `internal/app` | `App` 構造体・`bubbletea` ベースの TUIアダプタ・描画・`Coordinator` による状態管理 |
| `internal/app/layout` | TUI レイアウト計算・フォーカス状態・ステータスライン描画 |
| `internal/model` | ドメインモデル定義（PR・アクション種別など） |
| `internal/gh` | `gh` CLI の `exec` ラッパー。`ClientInterface` でモック可能 |
| `internal/config` | テーマ・キーバインド設定の読み込み |
| `internal/pr/diff` | PR unified diff の解析・表示・選択状態管理 |
| `internal/pr/list` | PR一覧パネル描画 |
| `internal/pr/overview` | PR概要パネル描画 |
| `internal/pr/detail` | PRの詳細情報表示 |
| `internal/pr/review` | レビュー保留ワークフロー（コメント・サマリー・submit/discard）の状態管理 |
| `internal/pr/action` | PRアクション（マージ等）パネル |
| `internal/help` | ヘルプパネル描画 |
| `pkg/gui/viewport` | 汎用 BubbleTea viewport ラッパー（ドメイン非依存） |
| `pkg/gui/textarea` | テキスト入力ウィジェット |
| `pkg/gui/ansi` | ANSIエスケープ処理ユーティリティ |
| `pkg/gui/widget` | 汎用 UI ウィジェット |

### テスト戦略

`gh` コマンドをモックするため `exec.Cmd` を差し替える仕組みを使用。各テストファイルに fake process のエントリーポイント関数を定義し、`execCommand` 差し替え経由で fake コマンドを注入する。
