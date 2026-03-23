package app

import (
	"strings"

	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/help"
	"github.com/rin2yh/lazygh/internal/model"
	"github.com/rin2yh/lazygh/internal/pr/diff"
	prhelp "github.com/rin2yh/lazygh/internal/pr/help"
	"github.com/rin2yh/lazygh/internal/pr/list"
	"github.com/rin2yh/lazygh/internal/pr/review"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

func (gui *Gui) render() string {
	isDiff := gui.coord.IsDiffMode()
	showDrawer := gui.review.ShouldShowDrawer()
	screen := layout.New(gui.coord.Width, gui.coord.Height, isDiff, showDrawer)
	focus := gui.focus
	statusLine := layout.Status{
		Fetching:  gui.coord.Overview.Fetching != model.FetchNone,
		DiffMode:  isDiff,
		Focus:     layout.Focus(focus),
		InputMode: gui.review.InputMode(),
		Keys:      gui.config.KeyBindings,
	}.String()

	leftInput := list.PanelInput{
		Repo:     gui.coord.Repo,
		Fetching: gui.coord.Fetching,
		Items:    gui.coord.Items,
		Selected: gui.coord.Selected,
		Filter:   gui.coord.Filter.Label(),
	}

	var rightLines []string
	if isDiff {
		rightLines = gui.currentDetailLines(screen, diff.ColorizeContent(gui.currentDiffContent()))
	} else {
		rightLines = gui.currentDetailLines(screen, gui.coord.Overview.Content)
	}
	rightInput := diff.PanelInput{
		DiffMode:      isDiff,
		OverviewTitle: "Overview",
		OverviewLines: rightLines,
	}
	if isDiff {
		rightInput.DiffFiles = gui.diff.Files()
		rightInput.DiffFileSelected = gui.diff.FileSelected()
		rightInput.DiffContentLines = diff.BuildContentLines(&gui.diff, gui.review.IsIndexWithinPendingRange)
	}

	leftLines := list.RenderLeft(leftInput, screen.RepoHeight, screen.PRHeight,
		func(f layout.Focus) bool { return focus == f },
		gui.style,
		screen.LeftWidth,
	)
	rightPanelLines := gui.renderRight(rightInput, screen, focus)

	lines := widget.JoinColumns(leftLines, screen.LeftWidth, rightPanelLines, screen.RightWidth, screen.MainHeight)
	drawerInput := gui.review.BuildDrawerInput(showDrawer)
	if drawerInput != nil && screen.DrawerHeight > 0 {
		drawerActive := focus == layout.FocusReviewDrawer
		for _, line := range review.RenderDrawer(*drawerInput, gui.style(drawerActive), screen.Width, screen.DrawerHeight) {
			lines = append(lines, widget.PadOrTrim(line, screen.Width))
		}
	}
	lines = append(lines, widget.PadOrTrim(statusLine, screen.Width))

	if gui.coord.FilterOpen {
		lines = applyFilterOverlay(lines, gui.coord.Filter, gui.coord.FilterCursor, screen.Width)
	}
	if gui.showHelp {
		sections := append(help.CommonSections(gui.config.KeyBindings), prhelp.Sections(gui.config.KeyBindings)...)
		lines = help.RenderOverlay(lines, sections, gui.config.KeyBindings, screen.Width)
	}

	var b strings.Builder
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func (gui *Gui) renderRight(input diff.PanelInput, screen layout.Screen, focus layout.Focus) []string {
	diffActive := focus == layout.FocusDiffContent
	if !input.DiffMode {
		return widget.FramePanel(input.OverviewTitle, input.OverviewLines, screen.RightWidth, screen.MainHeight, gui.style(diffActive))
	}
	filesWidth, diffWidth := layout.DiffSplitWidths(screen.RightWidth)
	if filesWidth == 0 {
		lines := diff.ContentLines(input, screen.MainHeight)
		if lines == nil {
			lines = input.OverviewLines
		}
		return widget.FramePanel("Diff", lines, screen.RightWidth, screen.MainHeight, gui.style(diffActive))
	}
	filesActive := focus == layout.FocusDiffFiles
	filesLines := diff.RenderFiles(input, gui.style(filesActive), filesWidth, screen.MainHeight)
	diffLines := diff.RenderContent(input, gui.style(diffActive), diffWidth, screen.MainHeight)
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

func (gui *Gui) currentDetailLines(dims layout.Screen, content string) []string {
	innerWidth := dims.RightWidth
	if gui.coord.IsDiffMode() {
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

func applyFilterOverlay(background []string, filter model.PRFilterMask, cursor int, screenWidth int) []string {
	panelLines, panelW := list.FilterPanelLines(filter, cursor)
	return widget.OverlayPanel(background, panelLines, panelW, screenWidth)
}
