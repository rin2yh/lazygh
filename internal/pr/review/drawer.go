package review

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rin2yh/lazygh/pkg/gui/widget"
	"github.com/rin2yh/lazygh/pkg/sanitize"
)

const (
	commentSummaryMaxLen      = 48
	commentStalePrefix        = "[stale] "
	CommentModeSingleLine     = "single-line"
	CommentModeRangeSelecting = "range-selecting"
)

type DrawerRange struct {
	Path string
	Line int
}

func (r DrawerRange) String() string {
	if r.Path == "" || r.Line <= 0 {
		return ""
	}
	return filepath.Base(r.Path) + ":" + strconv.Itoa(r.Line)
}

type DrawerComment struct {
	Path      string
	Line      int
	StartLine int
	Body      string
	Stale     bool
}

func (c DrawerComment) Summary() string {
	location := DrawerRange{Path: c.Path, Line: c.Line}.String()
	if c.StartLine > 0 {
		location = DrawerRange{Path: c.Path, Line: c.StartLine}.String() + "-" + strconv.Itoa(c.Line)
	}
	body := c.sanitize()
	if len(body) > commentSummaryMaxLen {
		body = body[:commentSummaryMaxLen] + "..."
	}
	summary := location + " " + body
	if c.Stale {
		summary = commentStalePrefix + summary
	}
	return summary
}

func (c DrawerComment) sanitize() string {
	return strings.TrimSpace(sanitize.SingleLine(c.Body))
}

// DrawerThreadComment holds display data for a single comment within a thread.
type DrawerThreadComment struct {
	Author string
	Body   string
}

// DrawerThread holds display data for an existing review thread.
type DrawerThread struct {
	Path       string
	Line       int
	DiffSide   string
	IsResolved bool
	IsOutdated bool
	Comments   []DrawerThreadComment
}

func (t DrawerThread) Summary() string {
	status := "open"
	if t.IsResolved {
		status = "resolved"
	} else if t.IsOutdated {
		status = "outdated"
	}
	location := DrawerRange{Path: t.Path, Line: t.Line}.String()
	replies := len(t.Comments)
	if replies == 0 {
		return location + " [" + status + "] (no comments)"
	}
	first := strings.TrimSpace(sanitize.SingleLine(t.Comments[0].Body))
	if len(first) > commentSummaryMaxLen {
		first = first[:commentSummaryMaxLen] + "..."
	}
	return location + " [" + status + "] " + first + " (" + strconv.Itoa(replies) + ")"
}

type Input struct {
	SummaryLines       []string
	CommentModeLabel   string
	EventLabel         string
	RangeStart         *DrawerRange
	AnchorConflict     bool
	Comments           []DrawerComment
	SelectedCommentIdx int
	Notice             string
	CommentInputLines  []string
	SummaryInputLines  []string
	Threads            []DrawerThread
	SelectedThreadIdx  int
	ThreadReplyLines   []string
}

func RenderDrawer(input Input, style widget.PanelStyle, width, height int) []string {
	if height <= 0 {
		return nil
	}
	inner := widget.InnerPanelHeight(height)
	lines := make([]string, 0, inner)
	if len(input.SummaryLines) == 0 {
		lines = append(lines, "Summary: (empty)")
	} else {
		lines = append(lines, "Summary:")
		for _, line := range input.SummaryLines {
			lines = append(lines, "  "+line)
		}
	}
	lines = append(lines, "Comment mode: "+input.CommentModeLabel)
	if input.EventLabel != "" {
		lines = append(lines, "Event: "+input.EventLabel)
	}
	if input.RangeStart != nil {
		if label := input.RangeStart.String(); label != "" {
			if input.AnchorConflict {
				lines = append(lines, "Range start: "+label+" [different file – range will be cleared]")
			} else {
				lines = append(lines, "Range start: "+label)
			}
		}
	}
	if len(input.Comments) == 0 {
		lines = append(lines, "Comments: none")
	} else {
		lines = append(lines, "Comments:")
		for i, comment := range input.Comments {
			prefix := "  - "
			if i == input.SelectedCommentIdx {
				prefix = "  > "
			}
			lines = append(lines, prefix+comment.Summary())
		}
	}
	if notice := strings.TrimSpace(input.Notice); notice != "" {
		lines = append(lines, "")
		lines = append(lines, "Notice: "+notice)
	}
	if len(input.CommentInputLines) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Comment Input [Ctrl+S save / Esc cancel]")
		lines = append(lines, input.CommentInputLines...)
	}
	if len(input.SummaryInputLines) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Summary Input [Ctrl+S save / Esc cancel]")
		lines = append(lines, input.SummaryInputLines...)
	}
	if len(input.Threads) == 0 {
		lines = append(lines, "")
		lines = append(lines, "Threads: none")
	} else {
		lines = append(lines, "")
		lines = append(lines, "Threads:")
		for i, thread := range input.Threads {
			prefix := "  - "
			if i == input.SelectedThreadIdx {
				prefix = "  > "
			}
			lines = append(lines, prefix+thread.Summary())
			if i == input.SelectedThreadIdx {
				for _, c := range thread.Comments {
					author := c.Author
					if author == "" {
						author = "unknown"
					}
					body := strings.TrimSpace(c.Body)
					for _, l := range strings.Split(body, "\n") {
						lines = append(lines, "    ["+author+"] "+l)
					}
				}
			}
		}
	}
	if len(input.ThreadReplyLines) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Reply Input [Ctrl+S save / Esc cancel]")
		lines = append(lines, input.ThreadReplyLines...)
	}
	return widget.FramePanel("Review", lines, width, height, style)
}
