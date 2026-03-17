package gui

import (
	"strings"

	"github.com/rin2yh/lazygh/internal/gh"
	guidiff "github.com/rin2yh/lazygh/internal/gui/diff"
	"github.com/rin2yh/lazygh/internal/gui/help"
	"github.com/rin2yh/lazygh/internal/gui/layout"
	"github.com/rin2yh/lazygh/internal/model"
	"github.com/rin2yh/lazygh/internal/pr"
	"github.com/rin2yh/lazygh/internal/review"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

func (gui *Gui) render() string {
	isDiff := gui.state.IsDiffMode()
	showDrawer := gui.review.ShouldShowDrawer()
	screen := layout.New(gui.state.Width, gui.state.Height, isDiff, showDrawer)
	focus := gui.renderFocus()
	statusLine := layout.Status{
		Loading:   gui.state.Detail.Loading != model.LoadingNone,
		DiffMode:  isDiff,
		Focus:     focus,
		InputMode: gui.review.InputMode(),
		Keys:      gui.config.KeyBindings,
	}.String()

	leftInput := pr.PanelInput{
		Repo:     gui.state.Repo,
		Fetching: gui.state.Fetching,
		Items:    gui.state.Items,
		Selected: gui.state.Selected,
		Filter:   gui.state.Filter.Label(),
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

	leftLines := pr.RenderLeft(leftInput, screen.RepoHeight, screen.PRHeight,
		func(f layout.Focus) bool { return focus == f },
		gui.style,
		screen.LeftWidth,
	)
	rightPanelLines := gui.renderRight(rightInput, screen, focus)

	lines := widget.JoinColumns(leftLines, screen.LeftWidth, rightPanelLines, screen.RightWidth, screen.MainHeight)
	drawerInput := gui.buildReviewDrawerInput(showDrawer)
	if drawerInput != nil && screen.DrawerHeight > 0 {
		drawerActive := focus == layout.FocusReviewDrawer
		for _, line := range review.RenderDrawer(*drawerInput, gui.style(drawerActive), screen.Width, screen.DrawerHeight) {
			lines = append(lines, widget.PadOrTrim(line, screen.Width))
		}
	}
	lines = append(lines, widget.PadOrTrim(statusLine, screen.Width))

	if gui.state.FilterOpen {
		lines = applyFilterOverlay(lines, gui.state.Filter, gui.state.FilterCursor, screen.Width)
	}
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

func (gui *Gui) buildReviewDrawerInput(showDrawer bool) *review.DrawerInput {
	if !showDrawer {
		return nil
	}
	inputMode := gui.review.InputMode()
	summary := gui.review.Summary()
	if inputMode == model.ReviewInputSummary {
		summary = gui.review.CurrentSummaryValue()
	}
	input := &review.DrawerInput{
		SummaryLines:     splitNonEmptyLines(summary),
		CommentModeLabel: review.CommentModeSingleLine,
		EventLabel:       gui.review.EventLabel(),
		Notice:           gui.review.Notice(),
	}
	if rs := gui.review.RangeStart(); rs != nil {
		input.CommentModeLabel = review.CommentModeRangeSelecting
		input.RangeStart = &review.DrawerRange{
			Path: rs.Path,
			Line: rs.Line,
		}
	}
	comments := gui.review.Comments()
	input.Comments = make([]review.DrawerComment, 0, len(comments))
	for _, comment := range comments {
		input.Comments = append(input.Comments, review.DrawerComment{
			Path:      comment.Path,
			Line:      comment.Line,
			StartLine: comment.StartLine,
			Body:      comment.Body,
		})
	}
	input.SelectedCommentIdx = gui.review.SelectedCommentIdx()
	if inputMode == model.ReviewInputComment {
		input.CommentInputLines = gui.review.CommentInputLines()
	}
	if inputMode == model.ReviewInputSummary {
		input.SummaryInputLines = gui.review.SummaryInputLines()
	}
	return input
}

func applyFilterOverlay(background []string, filter model.PRFilterMask, cursor int, screenWidth int) []string {
	panelLines, panelW := pr.FilterPanelLines(filter, cursor)
	return widget.OverlayPanel(background, panelLines, panelW, screenWidth)
}

func splitNonEmptyLines(content string) []string {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	return strings.Split(content, "\n")
}
