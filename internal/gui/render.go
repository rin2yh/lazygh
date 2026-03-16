package gui

import (
	"strings"

	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	guidiff "github.com/rin2yh/lazygh/internal/gui/diff"
	"github.com/rin2yh/lazygh/internal/gui/help"
	"github.com/rin2yh/lazygh/internal/gui/layout"
	"github.com/rin2yh/lazygh/internal/gui/prs"
	guireview "github.com/rin2yh/lazygh/internal/gui/review"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

func (gui *Gui) render() string {
	isDiff := gui.state.IsDiffMode()
	showDrawer := gui.review.ShouldShowDrawer()
	screen := layout.New(gui.state.Width, gui.state.Height, isDiff, showDrawer)
	focus := gui.renderFocus()
	statusLine := layout.Status{
		Loading:   gui.state.Detail.Loading != core.LoadingNone,
		DiffMode:  isDiff,
		Focus:     focus,
		InputMode: gui.state.Review.InputMode,
		Keys:      gui.config.KeyBindings,
	}.String()

	leftInput := prs.Input{
		Repo:       gui.state.List.Repo,
		PRsLoading: gui.state.List.PRsLoading,
		PRs:        gui.renderPRItems(),
		PRSelected: gui.state.List.PRsSelected,
	}

	var rightLines []string
	if isDiff {
		rightLines = gui.currentDetailLines(screen, guidiff.ColorizeContent(gui.currentDiffContent()))
	} else {
		rightLines = gui.currentDetailLines(screen, gui.state.Detail.Content)
	}
	rightInput := guidiff.PanelInput{
		DiffMode:      isDiff,
		OverviewTitle: "Overview",
		OverviewLines: rightLines,
	}
	if isDiff {
		rightInput.DiffFiles = gui.diff.Files()
		rightInput.DiffFileSelected = gui.diff.FileSelected()
		rightInput.DiffContentLines = gui.renderDiffContentLines()
	}

	leftLines := prs.RenderLeft(leftInput, screen.RepoHeight, screen.PRHeight,
		func(f layout.Focus) bool { return focus == f },
		gui.style,
		screen.LeftWidth,
	)
	rightPanelLines := gui.renderRight(rightInput, screen, focus)

	lines := widget.JoinColumns(leftLines, screen.LeftWidth, rightPanelLines, screen.RightWidth, screen.MainHeight)
	drawerInput := gui.buildReviewDrawerInput(showDrawer)
	if drawerInput != nil && screen.DrawerHeight > 0 {
		drawerActive := focus == layout.FocusReviewDrawer
		for _, line := range guireview.RenderDrawer(*drawerInput, gui.style(drawerActive), screen.Width, screen.DrawerHeight) {
			lines = append(lines, widget.PadOrTrim(line, screen.Width))
		}
	}
	lines = append(lines, widget.PadOrTrim(statusLine, screen.Width))

	if gui.showHelp {
		lines = help.RenderOverlay(lines, gui.config.KeyBindings, screen.Width)
	}

	var b strings.Builder
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func (gui *Gui) renderRight(input guidiff.PanelInput, screen layout.Screen, focus layout.Focus) []string {
	diffActive := focus == layout.FocusDiffContent
	if !input.DiffMode {
		return widget.FramePanel(input.OverviewTitle, input.OverviewLines, screen.RightWidth, screen.MainHeight, gui.style(diffActive))
	}
	filesWidth, diffWidth := layout.DiffSplitWidths(screen.RightWidth)
	if filesWidth == 0 {
		lines := guidiff.ContentLines(input, screen.MainHeight)
		if lines == nil {
			lines = input.OverviewLines
		}
		return widget.FramePanel("Diff", lines, screen.RightWidth, screen.MainHeight, gui.style(diffActive))
	}
	filesActive := focus == layout.FocusDiffFiles
	filesLines := guidiff.RenderFiles(input, gui.style(filesActive), filesWidth, screen.MainHeight)
	diffLines := guidiff.RenderContent(input, gui.style(diffActive), diffWidth, screen.MainHeight)
	return widget.JoinColumns(filesLines, filesWidth, diffLines, diffWidth, screen.MainHeight)
}

func (gui *Gui) style(active bool) widget.PanelStyle {
	if active {
		borderColor := widget.ResolveColorName(gui.config.Theme.ActiveBorderColor, "green")
		return widget.PanelStyle{BorderColor: borderColor, TitleColor: borderColor}
	}
	borderColor := widget.ResolveColorName(gui.config.Theme.InactiveBorderColor, "white")
	return widget.PanelStyle{BorderColor: borderColor}
}

func (gui *Gui) syncDetailViewport(width int, height int, content string) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	gui.detail.Sync(width, height, widget.WrapText(content, width))
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
	items := make([]string, 0, len(gui.state.List.PRs))
	for _, pr := range gui.state.List.PRs {
		items = append(items, core.FormatPRItem(pr))
	}
	return items
}

func (gui *Gui) currentDetailLines(dims layout.Screen, content string) []string {
	innerWidth := dims.RightWidth
	if gui.state.IsDiffMode() {
		filesWidth, diffWidth := layout.DiffSplitWidths(dims.RightWidth)
		if filesWidth > 0 {
			innerWidth = diffWidth
		}
	}
	innerWidth = dims.InnerWidth(innerWidth)
	bodyHeight := dims.InnerHeight(dims.MainHeight)
	if bodyHeight < 1 {
		bodyHeight = 1
	}
	gui.syncDetailViewport(innerWidth, bodyHeight, content)
	return strings.Split(gui.detail.View(), "\n")
}

func (gui *Gui) renderDiffContentLines() []guidiff.ContentLine {
	file, ok := gui.diff.CurrentFile()
	if !ok || len(file.Lines) == 0 {
		return nil
	}
	lineSelected := gui.diff.LineSelected()
	lines := make([]guidiff.ContentLine, 0, len(file.Lines))
	for idx, line := range file.Lines {
		lines = append(lines, guidiff.ContentLine{
			Location: gh.FormatDiffLineLocation(line),
			Text:     guidiff.ColorizeLine(line.Text),
			Selected: idx == lineSelected,
			InRange:  gui.review.IsIndexWithinPendingRange(line.Path, line.Commentable, idx),
		})
	}
	return lines
}

func (gui *Gui) buildReviewDrawerInput(showDrawer bool) *guireview.DrawerInput {
	if !showDrawer {
		return nil
	}
	summary := gui.state.Review.Summary
	if gui.state.Review.InputMode == core.ReviewInputSummary {
		summary = gui.review.CurrentSummaryValue()
	}
	input := &guireview.DrawerInput{
		SummaryLines:     splitNonEmptyLines(summary),
		CommentModeLabel: guireview.CommentModeSingleLine,
		Notice:           gui.state.Review.Notice,
	}
	if gui.state.Review.RangeStart != nil {
		input.CommentModeLabel = guireview.CommentModeRangeSelecting
		input.RangeStart = &guireview.DrawerRange{
			Path: gui.state.Review.RangeStart.Path,
			Line: gui.state.Review.RangeStart.Line,
		}
	}
	input.Comments = make([]guireview.DrawerComment, 0, len(gui.state.Review.Comments))
	for _, comment := range gui.state.Review.Comments {
		input.Comments = append(input.Comments, guireview.DrawerComment{
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
