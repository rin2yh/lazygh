package gui

import guireview "github.com/rin2yh/lazygh/internal/gui/review"

type panelFocus int

const (
	panelRepo panelFocus = iota
	panelPRs
	panelDiffFiles
	panelDiffContent
	panelReviewDrawer
)

func (gui *Gui) switchToOverview() bool {
	changed := gui.state.SwitchToOverview()
	if changed {
		gui.focus = panelPRs
	}
	return changed
}

func (gui *Gui) focusPRs() {
	gui.focus = panelPRs
}

func (gui *Gui) switchToDiff() bool {
	changed := gui.state.SwitchToDiff()
	if changed {
		gui.focus = panelDiffFiles
		gui.diffFiles = nil
		gui.diffFileSelected = 0
		gui.diffLineSelected = 0
	}
	return changed
}

func (gui *Gui) cycleFocus() {
	if !gui.state.IsDiffMode() {
		gui.focus = panelPRs
		return
	}

	order := gui.focusOrder()
	if len(order) == 0 {
		gui.focus = panelPRs
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

func (gui *Gui) focusOrder() []panelFocus {
	order := []panelFocus{panelRepo, panelPRs}
	if len(gui.diffFiles) > 0 {
		order = append(order, panelDiffFiles)
	}
	order = append(order, panelDiffContent)
	if gui.review.ShouldShowDrawer() {
		order = append(order, panelReviewDrawer)
	}
	return order
}

func (gui *Gui) setReviewFocus(target guireview.FocusTarget) {
	switch target {
	case guireview.FocusReviewDrawer:
		gui.focus = panelReviewDrawer
	default:
		gui.focus = panelDiffContent
	}
}
