// Package app は lazygh の BubbleTea アダプタ層である。
//
// ドメインパッケージ（review/、gh/、viewport/）がそれぞれのインターフェースを定義する。
// app → review の一方向インポートを維持して循環依存を避けている。
package app
