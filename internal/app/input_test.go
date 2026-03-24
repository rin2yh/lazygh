package app

import (
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/app/layout"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/pr"
	"github.com/rin2yh/lazygh/internal/pr/review"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestModelUpdate_VKeyTogglesRangeSelection(t *testing.T) {
	g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
	g.switchToDiff()
	g.updateDiffFiles(strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+new",
	}, "\n"))
	g.focus = layout.FocusDiffContent
	g.diff.SetLineSelected(5)

	m := &screen{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.review.RangeStart() == nil {
		t.Fatal("expected range start")
	}

	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if g.review.RangeStart() != nil {
		t.Fatal("expected range selection cleared")
	}
}

func TestModelUpdate_EnterKeyUsesRangeFlowAfterV(t *testing.T) {
	g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
	g.switchToDiff()
	g.updateDiffFiles(strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1,2 +1,2 @@",
		"-old",
		"+new",
	}, "\n"))
	g.focus = layout.FocusDiffContent
	g.diff.SetLineSelected(5)

	m := &screen{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.review.RangeStart() == nil {
		t.Fatal("expected range start")
	}

	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.review.InputMode() != review.InputComment {
		t.Fatalf("got %v, want %v", g.review.InputMode(), review.InputComment)
	}
}

func TestModelUpdate_EscCancelsCommentAndClearsRangeHighlight(t *testing.T) {
	g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
	g.switchToDiff()
	g.updateDiffFiles(strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1,2 +1,2 @@",
		"-old",
		"+new",
	}, "\n"))
	g.focus = layout.FocusDiffContent
	g.diff.SetLineSelected(5)
	m := &screen{gui: g}
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if g.review.RangeStart() == nil {
		t.Fatal("expected range start before cancel")
	}

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.review.RangeStart() != nil {
		t.Fatal("expected range start cleared after cancel")
	}
	if g.review.InputMode() != review.InputNone {
		t.Fatalf("got %v, want %v", g.review.InputMode(), review.InputNone)
	}
	if g.focus != layout.FocusDiffContent {
		t.Fatalf("got %v, want %v", g.focus, layout.FocusDiffContent)
	}
}

func TestModelUpdate_EscClearsRangeSelectionWithoutLeavingDiff(t *testing.T) {
	g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
	g.switchToDiff()
	g.updateDiffFiles(strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1,2 +1,2 @@",
		"-old",
		"+new",
	}, "\n"))
	g.focus = layout.FocusDiffContent
	g.diff.SetLineSelected(5)

	m := &screen{gui: g}
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if g.review.RangeStart() == nil {
		t.Fatal("expected range start before esc")
	}

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Fatal("did not expect command")
	}
	if g.review.RangeStart() != nil {
		t.Fatal("expected range start cleared")
	}
	if g.focus != layout.FocusDiffContent {
		t.Fatalf("got %v, want %v", g.focus, layout.FocusDiffContent)
	}
}

func TestModelUpdate_InputModeSubmitShortcutBypassesEditor(t *testing.T) {
	mc := &testmock.GHClient{}
	g, err := NewGui(config.Default(), NewCoordinator(), mc, mc)
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
	g.switchToDiff()
	rc := ReviewCtrl(g)
	rc.SetContext(1, "PR_kwDO123", "deadbeef", "PRR_kwDO456")
	rc.BeginCommentInput()
	rc.SetCommentValue("draft")

	m := &screen{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	if cmd == nil {
		t.Fatal("expected submit command")
	}
	if got := rc.CommentValue(); got != "draft" {
		t.Fatalf("got %q, want %q", got, "draft")
	}
	msg := cmd().(review.SubmittedMsg)
	if msg.ReviewID != "PRR_kwDO456" {
		t.Fatalf("got %q, want %q", msg.ReviewID, "PRR_kwDO456")
	}
	if len(mc.SubmittedReviews) != 1 {
		t.Fatalf("got %d submissions, want 1", len(mc.SubmittedReviews))
	}
}

func TestModelUpdate_InputModeDiscardShortcutBypassesEditor(t *testing.T) {
	mc := &testmock.GHClient{}
	g, err := NewGui(config.Default(), NewCoordinator(), mc, mc)
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
	g.switchToDiff()
	rc := ReviewCtrl(g)
	rc.SetContext(1, "PR_kwDO123", "deadbeef", "PRR_kwDO456")
	rc.BeginCommentInput()
	rc.SetCommentValue("draft")

	m := &screen{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})
	if cmd == nil {
		t.Fatal("expected discard command")
	}
	if got := rc.CommentValue(); got != "draft" {
		t.Fatalf("got %q, want %q", got, "draft")
	}
	_ = cmd().(review.DiscardedMsg)
	if len(mc.DeletedReviews) != 1 {
		t.Fatalf("got %d discards, want 1", len(mc.DeletedReviews))
	}
}

func TestModelUpdate_ReviewKeysIgnoredOutsideDiff(t *testing.T) {
	g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{testfactory.NewItem(1, "x")}, nil)
	m := &screen{gui: g}

	for _, key := range []tea.KeyMsg{
		{Type: tea.KeyEnter},
		{Type: tea.KeyRunes, Runes: []rune{'v'}},
		{Type: tea.KeyRunes, Runes: []rune{'R'}},
	} {
		_, cmd := m.Update(key)
		if cmd != nil {
			t.Fatal("did not expect command")
		}
		if g.review.InputMode() != review.InputNone {
			t.Fatalf("got %v, want %v", g.review.InputMode(), review.InputNone)
		}
		if g.review.RangeStart() != nil {
			t.Fatal("expected no range selection outside diff")
		}
	}
}
