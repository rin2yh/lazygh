// Package overview は PR 概要パネルのドメイン状態を所有するパッケージ。
package overview

// DetailMode represents the current detail panel display mode.
type DetailMode int

const (
	DetailModeOverview DetailMode = iota
	DetailModeDiff
)

// FetchKind represents the type of ongoing fetch operation.
type FetchKind int

const (
	FetchNone FetchKind = iota
	FetchingPRs
	FetchingDetail
	FetchingReview
)

// State holds overview panel display and fetching state.
type State struct {
	Mode     DetailMode
	Content  string
	Fetching FetchKind
}
