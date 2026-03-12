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
	tests := []struct {
		name    string
		loading bool
		want    string
	}{
		{
			name:    "idle",
			loading: false,
			want:    "[q]Quit  [j/k]Move  [enter]Reload detail",
		},
		{
			name:    "loading",
			loading: true,
			want:    "Loading...  | [q]Quit  [j/k]Move  [enter]Reload detail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatusLine(tt.loading)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatRepoLine(t *testing.T) {
	tests := []struct {
		name string
		repo string
		want string
	}{
		{
			name: "with repo",
			repo: "owner/repo",
			want: "owner/repo",
		},
		{
			name: "empty",
			repo: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatRepoLine(tt.repo)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
