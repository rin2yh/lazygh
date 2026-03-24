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

// NewState returns an initialized overview State with default mode.
func NewState() State {
	return State{
		Mode: DetailModeOverview,
	}
}

// State holds overview panel display and fetching state.
type State struct {
	Mode     DetailMode
	Content  string
	Fetching FetchKind
}

// SetMode sets the detail display mode.
func (s *State) SetMode(m DetailMode) { s.Mode = m }

// SetContent sets the overview content string.
func (s *State) SetContent(content string) { s.Content = content }

// SetFetching sets the fetch kind.
func (s *State) SetFetching(k FetchKind) { s.Fetching = k }
