package gui

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/gui/widget"
	"github.com/rin2yh/lazygh/internal/model"
)

func TestPRStatusPrefix(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{model.PRStatusOpen, widget.Colorize("O", "green")},
		{model.PRStatusDraft, widget.Colorize("D", "gray")},
		{model.PRStatusClosed, widget.Colorize("C", "red")},
		{model.PRStatusMerged, widget.Colorize("M", "purple")},
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
