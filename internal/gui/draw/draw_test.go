package draw

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/layout"
)

func TestRenderRightPanels_DiffModeHasFilesPanel(t *testing.T) {
	out := New(Input{
		Width:  60,
		Height: 10,
		Focus:  layout.FocusDiffFiles,
		Right: RightPanelsInput{
			DiffMode: true,
			DiffFiles: []gh.DiffFile{
				{Path: "a.txt", Status: gh.DiffFileStatusModified, Additions: 1, Deletions: 1},
			},
			DiffContentLines: []DiffContentLine{{Location: "+1", Text: "+new", Selected: true}},
		},
	}).String()

	lines := strings.Split(out, "\n")
	if len(lines) < 1 || !strings.Contains(lines[0], "Files") {
		t.Fatalf("line does not contain %q: %q", "Files", lines[0])
	}
}

func TestRenderRightPanels_NarrowDiffModeFallsBackToOverviewLines(t *testing.T) {
	lines := New(Input{
		Width:  24,
		Height: 10,
		Focus:  layout.FocusDiffContent,
		Right: RightPanelsInput{
			DiffMode:      true,
			OverviewLines: []string{"@@ -1 +1 @@", "+raw diff line"},
		},
	}).right(19, 9)

	joined := xansi.Strip(strings.Join(lines, "\n"))
	if !strings.Contains(joined, "+raw diff line") {
		t.Fatalf("narrow diff view did not render overview fallback: %q", joined)
	}
}

func TestRenderRightPanels_NarrowDiffModePrefersStructuredDiffLines(t *testing.T) {
	lines := New(Input{
		Width:  24,
		Height: 10,
		Focus:  layout.FocusDiffContent,
		Right: RightPanelsInput{
			DiffMode:         true,
			OverviewLines:    []string{"+raw diff line"},
			DiffContentLines: []DiffContentLine{{Location: "+1", Text: "+new", Selected: true}},
		},
	}).right(19, 9)

	joined := xansi.Strip(strings.Join(lines, "\n"))
	if !strings.Contains(joined, "+new") {
		t.Fatalf("narrow diff view did not render structured diff line: %q", joined)
	}
	if strings.Contains(joined, "+raw diff line") {
		t.Fatalf("narrow diff view unexpectedly rendered overview fallback: %q", joined)
	}
}

func TestRenderLeftPanelsSeparated(t *testing.T) {
	out := New(Input{
		Width:  80,
		Height: 10,
		Left: LeftPanelsInput{
			Repo:       "owner/repo",
			PRs:        []string{"PR #1 Fix bug"},
			PRSelected: 0,
		},
		Right: RightPanelsInput{
			OverviewTitle: "Overview",
			OverviewLines: []string{"detail"},
		},
	}).String()

	lines := strings.Split(out, "\n")
	if len(lines) < 5 {
		t.Fatalf("got %d lines, want at least 5", len(lines))
	}
	if !strings.HasPrefix(xansi.Strip(lines[0]), "┌ Repository ") {
		t.Fatalf("unexpected first line: %q", xansi.Strip(lines[0]))
	}
	if !strings.Contains(xansi.Strip(lines[4]), "PRs (Open/Draft)") {
		t.Fatalf("line does not contain expected title: %q", xansi.Strip(lines[4]))
	}
}

func TestRenderReviewDrawer_RendersMultilineSummaryWithoutBreakingLayout(t *testing.T) {
	lines := New(Input{
		Focus: layout.FocusReviewDrawer,
		Review: &ReviewDrawerInput{
			SummaryLines:     []string{"first line", "second line"},
			CommentModeLabel: "single-line",
		},
	}).drawer(40, 8)

	var foundHeader bool
	var foundFirst bool
	var foundSecond bool
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

func TestRenderCommentSummary(t *testing.T) {
	got := Comment{
		Path:      "a.txt",
		Line:      12,
		StartLine: 10,
		Body:      "hello\nworld",
	}.Summary()
	if got != "a.txt:10-12 hello world" {
		t.Fatalf("got %q", got)
	}
}
