package model

import "testing"

func TestReviewEventLabel(t *testing.T) {
	tests := []struct {
		event ReviewEvent
		label string
	}{
		{ReviewEventComment, "COMMENT"},
		{ReviewEventApprove, "APPROVE"},
		{ReviewEventRequestChanges, "REQUEST CHANGES"},
	}
	for _, tt := range tests {
		if got := tt.event.Label(); got != tt.label {
			t.Errorf("event=%v: got %q, want %q", tt.event, got, tt.label)
		}
	}
}
