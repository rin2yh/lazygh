package list

import "github.com/rin2yh/lazygh/internal/pr"

// NewState returns an initialized list State with default filter and empty items.
func NewState() State {
	return State{
		Items:  []pr.Item{},
		Filter: PRFilterOpen,
	}
}

// State holds PR list, selection, and filter state.
type State struct {
	Repo         string
	Items        []pr.Item
	Fetching     bool
	Selected     int
	Filter       PRFilterMask
	FilterOpen   bool
	FilterCursor int
}

// SelectedOverview returns the formatted overview string for the currently
// selected PR. Returns false if no item is selected.
func (ls *State) SelectedOverview() (string, bool) {
	if ls.Selected < 0 || ls.Selected >= len(ls.Items) {
		return "", false
	}
	return formatOverview(ls.Items[ls.Selected]), true
}

// SetFetching sets the fetching state.
func (ls *State) SetFetching(v bool) { ls.Fetching = v }

// SetRepo sets the repository name.
func (ls *State) SetRepo(repo string) { ls.Repo = repo }

// SetItems replaces the PR item list.
func (ls *State) SetItems(items []pr.Item) { ls.Items = items }

// SetSelected sets the selected index.
func (ls *State) SetSelected(i int) { ls.Selected = i }

// NavigateDown moves the selection down by one. Returns true if selection changed.
func (ls *State) NavigateDown() bool {
	if ls.Selected >= len(ls.Items)-1 {
		return false
	}
	ls.Selected++
	return true
}

// NavigateUp moves the selection up by one. Returns true if selection changed.
func (ls *State) NavigateUp() bool {
	if ls.Selected <= 0 {
		return false
	}
	ls.Selected--
	return true
}
