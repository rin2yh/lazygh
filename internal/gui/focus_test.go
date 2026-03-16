package gui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/model"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestCycleFocus_DiffMode(t *testing.T) {
	g, err := NewGui(config.Default(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []model.Item{testfactory.NewItem(1, "x")}, nil)
	g.switchToDiff()
	g.diff.SetFiles([]gh.DiffFile{{Path: "a.txt", Content: "x"}})

	if g.focus != panelDiffFiles {
		t.Fatalf("got %v, want %v", g.focus, panelDiffFiles)
	}

	g.cycleFocus()
	if g.focus != panelDiffContent {
		t.Fatalf("got %v, want %v", g.focus, panelDiffContent)
	}
	g.cycleFocus()
	if g.focus != panelRepo {
		t.Fatalf("got %v, want %v", g.focus, panelRepo)
	}
	g.cycleFocus()
	if g.focus != panelPRs {
		t.Fatalf("got %v, want %v", g.focus, panelPRs)
	}
	g.cycleFocus()
	if g.focus != panelDiffFiles {
		t.Fatalf("got %v, want %v", g.focus, panelDiffFiles)
	}
}

func TestModelUpdateFocusKeysInDiffMode(t *testing.T) {
	tests := []struct {
		name      string
		key       tea.KeyMsg
		files     []gh.DiffFile
		start     panelFocus
		wantFocus panelFocus
	}{
		{
			name:      "l moves repo to prs",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelRepo,
			wantFocus: panelPRs,
		},
		{
			name:      "l moves prs to files",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelPRs,
			wantFocus: panelDiffFiles,
		},
		{
			name:      "l moves files to diff",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelDiffFiles,
			wantFocus: panelDiffContent,
		},
		{
			name:      "h moves diff to files",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelDiffContent,
			wantFocus: panelDiffFiles,
		},
		{
			name:      "h moves files to prs",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelDiffFiles,
			wantFocus: panelPRs,
		},
		{
			name:      "h moves prs to repo",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelPRs,
			wantFocus: panelRepo,
		},
		{
			name:      "l moves diff to review",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelDiffContent,
			wantFocus: panelReviewDrawer,
		},
		{
			name:      "h moves review to diff",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelReviewDrawer,
			wantFocus: panelDiffContent,
		},
		{
			name:      "h stops at first panel",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelRepo,
			wantFocus: panelRepo,
		},
		{
			name:      "l stops at last panel",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelReviewDrawer,
			wantFocus: panelReviewDrawer,
		},
		{
			name:      "esc moves to prs",
			key:       tea.KeyMsg{Type: tea.KeyEsc},
			files:     []gh.DiffFile{{Path: "a.txt", Content: "x"}},
			start:     panelDiffContent,
			wantFocus: panelPRs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGui(config.Default(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.state.ApplyPRsResult("owner/repo", []model.Item{testfactory.NewItem(1, "x")}, nil)
			g.switchToDiff()
			g.diff.SetFiles(tt.files)
			g.focus = tt.start
			g.state.Review.DrawerOpen = true
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
		start     panelFocus
		wantFocus panelFocus
	}{
		{
			name:      "l moves repo to prs",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			start:     panelRepo,
			wantFocus: panelPRs,
		},
		{
			name:      "l moves prs to overview",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			start:     panelPRs,
			wantFocus: panelDiffContent,
		},
		{
			name:      "h moves overview to prs",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			start:     panelDiffContent,
			wantFocus: panelPRs,
		},
		{
			name:      "h moves prs to repo",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			start:     panelPRs,
			wantFocus: panelRepo,
		},
		{
			name:      "h stops at first panel",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			start:     panelRepo,
			wantFocus: panelRepo,
		},
		{
			name:      "l stops at last panel",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			start:     panelDiffContent,
			wantFocus: panelDiffContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGui(config.Default(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.state.ApplyPRsResult("owner/repo", []model.Item{testfactory.NewItem(1, "x")}, nil)
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
