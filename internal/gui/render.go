package gui

import (
	"strconv"
	"strings"

	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
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

	mainHeight := contentHeight
	drawerHeight := 0
	if gui.shouldShowReviewDrawer() {
		drawerHeight = contentHeight / 3
		if drawerHeight < 8 {
			drawerHeight = 8
		}
		if drawerHeight >= contentHeight {
			drawerHeight = contentHeight - 1
		}
		if drawerHeight < 0 {
			drawerHeight = 0
		}
		mainHeight = contentHeight - drawerHeight
	}
	if mainHeight < 1 {
		mainHeight = 1
	}

	leftLines := gui.renderLeftPanels(leftWidth, mainHeight)
	rightLines := gui.renderRightPanels(rightWidth, mainHeight)

	var b strings.Builder
	for i := 0; i < mainHeight; i++ {
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
	if drawerHeight > 0 {
		for _, line := range gui.renderReviewDrawer(w, drawerHeight) {
			b.WriteString(padOrTrim(line, w))
			b.WriteByte('\n')
		}
	}
	b.WriteString(padOrTrim(
		formatStatusLine(
			gui.state.Loading != core.LoadingNone,
			gui.state.IsDiffMode(),
			len(gui.state.PRs) > 0,
			gui.focus,
			len(gui.diffFiles) > 0,
			gui.shouldShowReviewDrawer(),
			gui.state.Review.InputMode,
		),
		w,
	))
	return b.String()
}

func (gui *Gui) renderRightPanels(width int, height int) []string {
	if !gui.state.IsDiffMode() {
		return gui.renderDetailPanel("Overview", gui.focus == panelDiffContent, width, height, gui.state.DetailContent)
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
	diffLines := gui.renderDiffContentPanel(diffWidth, height, coloredDiff)

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

func (gui *Gui) renderDiffContentPanel(width int, height int, content string) []string {
	if height <= 0 {
		return nil
	}
	file, ok := gui.currentDiffFile()
	if !ok || len(file.Lines) == 0 {
		return gui.renderDetailPanel("Diff", gui.focus == panelDiffContent, width, height, content)
	}
	innerHeight := height
	if height > 2 {
		innerHeight = height - 2
	}
	start := 0
	if gui.diffLineSelected >= innerHeight {
		start = gui.diffLineSelected - innerHeight + 1
	}
	lines := make([]string, 0, innerHeight)
	for i := 0; len(lines) < innerHeight; i++ {
		idx := start + i
		if idx >= len(file.Lines) {
			lines = append(lines, "")
			continue
		}
		line := colorizeDiffContent(file.Lines[idx].Text)
		prefix := "  "
		location := gh.FormatDiffLineLocation(file.Lines[idx])
		if location != "" {
			prefix = padOrTrim(location, 7) + " "
		}
		inRange := gui.isDiffLineWithinPendingRange(file.Lines[idx])
		renderedPrefix := prefix
		if inRange {
			renderedPrefix = highlightPendingRangeLine(prefix)
		}
		rendered := renderedPrefix + styleDiffContentLine(line, false)
		if idx == gui.diffLineSelected {
			if inRange {
				rendered = highlightLine(renderedPrefix) + line
			} else {
				rendered = highlightLine(prefix) + line
			}
		}
		lines = append(lines, rendered)
	}
	return gui.framePanel("Diff", gui.focus == panelDiffContent, lines, width, height)
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

	repoLines := gui.framePanel("Repository", gui.focus == panelRepo, gui.renderRepoPanel(repoInnerHeight), width, repoPanelHeight)
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

func highlightPendingRangeLine(s string) string {
	if s == "" {
		return s
	}
	restyled := strings.ReplaceAll(s, ansiReset, ansiReset+ansiCyan)
	return ansiCyan + restyled + ansiReset
}

func styleDiffContentLine(s string, inRange bool) string {
	_ = inRange
	return s
}

func appendReviewSummaryLines(lines []string, summary string) []string {
	summary = strings.TrimSpace(summary)
	if summary == "" {
		return append(lines, "Summary: (empty)")
	}
	lines = append(lines, "Summary:")
	for _, line := range strings.Split(summary, "\n") {
		lines = append(lines, "  "+line)
	}
	return lines
}

func (gui *Gui) renderReviewDrawer(width int, height int) []string {
	if height <= 0 {
		return nil
	}
	innerHeight := height
	if height > 2 {
		innerHeight = height - 2
	}
	lines := make([]string, 0, innerHeight)

	summary := gui.state.Review.Summary
	if gui.state.Review.InputMode == core.ReviewInputSummary {
		summary = gui.summaryEditor.Value()
	}
	lines = appendReviewSummaryLines(lines, summary)
	modeLabel := "single-line"
	if gui.state.Review.RangeStart != nil {
		modeLabel = "range-selecting"
	}
	lines = append(lines, "Comment mode: "+modeLabel)

	if gui.state.Review.RangeStart != nil {
		lines = append(lines, "Range start: "+gui.state.Review.RangeStart.Path+":"+strconv.Itoa(gui.state.Review.RangeStart.Line))
	}

	if len(gui.state.Review.Comments) == 0 {
		lines = append(lines, "Comments: none")
	} else {
		lines = append(lines, "Comments:")
		for _, comment := range gui.state.Review.Comments {
			lines = append(lines, "  - "+renderReviewCommentSummary(comment))
		}
	}

	if notice := strings.TrimSpace(gui.state.Review.Notice); notice != "" {
		lines = append(lines, "")
		lines = append(lines, "Notice: "+notice)
	}

	if gui.state.Review.InputMode == core.ReviewInputComment {
		lines = append(lines, "")
		lines = append(lines, "Comment Input [Ctrl+S save / Esc cancel]")
		lines = append(lines, strings.Split(gui.commentEditor.View(), "\n")...)
	}
	if gui.state.Review.InputMode == core.ReviewInputSummary {
		lines = append(lines, "")
		lines = append(lines, "Summary Input [Ctrl+S save / Esc cancel]")
		lines = append(lines, strings.Split(gui.summaryEditor.View(), "\n")...)
	}

	for len(lines) < innerHeight {
		lines = append(lines, "")
	}
	return gui.framePanel("Review", gui.focus == panelReviewDrawer, lines[:innerHeight], width, height)
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
