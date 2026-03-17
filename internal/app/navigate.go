package app

const pendingReviewBlockNotice = "Pending review exists. Submit with S or discard with X."

// navigateDown moves selection down, blocking if a pending review exists for
// the current PR. Returns true if the selection changed.
func (gui *Gui) navigateDown() bool {
	if gui.coord.BlocksPRSelectionChange() {
		gui.review.SetNotice(pendingReviewBlockNotice)
		return false
	}
	return gui.coord.NavigateDown()
}

// navigateUp moves selection up, blocking if a pending review exists for the
// current PR. Returns true if the selection changed.
func (gui *Gui) navigateUp() bool {
	if gui.coord.BlocksPRSelectionChange() {
		gui.review.SetNotice(pendingReviewBlockNotice)
		return false
	}
	return gui.coord.NavigateUp()
}
