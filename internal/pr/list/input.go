package list

// FilterKeyResult represents the outcome of handling a filter key press.
type FilterKeyResult int

const (
	FilterKeyNoop    FilterKeyResult = iota
	FilterKeyHandled                 // state updated, no further action needed
	FilterKeyApply                   // filter applied, begin PR fetch
)

// HandleFilterKey processes a key event while the filter panel is open.
// It returns a FilterKeyResult indicating whether a PR fetch is needed.
func (ls *ListState) HandleFilterKey(key string) FilterKeyResult {
	switch key {
	case "esc":
		ls.CloseFilterSelect()
		return FilterKeyHandled
	case "enter":
		ls.CloseFilterSelect()
		return FilterKeyApply
	case "j", "down":
		ls.MoveFilterCursor(1)
		return FilterKeyHandled
	case "k", "up":
		ls.MoveFilterCursor(-1)
		return FilterKeyHandled
	case " ":
		ls.ToggleFilterAtCursor()
		return FilterKeyHandled
	}
	return FilterKeyNoop
}

// OpenFilterSelect opens the filter selection panel.
func (ls *ListState) OpenFilterSelect() {
	ls.FilterOpen = true
	ls.FilterCursor = 0
}

// CloseFilterSelect closes the filter selection panel.
func (ls *ListState) CloseFilterSelect() {
	ls.FilterOpen = false
}

// MoveFilterCursor moves the filter cursor by dir steps (wraps around).
func (ls *ListState) MoveFilterCursor(dir int) {
	n := len(PRFilterOptions)
	ls.FilterCursor = (ls.FilterCursor + dir + n) % n
}

// ToggleFilterAtCursor toggles the filter option under the cursor.
// At least one option must remain selected.
func (ls *ListState) ToggleFilterAtCursor() {
	opt := PRFilterOptions[ls.FilterCursor]
	next := ls.Filter.Toggle(opt)
	if next == 0 {
		return
	}
	ls.Filter = next
}
