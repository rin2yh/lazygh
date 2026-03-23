package textarea

import (
	"strings"
	"testing"
)

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
	s := New("")
	s.Load("line1\nline2")
	lines := s.Lines()
	found1, found2 := false, false
	for _, l := range lines {
		if strings.Contains(l, "line1") {
			found1 = true
		}
		if strings.Contains(l, "line2") {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Errorf("Lines() = %v, want to contain line1 and line2", lines)
	}
}
