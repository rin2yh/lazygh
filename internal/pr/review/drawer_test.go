package review

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

func TestRenderCommentSummary(t *testing.T) {
	got := DrawerComment{
		Path:      "a.txt",
		Line:      12,
		StartLine: 10,
		Body:      "hello\nworld",
	}.Summary()
	if got != "a.txt:10-12 hello world" {
		t.Fatalf("got %q", got)
	}
}

func TestRenderReviewDrawer_RendersMultilineSummaryWithoutBreakingLayout(t *testing.T) {
	lines := RenderDrawer(DrawerInput{
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
