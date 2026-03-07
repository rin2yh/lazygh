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
	repos     []string
	issues    []gh.IssueItem
	prs       []gh.PRItem
	issueView string
	prView    string
	err       error

	reposCalled  chan struct{}
	itemsCalled  chan struct{}
	detailCalled chan struct{}

	releaseRepos  <-chan struct{}
	releaseItems  <-chan struct{}
	releaseDetail <-chan struct{}

	reposOnce  sync.Once
	itemsOnce  sync.Once
	detailOnce sync.Once
}

func (c *controlledClient) ListRepos() ([]string, error) {
	c.reposOnce.Do(func() { close(c.reposCalled) })
	if c.releaseRepos != nil {
		<-c.releaseRepos
	}
	if c.err != nil {
		return nil, c.err
	}
	return c.repos, nil
}

func (c *controlledClient) ListPRs(_ string) ([]gh.PRItem, error) {
	c.itemsOnce.Do(func() { close(c.itemsCalled) })
	if c.releaseItems != nil {
		<-c.releaseItems
	}
	if c.err != nil {
		return nil, c.err
	}
	return c.prs, nil
}

func (c *controlledClient) ListIssues(_ string) ([]gh.IssueItem, error) {
	c.itemsOnce.Do(func() { close(c.itemsCalled) })
	if c.releaseItems != nil {
		<-c.releaseItems
	}
	if c.err != nil {
		return nil, c.err
	}
	return c.issues, nil
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

func (c *controlledClient) ViewIssue(_ string, _ int) (string, error) {
	c.detailOnce.Do(func() { close(c.detailCalled) })
	if c.releaseDetail != nil {
		<-c.releaseDetail
	}
	if c.err != nil {
		return "", c.err
	}
	return c.issueView, nil
}

func TestProgramE2E_MainFlow(t *testing.T) {
	releaseRepos := make(chan struct{})
	releaseItems := make(chan struct{})
	releaseDetail := make(chan struct{})

	client := &controlledClient{
		repos:         []string{"owner/repo1"},
		issues:        []gh.IssueItem{{Number: 10, Title: "Issue one"}},
		prs:           []gh.PRItem{{Number: 1, Title: "Fix bug"}},
		issueView:     "Issue detail",
		reposCalled:   make(chan struct{}),
		itemsCalled:   make(chan struct{}),
		detailCalled:  make(chan struct{}),
		releaseRepos:  releaseRepos,
		releaseItems:  releaseItems,
		releaseDetail: releaseDetail,
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

	waitChan(t, client.reposCalled, "ListRepos was not called")
	if !g.state.ReposLoading {
		t.Fatal("repos panel should stay loading before release")
	}
	close(releaseRepos)

	waitUntil(t, func() bool {
		return g.state.ReposLoaded && len(g.state.Repos) == 1
	}, "repos were not loaded")

	p.Send(tea.KeyMsg{Type: tea.KeyEnter})
	waitChan(t, client.itemsCalled, "ListIssues/ListPRs was not called")
	close(releaseItems)

	waitUntil(t, func() bool {
		return len(g.state.Issues) == 1 && len(g.state.PRs) == 1
	}, "issues/prs were not loaded")

	p.Send(tea.KeyMsg{Type: tea.KeyTab})
	waitUntil(t, func() bool {
		return g.state.ActivePanel == PanelIssues
	}, "active panel did not move to issues")

	p.Send(tea.KeyMsg{Type: tea.KeyEnter})
	waitChan(t, client.detailCalled, "ViewIssue/ViewPR was not called")
	close(releaseDetail)
	waitUntil(t, func() bool {
		return g.state.DetailContent == "Issue detail"
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

func TestProgramE2E_ShowErrorWhenRepoLoadFails(t *testing.T) {
	client := &controlledClient{
		err:          errors.New("boom"),
		reposCalled:  make(chan struct{}),
		itemsCalled:  make(chan struct{}),
		detailCalled: make(chan struct{}),
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

	waitChan(t, client.reposCalled, "ListRepos was not called")
	waitUntil(t, func() bool {
		return strings.Contains(g.state.DetailContent, "Error loading repos")
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
