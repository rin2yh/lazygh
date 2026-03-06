package panels

import "testing"

func TestNormalizeDisplayText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "ASCII", in: "abc", want: "abc"},
		{name: "Japanese", in: "A日本B", want: "A日\x00本\x00B"},
		{name: "Control", in: "日\n本", want: "日\x00\n本\x00"},
	}

	for _, tt := range tests {
		got := normalizeDisplayText(tt.in)
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestRuneDisplayWidth(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want int
	}{
		{name: "ASCII", r: 'A', want: 1},
		{name: "CJK", r: '日', want: 2},
		{name: "Control", r: '\n', want: 1},
	}
	for _, tt := range tests {
		got := runeDisplayWidth(tt.r)
		if got != tt.want {
			t.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
		}
	}
}
