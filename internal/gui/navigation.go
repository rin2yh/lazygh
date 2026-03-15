package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/diff"
)

func (gui *Gui) CurrentDiffFile() (gh.DiffFile, bool) {
	return gui.currentDiffFile()
}

func (gui *Gui) CurrentDiffLine() (gh.DiffLine, bool) {
	return gui.currentDiffLine()
}

func (gui *Gui) CurrentLineIndex() int {
	return gui.diffLineSelected
}

func (gui *Gui) navigateDown() bool {
	return gui.state.NavigateDown()
}

func (gui *Gui) navigateUp() bool {
	return gui.state.NavigateUp()
}

func (gui *Gui) currentDiffFile() (gh.DiffFile, bool) {
	return diff.CurrentFile(gui.diffFiles, gui.diffFileSelected)
}

func (gui *Gui) currentDiffLine() (gh.DiffLine, bool) {
	file, ok := gui.currentDiffFile()
	if !ok {
		return gh.DiffLine{}, false
	}
	return diff.CurrentLine(file, gui.diffLineSelected)
}

func (gui *Gui) ensureDiffLineSelection() {
	file, ok := gui.currentDiffFile()
	if !ok {
		gui.diffLineSelected = 0
		return
	}
	gui.diffLineSelected = diff.EnsureLineSelection(file, gui.diffLineSelected)
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

func (gui *Gui) selectNextDiffLine(step int) bool {
	file, ok := gui.currentDiffFile()
	if !ok {
		return false
	}
	next, changed := diff.SelectNextLine(file, gui.diffLineSelected, step)
	if !changed {
		return false
	}
	gui.diffLineSelected = next
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

func (gui *Gui) selectPrevDiffLine(step int) bool {
	file, ok := gui.currentDiffFile()
	if !ok {
		return false
	}
	prev, changed := diff.SelectPrevLine(file, gui.diffLineSelected, step)
	if !changed {
		return false
	}
	gui.diffLineSelected = prev
	return true
}

func (gui *Gui) gotoFirstDiffLine() bool {
	next, changed := diff.GotoFirstLine(gui.diffLineSelected)
	if !changed {
		return false
	}
	gui.diffLineSelected = next
	return true
}

func (gui *Gui) gotoLastDiffLine() bool {
	file, ok := gui.currentDiffFile()
	if !ok {
		return false
	}
	next, changed := diff.GotoLastLine(file, gui.diffLineSelected)
	if !changed {
		return false
	}
	gui.diffLineSelected = next
	return true
}

func (gui *Gui) scrollDetailByKey(msg tea.KeyMsg) bool {
	if !gui.state.IsDiffMode() || gui.focus != panelDiffContent {
		return false
	}

	keys := gui.config.KeyBindings
	switch {
	case keys.Matches(msg, config.ActionPageUp):
		return gui.selectPrevDiffLine(gui.detailViewportHeight)
	case keys.Matches(msg, config.ActionPageDown):
		return gui.selectNextDiffLine(gui.detailViewportHeight)
	case keys.Matches(msg, config.ActionGoTop):
		return gui.gotoFirstDiffLine()
	case keys.Matches(msg, config.ActionGoBottom):
		return gui.gotoLastDiffLine()
	default:
		return false
	}
}

func (gui *Gui) scrollOverviewByKey(msg tea.KeyMsg) bool {
	if gui.state.IsDiffMode() || gui.focus != panelDiffContent {
		return false
	}
	switch {
	case gui.config.KeyBindings.Matches(msg, config.ActionPageDown),
		gui.config.KeyBindings.Matches(msg, config.ActionPageUp):
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
