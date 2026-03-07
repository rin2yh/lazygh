package gui

import (
	"errors"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
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

func repoItems(repos ...string) []core.Item {
	items := make([]core.Item, 0, len(repos))
	for _, repo := range repos {
		items = append(items, core.Item{Title: repo})
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
	g.state.PRs = []core.Item{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}}

	g.navigateDown()
	if g.state.PRsSelected != 1 {
		t.Fatalf("got %d, want 1", g.state.PRsSelected)
	}

	g.navigateUp()
	if g.state.PRsSelected != 0 {
		t.Fatalf("got %d, want 0", g.state.PRsSelected)
	}
}

func TestApplyReposResult(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.ReposLoading = true
	g.applyReposResult([]string{"owner/repo1", "owner/repo2"}, nil)
	if !g.state.ReposLoaded {
		t.Fatal("reposLoaded should be true")
	}
	if g.state.ReposLoading {
		t.Fatal("Loading should be false")
	}
	if len(g.state.Repos) != 2 {
		t.Fatalf("got %d, want 2", len(g.state.Repos))
	}
}

func TestApplyItemsResult(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.Repos = repoItems("owner/repo")
	g.state.IssuesLoading = true
	g.state.PRsLoading = true

	g.applyItemsResult(itemsLoadedMsg{
		repo:   "owner/repo",
		issues: []core.Item{{Number: 10, Title: "Issue one"}},
		prs:    []core.Item{{Number: 1, Title: "Fix bug"}},
	})

	if g.state.IssuesLoading || g.state.PRsLoading {
		t.Fatal("loading flags should be false")
	}
	if got := core.FormatIssueItem(g.state.Issues[0]); got != "Issue #10 Issue one" {
		t.Fatalf("unexpected issue row: %q", got)
	}
	if got := core.FormatPRItem(g.state.PRs[0]); got != "PR #1 Fix bug" {
		t.Fatalf("unexpected pr row: %q", got)
	}
}

func TestApplyItemsResult_MismatchRepoIgnored(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.Repos = repoItems("owner/repo")
	g.applyItemsResult(itemsLoadedMsg{
		repo:   "other/repo",
		issues: []core.Item{{Number: 10, Title: "Issue one"}},
		prs:    []core.Item{{Number: 1, Title: "Fix bug"}},
	})
	if len(g.state.Issues) != 0 || len(g.state.PRs) != 0 {
		t.Fatal("items should remain empty when repo mismatched")
	}
}

func TestApplyDetailResult(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.applyDetailResult(detailLoadedMsg{content: "hello"})
	if g.state.DetailContent != "hello" {
		t.Fatalf("got %q, want hello", g.state.DetailContent)
	}
}

func TestApplyDetailResult_Error(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.applyDetailResult(detailLoadedMsg{err: errors.New("boom")})
	if g.state.DetailContent == "" {
		t.Fatal("error message should be shown")
	}
}

func TestModelHandleEnterRepos(t *testing.T) {
	mc := &mockClient{issues: []gh.IssueItem{{Number: 1, Title: "i"}}, prs: []gh.PRItem{{Number: 2, Title: "p"}}}
	g := newTestGuiWithClient(mc)
	g.state.Repos = repoItems("owner/repo")
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
	g.state.Repos = repoItems("owner/repo")
	g.state.PRs = []core.Item{{Number: 1, Title: "x"}}
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

func TestNavigateDetailPanelScroll(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.ActivePanel = PanelDetail
	g.syncDetailViewport(20, 2, "line1\nline2\nline3\nline4")

	g.navigateDown()
	if g.detailViewport.YOffset <= 0 {
		t.Fatalf("expected detail viewport to scroll down, got offset=%d", g.detailViewport.YOffset)
	}

	g.navigateUp()
	if g.detailViewport.YOffset != 0 {
		t.Fatalf("expected detail viewport to scroll up, got offset=%d", g.detailViewport.YOffset)
	}
}
