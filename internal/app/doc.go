// Package app は lazygh の BubbleTea アダプタ層である。
//
// 各ファイルの責務:
//   - app.go: App 構造体とエントリーポイント
//   - gui.go: Gui 構造体（UI状態の集約）と BubbleTea プログラム起動
//   - screen.go: BubbleTea Model の実装（Init / Update / View）
//   - render.go: 画面描画ロジック
//   - coordinator.go: アプリ全体の状態遷移と機能間協調
//   - fetch.go: 非同期データロードコマンドと結果処理
//   - input.go: キー入力ディスパッチ・ナビゲーション・フィルター・レビュー入力の一元管理
//   - focus.go: フォーカス状態管理
//
// ドメインパッケージ（review/、gh/、viewport/）がそれぞれのインターフェースを定義する。
// app → review の一方向インポートを維持して循環依存を避けている。
package app
