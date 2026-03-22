package layout

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/pr/review"
)

func TestFormatStatusLine(t *testing.T) {
	tests := []struct {
		name      string
		fetching  bool
		diffMode  bool
		focus     Focus
		inputMode review.InputMode
		want      string
	}{
		{
			name:  "overview repo focus",
			focus: FocusRepo,
			want:  "[q]Quit [?]Help | [h/l]Panels [d]Diff",
		},
		{
			name:     "diff focus files",
			diffMode: true,
			focus:    FocusDiffFiles,
			want:     "[q]Quit [?]Help | [h/l]Panels [o]Overview",
		},
		{
			name:     "fetching",
			fetching: true,
			diffMode: true,
			focus:    FocusDiffContent,
			want:     "Fetching... | [q]Quit [?]Help | [h/l]Panels [o]Overview",
		},
		{
			name:      "review input comment",
			diffMode:  true,
			focus:     FocusReviewDrawer,
			inputMode: review.InputComment,
			want:      "[q]Quit [?]Help | [ctrl+s]Save Comment [esc]Cancel",
		},
		{
			name:      "review input summary",
			diffMode:  true,
			focus:     FocusReviewDrawer,
			inputMode: review.InputSummary,
			want:      "[q]Quit [?]Help | [ctrl+s]Save Summary [esc]Cancel",
		},
		{
			name:     "review drawer focus",
			diffMode: true,
			focus:    FocusReviewDrawer,
			want:     "[q]Quit [?]Help | [Review] [ctrl+r]Submit [X]Discard [esc]Cancel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Status{
				Fetching:  tt.fetching,
				DiffMode:  tt.diffMode,
				Focus:     tt.focus,
				InputMode: tt.inputMode,
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
	keys.SetBinding(config.ActionMoveUp, config.KeyBinding{Keys: []string{"p", "up"}})
	keys.SetBinding(config.ActionPanelNext, config.KeyBinding{Keys: []string{"n"}})
	keys.SetBinding(config.ActionReviewSummary, config.KeyBinding{Keys: []string{"r"}})

	got := Status{
		DiffMode:  true,
		Focus:     FocusPRs,
		InputMode: review.InputNone,
		Keys:      keys,
	}.String()

	want := "[q]Quit [?]Help | [h/n]Panels [o]Overview"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
