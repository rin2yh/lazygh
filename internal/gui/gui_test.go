package gui

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func newTestGuiWithClient(client gh.ClientInterface) *Gui {
	g, _ := NewGui(config.Default(), client)
	return g
}

func newTestGuiWithPRs(client gh.ClientInterface, prs ...core.Item) *Gui {
	g := newTestGuiWithClient(client)
	g.state.ApplyPRsResult("owner/repo", prs, nil)
	return g
}

func TestNavigatePRList(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, []core.Item{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}}...)

	g.navigateDown()
	if g.state.PRsSelected != 1 {
		t.Fatalf("got %d, want %d", g.state.PRsSelected, 1)
	}

	g.navigateUp()
	if g.state.PRsSelected != 0 {
		t.Fatalf("got %d, want %d", g.state.PRsSelected, 0)
	}
}

func TestApplyPRsResult(t *testing.T) {
	type want struct {
		repo   string
		prs    []core.Item
		detail string
	}

	tests := []struct {
		name string
		msg  prsLoadedMsg
		want want
	}{
		{
			name: "success",
			msg: prsLoadedMsg{
				repo: "owner/repo",
				prs:  []core.Item{{Number: 1, Title: "Fix bug", Status: "OPEN", Assignees: []string{"alice"}}},
			},
			want: want{
				repo:   "owner/repo",
				prs:    []core.Item{{Number: 1, Title: "Fix bug", Status: "OPEN", Assignees: []string{"alice"}}},
				detail: "PR #1 Fix bug\nStatus: OPEN\nAssignee: alice",
			},
		},
		{
			name: "empty",
			msg: prsLoadedMsg{
				repo: "owner/repo",
			},
			want: want{
				repo:   "owner/repo",
				prs:    nil,
				detail: "No pull requests",
			},
		},
		{
			name: "error",
			msg: prsLoadedMsg{
				err: errors.New("boom"),
			},
			want: want{
				repo:   "",
				prs:    nil,
				detail: "Error loading PRs: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithClient(&testmock.GHClient{})
			g.state.BeginLoadPRs()

			g.applyPRsResult(tt.msg)

			if g.state.PRsLoading {
				t.Fatal("expected PRsLoading=false")
			}
			if g.state.Loading != core.LoadingNone {
				t.Fatalf("got %v, want %v", g.state.Loading, core.LoadingNone)
			}
			if g.state.Repo != tt.want.repo {
				t.Fatalf("got %q, want %q", g.state.Repo, tt.want.repo)
			}
			if diff := cmp.Diff(tt.want.prs, g.state.PRs, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("prs mismatch (-want +got)\n%s", diff)
			}
			if g.state.DetailContent != tt.want.detail {
				t.Fatalf("got %q, want %q", g.state.DetailContent, tt.want.detail)
			}
		})
	}
}

func TestApplyDetailResult(t *testing.T) {
	type want struct {
		detail string
	}

	tests := []struct {
		name string
		msg  detailLoadedMsg
		want want
	}{
		{
			name: "success",
			msg: detailLoadedMsg{
				mode:    core.DetailModeOverview,
				number:  1,
				content: "hello",
			},
			want: want{
				detail: "hello",
			},
		},
		{
			name: "error",
			msg: detailLoadedMsg{
				mode:   core.DetailModeOverview,
				number: 1,
				err:    errors.New("boom"),
			},
			want: want{
				detail: "Error loading detail: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "Fix bug"})
			g.state.Loading = core.LoadingDetail

			g.applyDetailResult(tt.msg)

			if g.state.Loading != core.LoadingNone {
				t.Fatalf("got %v, want %v", g.state.Loading, core.LoadingNone)
			}
			if g.state.DetailContent != tt.want.detail {
				t.Fatalf("got %q, want %q", g.state.DetailContent, tt.want.detail)
			}
		})
	}
}

func TestApplyDetailResult_DiffUsesSanitizedContent(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "Fix bug"})
	g.switchToDiff()
	g.state.Loading = core.LoadingDetail

	raw := strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+ok\x1b[31mred",
	}, "\n")

	g.applyDetailResult(detailLoadedMsg{
		mode:    core.DetailModeDiff,
		number:  1,
		content: raw,
	})

	if strings.Contains(g.state.DetailContent, "\x1b") {
		t.Fatalf("detail content should be sanitized: %q", g.state.DetailContent)
	}
	if len(g.diffFiles) != 1 {
		t.Fatalf("got %d, want %d", len(g.diffFiles), 1)
	}
	if strings.Contains(g.diffFiles[0].Content, "\x1b") {
		t.Fatalf("diff file content should be sanitized: %q", g.diffFiles[0].Content)
	}
	if !strings.Contains(g.diffFiles[0].Content, "+ok[31mred") {
		t.Fatalf("unexpected diff content: %q", g.diffFiles[0].Content)
	}
}

