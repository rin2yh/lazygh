package list

import (
	"github.com/rin2yh/lazygh/internal/pr"
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

func TestSplitHeight(t *testing.T) {
	tests := []struct {
		total    int
		wantRepo int
		wantPR   int
	}{
		{10, 4, 6},
		{5, 4, 1},
		{2, 1, 1},
		{1, 0, 1},
	}
	for _, tt := range tests {
		repoH, prH := splitHeight(tt.total)
		if repoH != tt.wantRepo {
			t.Errorf("splitHeight(%d) repoH = %d, want %d", tt.total, repoH, tt.wantRepo)
		}
		if prH != tt.wantPR {
			t.Errorf("splitHeight(%d) prH = %d, want %d", tt.total, prH, tt.wantPR)
		}
	}
}

func TestStatusPrefix(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{pr.PRStatusOpen, widget.Colorize("O", "green")},
		{pr.PRStatusDraft, widget.Colorize("D", "gray")},
		{pr.PRStatusClosed, widget.Colorize("C", "red")},
		{pr.PRStatusMerged, widget.Colorize("M", "purple")},
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
	input := Input{
		Repo:     "owner/repo",
		Items:    []pr.Item{{Number: 1, Title: "Fix bug"}},
		Selected: 0,
		Filter:   "Open",
	}
	style := func(f layout.Focus) widget.PanelStyle {
		if f == layout.FocusRepo {
			return widget.PanelStyle{BorderColor: "green", TitleColor: "green"}
		}
		return widget.PanelStyle{BorderColor: "white"}
	}
	lines := RenderLeft(input, style, screen.LeftWidth, screen.MainHeight)
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
