package gui

import (
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
)

type panelStyle struct {
	borderColor string
	titleColor  string
}

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

func (gui *Gui) panelStyle(active bool) panelStyle {
	const (
		defaultActiveBorderColor   = "green"
		defaultInactiveBorderColor = "white"
	)

	borderColor := defaultInactiveBorderColor
	if active {
		borderColor = defaultActiveBorderColor
	}
	if gui != nil && gui.config != nil {
		if active {
			borderColor = resolveColorName(gui.config.Theme.ActiveBorderColor, defaultActiveBorderColor)
		} else {
			borderColor = resolveColorName(gui.config.Theme.InactiveBorderColor, defaultInactiveBorderColor)
		}
	}

	titleColor := ""
	if active {
		titleColor = borderColor
	}
	return panelStyle{
		borderColor: borderColor,
		titleColor:  titleColor,
	}
}

func (gui *Gui) framePanel(title string, active bool, content []string, width int, height int) []string {
	return framePanelWithStyle(title, content, width, height, gui.panelStyle(active))
}

func framePanel(title string, active bool, content []string, width int, height int) []string {
	_ = active
	return framePanelWithStyle(title, content, width, height, panelStyle{})
}

func framePanelWithStyle(title string, content []string, width int, height int, style panelStyle) []string {
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
	top := colorizeByName(strings.Repeat("─", innerWidth), style.borderColor)
	if strings.TrimSpace(title) != "" {
		rawLabel := formatPanelTitle(title)
		topLabel := rawLabel
		if strings.TrimSpace(style.titleColor) != "" {
			topLabel = colorizeByName(topLabel, style.titleColor)
		}
		labelWidth := xansi.StringWidth(rawLabel)
		if labelWidth > 0 {
			if labelWidth >= innerWidth {
				top = padOrTrim(topLabel, innerWidth)
			} else {
				top = topLabel + colorizeByName(strings.Repeat("─", innerWidth-labelWidth), style.borderColor)
			}
		}
	}
	leftBorder := colorizeByName("│", style.borderColor)
	rightBorder := colorizeByName("│", style.borderColor)
	lines = append(lines, colorizeByName("┌", style.borderColor)+top+colorizeByName("┐", style.borderColor))
	for i := 0; i < innerHeight; i++ {
		row := ""
		if i < len(content) {
			row = content[i]
		}
		lines = append(lines, leftBorder+padOrTrim(row, innerWidth)+rightBorder)
	}
	lines = append(lines, colorizeByName("└", style.borderColor)+colorizeByName(strings.Repeat("─", innerWidth), style.borderColor)+colorizeByName("┘", style.borderColor))
	return lines
}

func resolveColorName(name string, fallback string) string {
	if ansiCodeForColor(name) == "" {
		return fallback
	}
	return name
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
		return ansiRed
	case "green":
		return ansiGreen
	case "yellow":
		return ansiYellow
	case "blue":
		return ansiBlue
	case "magenta", "purple":
		return ansiPurple
	case "cyan":
		return ansiCyan
	case "white":
		return "\x1b[37m"
	case "gray", "grey", "brightblack", "bright-black":
		return ansiGray
	default:
		return ""
	}
}
