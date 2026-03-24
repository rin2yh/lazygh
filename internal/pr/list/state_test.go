package list

import (
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/pr"
)

func TestNavigateDown(t *testing.T) {
	tests := []struct {
		name    string
		ls      State
		wantOk  bool
		wantSel int
	}{
		{
			name:    "empty list",
			ls:      State{},
			wantOk:  false,
			wantSel: 0,
		},
		{
			name:    "at last item",
			ls:      State{items: []pr.Item{{}, {}}, selected: 1},
			wantOk:  false,
			wantSel: 1,
		},
		{
			name:    "normal move",
			ls:      State{items: []pr.Item{{}, {}}, selected: 0},
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
			if tt.ls.selected != tt.wantSel {
				t.Fatalf("Selected = %d, want %d", tt.ls.selected, tt.wantSel)
			}
		})
	}
}

func TestNavigateUp(t *testing.T) {
	tests := []struct {
		name    string
		ls      State
		wantOk  bool
		wantSel int
	}{
		{
			name:    "at top",
			ls:      State{items: []pr.Item{{}, {}}, selected: 0},
			wantOk:  false,
			wantSel: 0,
		},
		{
			name:    "normal move",
			ls:      State{items: []pr.Item{{}, {}}, selected: 1},
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
			if tt.ls.selected != tt.wantSel {
				t.Fatalf("Selected = %d, want %d", tt.ls.selected, tt.wantSel)
			}
		})
	}
}

func TestSelectedOverview(t *testing.T) {
	tests := []struct {
		name    string
		ls      *State
		wantOk  bool
		wantSub string
	}{
		{
			name:   "empty items",
			ls:     &State{},
			wantOk: false,
		},
		{
			name:   "selected out of range",
			ls:     &State{items: []pr.Item{{Number: 1, Title: "PR 1"}}, selected: 5},
			wantOk: false,
		},
		{
			name:   "selected negative",
			ls:     &State{items: []pr.Item{{Number: 1, Title: "PR 1"}}, selected: -1},
			wantOk: false,
		},
		{
			name:    "first item selected",
			ls:      &State{items: []pr.Item{{Number: 1, Title: "Fix bug", Status: "OPEN"}}, selected: 0},
			wantOk:  true,
			wantSub: "PR #1",
		},
		{
			name: "second item selected",
			ls: &State{
				items:    []pr.Item{{Number: 1, Title: "First"}, {Number: 2, Title: "Second", Status: "MERGED"}},
				selected: 1,
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
