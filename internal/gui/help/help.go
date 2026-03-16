package help

import (
	"fmt"
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

type section struct {
	title string
	rows  [][2]string // [key label, description]
}

func buildSections(keys config.KeyBindings) []section {
	return []section{
		{
			title: "Navigation",
			rows: [][2]string{
				{keys.MoveLabel(), "Move Up/Down"},
				{keys.PanelLabel(), "Panel Prev/Next"},
				{keys.FocusLabel(), "Cycle Focus"},
				{keys.PageLabel(), "Page Up/Down"},
				{keys.TopBottomLabel(), "Go Top/Bottom"},
				{keys.Label(config.ActionCancel), "Cancel / Close"},
			},
		},
		{
			title: "View",
			rows: [][2]string{
				{keys.DiffLabel(), "Show Diff"},
				{keys.OverviewLabel(), "Show Overview"},
				{keys.ReloadLabel(), "Reload PR"},
				{keys.QuitLabel(), "Quit"},
			},
		},
		{
			title: "Review",
			rows: [][2]string{
				{keys.RangeLabel(), "Select Range"},
				{keys.CommentLabel(), "Add Comment"},
				{keys.SummaryLabel(), "Edit Summary"},
				{keys.SaveLabel(), "Save Comment"},
				{keys.SubmitLabel(), "Submit Review"},
				{keys.DiscardLabel(), "Discard Review"},
			},
		},
	}
}

func buildPanelLines(keys config.KeyBindings, screenWidth int) ([]string, int) {
	content := buildContent(buildSections(keys))
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

func bgLeft(bg string, x int) string {
	return widget.PadOrTrim(xansi.Truncate(bg, x, ""), x)
}

func overlayLine(bg, panel string, startX, panelW, screenW int) string {
	return bgLeft(bg, startX) + widget.PadOrTrim(panel, panelW) + bgRight(startX+panelW, screenW)
}

func bgRight(endX, screenW int) string {
	if endX >= screenW {
		return ""
	}
	return strings.Repeat(" ", screenW-endX)
}

func buildContent(sections []section) []string {
	var lines []string
	for i, sec := range sections {
		if i > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, "  "+sec.title)
		lines = append(lines, "  "+strings.Repeat("─", 36))
		for _, row := range sec.rows {
			if row[0] == "" {
				continue
			}
			keyCol := fmt.Sprintf("[%s]", row[0])
			lines = append(lines, fmt.Sprintf("  %-10s %s", keyCol, row[1]))
		}
	}
	return lines
}
