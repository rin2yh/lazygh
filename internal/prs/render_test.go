package prs

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/gui/layout"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

func TestRenderLeftPanelsSeparated(t *testing.T) {
	screen := layout.New(80, 10, false, false)
	input := PanelInput{
		Repo:       "owner/repo",
		PRs:        []string{"PR #1 Fix bug"},
		PRSelected: 0,
		Filter:     "Open",
	}
	active := func(f layout.Focus) bool { return f == layout.FocusRepo }
	style := func(a bool) widget.PanelStyle {
		if a {
			return widget.PanelStyle{BorderColor: "green", TitleColor: "green"}
		}
		return widget.PanelStyle{BorderColor: "white"}
	}
	lines := RenderLeft(input, screen.RepoHeight, screen.PRHeight, active, style, screen.LeftWidth)
	if len(lines) < 5 {
		t.Fatalf("got %d lines, want at least 5", len(lines))
	}
	if !strings.HasPrefix(xansi.Strip(lines[0]), "┌ Repository ") {
		t.Fatalf("unexpected first line: %q", xansi.Strip(lines[0]))
	}
	if !strings.Contains(xansi.Strip(lines[4]), "PRs [Open]") {
		t.Fatalf("line does not contain expected title: %q", xansi.Strip(lines[4]))
	}
}
