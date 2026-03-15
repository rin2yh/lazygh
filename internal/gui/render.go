package gui

import (
	"strings"

	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/diff"
	"github.com/rin2yh/lazygh/internal/gui/draw"
	"github.com/rin2yh/lazygh/internal/gui/layout"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

func (gui *Gui) render() string {
	return draw.New(gui.buildRenderInput()).String()
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
	wrapped := widget.WrapText(content, width)
	if gui.detailViewportBody != wrapped {
		gui.detailViewport.SetContent(wrapped)
		gui.detailViewportBody = wrapped
		gui.detailViewport.GotoTop()
	}
}

func (gui *Gui) buildRenderInput() draw.Input {
	rightLines := gui.currentDetailLines(gui.state.DetailContent)
	if gui.state.IsDiffMode() {
		rightLines = gui.currentDetailLines(diff.ColorizeContent(gui.currentDiffContent()))
	}

	return draw.Input{
		Width:  gui.state.Width,
		Height: gui.state.Height,
		StatusLine: layout.Status{
			Loading:         gui.state.Loading != core.LoadingNone,
			DiffMode:        gui.state.IsDiffMode(),
			HasPR:           len(gui.state.PRs) > 0,
			Focus:           gui.renderFocus(),
			HasFiles:        len(gui.diffFiles) > 0,
			HasReviewDrawer: gui.review.ShouldShowDrawer(),
			InputMode:       gui.state.Review.InputMode,
			Keys:            gui.config.KeyBindings,
		}.String(),
		Theme: draw.Theme{
			ActiveBorderColor:   gui.config.Theme.ActiveBorderColor,
			InactiveBorderColor: gui.config.Theme.InactiveBorderColor,
		},
		Focus: gui.renderFocus(),
		Left: draw.LeftPanelsInput{
			Repo:       gui.state.Repo,
			PRsLoading: gui.state.PRsLoading,
			PRs:        gui.renderPRItems(),
			PRSelected: gui.state.PRsSelected,
		},
		Right: draw.RightPanelsInput{
			DiffMode:         gui.state.IsDiffMode(),
			OverviewTitle:    "Overview",
			OverviewLines:    rightLines,
			DiffFiles:        gui.diffFiles,
			DiffFileSelected: gui.diffFileSelected,
			DiffContentLines: gui.renderDiffContentLines(),
		},
		Review: gui.buildReviewDrawerInput(),
	}
}

func (gui *Gui) renderFocus() layout.Focus {
	switch gui.focus {
	case panelRepo:
		return layout.FocusRepo
	case panelPRs:
		return layout.FocusPRs
	case panelDiffFiles:
		return layout.FocusDiffFiles
	case panelReviewDrawer:
		return layout.FocusReviewDrawer
	default:
		return layout.FocusDiffContent
	}
}

func (gui *Gui) renderPRItems() []string {
	items := make([]string, 0, len(gui.state.PRs))
	for _, pr := range gui.state.PRs {
		items = append(items, core.FormatPRItem(pr))
	}
	return items
}

func (gui *Gui) currentDetailLines(content string) []string {
	dims := layout.New(gui.state.Width, gui.state.Height, gui.state.IsDiffMode(), gui.review.ShouldShowDrawer())
	innerWidth := dims.RightWidth
	if gui.state.IsDiffMode() && dims.RightWidth >= 20 {
		filesWidth := dims.RightWidth * 30 / 100
		if filesWidth < 16 {
			filesWidth = 16
		}
		if filesWidth > dims.RightWidth-10 {
			filesWidth = dims.RightWidth - 10
		}
		innerWidth = dims.RightWidth - filesWidth - 1
	}
	innerWidth = dims.InnerWidth(innerWidth)
	bodyHeight := dims.InnerHeight(dims.MainHeight)
	if bodyHeight < 1 {
		bodyHeight = 1
	}
	gui.syncDetailViewport(innerWidth, bodyHeight, content)
	return strings.Split(gui.detailViewport.View(), "\n")
}

func (gui *Gui) renderDiffContentLines() []draw.DiffContentLine {
	file, ok := gui.currentDiffFile()
	if !ok || len(file.Lines) == 0 {
		return nil
	}
	lines := make([]draw.DiffContentLine, 0, len(file.Lines))
	for idx, line := range file.Lines {
		lines = append(lines, draw.DiffContentLine{
			Location: gh.FormatDiffLineLocation(line),
			Text:     diff.ColorizeContent(line.Text),
			Selected: idx == gui.diffLineSelected,
			InRange:  gui.review.IsLineWithinPendingRange(line),
		})
	}
	return lines
}

func (gui *Gui) buildReviewDrawerInput() *draw.ReviewDrawerInput {
	if !gui.review.ShouldShowDrawer() {
		return nil
	}
	summary := gui.state.Review.Summary
	if gui.state.Review.InputMode == core.ReviewInputSummary {
		summary = gui.review.CurrentSummaryValue()
	}
	input := &draw.ReviewDrawerInput{
		SummaryLines:     splitNonEmptyLines(summary),
		CommentModeLabel: "single-line",
		Notice:           gui.state.Review.Notice,
	}
	if gui.state.Review.RangeStart != nil {
		input.CommentModeLabel = "range-selecting"
		input.RangeStart = &draw.Range{
			Path: gui.state.Review.RangeStart.Path,
			Line: gui.state.Review.RangeStart.Line,
		}
	}
	for _, comment := range gui.state.Review.Comments {
		input.Comments = append(input.Comments, draw.Comment{
			Path:      comment.Path,
			Line:      comment.Line,
			StartLine: comment.StartLine,
			Body:      comment.Body,
		})
	}
	if gui.state.Review.InputMode == core.ReviewInputComment {
		input.CommentInputLines = gui.review.CommentInputLines()
	}
	if gui.state.Review.InputMode == core.ReviewInputSummary {
		input.SummaryInputLines = gui.review.SummaryInputLines()
	}
	return input
}

func splitNonEmptyLines(content string) []string {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	return strings.Split(content, "\n")
}
