package core

import "testing"

func TestSanitizeSingleLine(t *testing.T) {
	in := "ok\x1b[31mred\x00\tline\nnext"
	got := sanitizeSingleLine(in)
	want := "ok[31mred line next"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestSanitizeMultiline(t *testing.T) {
	in := "title\x1b[31m\r\nbody\x00\nend"
	got := sanitizeMultiline(in)
	want := "title[31m\nbody\nend"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatItemsSanitizeTitle(t *testing.T) {
	issue := FormatIssueItem(Item{Number: 1, Title: "bad\x1b[31m\ntitle"})
	if issue != "Issue #1 bad[31m title" {
		t.Fatalf("unexpected issue format: %q", issue)
	}

	pr := FormatPRItem(Item{Number: 2, Title: "bad\x00title"})
	if pr != "PR #2 badtitle" {
		t.Fatalf("unexpected pr format: %q", pr)
	}
}
