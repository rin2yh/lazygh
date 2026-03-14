package gui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestModelInitLoadsPRs(t *testing.T) {
	mc := &testmock.GHClient{Repo: "owner/repo", PRs: []gh.PRItem{{Number: 2, Title: "p"}}}
	g := newTestGuiWithClient(mc)
	m := &model{gui: g}

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("expected init command")
	}
	msg := cmd().(prsLoadedMsg)
	if msg.err != nil {
		t.Fatalf("unexpected error: %v", msg.err)
	}
	if msg.repo != "owner/repo" {
		t.Fatalf("got %q, want %q", msg.repo, "owner/repo")
	}
	if len(msg.prs) != 1 {
		t.Fatalf("got %d, want %d", len(msg.prs), 1)
	}
}

func TestModelHandleDetailLoad(t *testing.T) {
	tests := []struct {
		name         string
		client       *testmock.GHClient
		pr           core.Item
		switchToDiff bool
		wantMode     core.DetailMode
		wantContent  string
		wantNumber   int
	}{
		{
			name:        "overview",
			client:      &testmock.GHClient{PRView: "detail"},
			pr:          core.Item{Number: 1, Title: "x"},
			wantMode:    core.DetailModeOverview,
			wantContent: "detail",
			wantNumber:  1,
		},
		{
			name:         "diff",
			client:       &testmock.GHClient{PRDiff: "diff"},
			pr:           core.Item{Number: 2, Title: "x"},
			switchToDiff: true,
			wantMode:     core.DetailModeDiff,
			wantContent:  "diff",
			wantNumber:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithPRs(tt.client, tt.pr)
			if tt.switchToDiff {
				g.switchToDiff()
			}
			m := &model{gui: g}

			cmd := m.handleDetailLoad()
			if cmd == nil {
				t.Fatal("expected detail load command")
			}
			msg := cmd().(detailLoadedMsg)
			if msg.err != nil {
				t.Fatalf("unexpected error: %v", msg.err)
			}
			if msg.content != tt.wantContent {
				t.Fatalf("got %q, want %q", msg.content, tt.wantContent)
			}
			if msg.mode != tt.wantMode {
				t.Fatalf("got %v, want %v", msg.mode, tt.wantMode)
			}
			if msg.number != tt.wantNumber {
				t.Fatalf("got %d, want %d", msg.number, tt.wantNumber)
			}
		})
	}
}

func TestToCorePRsMapsStatusAndAssignees(t *testing.T) {
	items := toCorePRs([]gh.PRItem{
		{
			Number:  1,
			Title:   "open",
			State:   "OPEN",
			IsDraft: false,
			Assignees: []gh.GHUser{
				{Login: "alice"},
				{Login: "bob"},
			},
		},
		{
			Number:  2,
			Title:   "draft",
			State:   "OPEN",
			IsDraft: true,
		},
	})

	if len(items) != 2 {
		t.Fatalf("got %d, want %d", len(items), 2)
	}
	if items[0].Status != "OPEN" {
		t.Fatalf("got %q, want %q", items[0].Status, "OPEN")
	}
	if strings.Join(items[0].Assignees, ",") != "alice,bob" {
		t.Fatalf("got %q, want %q", strings.Join(items[0].Assignees, ","), "alice,bob")
	}
	if items[1].Status != "DRAFT" {
		t.Fatalf("got %q, want %q", items[1].Status, "DRAFT")
	}
}

