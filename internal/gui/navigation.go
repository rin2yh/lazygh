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
	gui.diffLineSelected = 0
	gui.ensureDiffLineSelection()
	return true
}

func (gui *Gui) selectPrevDiffFile() bool {
	if len(gui.diffFiles) == 0 || gui.diffFileSelected <= 0 {
		return false
	}
	gui.diffFileSelected--
	gui.diffLineSelected = 0
	gui.ensureDiffLineSelection()
	return true
}

func (gui *Gui) scrollDetailByKey(msg tea.KeyMsg) bool {
	if !gui.state.IsDiffMode() || gui.focus != panelDiffContent {
		return false
	}

	switch msg.String() {
	case "pgdown", "f", " ", "pgup", "b":
		if msg.String() == "pgup" || msg.String() == "b" {
			return gui.selectPrevDiffLine(gui.detailViewportHeight)
		}
		return gui.selectNextDiffLine(gui.detailViewportHeight)
	case "home", "g":
		return gui.gotoFirstDiffLine()
	case "end", "G":
		return gui.gotoLastDiffLine()
	default:
		return false
	}
}

func (gui *Gui) scrollOverviewByKey(msg tea.KeyMsg) bool {
	if gui.state.IsDiffMode() || gui.focus != panelDiffContent {
		return false
	}
	switch msg.String() {
	case "pgdown", "f", " ", "pgup", "b":
		updated, _ := gui.detailViewport.Update(msg)
		gui.detailViewport = updated
		return true
	default:
		return false
	}
}

func (gui *Gui) scrollDetailDown() {
	if gui.state.IsDiffMode() {
		gui.selectNextDiffLine(1)
		return
	}
	gui.detailViewport.ScrollDown(1)
}

func (gui *Gui) scrollDetailUp() {
	if gui.state.IsDiffMode() {
		gui.selectPrevDiffLine(1)
		return
	}
	gui.detailViewport.ScrollUp(1)
}
