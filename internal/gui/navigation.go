package gui

import tea "github.com/charmbracelet/bubbletea"

func (gui *Gui) navigateDown() bool {
	return gui.state.NavigateDown()
}

func (gui *Gui) navigateUp() bool {
	return gui.state.NavigateUp()
}

func (gui *Gui) selectNextDiffFile() bool {
	if len(gui.diffFiles) == 0 || gui.diffFileSelected >= len(gui.diffFiles)-1 {
		return false
	}
	gui.diffFileSelected++
	return true
}

func (gui *Gui) selectPrevDiffFile() bool {
	if len(gui.diffFiles) == 0 || gui.diffFileSelected <= 0 {
		return false
	}
	gui.diffFileSelected--
	return true
}

func (gui *Gui) scrollDetailByKey(msg tea.KeyMsg) bool {
	if !gui.state.IsDiffMode() || gui.focus != panelDiffContent {
		return false
	}

	switch msg.String() {
	case "pgdown", "f", " ", "pgup", "b":
		updated, _ := gui.detailViewport.Update(msg)
		gui.detailViewport = updated
		return true
	case "home", "g":
		gui.detailViewport.GotoTop()
		return true
	case "end", "G":
		gui.detailViewport.GotoBottom()
		return true
	default:
		return false
	}
}

func (gui *Gui) scrollDetailDown() {
	gui.detailViewport.ScrollDown(1)
}

func (gui *Gui) scrollDetailUp() {
	gui.detailViewport.ScrollUp(1)
}
