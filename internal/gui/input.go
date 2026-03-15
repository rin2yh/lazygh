package gui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
)

func (s *screen) handleReviewInputKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	keys := s.gui.config.KeyBindings
	switch {
	case keys.Matches(msg, config.ActionReviewSubmit):
		return s.gui.review.HandleSubmit(), true
	case keys.Matches(msg, config.ActionReviewDiscard):
		return s.gui.review.HandleDiscard(), true
	case keys.Matches(msg, config.ActionReviewSave):
		if s.gui.state.Review.InputMode == core.ReviewInputComment {
			return s.gui.review.HandleCommentSave(), true
		}
	}
	if s.gui.review.HandleEditorKey(msg) {
		return nil, true
	}
	return nil, false
}

func (s *screen) handleKeyInput(msg tea.KeyMsg) tea.Cmd {
	keys := s.gui.config.KeyBindings
	switch {
	case keys.Matches(msg, config.ActionQuit):
		return tea.Quit
	case keys.Matches(msg, config.ActionCancel):
		return s.handleCancel()
	case keys.Matches(msg, config.ActionFocusNext):
		s.gui.cycleFocus()
		return nil
	case keys.Matches(msg, config.ActionMoveDown):
		return s.moveDown()
	case keys.Matches(msg, config.ActionMoveUp):
		return s.moveUp()
	case keys.Matches(msg, config.ActionPageDown),
		keys.Matches(msg, config.ActionPageUp),
		keys.Matches(msg, config.ActionGoTop),
		keys.Matches(msg, config.ActionGoBottom):
		s.gui.scrollDetailByKey(msg)
		s.gui.scrollOverviewByKey(msg)
		return nil
	case keys.Matches(msg, config.ActionPanelPrev):
		s.gui.moveFocus(-1)
		return nil
	case keys.Matches(msg, config.ActionPanelNext):
		s.gui.moveFocus(1)
		return nil
	case keys.Matches(msg, config.ActionShowOverview):
		s.gui.switchToOverview()
		return nil
	case keys.Matches(msg, config.ActionShowDiff):
		return s.showDiff()
	case keys.Matches(msg, config.ActionOpenSelected):
		return s.openSelectedPR()
	case keys.Matches(msg, config.ActionReviewRange):
		return s.startReviewRange()
	case keys.Matches(msg, config.ActionReviewComment):
		return s.startReviewComment()
	case keys.Matches(msg, config.ActionReviewSummary):
		return s.startReviewSummary()
	case keys.Matches(msg, config.ActionReviewSubmit):
		return s.gui.review.HandleSubmit()
	case keys.Matches(msg, config.ActionReviewDiscard):
		return s.gui.review.HandleDiscard()
	case keys.Matches(msg, config.ActionReviewClearComment):
		if s.gui.state.Review.InputMode == core.ReviewInputComment {
			s.gui.review.ClearCommentInput()
		}
		return nil
	default:
		return nil
	}
}

func (s *screen) handleCancel() tea.Cmd {
	if s.gui.state.Review.InputMode == core.ReviewInputNone && s.gui.state.Review.RangeStart != nil {
		s.gui.state.ClearReviewRangeStart()
		s.gui.state.SetReviewNotice("Range selection cleared.")
		s.gui.focus = panelDiffContent
		return nil
	}
	if s.gui.focus == panelReviewDrawer {
		s.gui.review.StopInput()
		s.gui.focus = panelDiffContent
		return nil
	}
	s.gui.focusPRs()
	return nil
}

func (s *screen) showDiff() tea.Cmd {
	if s.gui.switchToDiff() {
		return s.openSelectedPR()
	}
	return nil
}

func (s *screen) moveDown() tea.Cmd {
	if s.gui.state.IsDiffMode() {
		switch s.gui.focus {
		case panelPRs:
			if s.gui.navigateDown() {
				return s.openSelectedPR()
			}
			return nil
		case panelDiffFiles:
			s.gui.selectNextDiffFile()
			return nil
		case panelDiffContent:
			s.gui.scrollDetailDown()
			return nil
		case panelReviewDrawer:
			return nil
		}
		return nil
	}

	if s.gui.focus == panelPRs {
		s.gui.navigateDown()
	}
	return nil
}

func (s *screen) moveUp() tea.Cmd {
	if s.gui.state.IsDiffMode() {
		switch s.gui.focus {
		case panelPRs:
			if s.gui.navigateUp() {
				return s.openSelectedPR()
			}
			return nil
		case panelDiffFiles:
			s.gui.selectPrevDiffFile()
			return nil
		case panelDiffContent:
			s.gui.scrollDetailUp()
			return nil
		case panelReviewDrawer:
			return nil
		}
		return nil
	}

	if s.gui.focus == panelPRs {
		s.gui.navigateUp()
	}
	return nil
}

func (s *screen) openSelectedPR() tea.Cmd {
	action := s.gui.state.PlanEnter(s.gui.client != nil, os.Getenv("LAZYGH_DEBUG_DETAIL_TEXT"))
	switch action.Kind {
	case core.EnterLoadPRDiff:
		return s.loadDetailCmd(action.Repo, action.Number, core.DetailModeDiff)
	case core.EnterLoadPRDetail:
		return s.loadDetailCmd(action.Repo, action.Number, core.DetailModeOverview)
	default:
		return nil
	}
}

func (s *screen) startReviewRange() tea.Cmd {
	if !s.gui.state.IsDiffMode() {
		s.gui.state.SetReviewNotice("Review range selection is only available in diff view.")
		return nil
	}
	s.gui.review.ToggleRangeSelection()
	return nil
}

func (s *screen) startReviewComment() tea.Cmd {
	if !s.gui.state.IsDiffMode() {
		s.gui.state.SetReviewNotice("Review comments are only available in diff view.")
		return nil
	}
	s.gui.review.BeginCommentFlow()
	return nil
}

func (s *screen) startReviewSummary() tea.Cmd {
	if !s.gui.state.IsDiffMode() {
		s.gui.state.SetReviewNotice("Review summary is only available in diff view.")
		return nil
	}
	s.gui.review.BeginSummaryInput()
	return nil
}
