package core

import "testing"

func TestSanitize(t *testing.T) {
	tests := []struct {
		name string
		run  func(string) string
		in   string
		want string
	}{
		{
			name: "single line",
			run:  sanitizeSingleLine,
			in:   "ok\x1b[31mred\x00\tline\nnext",
			want: "ok[31mred line next",
		},
		{
			name: "multi line",
			run:  sanitizeMultiline,
			in:   "title\x1b[31m\r\nbody\x00\nend",
			want: "title[31m\nbody\nend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.run(tt.in); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatPRItemSanitizeTitle(t *testing.T) {
	pr := FormatPRItem(Item{Number: 2, Title: "bad\x00title"})
	if pr != "PR #2 badtitle" {
		t.Fatalf("unexpected pr format: %q", pr)
	}
}
