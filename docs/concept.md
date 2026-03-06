# lazygh コンセプト

## Why lazygh

lazydockerのようなスタンドアロンTUIアプリとして、GitHub操作をターミナルから快適に行えるツール。

### 既存ツールとの差別化

| 観点 | gh CLI | lazydocker的TUI | lazygh |
|------|--------|-----------------|--------|
| 操作性 | コマンド記憶が必要 | 直感的なナビゲーション | 直感的なナビゲーション |
| 並行処理 | 単発リクエスト | 対象依存 | goroutineで並列API呼び出し |
| 認証 | 自前実装 | 自前実装 | gh CLIに委譲 |
| エディタ統合 | なし | なし | 将来対応 |

## アーキテクチャ概要

```
lazygh (Go TUI)
    ↕ exec
gh CLI
    ↕
GitHub API
```

### 各レイヤーの責務

- **lazygh (Go TUI)**: UI描画・キーバインド・データ表示
- **gh CLI**: GitHub APIへのアクセス・認証管理
- **GitHub API**: データソース

参考UI: [lazydocker](https://github.com/jesseduffield/lazydocker)

## 設計方針

### gh CLI依存
GitHub API認証・トークン管理はgh CLIに委譲する。
独自のOAuth実装は持たず、`gh auth`で設定済みの認証情報を利用する。

### 読み取り優先
Phase 1は読み取り専用。まずナビゲーションの体験を磨き、書き込み操作は後から追加する。

### レイアウト

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

- 左カラム幅: 30%、右: 70%
- 左上: Repos、左下: Items
- 下部: ステータスバー
