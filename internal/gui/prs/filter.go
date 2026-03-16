package prs

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

var filterOptionLabels = []string{"Open", "Closed", "Merged"}

// FilterPanelLines builds the filter selection panel content and returns
// the framed lines and the panel width.
func FilterPanelLines(filter core.PRFilterMask, cursor int) ([]string, int) {
	content := buildFilterContent(filter, cursor)
	innerW := 0
	for _, line := range content {
		if w := len(line); w > innerW {
			innerW = w
		}
	}
	panelW := innerW + 4 // borders + padding
	panelH := len(content) + 2
	lines := widget.FramePanel("Filter",
		content,
		panelW,
		panelH,
		widget.PanelStyle{BorderColor: "yellow", TitleColor: "yellow"},
	)
	return lines, panelW
}

func buildFilterContent(filter core.PRFilterMask, cursor int) []string {
	lines := []string{""}
	for i, opt := range core.PRFilterOptions {
		check := "[ ]"
		if filter.Has(opt) {
			check = "[x]"
		}
		label := filterOptionLabels[i]
		var line string
		if i == cursor {
			line = fmt.Sprintf("  > %s %s", check, label)
		} else {
			line = fmt.Sprintf("    %s %s", check, label)
		}
		lines = append(lines, line)
	}
	lines = append(lines, "")
	lines = append(lines, "  space:toggle  enter:apply  esc:cancel")
	lines = append(lines, "")
	return lines
}
