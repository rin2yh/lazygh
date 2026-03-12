package gui

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
)

type controlledClient struct {
	repo   string
	prs    []gh.PRItem
	prView string
	err    error

	resolveCalled chan struct{}
	prsCalled     chan struct{}
	detailCalled  chan struct{}

	releaseResolve <-chan struct{}
	releasePRs     <-chan struct{}
	releaseDetail  <-chan struct{}

	resolveOnce sync.Once
	prsOnce     sync.Once
	detailOnce  sync.Once
}

func (c *controlledClient) ResolveCurrentRepo() (string, error) {
	c.resolveOnce.Do(func() { close(c.resolveCalled) })
	if c.releaseResolve != nil {
		<-c.releaseResolve
	}
	if c.err != nil {
		return "", c.err
	}
	return c.repo, nil
}

func (c *controlledClient) ListPRs(_ string) ([]gh.PRItem, error) {
	c.prsOnce.Do(func() { close(c.prsCalled) })
	if c.releasePRs != nil {
		<-c.releasePRs
	}
	if c.err != nil {
		return nil, c.err
	}
	return c.prs, nil
}

func (c *controlledClient) ViewPR(_ string, _ int) (string, error) {
	c.detailOnce.Do(func() { close(c.detailCalled) })
	if c.releaseDetail != nil {
		<-c.releaseDetail
	}
	if c.err != nil {
		return "", c.err
	}
	return c.prView, nil
}

func TestProgramE2E_MainFlow(t *testing.T) {
	releaseResolve := make(chan struct{})
	releasePRs := make(chan struct{})
	releaseDetail := make(chan struct{})

	client := &controlledClient{
		repo:           "owner/repo1",
		prs:            []gh.PRItem{{Number: 1, Title: "Fix bug"}},
		prView:         "PR detail",
		resolveCalled:  make(chan struct{}),
		prsCalled:      make(chan struct{}),
		detailCalled:   make(chan struct{}),
		releaseResolve: releaseResolve,
		releasePRs:     releasePRs,
		releaseDetail:  releaseDetail,
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

	waitChan(t, client.resolveCalled, "ResolveCurrentRepo was not called")
	if !g.state.PRsLoading {
		t.Fatal("prs panel should stay loading before release")
	}
	close(releaseResolve)

	waitChan(t, client.prsCalled, "ListPRs was not called")
	close(releasePRs)

	waitUntil(t, func() bool {
		return g.state.Repo == "owner/repo1" && len(g.state.PRs) == 1
	}, "prs were not loaded")

	p.Send(tea.KeyMsg{Type: tea.KeyEnter})
	waitChan(t, client.detailCalled, "ViewPR was not called")
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
	client := &controlledClient{
		err:           errors.New("boom"),
		resolveCalled: make(chan struct{}),
		prsCalled:     make(chan struct{}),
		detailCalled:  make(chan struct{}),
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

	waitChan(t, client.resolveCalled, "ResolveCurrentRepo was not called")
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
