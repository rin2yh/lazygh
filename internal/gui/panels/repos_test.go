package panels

import "testing"

func TestNewReposPanel(t *testing.T) {
	p := NewReposPanel()
	if p == nil {
		t.Fatal("NewReposPanel returned nil")
	}
	if p.Repos == nil {
		t.Error("Repos should not be nil")
	}
	if len(p.Repos) != 0 {
		t.Errorf("Repos length: got %d, want 0", len(p.Repos))
	}
	if p.Selected != 0 {
		t.Errorf("Selected: got %d, want 0", p.Selected)
	}
	if p.Loading {
		t.Error("Loading should be false")
	}
}

func TestCalcOriginY(t *testing.T) {
	tests := []struct {
		name     string
		selected int
		originY  int
		height   int
		want     int
	}{
		{"選択が表示範囲内", 3, 0, 5, 0},
		{"選択が下にはみ出す", 5, 0, 5, 1},
		{"選択が上にはみ出す", 0, 3, 5, 0},
		{"originY変更不要(下端ちょうど)", 4, 0, 5, 0},
		{"スクロール済みで範囲内", 5, 3, 5, 3},
		{"スクロール済みで上にはみ出す", 2, 3, 5, 2},
		{"スクロール済みで下にはみ出す", 9, 3, 5, 5},
	}
	for _, tt := range tests {
		got := calcOriginY(tt.selected, tt.originY, tt.height)
		if got != tt.want {
			t.Errorf("%s: calcOriginY(%d, %d, %d) = %d, want %d",
				tt.name, tt.selected, tt.originY, tt.height, got, tt.want)
		}
	}
}

func TestCalcCursorY(t *testing.T) {
	tests := []struct {
		name     string
		selected int
		originY  int
		height   int
		want     int
	}{
		{"通常", 3, 1, 5, 2},
		{"originと同じ", 4, 4, 5, 0},
		{"上にはみ出す", 2, 5, 5, 0},
		{"下にはみ出す", 12, 5, 5, 4},
		{"高さ0", 0, 0, 0, 0},
	}
	for _, tt := range tests {
		got := calcCursorY(tt.selected, tt.originY, tt.height)
		if got != tt.want {
			t.Errorf("%s: calcCursorY(%d, %d, %d) = %d, want %d",
				tt.name, tt.selected, tt.originY, tt.height, got, tt.want)
		}
	}
}
