package list

import (
	"github.com/rin2yh/lazygh/internal/pr"
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/gh"
)

func TestConvert(t *testing.T) {
	items := Convert([]gh.PRItem{
		{
			Number:  1,
			Title:   "open",
			State:   "OPEN",
			IsDraft: false,
			Assignees: []gh.GHUser{
				{Login: "alice"},
				{Login: "bob"},
			},
		},
		{
			Number:  2,
			Title:   "draft",
			State:   "OPEN",
			IsDraft: true,
		},
	}, PRFilterOpen)

	if len(items) != 2 {
		t.Fatalf("got %d, want %d", len(items), 2)
	}
	if items[0].Status != pr.PRStatusOpen {
		t.Fatalf("got %q, want %q", items[0].Status, pr.PRStatusOpen)
	}
	if strings.Join(items[0].Assignees, ",") != "alice,bob" {
		t.Fatalf("got %q, want %q", strings.Join(items[0].Assignees, ","), "alice,bob")
	}
	if items[1].Status != pr.PRStatusDraft {
		t.Fatalf("got %q, want %q", items[1].Status, pr.PRStatusDraft)
	}
}
