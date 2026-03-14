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

func TestComputeLeftPanels(t *testing.T) {
	got := New(120, 11, false, false)
	if got.RepoHeight != 4 {
		t.Fatalf("repo panel height: got %d, want %d", got.RepoHeight, 4)
	}
	if got.PRHeight != 6 {
		t.Fatalf("pr panel height: got %d, want %d", got.PRHeight, 6)
	}
}
