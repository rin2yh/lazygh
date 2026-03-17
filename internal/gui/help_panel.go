package gui

import (
	"fmt"
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/help"
	prhelp "github.com/rin2yh/lazygh/internal/pr/help"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

func renderHelpOverlay(background []string, keys config.KeyBindings, screenWidth int) []string {
	panelLines, panelW := buildHelpPanelLines(keys, screenWidth)
	return widget.OverlayPanel(background, panelLines, panelW, screenWidth)
}

func buildHelpPanelLines(keys config.KeyBindings, screenWidth int) ([]string, int) {
	sections := append(help.CommonSections(keys), prhelp.Sections(keys)...)
	content := buildHelpContent(sections)
	closeHint := fmt.Sprintf("Press [%s] or [%s] to close", keys.HelpLabel(), keys.Label(config.ActionCancel))

	panelContent := make([]string, 0, len(content)+4)
	panelContent = append(panelContent, "")
	panelContent = append(panelContent, content...)
	panelContent = append(panelContent, "")
	panelContent = append(panelContent, "  "+closeHint)
	panelContent = append(panelContent, "")

	innerW := 0
	for _, line := range panelContent {
		if w := xansi.StringWidth(line); w > innerW {
			innerW = w
		}
	}
	panelW := min(innerW+2, screenWidth-2) // FramePanel subtracts 2 for borders
	panelH := len(panelContent) + 2        // 2 border rows

	lines := widget.FramePanel("Keybindings", panelContent, panelW, panelH,
		widget.PanelStyle{BorderColor: "yellow", TitleColor: "yellow"},
	)
	return lines, panelW
}

func buildHelpContent(sections []help.Section) []string {
	var lines []string
	for i, sec := range sections {
		if i > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, "  "+sec.Title)
		lines = append(lines, "  "+strings.Repeat("─", 36))
		for _, row := range sec.Rows {
			if row[0] == "" {
				continue
			}
			keyCol := fmt.Sprintf("[%s]", row[0])
			lines = append(lines, fmt.Sprintf("  %-10s %s", keyCol, row[1]))
		}
	}
	return lines
}
