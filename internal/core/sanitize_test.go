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
			run:  SanitizeSingleLine,
			in:   "ok\x1b[31mred\x00\tline\nnext",
			want: "ok[31mred line next",
		},
		{
			name: "multi line",
			run:  SanitizeMultiline,
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
	if pr != "#2 badtitle" {
		t.Fatalf("unexpected pr format: %q", pr)
	}
}

func TestFormatPROverview(t *testing.T) {
	pr := FormatPROverview(Item{
		Number:    3,
		Title:     "bad\x00title",
		Status:    PRStatusDraft,
		Assignees: []string{"alice", "bob"},
	})
	want := "PR #3 badtitle\nStatus: DRAFT\nAssignee: alice (+1)"
	if pr != want {
		t.Fatalf("got %q, want %q", pr, want)
	}
}
