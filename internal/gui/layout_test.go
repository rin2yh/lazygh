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
			base:   "Repositories",
			active: true,
			want:   "> Repositories <",
		},
		{
			name:   "Inactive",
			base:   "Repositories",
			active: false,
			want:   " Repositories ",
		},
	}

	for _, tt := range tests {
		got := formatPanelTitle(tt.base, tt.active)
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestPanelDisplayName(t *testing.T) {
	tests := []struct {
		name  string
		panel PanelType
		want  string
	}{
		{name: "Repos", panel: PanelRepos, want: "Repositories"},
		{name: "Issues", panel: PanelIssues, want: "Issues"},
		{name: "PRs", panel: PanelPRs, want: "PRs"},
		{name: "Detail", panel: PanelDetail, want: "Detail"},
		{name: "Unknown", panel: PanelType(99), want: "Unknown"},
	}

	for _, tt := range tests {
		got := panelDisplayName(tt.panel)
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestFormatStatusLine(t *testing.T) {
	got := formatStatusLine(PanelIssues)

	if !strings.Contains(got, "Panel: Issues") {
		t.Errorf("status %q should contain %q", got, "Panel: Issues")
	}
	if !strings.Contains(got, "[q]Quit  [tab]Panel  [j/k]Navigate  [enter]Select") {
		t.Errorf("status %q should contain key guide", got)
	}
}

func TestFormatStatusLine_AllPanels(t *testing.T) {
	tests := []struct {
		name  string
		panel PanelType
		want  string
	}{
		{name: "Repos", panel: PanelRepos, want: "Panel: Repositories"},
		{name: "Issues", panel: PanelIssues, want: "Panel: Issues"},
		{name: "PRs", panel: PanelPRs, want: "Panel: PRs"},
		{name: "Detail", panel: PanelDetail, want: "Panel: Detail"},
	}

	for _, tt := range tests {
		got := formatStatusLine(tt.panel)
		if !strings.Contains(got, tt.want) {
			t.Errorf("%s: status %q should contain %q", tt.name, got, tt.want)
		}
	}
}

func TestStatusViewBounds_HasOneVisibleRow(t *testing.T) {
	_, y0, _, y1 := statusViewBounds(120, 40)
	visibleRows := y1 - y0 - 1
	if visibleRows != 1 {
		t.Errorf("visibleRows = %d, want %d", visibleRows, 1)
	}
}
