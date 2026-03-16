package help

import (
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

// RenderOverlay renders the help panel lines centered on the screen.
// It paints over the background lines in the center region.
func RenderOverlay(background []string, keys config.KeyBindings, screenWidth int) []string {
	panelLines, panelW := buildPanelLines(keys, screenWidth)
	return widget.OverlayPanel(background, panelLines, panelW, screenWidth)
}
