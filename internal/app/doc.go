// Package app is the BubbleTea adapter layer for lazygh.
//
// ドメインパッケージ（review/、gh/、viewport/）がそれぞれのインターフェースを定義する。
// review.Selection / review.AppState / review.PendingReviewClient が review/ 側に置かれているのは、
// app → review の一方向インポートを維持して循環依存を避けるためである。
package app
