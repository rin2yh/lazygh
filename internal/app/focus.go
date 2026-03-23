package app

import (
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/pr/review"
)

const pendingReviewBlockNotice = "Pending review exists. Submit with S or discard with X."

func (gui *Gui) navigate(fn func() bool) bool {
	if gui.coord.BlocksPRSelectionChange() {
		gui.review.SetNotice(pendingReviewBlockNotice)
		return false
	}
	return fn()
}

func (gui *Gui) navigateDown() bool { return gui.navigate(gui.coord.NavigateDown) }
func (gui *Gui) navigateUp() bool   { return gui.navigate(gui.coord.NavigateUp) }

func (gui *Gui) switchToOverview() bool {
	changed := gui.coord.SwitchToOverview()
	if changed {
		gui.focus = layout.FocusPRs
	}
	return changed
}

func (gui *Gui) focusPRs() {
	gui.focus = layout.FocusPRs
}

func (gui *Gui) switchToDiff() bool {
	changed := gui.coord.SwitchToDiff()
	if changed {
		gui.focus = layout.FocusDiffFiles
		gui.diff.Reset()
	}
	return changed
}

func (gui *Gui) cycleFocus() {
	if !gui.coord.IsDiffMode() {
		gui.focus = layout.FocusPRs
		return
	}

	order := gui.focusOrder()
	if len(order) == 0 {
		gui.focus = layout.FocusPRs
		return
	}
	for i, focus := range order {
		if focus == gui.focus {
			gui.focus = order[(i+1)%len(order)]
			return
		}
	}
	gui.focus = order[0]
}

func (gui *Gui) moveFocus(delta int) bool {
	if delta == 0 {
		return false
	}

	order := gui.focusOrder()
	if len(order) == 0 {
		return false
	}
	for i, focus := range order {
		if focus != gui.focus {
			continue
		}
		next := i + delta
		if next < 0 || next >= len(order) {
			return false
		}
		gui.focus = order[next]
		return true
	}
	return false
}

func (gui *Gui) focusOrder() []layout.Focus {
	order := []layout.Focus{layout.FocusRepo, layout.FocusPRs}
	if len(gui.diff.Files()) > 0 {
		order = append(order, layout.FocusDiffFiles)
	}
	order = append(order, layout.FocusDiffContent)
	if gui.review.ShouldShowDrawer() {
		order = append(order, layout.FocusReviewDrawer)
	}
	return order
}

// resetDiffFocusIfOnFiles moves focus off the files panel when there are no files to show.
func (gui *Gui) resetDiffFocusIfOnFiles() {
	if gui.focus == layout.FocusDiffFiles {
		gui.focus = layout.FocusDiffContent
	}
}

func (gui *Gui) setReviewFocus(target review.FocusTarget) {
	switch target {
	case review.FocusReviewDrawer:
		gui.focus = layout.FocusReviewDrawer
	default:
		gui.focus = layout.FocusDiffContent
	}
}
