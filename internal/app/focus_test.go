package app

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/app/layout"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/pr"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestNavigatePRList(t *testing.T) {
	g := mustNewGui(t, &testmock.GHClient{})
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "a"), testfactory.NewItem(2, "b")}, nil)

	g.navigateDown()
	if g.coord.Selected() != 1 {
		t.Fatalf("got %d, want %d", g.coord.Selected(), 1)
	}

	g.navigateUp()
	if g.coord.Selected() != 0 {
		t.Fatalf("got %d, want %d", g.coord.Selected(), 0)
	}
}

func TestCycleFocus_DiffMode(t *testing.T) {
	g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
	g.switchToDiff()
	g.diff.SetFiles([]gh.DiffFile{{Path: "a.txt", Content: "x"}})

	if g.focus != layout.FocusDiffFiles {
		t.Fatalf("got %v, want %v", g.focus, layout.FocusDiffFiles)
	}

	g.cycleFocus()
	if g.focus != layout.FocusDiffContent {
		t.Fatalf("got %v, want %v", g.focus, layout.FocusDiffContent)
	}
	g.cycleFocus()
	if g.focus != layout.FocusRepo {
		t.Fatalf("got %v, want %v", g.focus, layout.FocusRepo)
	}
	g.cycleFocus()
	if g.focus != layout.FocusPRs {
		t.Fatalf("got %v, want %v", g.focus, layout.FocusPRs)
	}
	g.cycleFocus()
	if g.focus != layout.FocusDiffFiles {
		t.Fatalf("got %v, want %v", g.focus, layout.FocusDiffFiles)
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
			g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
			g.switchToDiff()
			g.diff.SetFiles(tt.files)
			g.focus = tt.start
			reviewCtrl(g).OpenDrawer()
			m := &screen{gui: g}

			_, cmd := m.Update(tt.key)
			if cmd != nil {
				t.Fatal("did not expect command")
			}
			if g.focus != tt.wantFocus {
				t.Fatalf("got %v, want %v", g.focus, tt.wantFocus)
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
			g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
			g.focus = tt.start
			m := &screen{gui: g}

			_, cmd := m.Update(tt.key)
			if cmd != nil {
				t.Fatal("did not expect command")
			}
			if g.focus != tt.wantFocus {
				t.Fatalf("got %v, want %v", g.focus, tt.wantFocus)
			}
		})
	}
}
