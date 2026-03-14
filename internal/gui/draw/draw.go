package draw

import (
	"strconv"
	"strings"

	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/diff"
	"github.com/rin2yh/lazygh/internal/gui/layout"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

const (
	ansiReset   = "\x1b[0m"
	ansiReverse = "\x1b[7m"
	ansiCyan    = "\x1b[36m"
)

type Theme struct {
	ActiveBorderColor   string
	InactiveBorderColor string
}

type Input struct {
	Width      int
	Height     int
	StatusLine string
	Theme      Theme
	Focus      layout.Focus
	Left       LeftPanelsInput
	Right      RightPanelsInput
	Review     *ReviewDrawerInput
}

type LeftPanelsInput struct {
	Repo       string
	PRsLoading bool
	PRs        []string
	PRSelected int
}

type RightPanelsInput struct {
	DiffMode         bool
	OverviewTitle    string
	OverviewLines    []string
	DiffFiles        []gh.DiffFile
	DiffFileSelected int
	DiffContentLines []DiffContentLine
}

type DiffContentLine struct {
	Location string
	Text     string
	Selected bool
	InRange  bool
}

type Range struct {
	Path string
	Line int
}

func (r Range) String() string {
	if r.Path == "" || r.Line <= 0 {
		return ""
	}
	return r.Path + ":" + strconv.Itoa(r.Line)
}

type Comment struct {
	Path      string
	Line      int
	StartLine int
	Body      string
}

func (c Comment) Summary() string {
	location := c.Path + ":" + strconv.Itoa(c.Line)
	if c.StartLine > 0 {
		location = c.Path + ":" + strconv.Itoa(c.StartLine) + "-" + strconv.Itoa(c.Line)
	}
	body := c.sanitize()
	if len(body) > 48 {
		body = body[:48] + "..."
	}
	return location + " " + body
}

func (c Comment) sanitize() string {
	body := strings.ReplaceAll(c.Body, "\n", " ")
	body = strings.ReplaceAll(body, "\r", " ")
	return strings.TrimSpace(body)
}

type ReviewDrawerInput struct {
	SummaryLines      []string
	CommentModeLabel  string
	RangeStart        *Range
	Comments          []Comment
	Notice            string
	CommentInputLines []string
	SummaryInputLines []string
}

type View struct {
	input  Input
	screen layout.Screen
}

func New(input Input) View {
	return View{
		input:  input,
		screen: layout.New(input.Width, input.Height, input.Right.DiffMode, input.Review != nil),
	}
}

func (v View) String() string {
	leftLines := v.left(v.screen.LeftWidth, v.screen.MainHeight)
	rightLines := v.right(v.screen.RightWidth, v.screen.MainHeight)

	var b strings.Builder
	for i := 0; i < v.screen.MainHeight; i++ {
		left := ""
		if i < len(leftLines) {
			left = leftLines[i]
		}
		right := ""
		if i < len(rightLines) {
			right = rightLines[i]
		}
		b.WriteString(widget.PadOrTrim(left, v.screen.LeftWidth))
		b.WriteRune(' ')
		b.WriteString(widget.PadOrTrim(right, v.screen.RightWidth))
		b.WriteByte('\n')
	}
	if v.input.Review != nil && v.screen.DrawerHeight > 0 {
		for _, line := range v.drawer(v.screen.Width, v.screen.DrawerHeight) {
			b.WriteString(widget.PadOrTrim(line, v.screen.Width))
			b.WriteByte('\n')
		}
	}
	b.WriteString(widget.PadOrTrim(v.input.StatusLine, v.screen.Width))
	return b.String()
}

func (v View) left(width int, height int) []string {
	if height <= 0 {
		return nil
	}
	repoInnerHeight := v.screen.InnerHeight(v.screen.RepoHeight)
	prInnerHeight := v.screen.InnerHeight(v.screen.PRHeight)
	repoLines := v.panel("Repository", v.input.Focus == layout.FocusRepo, v.repo(repoInnerHeight), width, v.screen.RepoHeight)
	prLines := v.panel("PRs (Open/Draft)", v.input.Focus == layout.FocusPRs, v.pr(prInnerHeight), width, v.screen.PRHeight)
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

func (v View) right(width int, height int) []string {
	if !v.input.Right.DiffMode {
		return v.detail(v.input.Right.OverviewTitle, v.input.Focus == layout.FocusDiffContent, width, height, v.input.Right.OverviewLines)
	}
	if width < 20 {
		lines := v.diffContentLines(height)
		if len(v.input.Right.DiffContentLines) == 0 {
			lines = v.input.Right.OverviewLines
		}
		return v.detail("Diff", v.input.Focus == layout.FocusDiffContent, width, height, lines)
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

	filesLines := v.files(filesWidth, height)
	diffLines := v.diff(diffWidth, height)
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
		lines = append(lines, widget.PadOrTrim(left, filesWidth)+" "+widget.PadOrTrim(right, diffWidth))
	}
	return lines
}

func (v View) repo(height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)
	for len(lines) < height {
		if len(lines) == 0 {
			lines = append(lines, v.input.Left.Repo)
		} else {
			lines = append(lines, "")
		}
	}
	return lines
}

func (v View) pr(height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)
	if v.input.Left.PRsLoading {
		for len(lines) < height {
			lines = append(lines, "")
		}
		return lines
	}
	if len(v.input.Left.PRs) == 0 {
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
		if i >= len(v.input.Left.PRs) {
			lines = append(lines, "")
			continue
		}
		line := "  " + v.input.Left.PRs[i]
		if i == v.input.Left.PRSelected {
			line = v.highlight("> " + v.input.Left.PRs[i])
		}
		lines = append(lines, line)
	}
	return lines
}

