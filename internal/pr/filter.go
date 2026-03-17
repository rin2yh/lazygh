package pr

import (
	"fmt"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/model"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

// FilterPanelLines builds the filter selection panel content and returns
// the framed lines and the panel width.
func FilterPanelLines(filter model.PRFilterMask, cursor int) ([]string, int) {
	content := buildFilterContent(filter, cursor)
	innerW := 0
	for _, line := range content {
		if w := xansi.StringWidth(line); w > innerW {
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

func buildFilterContent(filter model.PRFilterMask, cursor int) []string {
	lines := []string{""}
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
		lines = append(lines, line)
	}
	lines = append(lines, "")
	lines = append(lines, "  space:toggle  enter:apply  esc:cancel")
	lines = append(lines, "")
	return lines
}
