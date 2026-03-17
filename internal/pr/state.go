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
