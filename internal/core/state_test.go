package core

import "testing"

func TestPRFilterMaskLabel(t *testing.T) {
	tests := []struct {
		mask  PRFilterMask
		label string
	}{
		{PRFilterOpen, "Open"},
		{PRFilterClosed, "Closed"},
		{PRFilterMerged, "Merged"},
		{PRFilterOpen | PRFilterClosed, "Open,Closed"},
		{PRFilterOpen | PRFilterMerged, "Open,Merged"},
		{PRFilterClosed | PRFilterMerged, "Closed,Merged"},
		{PRFilterOpen | PRFilterClosed | PRFilterMerged, "All"},
		{0, "None"},
	}
	for _, tt := range tests {
		if got := tt.mask.Label(); got != tt.label {
			t.Errorf("mask=%v: got %q, want %q", tt.mask, got, tt.label)
		}
	}
}

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
