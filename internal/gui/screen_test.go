package gui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/model"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestModelUpdate_JKMovesPRsOnlyWhenPRPanelFocusedInOverviewMode(t *testing.T) {
	tests := []struct {
		name       string
		key        tea.KeyMsg
		startFocus panelFocus
		startIndex int
		wantIndex  int
		wantCmd    bool
		wantDetail string
	}{
		{
			name:       "j on prs moves selection without reload",
			key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			startFocus: panelPRs,
			startIndex: 0,
			wantIndex:  1,
			wantCmd:    false,
			wantDetail: "PR #2 two\nStatus: OPEN\nAssignee: unassigned",
		},
		{
			name:       "j on repo does nothing",
			key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			startFocus: panelRepo,
			startIndex: 0,
			wantIndex:  0,
			wantCmd:    false,
			wantDetail: "PR #1 one\nStatus: \nAssignee: -",
		},
		{
			name:       "j on overview does nothing",
			key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			startFocus: panelDiffContent,
			startIndex: 0,
			wantIndex:  0,
			wantCmd:    false,
			wantDetail: "PR #1 one\nStatus: \nAssignee: -",
		},
	}

	prs := []model.Item{testfactory.NewItem(1, "one"), testfactory.NewItem(2, "two")}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := mustNewGui(t, &testmock.GHClient{})
			g.state.ApplyPRsResult("owner/repo", prs, nil)
			g.focus = tt.startFocus
			g.state.PRsSelected = tt.startIndex
			g.state.Detail.Content = "PR #1 one\nStatus: \nAssignee: -"
			m := &screen{gui: g}

			_, cmd := m.Update(tt.key)
			if (cmd != nil) != tt.wantCmd {
				t.Fatalf("cmd returned = %v, want %v", cmd != nil, tt.wantCmd)
			}
			if g.state.PRsSelected != tt.wantIndex {
				t.Fatalf("got selected %d, want %d", g.state.PRsSelected, tt.wantIndex)
			}
			if g.state.Detail.Content != tt.wantDetail {
				t.Fatalf("got detail %q, want %q", g.state.Detail.Content, tt.wantDetail)
			}
		})
	}
}

func TestModelUpdate_JKMovesPRsOnlyWhenPRPanelFocusedInDiffMode(t *testing.T) {
	prs := []model.Item{testfactory.NewItem(1, "one"), testfactory.NewItem(2, "two")}

	t.Run("j on prs returns reload command", func(t *testing.T) {
		client := &testmock.GHClient{PRDiff: "diff for two"}
		g := mustNewGui(t, client)
		g.state.ApplyPRsResult("owner/repo", prs, nil)
		g.switchToDiff()
		g.focus = panelPRs
		m := &screen{gui: g}

		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if cmd == nil {
			t.Fatal("expected reload command")
		}
		if g.state.PRsSelected != 1 {
			t.Fatalf("got selected %d, want %d", g.state.PRsSelected, 1)
		}

		msg := cmd().(detailLoadedMsg)
		if msg.err != nil {
			t.Fatalf("unexpected error: %v", msg.err)
		}
		if msg.number != 2 {
			t.Fatalf("got number %d, want %d", msg.number, 2)
		}
		if msg.content != "diff for two" {
			t.Fatalf("got content %q, want %q", msg.content, "diff for two")
		}
	})

	t.Run("j on repo does nothing", func(t *testing.T) {
		g := mustNewGui(t, &testmock.GHClient{PRDiff: "diff for two"})
		g.state.ApplyPRsResult("owner/repo", prs, nil)
		g.switchToDiff()
		g.focus = panelRepo
		m := &screen{gui: g}

		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if cmd != nil {
			t.Fatal("did not expect command")
		}
		if g.state.PRsSelected != 0 {
			t.Fatalf("got selected %d, want %d", g.state.PRsSelected, 0)
		}
	})

	t.Run("j on diff content scrolls without changing prs", func(t *testing.T) {
		g := mustNewGui(t, &testmock.GHClient{PRDiff: "diff for two"})
		g.state.ApplyPRsResult("owner/repo", prs, nil)
		g.switchToDiff()
		g.updateDiffFiles(strings.Join([]string{
			"diff --git a/a.txt b/a.txt",
			"--- a/a.txt",
			"+++ b/a.txt",
			"@@ -1,2 +1,2 @@",
			"-old",
			"+new",
		}, "\n"))
		g.focus = panelDiffContent
		g.diff.SetLineSelected(0)
		m := &screen{gui: g}

		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if cmd != nil {
			t.Fatal("did not expect command")
		}
		if g.state.PRsSelected != 0 {
			t.Fatalf("got selected %d, want %d", g.state.PRsSelected, 0)
		}
		if g.diff.LineSelected() == 0 {
			t.Fatal("expected diff line selection to move")
		}
	})
}
