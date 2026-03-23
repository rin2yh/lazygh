// Package app is the BubbleTea adapter layer for lazygh.
//
// インターフェースは利用側で定義する Go 慣習に従い interfaces.go にまとめている。
// review パッケージ側が必要とするインターフェース（Selection, AppState,
// PendingReviewClient）は循環インポートを避けるため internal/pr/review/interfaces.go
// に定義されている。
package app
