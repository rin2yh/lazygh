// Package detail は PR 詳細テキストを表示するスクロール可能なビューポートパッケージ。
//
// 分類: 機能固有（feature-specific）
//
// 現在は gui/ 配下に置かれているが、後続 issue で internal/detail/ へ昇格予定。
// 昇格時は bubbletea（viewport.Model, tea.KeyMsg）への依存を取り除き、
// gui/ レイヤーのアダプタがフレームワーク変換を担う構成に変更する。
// gui/ との境界は gui.DetailViewport インターフェースで定義されている。
package detail
