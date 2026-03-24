package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/pr/list"
	"github.com/rin2yh/lazygh/internal/pr/overview"
	"github.com/rin2yh/lazygh/internal/pr/review"
)

// --- メインキー入力ディスパッチャ ---

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
	if s.gui.review.InputMode() == review.InputNone && s.gui.review.HasRangeStart() {
		s.gui.review.ClearRangeStart()
		s.gui.review.Notify("Range selection cleared.")
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

// --- ナビゲーションアクション ---

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

func (s *screen) showDiff() tea.Cmd {
	if s.gui.switchToDiff() {
		return s.openSelectedPR()
	}
	return nil
}

func (s *screen) moveDown() tea.Cmd { return s.moveCursor(1) }
func (s *screen) moveUp() tea.Cmd   { return s.moveCursor(-1) }

func (s *screen) moveCursor(dir int) tea.Cmd {
	if s.gui.coord.IsDiffMode() {
		switch s.gui.focus {
		case layout.FocusPRs:
			navigate := s.gui.navigateDown
			if dir < 0 {
				navigate = s.gui.navigateUp
			}
			if navigate() {
				return s.openSelectedPR()
			}
			return nil
		case layout.FocusDiffFiles:
			if dir > 0 {
				s.gui.diff.SelectNextFile()
			} else {
				s.gui.diff.SelectPrevFile()
			}
			return nil
		case layout.FocusDiffContent:
			if dir > 0 {
				s.scrollDetailDown()
			} else {
				s.scrollDetailUp()
			}
			return nil
		case layout.FocusReviewDrawer:
			if dir > 0 {
				s.gui.review.SelectNextComment()
			} else {
				s.gui.review.SelectPrevComment()
			}
			return nil
		}
		return nil
	}

	if s.gui.focus == layout.FocusPRs {
		if dir > 0 {
			s.gui.navigateDown()
		} else {
			s.gui.navigateUp()
		}
	}
	return nil
}

func (s *screen) openSelectedPR() tea.Cmd {
	action := s.gui.coord.PlanEnter(s.gui.client != nil)
	switch action.Kind {
	case EnterLoadPRDiff:
		return s.loadDetailCmd(action.Repo, action.Number, overview.DetailModeDiff)
	case EnterLoadPRDetail:
		return s.loadDetailCmd(action.Repo, action.Number, overview.DetailModeOverview)
	default:
		return nil
	}
}

func (s *screen) handleDetailScrollAction(action config.Action) tea.Cmd {
	if s.gui.focus != layout.FocusDiffContent {
		return nil
	}

	if s.gui.coord.IsDiffMode() {
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
	case config.ActionGoTop:
		s.gui.detail.GotoTop()
	case config.ActionGoBottom:
		s.gui.detail.GotoBottom()
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
	if s.gui.coord.IsDiffMode() {
		s.gui.diff.SelectNextLine(1)
		return
	}
	s.gui.detail.ScrollDown(1)
}

func (s *screen) scrollDetailUp() {
	if s.gui.coord.IsDiffMode() {
		s.gui.diff.SelectPrevLine(1)
		return
	}
	s.gui.detail.ScrollUp(1)
}

// --- フィルター入力 ---

func (s *screen) handleFilterKey(msg tea.KeyMsg) tea.Cmd {
	if s.gui.coord.HandleFilterKey(msg.String()) == list.FilterKeyApply {
		s.gui.coord.BeginFetchPRs()
		return s.loadPRsCmd()
	}
	return nil
}

// --- レビュー入力 ---

func (s *screen) handleReviewInputKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	return s.gui.review.HandleInputKey(msg)
}

func (s *screen) handleReviewAction(action config.Action) tea.Cmd {
	switch action {
	case config.ActionReviewDeleteComment, config.ActionReviewEditComment:
		if s.gui.focus != layout.FocusReviewDrawer {
			return nil
		}
	}
	return s.gui.review.HandleAction(action)
}
