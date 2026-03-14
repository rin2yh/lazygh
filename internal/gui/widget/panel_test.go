package widget

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/google/go-cmp/cmp"
)

const ansiGreen = "\x1b[32m"

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
	colored := ansiGreen + "+10" + ansiReset
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

	if !strings.Contains(lines[0], ansiGreen+"┌") {
		t.Fatalf("top border is not active color: %q", lines[0])
	}
	if !strings.Contains(lines[0], ansiGreen+" Repo "+ansiReset) {
		t.Fatalf("title is not active color: %q", lines[0])
	}
	if strings.Contains(xansi.Strip(lines[0]), "> Repo <") {
		t.Fatalf("title should not use ascii emphasis: %q", xansi.Strip(lines[0]))
	}
}

func TestResolveColorName_InvalidFallsBack(t *testing.T) {
	if got := ResolveColorName("unknown-color", "green"); got != "green" {
		t.Fatalf("got %q, want %q", got, "green")
	}
	if got := ResolveColorName("invalid", "white"); got != "white" {
		t.Fatalf("got %q, want %q", got, "white")
	}
}
