package app_test

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/app"
	"github.com/rin2yh/lazygh/internal/app/layout"
	apptest "github.com/rin2yh/lazygh/internal/app/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/pr"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestNavigatePRList(t *testing.T) {
	g := apptest.NewGui(t, &testmock.GHClient{})
	app.GuiCoord(g).ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "a"), testfactory.NewItem(2, "b")}, nil)

	app.GuiNavigateDown(g)
	if app.GuiCoord(g).Selected != 1 {
		t.Fatalf("got %d, want %d", app.GuiCoord(g).Selected, 1)
	}

	app.GuiNavigateUp(g)
	if app.GuiCoord(g).Selected != 0 {
		t.Fatalf("got %d, want %d", app.GuiCoord(g).Selected, 0)
	}
}

func TestCycleFocus_DiffMode(t *testing.T) {
	g, err := app.NewGui(config.Default(), app.NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	app.GuiCoord(g).ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
	app.GuiSwitchToDiff(g)
	app.GuiDiff(g).SetFiles([]gh.DiffFile{{Path: "a.txt", Content: "x"}})

	if app.GuiFocus(g) != layout.FocusDiffFiles {
		t.Fatalf("got %v, want %v", app.GuiFocus(g), layout.FocusDiffFiles)
	}

	app.GuiCycleFocus(g)
	if app.GuiFocus(g) != layout.FocusDiffContent {
		t.Fatalf("got %v, want %v", app.GuiFocus(g), layout.FocusDiffContent)
	}
	app.GuiCycleFocus(g)
	if app.GuiFocus(g) != layout.FocusRepo {
		t.Fatalf("got %v, want %v", app.GuiFocus(g), layout.FocusRepo)
	}
	app.GuiCycleFocus(g)
	if app.GuiFocus(g) != layout.FocusPRs {
		t.Fatalf("got %v, want %v", app.GuiFocus(g), layout.FocusPRs)
	}
	app.GuiCycleFocus(g)
	if app.GuiFocus(g) != layout.FocusDiffFiles {
		t.Fatalf("got %v, want %v", app.GuiFocus(g), layout.FocusDiffFiles)
	}
}

func TestModelUpdateFocusKeysInDiffMode(t *testing.T) {
	tests := []struct {
		name      string
		key       tea.KeyMsg
		files     []gh.DiffFile
		start     layout.Focus
		wantFocus layout.Focus
	}{
		{
			name:      "l moves repo to prs",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusRepo,
			wantFocus: layout.FocusPRs,
		},
		{
			name:      "l moves prs to files",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusPRs,
			wantFocus: layout.FocusDiffFiles,
		},
		{
			name:      "l moves files to diff",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusDiffFiles,
			wantFocus: layout.FocusDiffContent,
		},
		{
			name:      "h moves diff to files",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusDiffContent,
			wantFocus: layout.FocusDiffFiles,
		},
		{
			name:      "h moves files to prs",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusDiffFiles,
			wantFocus: layout.FocusPRs,
		},
		{
			name:      "h moves prs to repo",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusPRs,
			wantFocus: layout.FocusRepo,
		},
		{
			name:      "l moves diff to review",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusDiffContent,
			wantFocus: layout.FocusReviewDrawer,
		},
		{
			name:      "h moves review to diff",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusReviewDrawer,
			wantFocus: layout.FocusDiffContent,
		},
		{
			name:      "h stops at first panel",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusRepo,
			wantFocus: layout.FocusRepo,
		},
		{
			name:      "l stops at last panel",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusReviewDrawer,
			wantFocus: layout.FocusReviewDrawer,
		},
		{
			name:      "esc moves to prs",
			key:       tea.KeyMsg{Type: tea.KeyEsc},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     layout.FocusDiffContent,
			wantFocus: layout.FocusPRs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := app.NewGui(config.Default(), app.NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			app.GuiCoord(g).ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
			app.GuiSwitchToDiff(g)
			app.GuiDiff(g).SetFiles(tt.files)
			app.SetGuiFocus(g, tt.start)
			app.ReviewCtrl(g).OpenDrawer()
			m := app.NewScreen(g)

			_, cmd := m.Update(tt.key)
			if cmd != nil {
				t.Fatal("did not expect command")
			}
			if app.GuiFocus(g) != tt.wantFocus {
				t.Fatalf("got %v, want %v", app.GuiFocus(g), tt.wantFocus)
			}
		})
	}
}

func TestModelUpdateFocusKeysInOverviewMode(t *testing.T) {
	tests := []struct {
		name      string
		key       tea.KeyMsg
		start     layout.Focus
		wantFocus layout.Focus
	}{
		{
			name:      "l moves repo to prs",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			start:     layout.FocusRepo,
			wantFocus: layout.FocusPRs,
		},
		{
			name:      "l moves prs to overview",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			start:     layout.FocusPRs,
			wantFocus: layout.FocusDiffContent,
		},
		{
			name:      "h moves overview to prs",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			start:     layout.FocusDiffContent,
			wantFocus: layout.FocusPRs,
		},
		{
			name:      "h moves prs to repo",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			start:     layout.FocusPRs,
			wantFocus: layout.FocusRepo,
		},
		{
			name:      "h stops at first panel",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			start:     layout.FocusRepo,
			wantFocus: layout.FocusRepo,
		},
		{
			name:      "l stops at last panel",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			start:     layout.FocusDiffContent,
			wantFocus: layout.FocusDiffContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := app.NewGui(config.Default(), app.NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			app.GuiCoord(g).ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
			app.SetGuiFocus(g, tt.start)
			m := app.NewScreen(g)

			_, cmd := m.Update(tt.key)
			if cmd != nil {
				t.Fatal("did not expect command")
			}
			if app.GuiFocus(g) != tt.wantFocus {
				t.Fatalf("got %v, want %v", app.GuiFocus(g), tt.wantFocus)
			}
		})
	}
}
