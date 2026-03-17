package prs

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/model"
)

func TestFormatPRItemSanitizeTitle(t *testing.T) {
	pr := FormatPRItem(model.Item{Number: 2, Title: "bad\x00title"})
	if pr != "#2 badtitle" {
		t.Fatalf("unexpected pr format: %q", pr)
	}
}

func TestFormatPROverview(t *testing.T) {
	pr := FormatPROverview(model.Item{
		Number:    3,
		Title:     "bad\x00title",
		Status:    model.PRStatusDraft,
		Assignees: []string{"alice", "bob"},
	})
	want := "PR #3 badtitle\nStatus: DRAFT\nAssignee: alice (+1)"
	if pr != want {
		t.Fatalf("got %q, want %q", pr, want)
	}
}
