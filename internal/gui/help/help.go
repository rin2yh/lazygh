package help

import "github.com/rin2yh/lazygh/internal/config"

// RenderOverlay renders the help panel lines centered on the screen.
// It paints over the background lines in the center region.
func RenderOverlay(background []string, keys config.KeyBindings, screenWidth int) []string {
	panelLines, panelW := buildPanelLines(keys, screenWidth)
	panelH := len(panelLines)

	startY := max(0, (len(background)-panelH)/2)
	startX := max(0, (screenWidth-panelW)/2)

	result := make([]string, len(background))
	copy(result, background)
	for i, line := range panelLines {
		y := startY + i
		if y >= 0 && y < len(result) {
			result[y] = overlayLine(result[y], line, startX, panelW, screenWidth)
		}
	}
	return result
}
