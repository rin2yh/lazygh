package diff

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

func TestRenderRightPanels_DiffModeHasFilesPanel(t *testing.T) {
	input := PanelInput{
		DiffMode: true,
		DiffFiles: []gh.DiffFile{
			{Path: "a.txt", Status: gh.DiffFileStatusModified, Additions: 1, Deletions: 1},
		},
		DiffContentLines: []ContentLine{{Location: "+1", Text: "+new", Selected: true}},
	}
	lines := RenderFiles(input, widget.PanelStyle{BorderColor: "white"}, 30, 9)
	if len(lines) < 1 || !strings.Contains(xansi.Strip(lines[0]), "Files") {
		t.Fatalf("line does not contain %q: %q", "Files", xansi.Strip(lines[0]))
	}
}

func TestRenderRightPanels_NarrowDiffModeFallsBackToOverviewLines(t *testing.T) {
	input := PanelInput{
		DiffMode:      true,
		OverviewLines: []string{"@@ -1 +1 @@", "+raw diff line"},
	}
	lines := ContentLines(input, 9)
	if lines != nil {
		t.Fatalf("expected nil ContentLines when no DiffContentLines, got: %v", lines)
	}
	rendered := RenderContent(input, widget.PanelStyle{BorderColor: "green", TitleColor: "green"}, 19, 9)
	joined := xansi.Strip(strings.Join(rendered, "\n"))
	if !strings.Contains(joined, "+raw diff line") {
		t.Fatalf("narrow diff view did not render overview fallback: %q", joined)
	}
}

func TestRenderRightPanels_NarrowDiffModePrefersStructuredDiffLines(t *testing.T) {
	input := PanelInput{
		DiffMode:         true,
		OverviewLines:    []string{"+raw diff line"},
		DiffContentLines: []ContentLine{{Location: "+1", Text: "+new", Selected: true}},
	}
	lines := ContentLines(input, 9)
	joined := xansi.Strip(strings.Join(lines, "\n"))
	if !strings.Contains(joined, "+new") {
		t.Fatalf("narrow diff view did not render structured diff line: %q", joined)
	}
	if strings.Contains(joined, "+raw diff line") {
		t.Fatalf("narrow diff view unexpectedly rendered overview fallback: %q", joined)
	}
}
