package gui

// navigateDown moves selection down, blocking if a pending review exists for
// the current PR. Returns true if the selection changed.
func (gui *Gui) navigateDown() bool {
	if gui.blocksPRSelectionChange() {
		gui.review.SetNotice("Pending review exists. Submit with S or discard with X.")
		return false
	}
	return gui.state.NavigateDown()
}

// navigateUp moves selection up, blocking if a pending review exists for the
// current PR. Returns true if the selection changed.
func (gui *Gui) navigateUp() bool {
	if gui.blocksPRSelectionChange() {
		gui.review.SetNotice("Pending review exists. Submit with S or discard with X.")
		return false
	}
	return gui.state.NavigateUp()
}

// blocksPRSelectionChange returns true when a pending review is open for the
// currently selected PR, preventing accidental navigation away.
func (gui *Gui) blocksPRSelectionChange() bool {
	item, ok := gui.state.SelectedPR()
	if !ok {
		return false
	}
	return gui.review.HasPendingReview() && gui.review.PRNumber() == item.Number
}
