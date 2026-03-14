package gui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestApplyReviewCommentResult_PersistsPendingReviewContextOnError(t *testing.T) {
	g, err := NewGui(config.Default(), &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "Fix bug")}, nil)
	g.state.Loading = core.LoadingReview

	g.applyReviewCommentResult(reviewCommentSavedMsg{
		prNumber: 1,
		ctx: gh.ReviewContext{
			PullRequestID: "PR_kwDO123",
			CommitOID:     "deadbeef",
		},
		reviewID: "PRR_kwDO456",
		comment: gh.ReviewComment{
			Path: "a.txt",
			Body: "body",
			Line: 1,
			Side: gh.DiffSideRight,
		},
		err: errors.New("add failed"),
	})

	if g.state.Review.ReviewID != "PRR_kwDO456" {
		t.Fatalf("got %q, want %q", g.state.Review.ReviewID, "PRR_kwDO456")
	}
	if g.state.Review.PullRequestID != "PR_kwDO123" {
		t.Fatalf("got %q, want %q", g.state.Review.PullRequestID, "PR_kwDO123")
	}
	if g.state.Review.CommitOID != "deadbeef" {
		t.Fatalf("got %q, want %q", g.state.Review.CommitOID, "deadbeef")
	}
	if len(g.state.Review.Comments) != 0 {
		t.Fatalf("got %d comments, want 0", len(g.state.Review.Comments))
	}
	if g.state.Review.Notice != "add failed" {
		t.Fatalf("got %q, want %q", g.state.Review.Notice, "add failed")
	}
}

func TestModelUpdate_VKeyTogglesRangeSelection(t *testing.T) {
	g, err := NewGui(config.Default(), &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "x")}, nil)
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

	m := &screen{gui: g}
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
	g, err := NewGui(config.Default(), &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "x")}, nil)
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

	m := &screen{gui: g}
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
	g, err := NewGui(config.Default(), &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "x")}, nil)
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
	m := &screen{gui: g}
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
	g, err := NewGui(config.Default(), &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "x")}, nil)
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

	m := &screen{gui: g}
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
	g, err := NewGui(config.Default(), mc)
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "x")}, nil)
	g.switchToDiff()
	g.state.SetReviewContext(1, "PR_kwDO123", "deadbeef", "PRR_kwDO456")
	g.state.BeginReviewCommentInput()
	g.commentEditor.SetValue("draft")

	m := &screen{gui: g}
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
	g, err := NewGui(config.Default(), mc)
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "x")}, nil)
	g.switchToDiff()
	g.state.SetReviewContext(1, "PR_kwDO123", "deadbeef", "PRR_kwDO456")
	g.state.BeginReviewCommentInput()
	g.commentEditor.SetValue("draft")

	m := &screen{gui: g}
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
	g, err := NewGui(config.Default(), &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "x")}, nil)
	m := &screen{gui: g}

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
