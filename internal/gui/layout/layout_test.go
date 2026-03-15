package layout

import "testing"

func TestComputeScreen(t *testing.T) {
	got := New(120, 40, true, true)
	if got.LeftWidth != 26 {
		t.Fatalf("left width: got %d, want %d", got.LeftWidth, 26)
	}
	if got.RightWidth != 93 {
		t.Fatalf("right width: got %d, want %d", got.RightWidth, 93)
	}
	if got.DrawerHeight <= 0 {
		t.Fatalf("drawer height should be positive: %d", got.DrawerHeight)
	}
	if got.MainHeight <= 0 {
		t.Fatalf("main height should be positive: %d", got.MainHeight)
	}
}

func TestDiffSplitWidths(t *testing.T) {
	tests := []struct {
		name       string
		totalWidth int
		wantFiles  int
		wantDiff   int
	}{
		{"too narrow to split", 19, 0, 19},
		{"zero width", 0, 0, 0},
		{"minimum split width", 20, 10, 9},
		{"normal width", 100, 30, 69},
		{"wide", 200, 60, 139},
		{"small split clamps to min files", 40, 16, 23},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, d := DiffSplitWidths(tt.totalWidth)
			if f != tt.wantFiles {
				t.Errorf("filesWidth = %d, want %d", f, tt.wantFiles)
			}
			if d != tt.wantDiff {
				t.Errorf("diffWidth = %d, want %d", d, tt.wantDiff)
			}
		})
	}
}

func TestComputeLeftPanels(t *testing.T) {
	got := New(120, 11, false, false)
	if got.RepoHeight != 4 {
		t.Fatalf("repo panel height: got %d, want %d", got.RepoHeight, 4)
	}
	if got.PRHeight != 6 {
		t.Fatalf("pr panel height: got %d, want %d", got.PRHeight, 6)
	}
}
