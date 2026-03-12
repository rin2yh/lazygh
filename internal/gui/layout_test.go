package gui

import (
	"strings"
	"testing"
)

func TestFormatPanelTitle(t *testing.T) {
	tests := []struct {
		name   string
		base   string
		active bool
		want   string
	}{
		{
			name:   "Active",
			base:   "PRs",
			active: true,
			want:   "> PRs <",
		},
		{
			name:   "Inactive",
			base:   "Detail",
			active: false,
			want:   " Detail ",
		},
	}

	for _, tt := range tests {
		got := formatPanelTitle(tt.base, tt.active)
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestFormatStatusLine(t *testing.T) {
	got := formatStatusLine("owner/repo")

	if !strings.Contains(got, "Repo: owner/repo") {
		t.Errorf("status %q should contain repo", got)
	}
	if !strings.Contains(got, "[q]Quit  [j/k]Move  [enter]Reload detail") {
		t.Errorf("status %q should contain key guide", got)
	}
}

func TestFormatStatusLine_Resolving(t *testing.T) {
	got := formatStatusLine("")
	if !strings.Contains(got, "Repo: (resolving...)") {
		t.Errorf("status %q should contain fallback repo", got)
	}
}
