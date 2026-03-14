package gui

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/core"
)

func TestFormatPanelTitle(t *testing.T) {
	tests := []struct {
		name string
		base string
		want string
	}{
		{
			name: "plain title",
			base: "Detail",
			want: " Detail ",
		},
	}

	for _, tt := range tests {
		got := formatPanelTitle(tt.base)
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
			name:     "overview repo focus with pr",
			loading:  false,
			diffMode: false,
			hasPR:    true,
			focus:    panelRepo,
			hasFiles: false,
			want:     "[q]Quit [enter]Reload | [Repo] [l]Next Panel [d]Diff",
		},
		{
			name:     "overview with pr",
			loading:  false,
			diffMode: false,
			hasPR:    true,
			focus:    panelPRs,
			hasFiles: false,
			want:     "[q]Quit [enter]Reload | [PRs] [h/l]Prev/Next Panel [j/k/↑/↓]Move [d]Diff",
		},
		{
			name:     "overview detail focus with pr",
			loading:  false,
			diffMode: false,
			hasPR:    true,
			focus:    panelDiffContent,
			hasFiles: false,
			want:     "[q]Quit [enter]Reload | [Overview] [h]Prev Panel [space/b]Page [enter]Reload [d]Diff",
		},
		{
			name:     "diff focus prs",
			loading:  false,
			diffMode: true,
			hasPR:    true,
			focus:    panelPRs,
			hasFiles: true,
			want:     "[q]Quit [enter]Reload | [tab]Focus [PRs] [h/l]Prev/Next Panel [j/k/↑/↓]Move [c/R]Review",
		},
		{
			name:     "diff focus repo",
			loading:  false,
			diffMode: true,
			hasPR:    true,
			focus:    panelRepo,
			hasFiles: true,
			want:     "[q]Quit [enter]Reload | [tab]Focus [Repo] [l]Next Panel [d]Diff",
		},
		{
			name:     "diff focus files",
			loading:  false,
			diffMode: true,
			hasPR:    true,
			focus:    panelDiffFiles,
			hasFiles: true,
			want:     "[q]Quit [enter]Reload | [tab]Focus [Files] [j/k/↑/↓]Move [h/l]Prev/Next Panel [o]Overview [v]Range [c]Comment",
		},
		{
			name:     "diff focus detail",
			loading:  false,
			diffMode: true,
			hasPR:    true,
			focus:    panelDiffContent,
			hasFiles: true,
			want:     "[q]Quit [enter]Reload | [tab]Focus [Diff] [j/k/↑/↓]Line [space/b]Page [g/G]Top/Bottom [h/l]Prev/Next Panel [v]Range [c]Comment [R]Summary [S]Submit [X]Discard [o]Overview",
		},
		{
			name:     "diff focus review drawer",
			loading:  false,
			diffMode: true,
			hasPR:    true,
			focus:    panelReviewDrawer,
			hasFiles: true,
			want:     "[q]Quit [enter]Reload | [h]Prev Panel [c]Comment [R]Summary [S]Submit [X]Discard [Esc]Diff",
		},
		{
			name:     "overview without pr",
			loading:  false,
			diffMode: false,
			hasPR:    false,
			focus:    panelRepo,
			hasFiles: false,
			want:     "[q]Quit | [Repo] [l]Next Panel [d]Diff",
		},
		{
			name:     "overview prs focus without pr",
			loading:  false,
			diffMode: false,
			hasPR:    false,
			focus:    panelPRs,
			hasFiles: false,
			want:     "[q]Quit | [l]Next Panel [d]Diff",
		},
		{
			name:     "overview detail focus without pr",
			loading:  false,
			diffMode: false,
			hasPR:    false,
			focus:    panelDiffContent,
			hasFiles: false,
			want:     "[q]Quit | [Overview] [h]Prev Panel [d]Diff",
		},
		{
			name:     "diff without pr",
			loading:  false,
			diffMode: true,
			hasPR:    false,
			focus:    panelDiffContent,
			hasFiles: false,
			want:     "[q]Quit | [o]Overview",
		},
		{
			name:     "loading",
			loading:  true,
			diffMode: true,
			hasPR:    false,
			focus:    panelDiffContent,
			hasFiles: false,
			want:     "Loading...  | [q]Quit | [o]Overview",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatusLine(tt.loading, tt.diffMode, tt.hasPR, tt.focus, tt.hasFiles, false, core.ReviewInputNone)
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
