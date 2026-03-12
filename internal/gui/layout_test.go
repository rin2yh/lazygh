package gui

import "testing"

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
	got := formatStatusLine(false)

	if got != "[q]Quit  [j/k]Move  [enter]Reload detail" {
		t.Errorf("got %q, want %q", got, "[q]Quit  [j/k]Move  [enter]Reload detail")
	}
}

func TestFormatStatusLine_Loading(t *testing.T) {
	got := formatStatusLine(true)

	if got != "Loading...  | [q]Quit  [j/k]Move  [enter]Reload detail" {
		t.Errorf("got %q, want %q", got, "Loading...  | [q]Quit  [j/k]Move  [enter]Reload detail")
	}
}

func TestFormatRepoLine(t *testing.T) {
	got := formatRepoLine("owner/repo")
	if got != "owner/repo" {
		t.Fatalf("got %q, want %q", got, "owner/repo")
	}
}

func TestFormatRepoLine_Resolving(t *testing.T) {
	got := formatRepoLine("")
	if got != "" {
		t.Fatalf("got %q, want empty string", got)
	}
}
