# lazygh コンセプト

## Why lazygh

### 既存プラグインとの差別化

既存のNeovim向けGitHubプラグイン（octo.nvim、gh.nvim）はすべてLua実装。
lazyghはGoバックエンドを採用することで以下の点で差別化する。

| 観点 | 既存プラグイン（Lua） | lazygh（Go + Lua） |
|------|----------------------|-------------------|
| 並行処理 | Neovimのスレッド制約あり | goroutineで並列API呼び出し |
| JSON解析 | 速度・型安全性が低い | 静的型付きで高速・安全 |
| テスタビリティ | Neovim依存でテストが難しい | Goで独立したユニットテスト可能 |
| 保守性 | 大規模化でコードが複雑になりやすい | 責務分離が明確 |

## アーキテクチャ概要

```
Neovim (Lua layer)
    ↕ jobstart / vim.system / RPC
Go binary (lazygh-server)
    ↕ exec
gh CLI
```

### 各レイヤーの責務

- **Lua layer**: UI描画・キーバインド・Neovim API連携
- **Go binary**: gh CLIの実行・JSON解析・データ変換・並行処理
- **gh CLI**: GitHub APIへのアクセス（認証もgh任せ）

Go binaryはgh CLIのラッパーとして機能し、複数APIの並列呼び出しやレスポンスの加工を担う。
Lua側はUIに専念し、GoバイナリとはJSON over stdioで通信する。

## 設計方針

### gh CLI依存
GitHub API認証・トークン管理はgh CLIに委譲する。
独自のOAuth実装は持たず、`gh auth`で設定済みの認証情報を利用する。

### 読み取り優先
Phase 1は読み取り専用。まずナビゲーションの体験を磨き、書き込み操作は後から追加する。

### シンプルなUI
Neovimのバッファ・ウィンドウ・フローティングウィンドウを活用したシンプルなUI。
独自のウィジェットは最小限にとどめる。
