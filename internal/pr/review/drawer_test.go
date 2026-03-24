package review

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

func TestRenderCommentSummary(t *testing.T) {
	tests := []struct {
		name string
		c    DrawerComment
		want string
	}{
		{
			name: "range comment",
			c:    DrawerComment{Path: "a.txt", Line: 12, StartLine: 10, Body: "hello\nworld"},
			want: "a.txt:10-12 hello world",
		},
		{
			name: "stale comment prepends marker",
			c:    DrawerComment{Path: "a.txt", Line: 5, Body: "fix this", Stale: true},
			want: "[stale] a.txt:5 fix this",
		},
		{
			name: "long path uses basename",
			c:    DrawerComment{Path: "internal/pr/review/types.go", Line: 3, Body: "ok"},
			want: "types.go:3 ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Summary(); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderDrawer_AnchorConflictWarning(t *testing.T) {
	lines := RenderDrawer(Input{
		CommentModeLabel: CommentModeRangeSelecting,
		RangeStart:       &DrawerRange{Path: "a.go", Line: 5},
		AnchorConflict:   true,
	}, widget.PanelStyle{}, 80, 10)

	var found bool
	for _, line := range lines {
		if strings.Contains(xansi.Strip(line), "different file") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected anchor conflict warning, got: %#v", lines)
	}
}

func TestRenderReviewDrawer_RendersMultilineSummaryWithoutBreakingLayout(t *testing.T) {
	lines := RenderDrawer(Input{
		SummaryLines:     []string{"first line", "second line"},
		CommentModeLabel: CommentModeSingleLine,
	}, widget.PanelStyle{BorderColor: "green", TitleColor: "green"}, 40, 8)

	var foundHeader, foundFirst, foundSecond bool
	for _, line := range lines {
		stripped := xansi.Strip(line)
		if strings.Contains(stripped, "Summary:") {
			foundHeader = true
		}
		if strings.Contains(stripped, "  first line") {
			foundFirst = true
		}
		if strings.Contains(stripped, "  second line") {
			foundSecond = true
		}
	}
	if !foundHeader || !foundFirst || !foundSecond {
		t.Fatalf("missing multiline summary rendering: %#v", lines)
	}
}
