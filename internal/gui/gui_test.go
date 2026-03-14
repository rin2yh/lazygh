package gui

import (
	"bytes"
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestGuiRun_LoadsPRsAndDetail(t *testing.T) {
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
		&screen{gui: g},
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
