package detail

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestStateSyncResetsViewportOnContentChange(t *testing.T) {
	state := NewState()
	state.Sync(20, 4, strings.Repeat("line\n", 20))
	state.ScrollDown(3)

	if state.YOffset() == 0 {
		t.Fatal("expected offset to move before sync")
	}

	state.Sync(20, 4, strings.Repeat("next\n", 20))
	if state.YOffset() != 0 {
		t.Fatalf("got %d, want %d", state.YOffset(), 0)
	}
}

func TestStateUpdateHandlesPageDown(t *testing.T) {
	state := NewState()
	state.Sync(20, 4, strings.Repeat("line\n", 30))

	if !state.Update(tea.KeyMsg{Type: tea.KeyPgDown}) {
		t.Fatal("expected update to handle page down")
	}
	if state.YOffset() == 0 {
		t.Fatal("expected offset to move")
	}
}

func TestStateScrollUpAndDown(t *testing.T) {
	state := NewState()
	state.Sync(20, 4, strings.Repeat("line\n", 30))

	state.ScrollDown(2)
	if state.YOffset() == 0 {
		t.Fatal("expected offset to move down")
	}

	state.ScrollUp(1)
	if state.YOffset() >= 2 {
		t.Fatalf("expected offset to decrease, got %d", state.YOffset())
	}
}
