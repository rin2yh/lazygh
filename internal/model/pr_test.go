package model

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
