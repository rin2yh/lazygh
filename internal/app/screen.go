package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/pr/overview"
	"github.com/rin2yh/lazygh/internal/pr/review"
)

type screen struct {
	gui *Gui
}

func (s *screen) Init() tea.Cmd {
	if s.gui.client == nil {
		return nil
	}
	s.gui.coord.BeginFetchPRs()
	return s.loadPRsCmd()
}

func (s *screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.gui.coord.SetWindowSize(msg.Width, msg.Height)
		return s, nil
	case prsLoadedMsg:
		s.gui.applyPRsResult(msg)
		return s, nil
	case detailLoadedMsg:
		cmd := s.gui.applyDetailResult(msg)
		if msg.mode == overview.DetailModeDiff && msg.err == nil {
			return s, tea.Batch(cmd, s.loadThreadsCmd(s.gui.coord.ListRepo(), msg.number))
		}
		return s, cmd
	case review.CommentSavedMsg:
		s.gui.review.CommentResult(msg)
		return s, nil
	case review.CommentDeletedMsg:
		s.gui.review.DeleteCommentResult(msg)
		return s, nil
	case review.CommentUpdatedMsg:
		s.gui.review.EditCommentResult(msg)
		return s, nil
	case review.SubmittedMsg:
		s.gui.review.SubmitResult(msg)
		return s, nil
	case review.DiscardedMsg:
		s.gui.review.DiscardResult(msg)
		return s, nil
	case threadsLoadedMsg:
		s.gui.applyThreadsResult(msg)
		return s, nil
	case review.ThreadReplyMsg:
		s.gui.review.ThreadReplyResult(msg)
		return s, nil
	case tea.KeyMsg:
		if s.gui.showHelp {
			s.gui.showHelp = false
			return s, nil
		}
		if s.gui.coord.FilterOpen {
			return s, s.handleFilterKey(msg)
		}
		if s.gui.review.IsInInputMode() {
			if cmd, handled := s.handleReviewInputKey(msg); handled {
				return s, cmd
			}
		}
		return s, s.handleKeyInput(msg)
	}
	return s, nil
}

func (s *screen) View() string {
	return s.gui.render()
}
