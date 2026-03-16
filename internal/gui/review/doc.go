// Package review は PR レビュー（コメント・サマリー・レンジ選択）の機能固有パッケージ。
//
// 分類: 機能固有（feature-specific）
//
// 現在は gui/ 配下に置かれているが、後続 issue で internal/review/ へ昇格予定。
// 昇格時は bubbletea（tea.Cmd, tea.KeyMsg）への依存を取り除き、
// gui/ レイヤーのアダプタがフレームワーク変換を担う構成に変更する。
// gui/ との境界は gui.ReviewController インターフェースで定義されている。
package review
