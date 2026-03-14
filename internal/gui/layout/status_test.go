package layout

import (
	"testing"

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
			want:  "[q]Quit [enter]Reload | [Repo] [l]Next Panel [d]Diff",
		},
		{
			name:     "diff focus files",
			diffMode: true,
			hasPR:    true,
			focus:    FocusDiffFiles,
			hasFiles: true,
			want:     "[q]Quit [enter]Reload | [tab]Focus [Files] [j/k/↑/↓]Move [h/l]Prev/Next Panel [o]Overview [v]Range [c]Comment",
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
			}.String()
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
