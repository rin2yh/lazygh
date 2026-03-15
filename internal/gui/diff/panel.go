package diff

import (
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

const diffLocationColWidth = 7

type ContentLine struct {
	Location string
	Text     string
	Selected bool
	InRange  bool
}

type PanelInput struct {
	DiffMode         bool
	OverviewTitle    string
	OverviewLines    []string
	DiffFiles        []gh.DiffFile
	DiffFileSelected int
	DiffContentLines []ContentLine
}

func RenderFiles(input PanelInput, style widget.PanelStyle, width, height int) []string {
	if len(input.DiffFiles) == 0 {
		return widget.FramePanel("Files", []string{"No changed files"}, width, height, style)
	}
	inner := widget.InnerPanelHeight(height)
	lines := make([]string, 0, inner)
	startIdx := startIndex(input.DiffFileSelected, inner)
	for i := 0; len(lines) < inner; i++ {
		idx := startIdx + i
		if idx >= len(input.DiffFiles) {
			lines = append(lines, "")
			continue
		}
		line := widget.ListItem(RenderFileListLine(input.DiffFiles[idx]), idx == input.DiffFileSelected)
		lines = append(lines, line)
	}
	return widget.FramePanel("Files", lines, width, height, style)
}

func RenderContent(input PanelInput, style widget.PanelStyle, width, height int) []string {
	if len(input.DiffContentLines) == 0 {
		return widget.FramePanel("Diff", input.OverviewLines, width, height, style)
	}
	inner := widget.InnerPanelHeight(height)
	lines := renderLines(input.DiffContentLines, inner)
	return widget.FramePanel("Diff", lines, width, height, style)
}

func ContentLines(input PanelInput, height int) []string {
	if len(input.DiffContentLines) == 0 {
		return nil
	}
	inner := widget.InnerPanelHeight(height)
	return renderLines(input.DiffContentLines, inner)
}

func renderLines(contentLines []ContentLine, inner int) []string {
	sel := selectedIndex(contentLines)
	startIdx := startIndex(sel, inner)
	lines := make([]string, 0, inner)
	for i := 0; len(lines) < inner; i++ {
		idx := startIdx + i
		if idx >= len(contentLines) {
			lines = append(lines, "")
			continue
		}
		lines = append(lines, diffLine(contentLines[idx]))
	}
	return lines
}

func diffLine(line ContentLine) string {
	prefix := "  "
	if line.Location != "" {
		prefix = widget.PadOrTrim(line.Location, diffLocationColWidth) + " "
	}
	if line.InRange {
		styledPrefix := colorize(ansiCyan, prefix)
		if line.Selected {
			return widget.Highlight(styledPrefix) + line.Text
		}
		return styledPrefix + line.Text
	}
	if line.Selected {
		return widget.Highlight(prefix) + line.Text
	}
	return prefix + line.Text
}

func selectedIndex(lines []ContentLine) int {
	for i, line := range lines {
		if line.Selected {
			return i
		}
	}
	return 0
}

func startIndex(sel, visible int) int {
	if visible <= 0 || sel < visible {
		return 0
	}
	return sel - visible + 1
}


