package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/review"
)

type screen struct {
	gui *Gui
}

func (s *screen) Init() tea.Cmd {
	if s.gui.client == nil {
		return nil
	}
	s.gui.state.BeginLoadPRs()
	return s.loadPRsCmd()
}

func (s *screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.gui.state.SetWindowSize(msg.Width, msg.Height)
		return s, nil
	case prsLoadedMsg:
		s.gui.applyPRsResult(msg)
		return s, nil
	case detailLoadedMsg:
		s.gui.applyDetailResult(msg)
		return s, nil
	case review.CommentSavedMsg:
		s.gui.review.ApplyCommentResult(msg)
		return s, nil
	case review.CommentDeletedMsg:
		s.gui.review.ApplyDeleteCommentResult(msg)
		return s, nil
	case review.CommentUpdatedMsg:
		s.gui.review.ApplyEditCommentResult(msg)
		return s, nil
	case review.SubmittedMsg:
		s.gui.review.ApplySubmitResult(msg)
		return s, nil
	case review.DiscardedMsg:
		s.gui.review.ApplyDiscardResult(msg)
		return s, nil
	case tea.KeyMsg:
		if s.gui.showHelp {
			s.gui.showHelp = false
			return s, nil
		}
		if s.gui.state.List.FilterOpen {
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
