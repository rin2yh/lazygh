package gui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestNavigatePRList(t *testing.T) {
	g, err := NewGui(config.Default(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "a"), testfactory.CoreItem(2, "b")}, nil)

	g.navigateDown()
	if g.state.PRsSelected != 1 {
		t.Fatalf("got %d, want %d", g.state.PRsSelected, 1)
	}

	g.navigateUp()
	if g.state.PRsSelected != 0 {
		t.Fatalf("got %d, want %d", g.state.PRsSelected, 0)
	}
}

func TestModelUpdate_JKMovesPRsOnlyWhenPRPanelFocusedInOverviewMode(t *testing.T) {
	tests := []struct {
		name       string
		key        tea.KeyMsg
		startFocus panelFocus
		startIndex int
		wantIndex  int
		wantCmd    bool
		wantDetail string
		client     *testmock.GHClient
	}{
		{
			name:       "j on prs moves selection without reload",
			key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			startFocus: panelPRs,
			startIndex: 0,
			wantIndex:  1,
			wantCmd:    false,
			wantDetail: "PR #2 two\nStatus: OPEN\nAssignee: unassigned",
			client:     &testmock.GHClient{},
		},
		{
			name:       "j on repo does nothing",
			key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			startFocus: panelRepo,
			startIndex: 0,
			wantIndex:  0,
			wantCmd:    false,
			wantDetail: "PR #1 one\nStatus: \nAssignee: -",
			client:     &testmock.GHClient{},
		},
		{
			name:       "j on overview does nothing",
			key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			startFocus: panelDiffContent,
			startIndex: 0,
			wantIndex:  0,
			wantCmd:    false,
			wantDetail: "PR #1 one\nStatus: \nAssignee: -",
			client:     &testmock.GHClient{},
		},
	}

	prs := []core.Item{testfactory.CoreItem(1, "one"), testfactory.CoreItem(2, "two")}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGui(config.Default(), tt.client, tt.client)
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.state.ApplyPRsResult("owner/repo", prs, nil)
			g.focus = tt.startFocus
			g.state.PRsSelected = tt.startIndex
			g.state.DetailContent = "PR #1 one\nStatus: \nAssignee: -"
			m := &screen{gui: g}

			_, cmd := m.Update(tt.key)
			if (cmd != nil) != tt.wantCmd {
				t.Fatalf("cmd returned = %v, want %v", cmd != nil, tt.wantCmd)
			}
			if g.state.PRsSelected != tt.wantIndex {
				t.Fatalf("got selected %d, want %d", g.state.PRsSelected, tt.wantIndex)
			}
			if g.state.DetailContent != tt.wantDetail {
				t.Fatalf("got detail %q, want %q", g.state.DetailContent, tt.wantDetail)
			}
		})
	}
}

func TestModelUpdate_JKMovesPRsOnlyWhenPRPanelFocusedInDiffMode(t *testing.T) {
	prs := []core.Item{testfactory.CoreItem(1, "one"), testfactory.CoreItem(2, "two")}

	t.Run("j on prs returns reload command", func(t *testing.T) {
		client := &testmock.GHClient{PRDiff: "diff for two"}
		g, err := NewGui(config.Default(), client, client)
		if err != nil {
			t.Fatalf("NewGui failed: %v", err)
		}
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
		g, err := NewGui(config.Default(), &testmock.GHClient{PRDiff: "diff for two"}, &testmock.GHClient{PRDiff: "diff for two"})
		if err != nil {
			t.Fatalf("NewGui failed: %v", err)
		}
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
		g, err := NewGui(config.Default(), &testmock.GHClient{PRDiff: "diff for two"}, &testmock.GHClient{PRDiff: "diff for two"})
		if err != nil {
			t.Fatalf("NewGui failed: %v", err)
		}
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
		g.diffLineSelected = 0
		m := &screen{gui: g}

		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if cmd != nil {
			t.Fatal("did not expect command")
		}
		if g.state.PRsSelected != 0 {
			t.Fatalf("got selected %d, want %d", g.state.PRsSelected, 0)
		}
		if g.diffLineSelected == 0 {
			t.Fatal("expected diff line selection to move")
		}
	})
}
