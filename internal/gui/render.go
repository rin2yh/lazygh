package gui

import (
	"strings"

	"github.com/rin2yh/lazygh/internal/core"
)

func (gui *Gui) render() string {
	w := gui.state.Width
	h := gui.state.Height
	if w <= 0 {
		w = 120
	}
	if h <= 0 {
		h = 40
	}

	leftRatio := 26
	if gui.state.IsDiffMode() {
		leftRatio = 22
	}
	leftWidth := w * leftRatio / 100
	if leftWidth < 1 {
		leftWidth = 1
	}
	if leftWidth > w-2 {
		leftWidth = w - 2
	}
	rightWidth := w - leftWidth - 1
	if rightWidth < 1 {
		rightWidth = 1
	}

	contentHeight := h - 1
	if contentHeight < 1 {
		contentHeight = 1
	}

	leftLines := gui.renderLeftPanels(leftWidth, contentHeight)
	rightLines := gui.renderRightPanels(rightWidth, contentHeight)

	var b strings.Builder
	for i := 0; i < contentHeight; i++ {
		left := ""
		if i < len(leftLines) {
			left = leftLines[i]
		}
		right := ""
		if i < len(rightLines) {
			right = rightLines[i]
		}
		b.WriteString(padOrTrim(left, leftWidth))
		b.WriteRune(' ')
		b.WriteString(padOrTrim(right, rightWidth))
		b.WriteByte('\n')
	}
	b.WriteString(padOrTrim(
		formatStatusLine(
			gui.state.Loading != core.LoadingNone,
			gui.state.IsDiffMode(),
			len(gui.state.PRs) > 0,
			gui.focus,
			len(gui.diffFiles) > 0,
		),
		w,
	))
	return b.String()
}

func (gui *Gui) renderRightPanels(width int, height int) []string {
	if !gui.state.IsDiffMode() {
		return gui.renderDetailPanel("Overview", false, width, height, gui.state.DetailContent)
	}
	coloredDiff := colorizeDiffContent(gui.currentDiffContent())

	if width < 20 {
		return gui.renderDetailPanel("Diff", gui.focus == panelDiffContent, width, height, coloredDiff)
	}

	filesWidth := width * 30 / 100
	if filesWidth < 16 {
		filesWidth = 16
	}
	if filesWidth > width-10 {
		filesWidth = width - 10
	}
	diffWidth := width - filesWidth - 1
	if diffWidth < 1 {
		diffWidth = 1
	}

	filesLines := gui.renderDiffFilesPanel(filesWidth, height)
	diffLines := gui.renderDetailPanel("Diff", gui.focus == panelDiffContent, diffWidth, height, coloredDiff)

	lines := make([]string, 0, height)
	for i := 0; i < height; i++ {
		left := ""
		if i < len(filesLines) {
			left = filesLines[i]
		}
		right := ""
		if i < len(diffLines) {
			right = diffLines[i]
		}
		lines = append(lines, padOrTrim(left, filesWidth)+" "+padOrTrim(right, diffWidth))
	}
	return lines
}

func (gui *Gui) renderPRPanel(height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)

	if gui.state.PRsLoading {
		for len(lines) < height {
			lines = append(lines, "")
		}
		return lines
	}

	if len(gui.state.PRs) == 0 {
		for len(lines) < height {
			if len(lines) == 0 {
				lines = append(lines, "No pull requests")
			} else {
				lines = append(lines, "")
			}
		}
		return lines
	}

	for i := 0; len(lines) < height; i++ {
		if i >= len(gui.state.PRs) {
			lines = append(lines, "")
			continue
		}
		prefix := "  "
		line := prefix + core.FormatPRItem(gui.state.PRs[i])
		if i == gui.state.PRsSelected {
			prefix = "> "
			line = highlightLine(prefix + core.FormatPRItem(gui.state.PRs[i]))
		}
		lines = append(lines, line)
	}
	return lines
}

func (gui *Gui) renderRepoPanel(height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)
	for len(lines) < height {
		if len(lines) == 0 {
			lines = append(lines, formatRepoLine(gui.state.Repo))
		} else {
			lines = append(lines, "")
		}
	}
	return lines
}

