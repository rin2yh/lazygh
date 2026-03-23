package list

import (
	"github.com/rin2yh/lazygh/internal/model"
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

func TestStatusPrefix(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{model.PRStatusOpen, widget.Colorize("O", "green")},
		{model.PRStatusDraft, widget.Colorize("D", "gray")},
		{model.PRStatusClosed, widget.Colorize("C", "red")},
		{model.PRStatusMerged, widget.Colorize("M", "purple")},
		{"", widget.Colorize("O", "green")},
		{"UNKNOWN", widget.Colorize("O", "green")},
	}
	for _, tt := range tests {
		got := statusPrefix(tt.status)
		if got != tt.want {
			t.Errorf("statusPrefix(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestRenderLeftPanelsSeparated(t *testing.T) {
	screen := layout.New(80, 10, false, false)
	input := PanelInput{
		Repo:     "owner/repo",
		Items:    []model.Item{{Number: 1, Title: "Fix bug"}},
		Selected: 0,
		Filter:   "Open",
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
	if !strings.Contains(xansi.Strip(lines[4]), "PR [Open]") {
		t.Fatalf("line does not contain expected title: %q", xansi.Strip(lines[4]))
	}
}
