package review

import "testing"

func TestEventLabel(t *testing.T) {
	tests := []struct {
		event Event
		label string
	}{
		{EventComment, "COMMENT"},
		{EventApprove, "APPROVE"},
		{EventRequestChanges, "REQUEST CHANGES"},
	}
	for _, tt := range tests {
		if got := tt.event.Label(); got != tt.label {
			t.Errorf("event=%v: got %q, want %q", tt.event, got, tt.label)
		}
	}
}
