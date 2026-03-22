package widget

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/google/go-cmp/cmp"
	"github.com/rin2yh/lazygh/pkg/gui/ansi"
)

func TestFramePanel(t *testing.T) {
	got := FramePanel("Repo", []string{"body"}, 10, 3, PanelStyle{})
	want := []string{
		"┌ Repo ──┐",
		"│body    │",
		"└────────┘",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("frame mismatch (-want +got)\n%s", diff)
	}
}

func TestFramePanelFallsBackWhenTooSmall(t *testing.T) {
	got := FramePanel("Repo", []string{"x"}, 1, 2, PanelStyle{})
	want := []string{"x", ""}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("frame mismatch (-want +got)\n%s", diff)
	}
}

func TestPadOrTrimHandlesANSI(t *testing.T) {
	colored := ansi.Green + "+10" + ansi.Reset
	got := PadOrTrim(colored, 4)
	if !strings.Contains(got, colored) {
		t.Fatalf("result does not contain colored text: %q", got)
	}
}

func TestFramePanel_ActiveUsesConfiguredColors(t *testing.T) {
	lines := FramePanel("Repo", []string{"body"}, 10, 3, PanelStyle{
		BorderColor: "green",
		TitleColor:  "green",
	})

	if !strings.Contains(lines[0], ansi.Green+"┌") {
		t.Fatalf("top border is not active color: %q", lines[0])
	}
	if !strings.Contains(lines[0], ansi.Green+" Repo "+ansi.Reset) {
		t.Fatalf("title is not active color: %q", lines[0])
	}
	if strings.Contains(xansi.Strip(lines[0]), "> Repo <") {
		t.Fatalf("title should not use ascii emphasis: %q", xansi.Strip(lines[0]))
	}
}

func TestOverlayPanel(t *testing.T) {
	// screenW=10, panelW=2 → startX=(10-2)/2=4
	// left="abcd"(4), panel="XY"(2), right="    "(4)
	bg := []string{"abcdefghij"}
	panel := []string{"XY"}
	got := OverlayPanel(bg, panel, 2, 10)
	want := "abcdXY    "
	if got[0] != want {
		t.Fatalf("got %q, want %q", got[0], want)
	}
}

func TestResolveColorName_InvalidFallsBack(t *testing.T) {
	tests := []struct {
		color    string
		fallback string
	}{
		{"unknown-color", "green"},
		{"invalid", "white"},
	}
	for _, tt := range tests {
		if got := ResolveColorName(tt.color, tt.fallback); got != tt.fallback {
			t.Errorf("ResolveColorName(%q, %q) = %q, want %q", tt.color, tt.fallback, got, tt.fallback)
		}
	}
}
