package gui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestScrollDetailByKey(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*Gui)
		key            tea.KeyMsg
		wantHandled    bool
		wantOffsetMove bool
	}{
		{
			name: "diff mode page down",
			setup: func(g *Gui) {
				g.switchToDiff()
				g.focus = panelDiffContent
			},
			key:            tea.KeyMsg{Type: tea.KeyPgDown},
			wantHandled:    true,
			wantOffsetMove: true,
		},
		{
			name:           "overview mode page down",
			setup:          func(_ *Gui) {},
			key:            tea.KeyMsg{Type: tea.KeyPgDown},
			wantHandled:    false,
			wantOffsetMove: false,
		},
		{
			name: "diff mode d key",
			setup: func(g *Gui) {
				g.switchToDiff()
				g.focus = panelDiffContent
			},
			key:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantHandled:    false,
			wantOffsetMove: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
			tt.setup(g)

			g.syncDetailViewport(20, 4, strings.Repeat("line\n", 30))
			before := g.detailViewport.YOffset

			handled := g.scrollDetailByKey(tt.key)
			if handled != tt.wantHandled {
				t.Fatalf("got %v, want %v", handled, tt.wantHandled)
			}
			if tt.wantOffsetMove {
				if g.detailViewport.YOffset <= before {
					t.Fatalf("expected offset to increase, before=%d after=%d", before, g.detailViewport.YOffset)
				}
				return
			}
			if g.detailViewport.YOffset != before {
				t.Fatalf("got %d, want %d", g.detailViewport.YOffset, before)
			}
		})
	}
}
