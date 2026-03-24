package list

import (
	"github.com/rin2yh/lazygh/internal/pr"
	"strings"
	"testing"
)

func TestSelectedOverview(t *testing.T) {
	tests := []struct {
		name    string
		ls      *ListState
		wantOk  bool
		wantSub string
	}{
		{
			name:   "empty items",
			ls:     &ListState{},
			wantOk: false,
		},
		{
			name:   "selected out of range",
			ls:     &ListState{Items: []pr.Item{{Number: 1, Title: "PR 1"}}, Selected: 5},
			wantOk: false,
		},
		{
			name:   "selected negative",
			ls:     &ListState{Items: []pr.Item{{Number: 1, Title: "PR 1"}}, Selected: -1},
			wantOk: false,
		},
		{
			name:    "first item selected",
			ls:      &ListState{Items: []pr.Item{{Number: 1, Title: "Fix bug", Status: "OPEN"}}, Selected: 0},
			wantOk:  true,
			wantSub: "PR #1",
		},
		{
			name: "second item selected",
			ls: &ListState{
				Items:    []pr.Item{{Number: 1, Title: "First"}, {Number: 2, Title: "Second", Status: "MERGED"}},
				Selected: 1,
			},
			wantOk:  true,
			wantSub: "PR #2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.ls.SelectedOverview()
			if ok != tt.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOk)
			}
			if tt.wantOk && !strings.Contains(got, tt.wantSub) {
				t.Fatalf("overview %q does not contain %q", got, tt.wantSub)
			}
		})
	}
}
