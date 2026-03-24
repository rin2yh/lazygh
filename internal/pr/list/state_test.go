package list

import (
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/pr"
)

func TestNavigateDown(t *testing.T) {
	tests := []struct {
		name    string
		ls      ListState
		wantOk  bool
		wantSel int
	}{
		{
			name:    "empty list",
			ls:      ListState{},
			wantOk:  false,
			wantSel: 0,
		},
		{
			name:    "at last item",
			ls:      ListState{Items: []pr.Item{{}, {}}, Selected: 1},
			wantOk:  false,
			wantSel: 1,
		},
		{
			name:    "normal move",
			ls:      ListState{Items: []pr.Item{{}, {}}, Selected: 0},
			wantOk:  true,
			wantSel: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := tt.ls.NavigateDown()
			if ok != tt.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOk)
			}
			if tt.ls.Selected != tt.wantSel {
				t.Fatalf("Selected = %d, want %d", tt.ls.Selected, tt.wantSel)
			}
		})
	}
}

func TestNavigateUp(t *testing.T) {
	tests := []struct {
		name    string
		ls      ListState
		wantOk  bool
		wantSel int
	}{
		{
			name:    "at top",
			ls:      ListState{Items: []pr.Item{{}, {}}, Selected: 0},
			wantOk:  false,
			wantSel: 0,
		},
		{
			name:    "normal move",
			ls:      ListState{Items: []pr.Item{{}, {}}, Selected: 1},
			wantOk:  true,
			wantSel: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := tt.ls.NavigateUp()
			if ok != tt.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOk)
			}
			if tt.ls.Selected != tt.wantSel {
				t.Fatalf("Selected = %d, want %d", tt.ls.Selected, tt.wantSel)
			}
		})
	}
}

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
