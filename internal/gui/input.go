package gui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
)

func (s *screen) handleReviewInputKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "S":
		return s.handleReviewSubmit(), true
	case "X":
		return s.handleReviewDiscard(), true
	case "ctrl+s":
		if s.gui.state.Review.InputMode == core.ReviewInputComment {
			return s.handleReviewCommentSave(), true
		}
	}
	if s.gui.handleReviewEditorKey(msg) {
		return nil, true
	}
	return nil, false
}

func (s *screen) handleKeyInput(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "esc":
		return s.handleCancel()
	case "tab":
		s.gui.cycleFocus()
		return nil
	case "j", "down":
		return s.moveDown()
	case "k", "up":
		return s.moveUp()
	case "pgdown", "f", " ", "pgup", "b", "home", "g", "end", "G":
		s.gui.scrollDetailByKey(msg)
		s.gui.scrollOverviewByKey(msg)
		return nil
	case "h":
		s.gui.moveFocus(-1)
		return nil
	case "l":
		s.gui.moveFocus(1)
		return nil
	case "o":
		s.gui.switchToOverview()
		return nil
	case "d":
		return s.showDiff()
	case "enter":
		return s.openSelectedPR()
	case "v":
		return s.startReviewRange()
	case "c":
		return s.startReviewComment()
	case "R":
		return s.startReviewSummary()
	case "S":
		return s.handleReviewSubmit()
	case "X":
		return s.handleReviewDiscard()
	case "x":
		if s.gui.state.Review.InputMode == core.ReviewInputComment {
			s.gui.commentEditor.SetValue("")
			s.gui.state.SetReviewNotice("Comment input cleared.")
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
		s.gui.stopReviewInput()
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
	s.gui.toggleReviewRangeSelection()
	return nil
}

func (s *screen) startReviewComment() tea.Cmd {
	if !s.gui.state.IsDiffMode() {
		s.gui.state.SetReviewNotice("Review comments are only available in diff view.")
		return nil
	}
	s.gui.beginReviewCommentFlow()
	return nil
}

func (s *screen) startReviewSummary() tea.Cmd {
	if !s.gui.state.IsDiffMode() {
		s.gui.state.SetReviewNotice("Review summary is only available in diff view.")
		return nil
	}
	s.gui.beginReviewSummaryInput()
	return nil
}
