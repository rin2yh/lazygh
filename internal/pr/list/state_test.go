package list

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/model"
)

func TestSelectedOverview_Empty(t *testing.T) {
	ls := &ListState{}

	_, ok := ls.SelectedOverview()
	if ok {
		t.Fatal("expected false when Items is empty")
	}
}

func TestSelectedOverview_OutOfRange(t *testing.T) {
	ls := &ListState{
		Items:    []model.Item{{Number: 1, Title: "PR 1"}},
		Selected: 5,
	}

	_, ok := ls.SelectedOverview()
	if ok {
		t.Fatal("expected false when Selected is out of range")
	}
}

func TestSelectedOverview_NegativeIndex(t *testing.T) {
	ls := &ListState{
		Items:    []model.Item{{Number: 1, Title: "PR 1"}},
		Selected: -1,
	}

	_, ok := ls.SelectedOverview()
	if ok {
		t.Fatal("expected false when Selected is negative")
	}
}

func TestSelectedOverview_Valid(t *testing.T) {
	tests := []struct {
		name     string
		items    []model.Item
		selected int
		wantSub  string
	}{
		{
			name:     "first item selected",
			items:    []model.Item{{Number: 1, Title: "Fix bug", Status: "OPEN"}},
			selected: 0,
			wantSub:  "PR #1",
		},
		{
			name: "second item selected",
			items: []model.Item{
				{Number: 1, Title: "First"},
				{Number: 2, Title: "Second", Status: "MERGED"},
			},
			selected: 1,
			wantSub:  "PR #2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &ListState{Items: tt.items, Selected: tt.selected}

			got, ok := ls.SelectedOverview()
			if !ok {
				t.Fatal("expected true for valid selection")
			}
			if len(got) == 0 {
				t.Fatal("expected non-empty overview")
			}
			// Check the PR number appears in the output
			for _, r := range []string{tt.wantSub} {
				found := false
				for i := 0; i <= len(got)-len(r); i++ {
					if got[i:i+len(r)] == r {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("overview %q does not contain %q", got, r)
				}
			}
		})
	}
}
