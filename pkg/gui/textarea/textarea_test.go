package textarea

import (
	"strings"
	"testing"
)

func newFocused(placeholder string) State {
	s := New(placeholder)
	s.Focus()
	return s
}

func TestNew(t *testing.T) {
	s := New("placeholder text")
	if got := s.Text(); got != "" {
		t.Errorf("Text() = %q, want empty", got)
	}
}

func TestLoadAndText(t *testing.T) {
	s := New("")
	s.Load("hello")
	if got := s.Text(); got != "hello" {
		t.Errorf("Text() = %q, want %q", got, "hello")
	}
}

func TestClear(t *testing.T) {
	s := New("")
	s.Load("some content")
	s.Clear()
	if got := s.Text(); got != "" {
		t.Errorf("after Clear(), Text() = %q, want empty", got)
	}
}

func TestLines(t *testing.T) {
	s := newFocused("")
	s.Load("line1\nline2")
	lines := s.Lines()
	if len(lines) == 0 {
		t.Fatal("Lines() returned empty slice")
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "line1") || !strings.Contains(joined, "line2") {
		t.Errorf("Lines() = %v, want to contain line1 and line2", lines)
	}
}

func TestViewMatchesLines(t *testing.T) {
	s := newFocused("")
	s.Load("abc")
	if got, want := s.View(), strings.Join(s.Lines(), "\n"); got != want {
		t.Errorf("View() != strings.Join(Lines(), newline)\nView=%q\nJoined=%q", got, want)
	}
}
