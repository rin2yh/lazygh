---
name: tech-debt-to-issues
description: コードベースの技術負債を調査し、GitHub Issues を作成してマイルストーンとラベルを付与する。「技術負債を調査」「負債をissueにして」「tech debtをトリアージ」がトリガー。
---

# Tech Debt to Issues

## Overview

コードベースを静的調査し、発見した技術負債を GitHub Issues として登録する。
マイルストーンへのアサインとラベル付けまで一括で行う。

## When to Use

- 「技術負債を調査して issue にして」
- 「tech debt をマイルストーンに追加して」
- リファクタリングフェーズの開始前にバックログを整理したいとき

## Instructions

### Step 1: リポジトリ情報を取得

```bash
gh api repos/{owner}/{repo}/milestones
gh api repos/{owner}/{repo}/labels
```

対象マイルストーンの **number**（整数）を控える。`--milestone` フラグはタイトル指定が効かないことがあるため、必ず number を使う。

### Step 2: 技術負債の調査

Explore エージェントを使ってコードベースを調査する。以下の観点を確認：

1. TODO / FIXME コメント
2. 重複コード・コピペ
3. 複雑すぎる関数（行数・サイクロマティック複雑度）
4. テスト不足な箇所
5. エラーハンドリングの漏れ（戻り値の `_` 無視など）
6. マジックナンバー
7. パッケージ間の密結合

### Step 3: Issues を作成

`--milestone` フラグは **使わず**、まず issue を作成する。

```bash
gh issue create \
  --repo {owner}/{repo} \
  --title "{prefix}: {タイトル}" \
  --body "..."
```

作成後、API で milestone をアサインする：

```bash
gh api repos/{owner}/{repo}/issues/{number} -X PATCH -f milestone={milestone_number}
```

### Step 4: ラベルを付与

複数ラベルは `-f "labels[]="` を繰り返す：

```bash
gh api repos/{owner}/{repo}/issues/{number}/labels \
  -X POST \
  -f "labels[]=tech-debt" \
  -f "labels[]=area: gui" \
  -f "labels[]=priority: m" \
  -f "labels[]=size: s"
```

複数 issue へのラベル付けはバックグラウンド並列で実行し、最後に `wait` する。

### Step 5: 結果を報告

作成した issue 一覧を表形式で報告する。

## Label Selection Guide

| ラベル種別 | 選び方 |
|-----------|--------|
| `tech-debt` | 技術負債 issue には必ず付与 |
| `area: *` | 影響パッケージに合わせて選択（gui / core / gh / test / ci / config） |
| `priority: xl/l/m/s/xs` | UI 安定性・バグに直結 → l、保守性改善 → m、スタイル → s |
| `size: xl/l/m/s/xs` | 変更規模の見積もりに応じて選択 |

## Examples

### 入力
「現状の技術負債を調査し、フェーズ1.5のマイルストーンに追加して」

### 出力

| # | タイトル | ラベル |
|---|---------|--------|
| #36 | refactor: draw.go を責務ごとに分割する | `tech-debt` `area: gui` `priority: l` `size: l` |
| #37 | fix: bubbles Update() 戻り値エラーを無視している箇所を修正 | `tech-debt` `area: gui` `priority: l` `size: s` |

## Guidelines

- `gh issue create --milestone` はマイルストーン名でなく番号が必要なケースがある。必ず API で PATCH する
- issue は1つの問題・1つの対応方針にとどめる（複数問題を1 issue にまとめない）
- 優先度が判断できない場合はユーザーに確認する
- 既存 issue と重複しないか、作成前に `gh issue list` で確認することを推奨
