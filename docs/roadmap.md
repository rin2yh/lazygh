# lazygh ロードマップ

## ナビゲーションフロー（v1）

```
現在レポジトリ解決
  └→ PR一覧（Open/Draft）
      └→ PR詳細（Overview）
          └→ Diff表示（Unified）
              └→ pending review 作成
                  ├→ 行コメント / 範囲コメント追加
                  ├→ レビュー概要入力
                  └→ Comment として submit / discard
```

---

## Phase 1: MVP（現在レポ専用PRレビュー）

| 機能 | 詳細 |
|------|------|
| 現在レポ解決 | 起動時に対象レポを1つ決定 |
| PR一覧 | Open + Draft のPRを表示（番号・タイトル・作者など） |
| PR詳細 | タイトル・本文を表示 |
| PR差分 | `gh pr diff` ベースのUnified diff表示 |
| レビュー下書き | pending review の作成・追記・送信・破棄 |
| レビュー入力 | 行コメント、範囲コメント、レビュー概要入力 |
| 操作体系 | `h/j/k/l` + `Enter` + `v` + `c` + `R` + `S` + `X` + `q` |

### Phase1 対応済み

- Unified Diff表示（`h/l` 切替、Diff表示、エラー表示）: 2026-03-12
- pending review ベースの行コメント / 範囲コメント / レビュー概要入力: 2026-03-14
- review drawer と submit / discard フロー、入力中ショートカット整理: 2026-03-14

---

## Phase 2: レビュー強化

| 機能 | 詳細 |
|------|------|
| レビュー判定 | Approve / Request changes / Comment |
| 既存コメント表示 | 既存 review thread の表示、返信、解決状態の確認 |
| コメント編集系 | pending review comment の編集 / 削除 / 並び替え |
| review UX | anchor表示改善、ショートカットヘルプ、衝突時の回復導線 |
| フィルタ拡張 | Closed / Merged の表示切替 |

---

## Phase 3: 拡張機能

| 機能 | 詳細 |
|------|------|
| 複数レポ切替 | レポ一覧から対象レポを切替 |
| GitHub Actions | ワークフロー一覧・ログ表示 |
| 通知 | GitHub通知の閲覧・既読管理 |
