package list

import (
	"fmt"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/model"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

// FilterPanelLines builds the filter selection panel content and returns
// the framed lines and the panel width.
func FilterPanelLines(filter model.PRFilterMask, cursor int) ([]string, int) {
	content, innerW := buildFilterContent(filter, cursor)
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

func buildFilterContent(filter model.PRFilterMask, cursor int) ([]string, int) {
	lines := []string{""}
	maxW := 0
	for i, opt := range model.PRFilterOptions {
		check := "[ ]"
		if filter.Has(opt) {
			check = "[x]"
		}
		var line string
		if i == cursor {
			line = fmt.Sprintf("  > %s %s", check, opt.Label())
		} else {
			line = fmt.Sprintf("    %s %s", check, opt.Label())
		}
		if w := xansi.StringWidth(line); w > maxW {
			maxW = w
		}
		lines = append(lines, line)
	}
	footer := "  space:toggle  enter:apply  esc:cancel"
	if w := xansi.StringWidth(footer); w > maxW {
		maxW = w
	}
	lines = append(lines, "", footer, "")
	return lines, maxW
}
