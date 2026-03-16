package gui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/model"
)

func (s *screen) handleReviewInputKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	action, ok := s.gui.config.KeyBindings.ActionFor(msg)
	if ok {
		switch action {
		case config.ActionReviewSubmit:
			return s.gui.review.HandleSubmit(), true
		case config.ActionReviewDiscard:
			return s.gui.review.HandleDiscard(), true
		case config.ActionReviewSave:
			if s.gui.state.Review.InputMode == model.ReviewInputComment {
				if s.gui.review.IsEditingComment() {
					return s.gui.review.HandleEditCommentSave(), true
				}
				return s.gui.review.HandleCommentSave(), true
			}
		}
	}
	if cmd, handled := s.gui.review.HandleEditorKey(msg); handled {
		return cmd, true
	}
	return nil, false
}

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
		s.gui.state.OpenFilterSelect()
		return nil, true
	default:
		return nil, false
	}
}

func (s *screen) handleNavigationAction(action config.Action) (tea.Cmd, bool) {
	switch action {
	case config.ActionMoveDown:
		return s.moveDown(), true
	case config.ActionMoveUp:
		return s.moveUp(), true
	case config.ActionPageDown, config.ActionPageUp, config.ActionGoTop, config.ActionGoBottom:
		return s.handleDetailScrollAction(action), true
	default:
		return nil, false
	}
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
		return s.gui.review.HandleSubmit()
	case config.ActionReviewDiscard:
		return s.gui.review.HandleDiscard()
	case config.ActionReviewClearComment:
		if s.gui.state.Review.InputMode == model.ReviewInputComment {
			s.gui.review.ClearCommentInput()
		}
	case config.ActionReviewEvent:
		if s.gui.state.IsDiffMode() {
			s.gui.review.CycleReviewEvent()
		}
	case config.ActionReviewDeleteComment:
		if s.gui.focus == panelReviewDrawer {
			return s.gui.review.HandleDeleteComment()
		}
	case config.ActionReviewEditComment:
		if s.gui.focus == panelReviewDrawer {
			s.gui.review.BeginEditComment()
		}
	}
	return nil
}

func (s *screen) handleFilterKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		s.gui.state.CloseFilterSelect()
		return nil
	case "enter":
		s.gui.state.CloseFilterSelect()
		s.gui.state.BeginLoadPRs()
		return s.loadPRsCmd()
	case "j", "down":
		s.gui.state.MoveFilterCursor(1)
		return nil
	case "k", "up":
		s.gui.state.MoveFilterCursor(-1)
		return nil
	case " ":
		s.gui.state.ToggleFilterAtCursor()
		return nil
	}
	return nil
}

func (s *screen) handleCancel() tea.Cmd {
	if s.gui.state.Review.InputMode == model.ReviewInputNone && s.gui.state.Review.RangeStart != nil {
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

func (s *screen) moveDown() tea.Cmd { return s.moveCursor(1) }
func (s *screen) moveUp() tea.Cmd   { return s.moveCursor(-1) }

// moveCursor moves the cursor in the given direction (1 = down, -1 = up).
func (s *screen) moveCursor(dir int) tea.Cmd {
	if s.gui.state.IsDiffMode() {
		switch s.gui.focus {
		case panelPRs:
			navigate := s.gui.navigateDown
			if dir < 0 {
				navigate = s.gui.navigateUp
			}
			if navigate() {
				return s.openSelectedPR()
			}
			return nil
		case panelDiffFiles:
			if dir > 0 {
				s.gui.diff.SelectNextFile()
			} else {
				s.gui.diff.SelectPrevFile()
			}
			return nil
		case panelDiffContent:
			if dir > 0 {
				s.scrollDetailDown()
			} else {
				s.scrollDetailUp()
			}
			return nil
		case panelReviewDrawer:
			if dir > 0 {
				s.gui.review.SelectNextComment()
			} else {
				s.gui.review.SelectPrevComment()
			}
			return nil
		}
		return nil
	}

	if s.gui.focus == panelPRs {
		if dir > 0 {
			s.gui.navigateDown()
		} else {
			s.gui.navigateUp()
		}
	}
	return nil
}

func (s *screen) openSelectedPR() tea.Cmd {
	action := s.gui.state.PlanEnter(s.gui.client != nil, os.Getenv("LAZYGH_DEBUG_DETAIL_TEXT"))
	switch action.Kind {
	case model.EnterLoadPRDiff:
		return s.loadDetailCmd(action.Repo, action.Number, model.DetailModeDiff)
	case model.EnterLoadPRDetail:
		return s.loadDetailCmd(action.Repo, action.Number, model.DetailModeOverview)
	default:
		return nil
	}
}

func (s *screen) handleDetailScrollAction(action config.Action) tea.Cmd {
	if s.gui.focus != panelDiffContent {
		return nil
	}

	if s.gui.state.IsDiffMode() {
		switch action {
		case config.ActionPageUp:
			s.gui.diff.SelectPrevLine(s.gui.detail.Height())
		case config.ActionPageDown:
			s.gui.diff.SelectNextLine(s.gui.detail.Height())
		case config.ActionGoTop:
			s.gui.diff.GotoFirstLine()
		case config.ActionGoBottom:
			s.gui.diff.GotoLastLine()
		}
		return nil
	}

	switch action {
	case config.ActionPageDown, config.ActionPageUp:
		key, ok := primaryKeyMsg(s.gui.config.KeyBindings.Binding(action))
		if !ok {
			return nil
		}
		_, cmd := s.gui.detail.Update(key)
		return cmd
	}
	return nil
}

func primaryKeyMsg(binding config.KeyBinding) (tea.KeyMsg, bool) {
	if len(binding.Keys) == 0 {
		return tea.KeyMsg{}, false
	}
	switch binding.Keys[0] {
	case "pgdown":
		return tea.KeyMsg{Type: tea.KeyPgDown}, true
	case "pgup":
		return tea.KeyMsg{Type: tea.KeyPgUp}, true
	case " ":
		return tea.KeyMsg{Type: tea.KeySpace}, true
	case "b":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}, true
	case "f":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}, true
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(binding.Keys[0])}, true
	}
}

func (s *screen) scrollDetailDown() {
	if s.gui.state.IsDiffMode() {
		s.gui.diff.SelectNextLine(1)
		return
	}
	s.gui.detail.ScrollDown(1)
}

func (s *screen) scrollDetailUp() {
	if s.gui.state.IsDiffMode() {
		s.gui.diff.SelectPrevLine(1)
		return
	}
	s.gui.detail.ScrollUp(1)
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