func (v View) files(width int, height int) []string {
	if height <= 0 {
		return nil
	}
	inner := v.screen.InnerHeight(height)
	lines := make([]string, 0, inner)
	if len(v.input.Right.DiffFiles) == 0 {
		for len(lines) < inner {
			if len(lines) == 0 {
				lines = append(lines, "No changed files")
			} else {
				lines = append(lines, "")
			}
		}
		return v.panel("Files", v.input.Focus == layout.FocusDiffFiles, lines, width, height)
	}

	start := v.start(v.input.Right.DiffFileSelected, inner)
	for i := 0; len(lines) < inner; i++ {
		idx := start + i
		if idx >= len(v.input.Right.DiffFiles) {
			lines = append(lines, "")
			continue
		}
		line := "  " + diff.RenderFileListLine(v.input.Right.DiffFiles[idx])
		if idx == v.input.Right.DiffFileSelected {
			line = v.highlight("> " + diff.RenderFileListLine(v.input.Right.DiffFiles[idx]))
		}
		lines = append(lines, line)
	}
	return v.panel("Files", v.input.Focus == layout.FocusDiffFiles, lines, width, height)
}

func (v View) diff(width int, height int) []string {
	if len(v.input.Right.DiffContentLines) == 0 {
		return v.detail("Diff", v.input.Focus == layout.FocusDiffContent, width, height, v.input.Right.OverviewLines)
	}
	inner := v.screen.InnerHeight(height)
	start := v.start(v.selected(), inner)
	lines := make([]string, 0, inner)
	for i := 0; len(lines) < inner; i++ {
		idx := start + i
		if idx >= len(v.input.Right.DiffContentLines) {
			lines = append(lines, "")
			continue
		}
		lines = append(lines, v.diffLine(v.input.Right.DiffContentLines[idx]))
	}
	return v.panel("Diff", v.input.Focus == layout.FocusDiffContent, lines, width, height)
}

func (v View) detail(title string, active bool, width int, height int, content []string) []string {
	if height <= 0 {
		return nil
	}
	inner := v.screen.InnerHeight(height)
	lines := make([]string, 0, inner)
	for _, line := range content {
		if len(lines) >= inner {
			break
		}
		lines = append(lines, line)
	}
	for len(lines) < inner {
		lines = append(lines, "")
	}
	return v.panel(title, active, lines, width, height)
}

