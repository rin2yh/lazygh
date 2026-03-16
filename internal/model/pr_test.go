package model

import "testing"

func TestFormatPRItemSanitizeTitle(t *testing.T) {
	pr := FormatPRItem(Item{Number: 2, Title: "bad\x00title"})
	if pr != "#2 badtitle" {
		t.Fatalf("unexpected pr format: %q", pr)
	}
}

func TestFormatPROverview(t *testing.T) {
	pr := FormatPROverview(Item{
		Number:    3,
		Title:     "bad\x00title",
		Status:    PRStatusDraft,
		Assignees: []string{"alice", "bob"},
	})
	want := "PR #3 badtitle\nStatus: DRAFT\nAssignee: alice (+1)"
	if pr != want {
		t.Fatalf("got %q, want %q", pr, want)
	}
}

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
