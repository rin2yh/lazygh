// Package detail は PR 詳細テキストを表示するスクロール可能なビューポートパッケージ。
//
// 分類: 機能固有（feature-specific）
//
// gui/ レイヤーのアダプタとして bubbletea（viewport.Model, tea.KeyMsg）への依存を担う。
// ドメイン状態（Mode, Content, Loading）は internal/pr/overview/ が所有する。
// gui/ との境界は gui.DetailViewport インターフェースで定義されている。
package detail