func TestModelHandleLKeyShowsOverviewFromPRsInDiffMode(t *testing.T) {
	mc := &testmock.GHClient{PRView: "overview"}
	g := newTestGuiWithPRs(mc, core.Item{Number: 1, Title: "x"})
	g.switchToDiff()
	g.focus = panelPRs

	m := &model{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if g.state.IsDiffMode() {
		t.Fatal("expected overview mode")
	}
	if cmd == nil {
		t.Fatal("expected detail load command")
	}

	msg := cmd().(detailLoadedMsg)
	if msg.err != nil {
		t.Fatalf("unexpected error: %v", msg.err)
	}
	if msg.mode != core.DetailModeOverview {
		t.Fatalf("got %v, want %v", msg.mode, core.DetailModeOverview)
	}
	if msg.number != 1 {
		t.Fatalf("got %d, want %d", msg.number, 1)
	}
	if msg.content != "overview" {
		t.Fatalf("got %q, want %q", msg.content, "overview")
	}
}

func TestModelUpdateFocusKeysInDiffMode(t *testing.T) {
	tests := []struct {
		name      string
		key       tea.KeyMsg
		start     panelFocus
		wantFocus panelFocus
	}{
		{
			name:      "l moves files to diff",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			start:     panelDiffFiles,
			wantFocus: panelDiffContent,
		},
		{
			name:      "h moves diff to files",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			start:     panelDiffContent,
			wantFocus: panelDiffFiles,
		},
		{
			name:      "esc moves to prs",
			key:       tea.KeyMsg{Type: tea.KeyEsc},
			start:     panelDiffContent,
			wantFocus: panelPRs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
			g.switchToDiff()
			g.diffFiles = []gh.DiffFile{{Path: "a.txt", Content: "x"}}
			g.focus = tt.start
			m := &model{gui: g}

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

func TestModelUpdate_VKeyTogglesRangeSelection(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
	g.switchToDiff()
	g.updateDiffFiles(strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+new",
	}, "\n"))
	g.focus = panelDiffContent
	g.diffLineSelected = 5

	m := &model{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.state.Review.RangeStart == nil {
		t.Fatal("expected range start")
	}

	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if g.state.Review.RangeStart != nil {
		t.Fatal("expected range selection cleared")
	}
}

func TestModelUpdate_CKeyUsesRangeFlowAfterV(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
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
	g.diffLineSelected = 5

	m := &model{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.state.Review.RangeStart == nil {
		t.Fatal("expected range start")
	}

	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.state.Review.InputMode != core.ReviewInputComment {
		t.Fatalf("got %v, want %v", g.state.Review.InputMode, core.ReviewInputComment)
	}
}

func TestModelUpdate_EscCancelsCommentAndClearsRangeHighlight(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
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
	g.diffLineSelected = 5
	m := &model{gui: g}
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if g.state.Review.RangeStart == nil {
		t.Fatal("expected range start before cancel")
	}

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.state.Review.RangeStart != nil {
		t.Fatal("expected range start cleared after cancel")
	}
	if g.state.Review.InputMode != core.ReviewInputNone {
		t.Fatalf("got %v, want %v", g.state.Review.InputMode, core.ReviewInputNone)
	}
	if g.focus != panelDiffContent {
		t.Fatalf("got %v, want %v", g.focus, panelDiffContent)
	}
}

func TestModelUpdate_EscClearsRangeSelectionWithoutLeavingDiff(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
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
	g.diffLineSelected = 5

	m := &model{gui: g}
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if g.state.Review.RangeStart == nil {
		t.Fatal("expected range start before esc")
	}

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.state.Review.RangeStart != nil {
		t.Fatal("expected range start cleared")
	}
	if g.focus != panelDiffContent {
		t.Fatalf("got %v, want %v", g.focus, panelDiffContent)
	}
}

func TestModelUpdate_InputModeSubmitShortcutBypassesEditor(t *testing.T) {
	mc := &testmock.GHClient{}
	g := newTestGuiWithPRs(mc, core.Item{Number: 1, Title: "x"})
	g.switchToDiff()
	g.state.SetReviewContext(1, "PR_kwDO123", "deadbeef", "PRR_kwDO456")
	g.state.BeginReviewCommentInput()
	g.commentEditor.SetValue("draft")

	m := &model{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'S'}})
	if cmd == nil {
		t.Fatal("expected submit command")
	}
	if got := g.commentEditor.Value(); got != "draft" {
		t.Fatalf("got %q, want %q", got, "draft")
	}
	msg := cmd().(reviewSubmittedMsg)
	if msg.reviewID != "PRR_kwDO456" {
		t.Fatalf("got %q, want %q", msg.reviewID, "PRR_kwDO456")
	}
	if len(mc.SubmittedReviews) != 1 {
		t.Fatalf("got %d submissions, want 1", len(mc.SubmittedReviews))
	}
}

func TestModelUpdate_InputModeDiscardShortcutBypassesEditor(t *testing.T) {
	mc := &testmock.GHClient{}
	g := newTestGuiWithPRs(mc, core.Item{Number: 1, Title: "x"})
	g.switchToDiff()
	g.state.SetReviewContext(1, "PR_kwDO123", "deadbeef", "PRR_kwDO456")
	g.state.BeginReviewCommentInput()
	g.commentEditor.SetValue("draft")

	m := &model{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})
	if cmd == nil {
		t.Fatal("expected discard command")
	}
	if got := g.commentEditor.Value(); got != "draft" {
		t.Fatalf("got %q, want %q", got, "draft")
	}
	_ = cmd().(reviewDiscardedMsg)
	if len(mc.DeletedReviews) != 1 {
		t.Fatalf("got %d discards, want 1", len(mc.DeletedReviews))
	}
}

func TestModelUpdate_ReviewKeysIgnoredOutsideDiff(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
	m := &model{gui: g}

	for _, key := range []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'c'}},
		{Type: tea.KeyRunes, Runes: []rune{'v'}},
		{Type: tea.KeyRunes, Runes: []rune{'R'}},
	} {
		_, cmd := m.Update(key)
		if cmd != nil {
			t.Fatal("did not expect command")
		}
		if g.state.Review.InputMode != core.ReviewInputNone {
			t.Fatalf("got %v, want %v", g.state.Review.InputMode, core.ReviewInputNone)
		}
		if g.state.Review.RangeStart != nil {
			t.Fatal("expected no range selection outside diff")
		}
	}
}
