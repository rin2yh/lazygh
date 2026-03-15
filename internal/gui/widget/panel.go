package widget

import (
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
)

const (
	ansiReset   = "\x1b[0m"
	ansiReverse = "\x1b[7m"
)

// Highlight applies reverse-video ANSI styling to s, preserving any inner resets.
func Highlight(s string) string {
	if s == "" {
		return s
	}
	restyled := strings.ReplaceAll(s, ansiReset, ansiReset+ansiReverse)
	return ansiReverse + restyled + ansiReset
}

// ListItem formats a list entry with selection indicator.
func ListItem(text string, selected bool) string {
	if selected {
		return Highlight("> " + text)
	}
	return "  " + text
}

type PanelStyle struct {
	BorderColor string
	TitleColor  string
}

// InnerPanelHeight returns the content height inside a bordered panel.
func InnerPanelHeight(height int) int {
	if height > 2 {
		return height - 2
	}
	return height
}

func PadOrTrim(s string, width int) string {
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

// JoinColumns horizontally joins two padded column slices into one slice of combined lines.
func JoinColumns(left []string, leftW int, right []string, rightW int, height int) []string {
	lines := make([]string, 0, height)
	for i := 0; i < height; i++ {
		l := ""
		if i < len(left) {
			l = left[i]
		}
		r := ""
		if i < len(right) {
			r = right[i]
		}
		lines = append(lines, PadOrTrim(l, leftW)+" "+PadOrTrim(r, rightW))
	}
	return lines
}

func FramePanel(title string, content []string, width int, height int, style PanelStyle) []string {
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
	top := colorizeByName(strings.Repeat("─", innerWidth), style.BorderColor)
	if strings.TrimSpace(title) != "" {
		rawLabel := formatPanelTitle(title)
		topLabel := rawLabel
		if strings.TrimSpace(style.TitleColor) != "" {
			topLabel = colorizeByName(topLabel, style.TitleColor)
		}
		labelWidth := xansi.StringWidth(rawLabel)
		if labelWidth > 0 {
			if labelWidth >= innerWidth {
				top = PadOrTrim(topLabel, innerWidth)
			} else {
				top = topLabel + colorizeByName(strings.Repeat("─", innerWidth-labelWidth), style.BorderColor)
			}
		}
	}
	leftBorder := colorizeByName("│", style.BorderColor)
	rightBorder := colorizeByName("│", style.BorderColor)
	lines = append(lines, colorizeByName("┌", style.BorderColor)+top+colorizeByName("┐", style.BorderColor))
	for i := 0; i < innerHeight; i++ {
		row := ""
		if i < len(content) {
			row = content[i]
		}
		lines = append(lines, leftBorder+PadOrTrim(row, innerWidth)+rightBorder)
	}
	lines = append(lines, colorizeByName("└", style.BorderColor)+colorizeByName(strings.Repeat("─", innerWidth), style.BorderColor)+colorizeByName("┘", style.BorderColor))
	return lines
}

func ResolveColorName(name string, fallback string) string {
	if ansiCodeForColor(name) == "" {
		return fallback
	}
	return name
}

func formatPanelTitle(base string) string {
	return " " + base + " "
}

func colorizeByName(s string, colorName string) string {
	code := ansiCodeForColor(colorName)
	if code == "" || s == "" {
		return s
	}
	return code + s + ansiReset
}

func ansiCodeForColor(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "black":
		return "\x1b[30m"
	case "red":
		return "\x1b[31m"
	case "green":
		return "\x1b[32m"
	case "yellow":
		return "\x1b[33m"
	case "blue":
		return "\x1b[34m"
	case "magenta", "purple":
		return "\x1b[35m"
	case "cyan":
		return "\x1b[36m"
	case "white":
		return "\x1b[37m"
	case "gray", "grey", "brightblack", "bright-black":
		return "\x1b[90m"
	default:
		return ""
	}
}
