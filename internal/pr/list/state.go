package list

import "github.com/rin2yh/lazygh/internal/pr"

// NewState returns an initialized list State with default filter and empty items.
func NewState() State {
	return State{
		items:  []pr.Item{},
		Filter: PRFilterOpen,
	}
}

// State holds PR list, selection, and filter state.
type State struct {
	repo     string
	items    []pr.Item
	loading  bool
	selected int

	Filter       PRFilterMask
	FilterOpen   bool
	FilterCursor int
}

func (ls *State) Repo() string     { return ls.repo }
func (ls *State) Items() []pr.Item { return ls.items }
func (ls *State) Selected() int    { return ls.selected }

// IsFetching reports whether a PR fetch is in progress.
func (ls *State) IsFetching() bool { return ls.loading }

// Load replaces the PR list with the given repo and items, and resets the selection to 0.
func (ls *State) Load(repo string, items []pr.Item) {
	ls.repo = repo
	ls.items = items
	ls.selected = 0
}

func (ls *State) StartLoading() { ls.loading = true }
func (ls *State) StopLoading()  { ls.loading = false }

// SelectedOverview returns the formatted overview string for the currently
// selected PR. Returns false if no item is selected.
func (ls *State) SelectedOverview() (string, bool) {
	if ls.selected < 0 || ls.selected >= len(ls.items) {
		return "", false
	}
	return formatOverview(ls.items[ls.selected]), true
}

// NavigateDown moves the selection down by one. Returns true if selection changed.
func (ls *State) NavigateDown() bool {
	if ls.selected >= len(ls.items)-1 {
		return false
	}
	ls.selected++
	return true
}

// NavigateUp moves the selection up by one. Returns true if selection changed.
func (ls *State) NavigateUp() bool {
	if ls.selected <= 0 {
		return false
	}
	ls.selected--
	return true
}
