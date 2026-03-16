package review

import (
	"strconv"
	"strings"

	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

const (
	commentSummaryMaxLen      = 48
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
	return r.Path + ":" + strconv.Itoa(r.Line)
}

type DrawerComment struct {
	Path      string
	Line      int
	StartLine int
	Body      string
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
	return location + " " + body
}

func (c DrawerComment) sanitize() string {
	return strings.TrimSpace(core.SanitizeSingleLine(c.Body))
}

type DrawerInput struct {
	SummaryLines       []string
	CommentModeLabel   string
	EventLabel         string
	RangeStart         *DrawerRange
	Comments           []DrawerComment
	SelectedCommentIdx int
	Notice             string
	CommentInputLines  []string
	SummaryInputLines  []string
}

func RenderDrawer(input DrawerInput, style widget.PanelStyle, width, height int) []string {
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
			lines = append(lines, "Range start: "+label)
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
	return widget.FramePanel("Review", lines, width, height, style)
}
