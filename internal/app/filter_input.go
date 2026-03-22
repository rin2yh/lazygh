package app

import tea "github.com/charmbracelet/bubbletea"

func (s *screen) handleFilterKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		s.gui.coord.CloseFilterSelect()
		return nil
	case "enter":
		s.gui.coord.CloseFilterSelect()
		s.gui.coord.BeginFetchPRs()
		return s.loadPRsCmd()
	case "j", "down":
		s.gui.coord.MoveFilterCursor(1)
		return nil
	case "k", "up":
		s.gui.coord.MoveFilterCursor(-1)
		return nil
	case " ":
		s.gui.coord.ToggleFilterAtCursor()
		return nil
	}
	return nil
}
