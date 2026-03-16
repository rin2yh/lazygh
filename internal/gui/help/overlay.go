package help

import (
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

func overlayLine(bg, panel string, startX, panelW, screenW int) string {
	return bgLeft(bg, startX) + widget.PadOrTrim(panel, panelW) + bgRight(startX+panelW, screenW)
}

func bgLeft(bg string, x int) string {
	return widget.PadOrTrim(xansi.Truncate(bg, x, ""), x)
}

func bgRight(endX, screenW int) string {
	if endX >= screenW {
		return ""
	}
	return strings.Repeat(" ", screenW-endX)
}
