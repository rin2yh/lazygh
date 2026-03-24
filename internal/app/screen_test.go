package app_test

import (
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/app"
	"github.com/rin2yh/lazygh/internal/app/layout"
	apptest "github.com/rin2yh/lazygh/internal/app/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/pr"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestModelUpdate_JKMovesPRsOnlyWhenPRPanelFocusedInOverviewMode(t *testing.T) {
	tests := []struct {
		name       string
		key        tea.KeyMsg
		startFocus layout.Focus
		startIndex int
		wantIndex  int
		wantCmd    bool
		wantDetail string
	}{
		{
			name:       "j on prs moves selection without reload",
			key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			startFocus: layout.FocusPRs,
			startIndex: 0,
			wantIndex:  1,
			wantCmd:    false,
			wantDetail: "PR #2 two\nStatus: OPEN\nAssignee: unassigned",
		},
		{
			name:       "j on repo does nothing",
			key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			startFocus: layout.FocusRepo,
			startIndex: 0,
			wantIndex:  0,
			wantCmd:    false,
			wantDetail: "PR #1 one\nStatus: \nAssignee: -",
		},
		{
			name:       "j on overview does nothing",
			key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			startFocus: layout.FocusDiffContent,
			startIndex: 0,
			wantIndex:  0,
			wantCmd:    false,
			wantDetail: "PR #1 one\nStatus: \nAssignee: -",
		},
	}

	prs := []pr.Item{testfactory.NewItem(1, "one"), testfactory.NewItem(2, "two")}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := apptest.NewGui(t, &testmock.GHClient{})
			app.GuiCoord(g).ApplyPRsResult("owner/repo", prs, nil)
			app.SetGuiFocus(g, tt.startFocus)
			app.GuiCoord(g).Selected = tt.startIndex
			app.GuiCoord(g).Overview.Content = "PR #1 one\nStatus: \nAssignee: -"
			m := app.NewScreen(g)

			_, cmd := m.Update(tt.key)
			if (cmd != nil) != tt.wantCmd {
				t.Fatalf("cmd returned = %v, want %v", cmd != nil, tt.wantCmd)
			}
			if app.GuiCoord(g).Selected != tt.wantIndex {
				t.Fatalf("got selected %d, want %d", app.GuiCoord(g).Selected, tt.wantIndex)
			}
			if app.GuiCoord(g).Overview.Content != tt.wantDetail {
				t.Fatalf("got detail %q, want %q", app.GuiCoord(g).Overview.Content, tt.wantDetail)
			}
		})
	}
}

func TestModelUpdate_JKMovesPRsOnlyWhenPRPanelFocusedInDiffMode(t *testing.T) {
	prs := []pr.Item{testfactory.NewItem(1, "one"), testfactory.NewItem(2, "two")}

	t.Run("j on prs returns reload command", func(t *testing.T) {
		client := &testmock.GHClient{PRDiff: "diff for two"}
		g := apptest.NewGui(t, client)
		app.GuiCoord(g).ApplyPRsResult("owner/repo", prs, nil)
		app.GuiSwitchToDiff(g)
		app.SetGuiFocus(g, layout.FocusPRs)
		m := app.NewScreen(g)

		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if cmd == nil {
			t.Fatal("expected reload command")
		}
		if app.GuiCoord(g).Selected != 1 {
			t.Fatalf("got selected %d, want %d", app.GuiCoord(g).Selected, 1)
		}

		number, content, _, err := app.CastDetailLoadedMsg(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if number != 2 {
			t.Fatalf("got number %d, want %d", number, 2)
		}
		if content != "diff for two" {
			t.Fatalf("got content %q, want %q", content, "diff for two")
		}
	})

	t.Run("j on repo does nothing", func(t *testing.T) {
		g := apptest.NewGui(t, &testmock.GHClient{PRDiff: "diff for two"})
		app.GuiCoord(g).ApplyPRsResult("owner/repo", prs, nil)
		app.GuiSwitchToDiff(g)
		app.SetGuiFocus(g, layout.FocusRepo)
		m := app.NewScreen(g)

		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if cmd != nil {
			t.Fatal("did not expect command")
		}
		if app.GuiCoord(g).Selected != 0 {
			t.Fatalf("got selected %d, want %d", app.GuiCoord(g).Selected, 0)
		}
	})

	t.Run("j on diff content scrolls without changing prs", func(t *testing.T) {
		g := apptest.NewGui(t, &testmock.GHClient{PRDiff: "diff for two"})
		app.GuiCoord(g).ApplyPRsResult("owner/repo", prs, nil)
		app.GuiSwitchToDiff(g)
		app.GuiUpdateDiffFiles(g, strings.Join([]string{
			"diff --git a/a.txt b/a.txt",
			"--- a/a.txt",
			"+++ b/a.txt",
			"@@ -1,2 +1,2 @@",
			"-old",
			"+new",
		}, "\n"))
		app.SetGuiFocus(g, layout.FocusDiffContent)
		app.GuiDiff(g).SetLineSelected(0)
		m := app.NewScreen(g)

		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if cmd != nil {
			t.Fatal("did not expect command")
		}
		if app.GuiCoord(g).Selected != 0 {
			t.Fatalf("got selected %d, want %d", app.GuiCoord(g).Selected, 0)
		}
		if app.GuiDiff(g).LineSelected() == 0 {
			t.Fatal("expected diff line selection to move")
		}
	})
}
