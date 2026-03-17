package prs

import "github.com/rin2yh/lazygh/internal/model"

// ListState holds PR list, selection, and filter state.
type ListState struct {
	Repo         string
	PRs          []model.Item
	PRsLoading   bool
	PRsSelected  int
	Filter       model.PRFilterMask
	FilterOpen   bool
	FilterCursor int
}
