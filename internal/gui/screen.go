package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
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
	case reviewCommentSavedMsg:
		s.gui.applyReviewCommentResult(msg)
		return s, nil
	case reviewSubmittedMsg:
		s.gui.applyReviewSubmitResult(msg)
		return s, nil
	case reviewDiscardedMsg:
		s.gui.applyReviewDiscardResult(msg)
		return s, nil
	case tea.KeyMsg:
		if s.gui.state.Review.InputMode != core.ReviewInputNone {
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
