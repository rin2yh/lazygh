package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/model"
)

func (s *screen) handleKeyInput(msg tea.KeyMsg) tea.Cmd {
	action, ok := s.gui.config.KeyBindings.ActionFor(msg)
	if !ok {
		return nil
	}

	if cmd, handled := s.handleGlobalAction(action); handled {
		return cmd
	}
	if cmd, handled := s.handleNavigationAction(action); handled {
		return cmd
	}
	return s.handleReviewAction(action)
}

func (s *screen) handleGlobalAction(action config.Action) (tea.Cmd, bool) {
	switch action {
	case config.ActionShowHelp:
		s.gui.showHelp = !s.gui.showHelp
		return nil, true
	case config.ActionQuit:
		return tea.Quit, true
	case config.ActionCancel:
		return s.handleCancel(), true
	case config.ActionFocusNext:
		s.gui.cycleFocus()
		return nil, true
	case config.ActionPanelPrev:
		s.gui.moveFocus(-1)
		return nil, true
	case config.ActionPanelNext:
		s.gui.moveFocus(1)
		return nil, true
	case config.ActionShowOverview:
		s.gui.switchToOverview()
		return nil, true
	case config.ActionShowDiff:
		return s.showDiff(), true
	case config.ActionOpenSelected:
		return s.openSelectedPR(), true
	case config.ActionFilterPRs:
		s.gui.coord.OpenFilterSelect()
		return nil, true
	default:
		return nil, false
	}
}

func (s *screen) handleCancel() tea.Cmd {
	if s.gui.review.InputMode() == model.ReviewInputNone && s.gui.review.HasRangeStart() {
		s.gui.review.ClearRangeStart()
		s.gui.review.SetNotice("Range selection cleared.")
		s.gui.focus = layout.FocusDiffContent
		return nil
	}
	if s.gui.focus == layout.FocusReviewDrawer {
		s.gui.review.StopInput()
		s.gui.focus = layout.FocusDiffContent
		return nil
	}
	s.gui.focusPRs()
	return nil
}
