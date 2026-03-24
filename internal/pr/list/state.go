package list

import "github.com/rin2yh/lazygh/internal/pr"

// ListState holds PR list, selection, and filter state.
type ListState struct {
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
func (ls *ListState) SelectedOverview() (string, bool) {
	if ls.Selected < 0 || ls.Selected >= len(ls.Items) {
		return "", false
	}
	return formatOverview(ls.Items[ls.Selected]), true
}

// NavigateDown moves the selection down by one. Returns true if selection changed.
func (ls *ListState) NavigateDown() bool {
	if ls.Selected >= len(ls.Items)-1 {
		return false
	}
	ls.Selected++
	return true
}

// NavigateUp moves the selection up by one. Returns true if selection changed.
func (ls *ListState) NavigateUp() bool {
	if ls.Selected <= 0 {
		return false
	}
	ls.Selected--
	return true
}