func (gui *Gui) renderLeftPanels(width int, height int) []string {
	if height <= 0 {
		return nil
	}

	repoPanelHeight := 4
	if height < repoPanelHeight+1 {
		repoPanelHeight = height / 2
	}
	if repoPanelHeight < 1 {
		repoPanelHeight = 1
	}
	prPanelHeight := height - repoPanelHeight
	if prPanelHeight < 1 {
		prPanelHeight = 1
		repoPanelHeight = height - prPanelHeight
	}

	repoInnerHeight := repoPanelHeight
	if repoPanelHeight > 2 {
		repoInnerHeight = repoPanelHeight - 2
	}
	prInnerHeight := prPanelHeight
	if prPanelHeight > 2 {
		prInnerHeight = prPanelHeight - 2
	}

	repoLines := gui.framePanel("Repository", false, gui.renderRepoPanel(repoInnerHeight), width, repoPanelHeight)
	prLines := gui.framePanel("PRs (Open/Draft)", gui.focus == panelPRs, gui.renderPRPanel(prInnerHeight), width, prPanelHeight)

	lines := make([]string, 0, height)
	lines = append(lines, repoLines...)
	lines = append(lines, prLines...)
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return lines
}

func (gui *Gui) renderDiffFilesPanel(width int, height int) []string {
	if height <= 0 {
		return nil
	}

	innerHeight := height
	if height > 2 {
		innerHeight = height - 2
	}
	lines := make([]string, 0, innerHeight)

	if len(gui.diffFiles) == 0 {
		for len(lines) < innerHeight {
			if len(lines) == 0 {
				lines = append(lines, "No changed files")
			} else {
				lines = append(lines, "")
			}
		}
		return gui.framePanel("Files", gui.focus == panelDiffFiles, lines, width, height)
	}

	start := 0
	if gui.diffFileSelected >= innerHeight {
		start = gui.diffFileSelected - innerHeight + 1
	}
	for i := 0; len(lines) < innerHeight; i++ {
		idx := start + i
		if idx >= len(gui.diffFiles) {
			lines = append(lines, "")
			continue
		}
		prefix := "  "
		line := prefix + renderDiffFileListLine(gui.diffFiles[idx])
		if idx == gui.diffFileSelected {
			prefix = "> "
			line = highlightLine(prefix + renderDiffFileListLine(gui.diffFiles[idx]))
		}
		lines = append(lines, line)
	}
	return gui.framePanel("Files", gui.focus == panelDiffFiles, lines, width, height)
}

func (gui *Gui) renderDetailPanel(title string, active bool, width int, height int, content string) []string {
	if height <= 0 {
		return nil
	}

	innerWidth := width
	if width > 2 {
		innerWidth = width - 2
	}
	innerHeight := height
	if height > 2 {
		innerHeight = height - 2
	}
	bodyHeight := innerHeight
	if bodyHeight < 1 {
		bodyHeight = 1
	}
	gui.syncDetailViewport(innerWidth, bodyHeight, content)

	lines := make([]string, 0, innerHeight)
	for _, line := range strings.Split(gui.detailViewport.View(), "\n") {
		if len(lines) >= innerHeight {
			break
		}
		lines = append(lines, line)
	}
	for len(lines) < innerHeight {
		lines = append(lines, "")
	}
	return gui.framePanel(title, active, lines, width, height)
}

func highlightLine(s string) string {
	if s == "" {
		return s
	}
	// keep reverse-video active even when the line contains inner color resets.
	restyled := strings.ReplaceAll(s, ansiReset, ansiReset+ansiReverse)
	return ansiReverse + restyled + ansiReset
}

func (gui *Gui) syncDetailViewport(width int, height int, content string) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if gui.detailViewportWidth != width || gui.detailViewportHeight != height {
		gui.detailViewport.Width = width
		gui.detailViewport.Height = height
		gui.detailViewportWidth = width
		gui.detailViewportHeight = height
	}
	wrapped := wrapText(content, width)
	if gui.detailViewportBody != wrapped {
		gui.detailViewport.SetContent(wrapped)
		gui.detailViewportBody = wrapped
		gui.detailViewport.GotoTop()
	}
}
