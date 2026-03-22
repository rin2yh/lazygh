package widget

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/pkg/gui/ansi"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
		want    string
	}{
		{
			name:    "wrap long line",
			content: "abcdefghij",
			width:   4,
			want:    "abcd\nefgh\nij",
		},
		{
			name:    "keep existing line breaks",
			content: "abcde\nfghij",
			width:   3,
			want:    "abc\nde\nfgh\nij",
		},
		{
			name:    "no wrap when width is enough",
			content: "abc",
			width:   10,
			want:    "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WrapText(tt.content, tt.width); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWrapTextWithANSI(t *testing.T) {
	got := WrapText(ansi.Green+"abcdef"+ansi.Reset, 3)
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d, want %d", len(lines), 2)
	}
	if xansi.Strip(lines[0]) != "abc" {
		t.Fatalf("got %q, want %q", xansi.Strip(lines[0]), "abc")
	}
	if xansi.StringWidth(lines[0]) != 3 {
		t.Fatalf("got %d, want %d", xansi.StringWidth(lines[0]), 3)
	}
	if xansi.Strip(lines[1]) != "def" {
		t.Fatalf("got %q, want %q", xansi.Strip(lines[1]), "def")
	}
	if xansi.StringWidth(lines[1]) != 3 {
		t.Fatalf("got %d, want %d", xansi.StringWidth(lines[1]), 3)
	}
}
