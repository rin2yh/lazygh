package config

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Theme.ActiveBorderColor != "green" {
		t.Fatalf("got %q, want %q", cfg.Theme.ActiveBorderColor, "green")
	}
	if cfg.Theme.InactiveBorderColor != "white" {
		t.Fatalf("got %q, want %q", cfg.Theme.InactiveBorderColor, "white")
	}
	if got := cfg.KeyBindings.Binding(ActionQuit).Keys; len(got) != 2 || got[0] != "q" || got[1] != "ctrl+c" {
		t.Fatalf("got %v, want [q ctrl+c]", got)
	}
	if got := cfg.KeyBindings.Binding(ActionMoveDown).Keys; len(got) != 2 || got[0] != "j" || got[1] != "down" {
		t.Fatalf("got %v, want [j down]", got)
	}
	if got := cfg.KeyBindings.Binding(ActionOpenSelected).Keys; len(got) != 1 || got[0] != "r" {
		t.Fatalf("got %v, want [r]", got)
	}
	if got := cfg.KeyBindings.Binding(ActionReviewComment).Keys; len(got) != 1 || got[0] != "enter" {
		t.Fatalf("got %v, want [enter]", got)
	}
}

func TestKeyBindingsMatches(t *testing.T) {
	keys := Default().KeyBindings

	tests := []struct {
		name   string
		msg    tea.KeyMsg
		action Action
		want   bool
	}{
		{name: "quit rune", msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, action: ActionQuit, want: true},
		{name: "quit ctrl", msg: tea.KeyMsg{Type: tea.KeyCtrlC}, action: ActionQuit, want: true},
		{name: "move down rune", msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}, action: ActionMoveDown, want: true},
		{name: "move down arrow", msg: tea.KeyMsg{Type: tea.KeyDown}, action: ActionMoveDown, want: true},
		{name: "save review", msg: tea.KeyMsg{Type: tea.KeyCtrlS}, action: ActionReviewSave, want: true},
		{name: "submit review ctrl+r", msg: tea.KeyMsg{Type: tea.KeyCtrlR}, action: ActionReviewSubmit, want: true},
		{name: "submit review S no longer bound", msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'S'}}, action: ActionReviewSubmit, want: false},
		{name: "go bottom uppercase only", msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}, action: ActionGoBottom, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keys.Matches(tt.msg, tt.action); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyBindingsLabels(t *testing.T) {
	keys := Default().KeyBindings

	if got := keys.QuitLabel(); got != "q" {
		t.Fatalf("got %q, want %q", got, "q")
	}
	if got := keys.MoveLabel(); got != "j/k/↑/↓" {
		t.Fatalf("got %q, want %q", got, "j/k/↑/↓")
	}
	if got := keys.PageLabel(); got != "space/b" {
		t.Fatalf("got %q, want %q", got, "space/b")
	}
	if got := keys.TopBottomLabel(); got != "g/G" {
		t.Fatalf("got %q, want %q", got, "g/G")
	}
	if got := keys.ReviewModeLabel(); got != "enter/R" {
		t.Fatalf("got %q, want %q", got, "enter/R")
	}
}

func TestKeyBindingsLabelsFollowCustomBindings(t *testing.T) {
	keys := Default().KeyBindings
	keys.SetBinding(ActionMoveUp, KeyBinding{Keys: []string{"p", "up"}})
	keys.SetBinding(ActionPanelNext, KeyBinding{Keys: []string{"n"}})
	keys.SetBinding(ActionPageUp, KeyBinding{Keys: []string{"u"}})
	keys.SetBinding(ActionGoBottom, KeyBinding{Keys: []string{"B"}})
	keys.SetBinding(ActionReviewSummary, KeyBinding{Keys: []string{"r"}})

	if got := keys.MoveLabel(); got != "j/p/↑/↓" {
		t.Fatalf("got %q, want %q", got, "j/p/↑/↓")
	}
	if got := keys.PanelLabel(); got != "h/n" {
		t.Fatalf("got %q, want %q", got, "h/n")
	}
	if got := keys.PageLabel(); got != "space/u" {
		t.Fatalf("got %q, want %q", got, "space/u")
	}
	if got := keys.TopBottomLabel(); got != "g/B" {
		t.Fatalf("got %q, want %q", got, "g/B")
	}
	if got := keys.ReviewModeLabel(); got != "enter/r" {
		t.Fatalf("got %q, want %q", got, "enter/r")
	}
}

func TestKeyBindingsLabelsDeduplicate(t *testing.T) {
	keys := Default().KeyBindings
	keys.SetBinding(ActionMoveDown, KeyBinding{Keys: []string{"j", "down"}})
	keys.SetBinding(ActionMoveUp, KeyBinding{Keys: []string{"j", "up"}})

	if got := keys.MoveLabel(); got != "j/↑/↓" {
		t.Fatalf("got %q, want %q", got, "j/↑/↓")
	}
}

func TestKeyBindingsActionFor(t *testing.T) {
	keys := Default().KeyBindings

	action, ok := keys.ActionFor(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	if !ok {
		t.Fatal("expected action")
	}
	if action != ActionShowDiff {
		t.Fatalf("got %v, want %v", action, ActionShowDiff)
	}
}
