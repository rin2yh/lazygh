package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/pr/diff"
	"github.com/rin2yh/lazygh/internal/pr/overview"
	"github.com/rin2yh/lazygh/internal/pr/review"
)

// ReviewCtrl returns the underlying *review.Controller from g for testing.
func ReviewCtrl(g *Gui) *review.Controller {
	return g.review.(*review.Controller)
}

// NewScreen creates a screen wrapping g for testing.
func NewScreen(g *Gui) *screen {
	return &screen{gui: g}
}

// GuiCoord returns the Coordinator for testing.
func GuiCoord(g *Gui) *Coordinator {
	return g.coord
}

// GuiFocus returns the current focus for testing.
func GuiFocus(g *Gui) layout.Focus {
	return g.focus
}

// SetGuiFocus sets the current focus for testing.
func SetGuiFocus(g *Gui, f layout.Focus) {
	g.focus = f
}

// GuiDiff returns a pointer to the diff selection for testing.
func GuiDiff(g *Gui) *diff.Selection {
	return &g.diff
}

// GuiNavigateDown calls g.navigateDown() for testing.
func GuiNavigateDown(g *Gui) bool {
	return g.navigateDown()
}

// GuiNavigateUp calls g.navigateUp() for testing.
func GuiNavigateUp(g *Gui) bool {
	return g.navigateUp()
}

// GuiCycleFocus calls g.cycleFocus() for testing.
func GuiCycleFocus(g *Gui) {
	g.cycleFocus()
}

// GuiSwitchToDiff calls g.switchToDiff() for testing.
func GuiSwitchToDiff(g *Gui) bool {
	return g.switchToDiff()
}

// GuiUpdateDiffFiles calls g.updateDiffFiles() for testing.
func GuiUpdateDiffFiles(g *Gui, content string) {
	g.updateDiffFiles(content)
}

// CastDetailLoadedMsg calls cmd() and returns the fields of the resulting detailLoadedMsg.
func CastDetailLoadedMsg(cmd tea.Cmd) (number int, content string, mode overview.DetailMode, err error) {
	msg := cmd().(detailLoadedMsg)
	return msg.number, msg.content, msg.mode, msg.err
}
