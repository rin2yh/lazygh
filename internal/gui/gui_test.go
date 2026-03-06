package gui

import (
	"errors"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/panels"
)

type mockClient struct {
	repos     []string
	prs       []gh.PRItem
	issues    []gh.IssueItem
	prView    string
	issueView string
	err       error
}

func (m *mockClient) ListRepos() ([]string, error)                { return m.repos, m.err }
func (m *mockClient) ListPRs(_ string) ([]gh.PRItem, error)       { return m.prs, m.err }
func (m *mockClient) ListIssues(_ string) ([]gh.IssueItem, error) { return m.issues, m.err }
func (m *mockClient) ViewPR(_ string, _ int) (string, error)      { return m.prView, m.err }
func (m *mockClient) ViewIssue(_ string, _ int) (string, error)   { return m.issueView, m.err }

func newTestGuiWithClient(client gh.ClientInterface) *Gui {
	g, _ := NewGui(config.Default(), client)
	return g
}

func repoItems(repos ...string) []panels.Item {
	items := make([]panels.Item, 0, len(repos))
	for _, repo := range repos {
		items = append(items, panels.Item{Title: repo})
	}
	return items
}

func TestPanelCycle(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.nextPanel()
	if g.state.ActivePanel != PanelIssues {
		t.Fatalf("got %v, want %v", g.state.ActivePanel, PanelIssues)
	}
	g.prevPanel()
	if g.state.ActivePanel != PanelRepos {
		t.Fatalf("got %v, want %v", g.state.ActivePanel, PanelRepos)
	}
}

func TestNavigateListPanels(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.ActivePanel = PanelPRs
	g.panels.PRs.Items = []panels.Item{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}}

	g.navigateDown()
	if g.panels.PRs.Selected != 1 {
		t.Fatalf("got %d, want 1", g.panels.PRs.Selected)
	}

	g.navigateUp()
	if g.panels.PRs.Selected != 0 {
		t.Fatalf("got %d, want 0", g.panels.PRs.Selected)
	}
}

func TestApplyReposResult(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.panels.Repos.Loading = true
	g.applyReposResult([]string{"owner/repo1", "owner/repo2"}, nil)
	if !g.reposLoaded {
		t.Fatal("reposLoaded should be true")
	}
	if g.panels.Repos.Loading {
		t.Fatal("Loading should be false")
	}
	if len(g.panels.Repos.Items) != 2 {
		t.Fatalf("got %d, want 2", len(g.panels.Repos.Items))
	}
}

func TestApplyItemsResult(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.panels.Repos.Items = repoItems("owner/repo")
	g.panels.Issues.Loading = true
	g.panels.PRs.Loading = true

	g.applyItemsResult(itemsLoadedMsg{
		repo:   "owner/repo",
		issues: []gh.IssueItem{{Number: 10, Title: "Issue one"}},
		prs:    []gh.PRItem{{Number: 1, Title: "Fix bug"}},
	})

	if g.panels.Issues.Loading || g.panels.PRs.Loading {
		t.Fatal("loading flags should be false")
	}
	if got := g.panels.Issues.Format(g.panels.Issues.Items[0]); got != "Issue #10 Issue one" {
		t.Fatalf("unexpected issue row: %q", got)
	}
	if got := g.panels.PRs.Format(g.panels.PRs.Items[0]); got != "PR #1 Fix bug" {
		t.Fatalf("unexpected pr row: %q", got)
	}
}

func TestApplyItemsResult_MismatchRepoIgnored(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.panels.Repos.Items = repoItems("owner/repo")
	g.applyItemsResult(itemsLoadedMsg{
		repo:   "other/repo",
		issues: []gh.IssueItem{{Number: 10, Title: "Issue one"}},
		prs:    []gh.PRItem{{Number: 1, Title: "Fix bug"}},
	})
	if len(g.panels.Issues.Items) != 0 || len(g.panels.PRs.Items) != 0 {
		t.Fatal("items should remain empty when repo mismatched")
	}
}

func TestApplyDetailResult(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.applyDetailResult(detailLoadedMsg{content: "hello"})
	if g.panels.Detail.Content != "hello" {
		t.Fatalf("got %q, want hello", g.panels.Detail.Content)
	}
}

func TestApplyDetailResult_Error(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.applyDetailResult(detailLoadedMsg{err: errors.New("boom")})
	if g.panels.Detail.Content == "" {
		t.Fatal("error message should be shown")
	}
}

func TestModelHandleEnterRepos(t *testing.T) {
	mc := &mockClient{issues: []gh.IssueItem{{Number: 1, Title: "i"}}, prs: []gh.PRItem{{Number: 2, Title: "p"}}}
	g := newTestGuiWithClient(mc)
	g.panels.Repos.Items = repoItems("owner/repo")
	m := &model{gui: g}
	cmd := m.handleEnter()
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd().(itemsLoadedMsg)
	if msg.err != nil || len(msg.issues) != 1 || len(msg.prs) != 1 {
		t.Fatal("unexpected itemsLoadedMsg")
	}
}

func TestModelHandleEnterDetail(t *testing.T) {
	mc := &mockClient{prView: "detail"}
	g := newTestGuiWithClient(mc)
	g.state.ActivePanel = PanelPRs
	g.panels.Repos.Items = repoItems("owner/repo")
	g.panels.PRs.Items = []panels.Item{{Number: 1, Title: "x"}}
	m := &model{gui: g}
	cmd := m.handleEnter()
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd().(detailLoadedMsg)
	if msg.err != nil || msg.content != "detail" {
		t.Fatalf("unexpected msg: %+v", msg)
	}
}
