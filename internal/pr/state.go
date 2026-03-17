package pr

import "github.com/rin2yh/lazygh/internal/model"

// ListState holds PR list, selection, and filter state.
type ListState struct {
	Repo         string
	Items        []model.Item
	Fetching     bool
	Selected     int
	Filter       model.PRFilterMask
	FilterOpen   bool
	FilterCursor int
}

// SelectedOverview returns the formatted overview string for the currently
// selected PR, or empty string if no item is selected.
func (ls *ListState) SelectedOverview() string {
	if ls.Selected < 0 || ls.Selected >= len(ls.Items) {
		return ""
	}
	return formatOverview(ls.Items[ls.Selected])
}
