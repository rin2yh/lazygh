# lazygh ロードマップ

## ナビゲーションフロー（v1）

```
現在レポジトリ解決
  └→ PR一覧（Open/Draft）
      └→ PR詳細（Overview）
          └→ Diff表示（Unified）
              └→ PR全体コメント投稿
```

---

## Phase 1: MVP（現在レポ専用PRレビュー）

| 機能 | 詳細 |
|------|------|
| 現在レポ解決 | 起動時に対象レポを1つ決定 |
| PR一覧 | Open + Draft のPRを表示（番号・タイトル・作者など） |
| PR詳細 | タイトル・本文を表示 |
| PR差分 | `gh pr diff` ベースのUnified diff表示 |
| PRコメント投稿 | PR全体コメントの投稿 |
| 操作体系 | `h/j/k/l` + `Enter` + `c` + `q` |

### Phase1 未達項目（別タスク管理）

#### 1) Unified Diff表示

- 要件: `gh pr diff` ベースで選択PRのUnified diffを表示できる
- 受け入れ条件:
  - `h/l` で Overview / Diff を切り替えできる
  - Diffモード時に選択PRの `gh pr diff` 結果を詳細ペインに表示できる
  - エラー時は詳細ペインにエラー内容を表示する

#### 2) PR全体コメント投稿

- 要件: 選択PRへ全体コメントを投稿できる
- 受け入れ条件:
  - `c` でコメント入力フローを開始できる
  - 投稿成功/失敗がUI上で判別できる
  - 投稿先は現在選択中のPRである

#### 実施順

1. Unified Diff表示
2. PR全体コメント投稿

---

## Phase 2: レビュー強化

| 機能 | 詳細 |
|------|------|
| 行コメント | diff行単位のコメント投稿 |
| レビュー判定 | Approve / Request changes / Comment |
| フィルタ拡張 | Closed / Merged の表示切替 |

---

## Phase 3: 拡張機能

| 機能 | 詳細 |
|------|------|
| 複数レポ切替 | レポ一覧から対象レポを切替 |
| GitHub Actions | ワークフロー一覧・ログ表示 |
| 通知 | GitHub通知の閲覧・既読管理 |
