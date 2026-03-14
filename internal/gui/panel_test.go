package gui

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/google/go-cmp/cmp"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestFramePanel(t *testing.T) {
	got := framePanel("Repo", false, []string{"body"}, 10, 3)
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
	got := framePanel("Repo", false, []string{"x"}, 1, 2)
	want := []string{"x", ""}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("frame mismatch (-want +got)\n%s", diff)
	}
}

func TestPadOrTrimHandlesANSI(t *testing.T) {
	colored := ansiGreen + "+10" + ansiReset
	got := padOrTrim(colored, 4)
	if !strings.Contains(got, colored) {
		t.Fatalf("result does not contain colored text: %q", got)
	}
}

func TestGuiFramePanel_ActiveUsesConfiguredColors(t *testing.T) {
	g := newTestGuiWithClient(&testmock.GHClient{})
	lines := g.framePanel("Repo", true, []string{"body"}, 10, 3)

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

func TestGuiFramePanel_InvalidThemeColorFallsBack(t *testing.T) {
	g := newTestGuiWithClient(&testmock.GHClient{})
	g.config.Theme.ActiveBorderColor = "unknown-color"
	g.config.Theme.InactiveBorderColor = "invalid"

	active := g.framePanel("Repo", true, []string{"body"}, 10, 3)
	inactive := g.framePanel("Repo", false, []string{"body"}, 10, 3)

	if !strings.Contains(active[0], ansiGreen+"┌") {
		t.Fatalf("active border should fallback to green: %q", active[0])
	}
	if !strings.Contains(inactive[0], "\x1b[37m┌") {
		t.Fatalf("inactive border should fallback to white: %q", inactive[0])
	}
}
