package list

import (
	"github.com/rin2yh/lazygh/internal/pr"
	"testing"
)

func TestFormatItemSanitizeTitle(t *testing.T) {
	pr := formatItem(pr.Item{Number: 2, Title: "bad\x00title"})
	if pr != "#2 badtitle" {
		t.Fatalf("unexpected pr format: %q", pr)
	}
}

func TestFormatOverview(t *testing.T) {
	pr := formatOverview(pr.Item{
		Number:    3,
		Title:     "bad\x00title",
		Status:    pr.PRStatusDraft,
		Assignees: []string{"alice", "bob"},
	})
	want := "PR #3 badtitle\nStatus: DRAFT\nAssignee: alice (+1)"
	if pr != want {
		t.Fatalf("got %q, want %q", pr, want)
	}
}
