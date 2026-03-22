package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/review"
)

func (s *screen) handleReviewInputKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	action, ok := s.gui.config.KeyBindings.ActionFor(msg)
	if ok {
		switch action {
		case config.ActionReviewSubmit:
			return s.gui.review.Submit(), true
		case config.ActionReviewDiscard:
			return s.gui.review.Discard(), true
		case config.ActionReviewSave:
			if s.gui.review.InputMode() == review.InputComment {
				if s.gui.review.IsEditingComment() {
					return s.gui.review.SaveEditComment(), true
				}
				return s.gui.review.SaveComment(), true
			}
		}
	}
	if cmd, handled := s.gui.review.EditorKey(msg); handled {
		return cmd, true
	}
	return nil, false
}

func (s *screen) handleReviewAction(action config.Action) tea.Cmd {
	switch action {
	case config.ActionReviewRange:
		return s.startReviewRange()
	case config.ActionReviewComment:
		return s.startReviewComment()
	case config.ActionReviewSummary:
		return s.startReviewSummary()
	case config.ActionReviewSubmit:
		return s.gui.review.Submit()
	case config.ActionReviewDiscard:
		return s.gui.review.Discard()
	case config.ActionReviewClearComment:
		if s.gui.review.InputMode() == review.InputComment {
			s.gui.review.ClearCommentInput()
		}
	case config.ActionReviewEvent:
		if s.gui.coord.IsDiffMode() {
			s.gui.review.CycleReviewEvent()
		}
	case config.ActionReviewDeleteComment:
		if s.gui.focus == layout.FocusReviewDrawer {
			return s.gui.review.DeleteComment()
		}
	case config.ActionReviewEditComment:
		if s.gui.focus == layout.FocusReviewDrawer {
			s.gui.review.EditComment()
		}
	}
	return nil
}

func (s *screen) requireDiffMode(notice string, fn func()) tea.Cmd {
	if !s.gui.coord.IsDiffMode() {
		s.gui.review.SetNotice(notice)
		return nil
	}
	fn()
	return nil
}

func (s *screen) startReviewRange() tea.Cmd {
	return s.requireDiffMode("Review range selection is only available in diff view.", s.gui.review.ToggleRangeSelection)
}

func (s *screen) startReviewComment() tea.Cmd {
	return s.requireDiffMode("Review comments are only available in diff view.", s.gui.review.BeginCommentFlow)
}

func (s *screen) startReviewSummary() tea.Cmd {
	return s.requireDiffMode("Review summary is only available in diff view.", s.gui.review.BeginSummaryInput)
}
