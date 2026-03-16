package gui

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gui/widget"
)

func TestPRStatusPrefix(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{core.PRStatusOpen, widget.Colorize("O", "green")},
		{core.PRStatusDraft, widget.Colorize("D", "gray")},
		{core.PRStatusClosed, widget.Colorize("C", "red")},
		{core.PRStatusMerged, widget.Colorize("M", "purple")},
		{"", widget.Colorize("O", "green")},
		{"UNKNOWN", widget.Colorize("O", "green")},
	}
	for _, tt := range tests {
		got := prStatusPrefix(tt.status)
		if got != tt.want {
			t.Errorf("prStatusPrefix(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}