func TestUpdateDiffFiles(t *testing.T) {
	g := newTestGuiWithClient(&testmock.GHClient{})
	diff := strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+new",
		"diff --git a/b.txt b/b.txt",
		"--- a/b.txt",
		"+++ b/b.txt",
		"@@ -1 +1 @@",
		"-x",
		"+y",
	}, "\n")

	g.updateDiffFiles(diff)
	want := []gh.DiffFile{
		{Path: "a.txt", Content: strings.Join([]string{
			"diff --git a/a.txt b/a.txt",
			"--- a/a.txt",
			"+++ b/a.txt",
			"@@ -1 +1 @@",
			"-old",
			"+new",
		}, "\n"), Status: gh.DiffFileStatusModified, Additions: 1, Deletions: 1},
		{Path: "b.txt", Content: strings.Join([]string{
			"diff --git a/b.txt b/b.txt",
			"--- a/b.txt",
			"+++ b/b.txt",
			"@@ -1 +1 @@",
			"-x",
			"+y",
		}, "\n"), Status: gh.DiffFileStatusModified, Additions: 1, Deletions: 1},
	}
	if diff := cmp.Diff(want, g.diffFiles); diff != "" {
		t.Fatalf("diffFiles mismatch (-want +got)\n%s", diff)
	}

	g.diffFileSelected = 1
	g.updateDiffFiles(diff)
	if g.diffFileSelected != 1 {
		t.Fatalf("got %d, want %d", g.diffFileSelected, 1)
	}
}

func TestCycleFocus_DiffMode(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
	g.switchToDiff()
	g.diffFiles = []gh.DiffFile{{Path: "a.txt", Content: "x"}}

	if g.focus != panelDiffFiles {
		t.Fatalf("got %v, want %v", g.focus, panelDiffFiles)
	}

	g.cycleFocus()
	if g.focus != panelDiffContent {
		t.Fatalf("got %v, want %v", g.focus, panelDiffContent)
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

func TestProgramE2E_MainFlow(t *testing.T) {
	releaseResolve := make(chan struct{})
	releasePRs := make(chan struct{})
	releaseDetail := make(chan struct{})

	client := &testmock.ControlledGHClient{
		Repo:           "owner/repo1",
		PRs:            []gh.PRItem{{Number: 1, Title: "Fix bug"}},
		PRView:         "PR detail",
		ResolveCalled:  make(chan struct{}),
		PRsCalled:      make(chan struct{}),
		DetailCalled:   make(chan struct{}),
		ReleaseResolve: releaseResolve,
		ReleasePRs:     releasePRs,
		ReleaseDetail:  releaseDetail,
	}
	g, err := NewGui(config.Default(), client)
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}

	p := tea.NewProgram(
		&model{gui: g},
		tea.WithInput(bytes.NewBuffer(nil)),
		tea.WithOutput(io.Discard),
		tea.WithoutSignals(),
	)

	runDone := make(chan error, 1)
	go func() {
		_, runErr := p.Run()
		runDone <- runErr
	}()

	waitChan(t, client.ResolveCalled, "ResolveCurrentRepo was not called")
	if !g.state.PRsLoading {
		t.Fatal("prs panel should stay loading before release")
	}
	close(releaseResolve)

	waitChan(t, client.PRsCalled, "ListPRs was not called")
	close(releasePRs)

	waitUntil(t, func() bool {
		return g.state.Repo == "owner/repo1" && len(g.state.PRs) == 1
	}, "prs were not loaded")

	p.Send(tea.KeyMsg{Type: tea.KeyEnter})
	waitChan(t, client.DetailCalled, "ViewPR was not called")
	close(releaseDetail)
	waitUntil(t, func() bool {
		return g.state.DetailContent == "PR detail"
	}, "detail content was not updated")

	p.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	select {
	case runErr := <-runDone:
		if runErr != nil {
			t.Fatalf("program failed: %v", runErr)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("program did not quit")
	}
}

func TestProgramE2E_ShowErrorWhenInitialLoadFails(t *testing.T) {
	client := &testmock.ControlledGHClient{
		Err:           errors.New("boom"),
		ResolveCalled: make(chan struct{}),
		PRsCalled:     make(chan struct{}),
		DetailCalled:  make(chan struct{}),
	}
	g, err := NewGui(config.Default(), client)
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}

	p := tea.NewProgram(
		&model{gui: g},
		tea.WithInput(bytes.NewBuffer(nil)),
		tea.WithOutput(io.Discard),
		tea.WithoutSignals(),
	)

	runDone := make(chan error, 1)
	go func() {
		_, runErr := p.Run()
		runDone <- runErr
	}()

	waitChan(t, client.ResolveCalled, "ResolveCurrentRepo was not called")
	waitUntil(t, func() bool {
		return strings.Contains(g.state.DetailContent, "Error loading PRs")
	}, "error message was not shown")

	p.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	select {
	case runErr := <-runDone:
		if runErr != nil {
			t.Fatalf("program failed: %v", runErr)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("program did not quit")
	}
}

func waitChan(t *testing.T, ch <-chan struct{}, msg string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(2 * time.Second):
		t.Fatal(msg)
	}
}

func waitUntil(t *testing.T, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal(msg)
}
