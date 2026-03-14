package gui

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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
