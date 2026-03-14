package gui

import (
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

func padOrTrim(s string, width int) string {
	if width <= 0 {
		return ""
	}
	out := xansi.Truncate(s, width, "")
	col := xansi.StringWidth(out)
	if col < width {
		out += strings.Repeat(" ", width-col)
	}
	return out
}

func framePanel(title string, active bool, content []string, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	if width < 2 || height < 3 {
		lines := make([]string, 0, height)
		for i := 0; i < height; i++ {
			if i < len(content) {
				lines = append(lines, content[i])
			} else {
				lines = append(lines, "")
			}
		}
		return lines
	}

	innerWidth := width - 2
	innerHeight := height - 2
	lines := make([]string, 0, height)
	top := strings.Repeat("─", innerWidth)
	if strings.TrimSpace(title) != "" {
		topLabel := formatPanelTitle(title, active)
		labelWidth := runewidth.StringWidth(topLabel)
		if labelWidth > 0 {
			if labelWidth >= innerWidth {
				top = padOrTrim(topLabel, innerWidth)
			} else {
				top = topLabel + strings.Repeat("─", innerWidth-labelWidth)
			}
		}
	}
	lines = append(lines, "┌"+top+"┐")
	for i := 0; i < innerHeight; i++ {
		row := ""
		if i < len(content) {
			row = content[i]
		}
		lines = append(lines, "│"+padOrTrim(row, innerWidth)+"│")
	}
	lines = append(lines, "└"+strings.Repeat("─", innerWidth)+"┘")
	return lines
}
