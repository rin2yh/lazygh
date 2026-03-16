package gui

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/core"
)

func TestPRStatusPrefix(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{core.PRStatusOpen, prPrefixOpen},
		{core.PRStatusDraft, prPrefixDraft},
		{core.PRStatusClosed, prPrefixClosed},
		{core.PRStatusMerged, prPrefixMerged},
		{"", prPrefixOpen},
		{"UNKNOWN", prPrefixOpen},
	}
	for _, tt := range tests {
		got := prStatusPrefix(tt.status)
		if got != tt.want {
			t.Errorf("prStatusPrefix(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}
