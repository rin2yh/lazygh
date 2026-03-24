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
	return State{mode: DetailModeOverview}
}

// State holds overview panel display and fetching state.
type State struct {
	mode     DetailMode
	content  string
	fetching FetchKind
}

// Mode returns the current detail display mode.
func (s *State) Mode() DetailMode { return s.mode }

// Content returns the current overview content string.
func (s *State) Content() string { return s.content }

// FetchKind returns the kind of fetch currently in progress.
func (s *State) FetchKind() FetchKind { return s.fetching }

// IsFetching reports whether any fetch is in progress.
func (s *State) IsFetching() bool { return s.fetching != FetchNone }

// StartFetching marks that a fetch of kind k is in progress.
func (s *State) StartFetching(k FetchKind) { s.fetching = k }

// StopFetching marks the current fetch as complete.
func (s *State) StopFetching() { s.fetching = FetchNone }

// ShowContent updates the displayed content without affecting fetch state.
// Use for preview updates (e.g. selection change).
func (s *State) ShowContent(c string) { s.content = c }

// LoadResult sets content from a completed fetch and clears fetching state.
func (s *State) LoadResult(c string) { s.fetching = FetchNone; s.content = c }

// EnterOverviewMode switches to overview mode and clears fetching state.
func (s *State) EnterOverviewMode() { s.mode = DetailModeOverview; s.fetching = FetchNone }

// EnterDiffMode switches to diff mode and clears fetching state.
func (s *State) EnterDiffMode() { s.mode = DetailModeDiff; s.fetching = FetchNone }
