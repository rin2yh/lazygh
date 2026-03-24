package list

import (
	"strings"
	"testing"
)

func TestFilterPanelLines(t *testing.T) {
	tests := []struct {
		name        string
		filter      PRFilterMask
		cursor      int
		wantContain []string
	}{
		{
			name:        "returns lines and positive width",
			filter:      PRFilterOpen,
			cursor:      0,
			wantContain: []string{"Open"},
		},
		{
			name:        "all filter option labels present",
			filter:      PRFilterOpen | PRFilterMerged,
			cursor:      0,
			wantContain: []string{"Open", "Closed", "Merged"},
		},
		{
			name:        "enabled filter shows checked marker",
			filter:      PRFilterOpen,
			cursor:      0,
			wantContain: []string{"[x]", "[ ]"},
		},
		{
			name:        "cursor row is marked",
			filter:      PRFilterOpen,
			cursor:      1,
			wantContain: []string{">"},
		},
		{
			name:        "footer hints present",
			filter:      PRFilterOpen,
			cursor:      0,
			wantContain: []string{"toggle", "apply", "cancel"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, width := FilterPanelLines(tt.filter, tt.cursor)

			if len(lines) == 0 {
				t.Fatal("expected non-empty lines")
			}
			if width <= 0 {
				t.Fatalf("expected positive width, got %d", width)
			}

			joined := strings.Join(lines, "\n")
			for _, s := range tt.wantContain {
				if !strings.Contains(joined, s) {
					t.Errorf("expected %q in output", s)
				}
			}
		})
	}
}

func TestBuildFilterContent(t *testing.T) {
	tests := []struct {
		name        string
		filter      PRFilterMask
		wantChecked int
	}{
		{
			name:        "all enabled",
			filter:      PRFilterOpen | PRFilterClosed | PRFilterMerged,
			wantChecked: len(PRFilterOptions),
		},
		{
			name:        "none enabled",
			filter:      0,
			wantChecked: 0,
		},
		{
			name:        "one enabled",
			filter:      PRFilterOpen,
			wantChecked: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, maxW := buildFilterContent(tt.filter, -1)

			if maxW <= 0 {
				t.Fatalf("expected positive maxW, got %d", maxW)
			}

			got := 0
			for _, l := range lines {
				got += strings.Count(l, "[x]")
			}
			if got != tt.wantChecked {
				t.Fatalf("checked count = %d, want %d", got, tt.wantChecked)
			}
		})
	}
}
