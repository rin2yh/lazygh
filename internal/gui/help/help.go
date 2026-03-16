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

// RenderOverlay renders the help panel lines centered on the screen.
// It paints over the background lines in the center region.
func RenderOverlay(background []string, keys config.KeyBindings, screenWidth int) []string {
	sections := buildSections(keys)
	content := buildContent(sections)
	closeHint := fmt.Sprintf("Press [%s] or [%s] to close", keys.HelpLabel(), keys.Label(config.ActionCancel))

	// Build panel lines
	panelContent := make([]string, 0, len(content)+2)
	panelContent = append(panelContent, "")
	panelContent = append(panelContent, content...)
	panelContent = append(panelContent, "")
	panelContent = append(panelContent, "  "+closeHint)
	panelContent = append(panelContent, "")

	// Determine panel dimensions
	innerW := 0
	for _, line := range panelContent {
		w := xansi.StringWidth(line)
		if w > innerW {
			innerW = w
		}
	}
	panelW := innerW + 2 // FramePanel subtracts 2 for borders; content already has 2-space indent
	if panelW > screenWidth-2 {
		panelW = screenWidth - 2
	}
	panelH := len(panelContent) + 2 // 2 border rows

	panelLines := widget.FramePanel("Keybindings", panelContent, panelW, panelH,
		widget.PanelStyle{BorderColor: "yellow", TitleColor: "yellow"},
	)

	screenH := len(background)
	startY := (screenH - panelH) / 2
	if startY < 0 {
		startY = 0
	}
	startX := (screenWidth - panelW) / 2
	if startX < 0 {
		startX = 0
	}

	result := make([]string, len(background))
	copy(result, background)

	for i, line := range panelLines {
		y := startY + i
		if y < 0 || y >= len(result) {
			continue
		}
		result[y] = overlayLine(result[y], line, startX, panelW, screenWidth)
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