func (v View) drawer(width int, height int) []string {
	if v.input.Review == nil || height <= 0 {
		return nil
	}
	inner := v.screen.InnerHeight(height)
	lines := make([]string, 0, inner)
	if len(v.input.Review.SummaryLines) == 0 {
		lines = append(lines, "Summary: (empty)")
	} else {
		lines = append(lines, "Summary:")
		for _, line := range v.input.Review.SummaryLines {
			lines = append(lines, "  "+line)
		}
	}
	lines = append(lines, "Comment mode: "+v.input.Review.CommentModeLabel)
	if v.input.Review.RangeStart != nil {
		if label := v.input.Review.RangeStart.String(); label != "" {
			lines = append(lines, "Range start: "+label)
		}
	}
	if len(v.input.Review.Comments) == 0 {
		lines = append(lines, "Comments: none")
	} else {
		lines = append(lines, "Comments:")
		for _, comment := range v.input.Review.Comments {
			lines = append(lines, "  - "+comment.Summary())
		}
	}
	if notice := strings.TrimSpace(v.input.Review.Notice); notice != "" {
		lines = append(lines, "")
		lines = append(lines, "Notice: "+notice)
	}
	if len(v.input.Review.CommentInputLines) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Comment Input [Ctrl+S save / Esc cancel]")
		lines = append(lines, v.input.Review.CommentInputLines...)
	}
	if len(v.input.Review.SummaryInputLines) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Summary Input [Ctrl+S save / Esc cancel]")
		lines = append(lines, v.input.Review.SummaryInputLines...)
	}
	for len(lines) < inner {
		lines = append(lines, "")
	}
	return v.panel("Review", v.input.Focus == layout.FocusReviewDrawer, lines[:inner], width, height)
}

func (v View) panel(title string, active bool, content []string, width int, height int) []string {
	return widget.FramePanel(title, content, width, height, v.style(active))
}

func (v View) style(active bool) widget.PanelStyle {
	borderColor := "white"
	if active {
		borderColor = widget.ResolveColorName(v.input.Theme.ActiveBorderColor, "green")
		return widget.PanelStyle{BorderColor: borderColor, TitleColor: borderColor}
	}
	borderColor = widget.ResolveColorName(v.input.Theme.InactiveBorderColor, "white")
	return widget.PanelStyle{BorderColor: borderColor}
}

func (v View) diffContentLines(height int) []string {
	inner := v.screen.InnerHeight(height)
	if len(v.input.Right.DiffContentLines) == 0 {
		return nil
	}
	start := v.start(v.selected(), inner)
	out := make([]string, 0, inner)
	for i := 0; len(out) < inner; i++ {
		idx := start + i
		if idx >= len(v.input.Right.DiffContentLines) {
			out = append(out, "")
			continue
		}
		out = append(out, v.diffLine(v.input.Right.DiffContentLines[idx]))
	}
	return out
}

func (v View) diffLine(line DiffContentLine) string {
	prefix := "  "
	if line.Location != "" {
		prefix = widget.PadOrTrim(line.Location, 7) + " "
	}
	renderedPrefix := prefix
	if line.InRange {
		renderedPrefix = v.rangeHighlight(prefix)
	}
	rendered := renderedPrefix + line.Text
	if line.Selected {
		if line.InRange {
			return v.highlight(renderedPrefix) + line.Text
		}
		return v.highlight(prefix) + line.Text
	}
	return rendered
}

func (v View) selected() int {
	for i, line := range v.input.Right.DiffContentLines {
		if line.Selected {
			return i
		}
	}
	return 0
}

func (v View) start(selected int, visible int) int {
	if visible <= 0 || selected < visible {
		return 0
	}
	return selected - visible + 1
}

func (v View) highlight(s string) string {
	if s == "" {
		return s
	}
	restyled := strings.ReplaceAll(s, ansiReset, ansiReset+ansiReverse)
	return ansiReverse + restyled + ansiReset
}

func (v View) rangeHighlight(s string) string {
	if s == "" {
		return s
	}
	restyled := strings.ReplaceAll(s, ansiReset, ansiReset+ansiCyan)
	return ansiCyan + restyled + ansiReset
}
