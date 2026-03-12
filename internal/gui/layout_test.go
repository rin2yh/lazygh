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
		name     string
		loading  bool
		diffMode bool
		hasPR    bool
		focus    panelFocus
		hasFiles bool
		want     string
	}{
		{
			name:     "overview with pr",
			loading:  false,
			diffMode: false,
			hasPR:    true,
			focus:    panelPRs,
			hasFiles: false,
			want:     "[PRs] [j/k/↑/↓]Move [enter]Reload [d]Diff [q]Quit",
		},
		{
			name:     "diff focus prs",
			loading:  false,
			diffMode: true,
			hasPR:    true,
			focus:    panelPRs,
			hasFiles: true,
			want:     "[tab]Focus [PRs] [j/k/↑/↓]Move [l]Overview [enter]Reload [q]Quit",
		},
		{
			name:     "diff focus files",
			loading:  false,
			diffMode: true,
			hasPR:    true,
			focus:    panelDiffFiles,
			hasFiles: true,
			want:     "[tab]Focus [Files] [j/k/↑/↓]Move [l]Diff [o]Overview [q]Quit",
		},
		{
			name:     "diff focus detail",
			loading:  false,
			diffMode: true,
			hasPR:    true,
			focus:    panelDiffContent,
			hasFiles: true,
			want:     "[tab]Focus [Diff] [j/k/↑/↓]Line [space/b]Page [g/G]Top/Bottom [h]Files [enter]Reload [o]Overview [q]Quit",
		},
		{
			name:     "overview without pr",
			loading:  false,
			diffMode: false,
			hasPR:    false,
			focus:    panelPRs,
			hasFiles: false,
			want:     "[d]Diff [q]Quit",
		},
		{
			name:     "diff without pr",
			loading:  false,
			diffMode: true,
			hasPR:    false,
			focus:    panelDiffContent,
			hasFiles: false,
			want:     "[o]Overview [q]Quit",
		},
		{
			name:     "loading",
			loading:  true,
			diffMode: true,
			hasPR:    false,
			focus:    panelDiffContent,
			hasFiles: false,
			want:     "Loading...  | [o]Overview [q]Quit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatusLine(tt.loading, tt.diffMode, tt.hasPR, tt.focus, tt.hasFiles)
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
