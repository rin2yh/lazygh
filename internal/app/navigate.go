package app

const pendingReviewBlockNotice = "Pending review exists. Submit with S or discard with X."

func (gui *Gui) navigate(fn func() bool) bool {
	if gui.coord.BlocksPRSelectionChange() {
		gui.review.SetNotice(pendingReviewBlockNotice)
		return false
	}
	return fn()
}

// navigateDown moves selection down, blocking if a pending review exists for
// the current PR. Returns true if the selection changed.
func (gui *Gui) navigateDown() bool { return gui.navigate(gui.coord.NavigateDown) }

// navigateUp moves selection up, blocking if a pending review exists for the
// current PR. Returns true if the selection changed.
func (gui *Gui) navigateUp() bool { return gui.navigate(gui.coord.NavigateUp) }
