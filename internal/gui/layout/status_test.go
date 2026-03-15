package layout

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
)

func TestFormatStatusLine(t *testing.T) {
	tests := []struct {
		name     string
		loading  bool
		diffMode bool
		hasPR    bool
		focus    Focus
		hasFiles bool
		want     string
	}{
		{
			name:  "overview repo focus with pr",
			hasPR: true,
			focus: FocusRepo,
			want:  "[q]Quit [r]Reload | [Repo] [l]Next Panel [d]Diff",
		},
		{
			name:     "diff focus files",
			diffMode: true,
			hasPR:    true,
			focus:    FocusDiffFiles,
			hasFiles: true,
			want:     "[q]Quit [r]Reload | [tab]Focus [Files] [j/k/↑/↓]Move [h/l]Prev/Next Panel [o]Overview [v]Range [enter]Comment",
		},
		{
			name:     "loading",
			loading:  true,
			diffMode: true,
			focus:    FocusDiffContent,
			want:     "Loading...  | [q]Quit | [o]Overview",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Status{
				Loading:   tt.loading,
				DiffMode:  tt.diffMode,
				HasPR:     tt.hasPR,
				Focus:     tt.focus,
				HasFiles:  tt.hasFiles,
				InputMode: core.ReviewInputNone,
				Keys:      config.Default().KeyBindings,
			}.String()
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatStatusLine_UsesCustomBindings(t *testing.T) {
	keys := config.Default().KeyBindings
	keys.MoveUp = config.KeyBinding{Keys: []string{"p", "up"}}
	keys.PanelNext = config.KeyBinding{Keys: []string{"n"}}
	keys.ReviewSummary = config.KeyBinding{Keys: []string{"r"}}

	got := Status{
		DiffMode:  true,
		HasPR:     true,
		Focus:     FocusPRs,
		InputMode: core.ReviewInputNone,
		Keys:      keys,
	}.String()

	want := "[q]Quit [r]Reload | [tab]Focus [PRs] [h/n]Prev/Next Panel [j/p/↑/↓]Move [enter/r]Review"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
