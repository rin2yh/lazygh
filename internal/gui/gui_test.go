package gui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/panels"
)

// mockClient は gh.ClientInterface のテスト用実装。
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

func newTestGui() *Gui {
	return &Gui{
		g:      nil,
		config: config.Default(),
		state:  &State{ActivePanel: PanelRepos},
		panels: &Panels{
			Repos:  panels.NewItemsPanel(panels.FormatRepoItem, true),
			Issues: panels.NewItemsPanel(panels.FormatIssueItem, false),
			PRs:    panels.NewItemsPanel(panels.FormatPRItem, false),
			Detail: panels.NewDetailPanel(),
		},
	}
}

func newTestGuiWithClient(client gh.ClientInterface) *Gui {
	g := newTestGui()
	g.client = client
	return g
}

func repoItems(repos ...string) []panels.Item {
	items := make([]panels.Item, 0, len(repos))
	for _, repo := range repos {
		items = append(items, panels.Item{Title: repo})
	}
	return items
}

func TestActiveViewName(t *testing.T) {
	tests := []struct {
		panel PanelType
		want  string
	}{
		{PanelRepos, "repos"},
		{PanelIssues, "issues"},
		{PanelPRs, "prs"},
		{PanelDetail, "detail"},
	}
	for _, tt := range tests {
		g := newTestGui()
		g.state.ActivePanel = tt.panel
		got := g.activeViewName()
		if got != tt.want {
			t.Errorf("panel=%d: got %q, want %q", tt.panel, got, tt.want)
		}
	}
}

func TestNextPanel_Cycle(t *testing.T) {
	tests := []struct {
		before PanelType
		after  PanelType
	}{
		{PanelRepos, PanelIssues},
		{PanelIssues, PanelPRs},
		{PanelPRs, PanelDetail},
		{PanelDetail, PanelRepos},
	}
	for _, tt := range tests {
		g := newTestGui()
		g.state.ActivePanel = tt.before
		_ = g.nextPanel(nil, nil)
		if g.state.ActivePanel != tt.after {
			t.Errorf("after %d: got %d, want %d", tt.before, g.state.ActivePanel, tt.after)
		}
	}
}

func TestNavigateDown_Repos(t *testing.T) {
	g := newTestGui()
	g.state.ActivePanel = PanelRepos
	g.panels.Repos.Items = repoItems("a", "b", "c")
	_ = g.navigateDown(nil, "repos")

	if g.panels.Repos.Selected != 1 {
		t.Errorf("repos: got %d, want %d", g.panels.Repos.Selected, 1)
	}
}

func TestNavigateDown_ListPanels(t *testing.T) {
	tests := []struct {
		name  string
		view  string
		panel PanelType
		items []panels.Item
	}{
		{
			name:  "Issues",
			view:  "issues",
			panel: PanelIssues,
			items: []panels.Item{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}},
		},
		{
			name:  "PRs",
			view:  "prs",
			panel: PanelPRs,
			items: []panels.Item{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}},
		},
	}

	for _, tt := range tests {
		g := newTestGui()
		g.state.ActivePanel = tt.panel
		switch tt.panel {
		case PanelIssues:
			g.panels.Issues.Items = tt.items
		case PanelPRs:
			g.panels.PRs.Items = tt.items
		}

		_ = g.navigateDown(nil, tt.view)

		var got int
		switch tt.panel {
		case PanelIssues:
			got = g.panels.Issues.Selected
		case PanelPRs:
			got = g.panels.PRs.Selected
		}
		if got != 1 {
			t.Errorf("%s: got %d, want %d", tt.name, got, 1)
		}
	}
}

func TestNavigateUp_Repos(t *testing.T) {
	g := newTestGui()
	g.state.ActivePanel = PanelRepos
	g.panels.Repos.Items = repoItems("a", "b", "c")
	g.panels.Repos.Selected = 2
	_ = g.navigateUp(nil, "repos")

	if g.panels.Repos.Selected != 1 {
		t.Errorf("repos: got %d, want %d", g.panels.Repos.Selected, 1)
	}
}

func TestNavigateUp_ListPanels(t *testing.T) {
	tests := []struct {
		name  string
		view  string
		panel PanelType
		items []panels.Item
	}{
		{
			name:  "Issues",
			view:  "issues",
			panel: PanelIssues,
			items: []panels.Item{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}},
		},
		{
			name:  "PRs",
			view:  "prs",
			panel: PanelPRs,
			items: []panels.Item{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}},
		},
	}

	for _, tt := range tests {
		g := newTestGui()
		g.state.ActivePanel = tt.panel
		switch tt.panel {
		case PanelIssues:
			g.panels.Issues.Items = tt.items
			g.panels.Issues.Selected = 1
		case PanelPRs:
			g.panels.PRs.Items = tt.items
			g.panels.PRs.Selected = 1
		}

		_ = g.navigateUp(nil, tt.view)

		var got int
		switch tt.panel {
		case PanelIssues:
			got = g.panels.Issues.Selected
		case PanelPRs:
			got = g.panels.PRs.Selected
		}
		if got != 0 {
			t.Errorf("%s: got %d, want %d", tt.name, got, 0)
		}
	}
}

// --- renderPanel ---

func TestRenderPanel_NilGui(t *testing.T) {
	g := newTestGui() // g.g == nil
	// panic しないことを確認
	g.renderPanel("repos")
	g.renderPanel("issues")
	g.renderPanel("prs")
	g.renderPanel("detail")
}

// --- applyReposResult ---

func TestApplyReposResult_Success(t *testing.T) {
	g := newTestGui()
	g.panels.Repos.Loading = true
	_ = g.applyReposResult([]string{"owner/repo1", "owner/repo2"}, nil)
	if g.panels.Repos.Loading {
		t.Error("Loading should be false after success")
	}
	if len(g.panels.Repos.Items) != 2 {
		t.Fatalf("got %d repos, want 2", len(g.panels.Repos.Items))
	}
	if got := g.panels.Repos.Format(g.panels.Repos.Items[0]); got != "owner/repo1" {
		t.Errorf("repos[0] = %q, want %q", got, "owner/repo1")
	}
	if !g.reposLoaded {
		t.Error("reposLoaded should be true")
	}
}

func TestApplyReposResult_Error(t *testing.T) {
	g := newTestGui()
	g.panels.Repos.Loading = true
	_ = g.applyReposResult(nil, fmt.Errorf("api error"))
	if g.panels.Repos.Loading {
		t.Error("Loading should be false after error")
	}
	if !strings.Contains(g.panels.Detail.Content, "api error") {
		t.Errorf("detail content %q should contain error", g.panels.Detail.Content)
	}
}

// --- loadRepos ---

func TestLoadRepos_PopulatesPanel(t *testing.T) {
	mc := &mockClient{repos: []string{"owner/repo1", "owner/repo2"}}
	g := newTestGuiWithClient(mc)
	if err := g.loadRepos(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.panels.Repos.Items) != 2 {
		t.Fatalf("got %d repos, want 2", len(g.panels.Repos.Items))
	}
	if got := g.panels.Repos.Format(g.panels.Repos.Items[0]); got != "owner/repo1" {
		t.Errorf("repos[0] = %q, want %q", got, "owner/repo1")
	}
	if g.panels.Repos.Selected != 0 {
		t.Errorf("Selected = %d, want 0", g.panels.Repos.Selected)
	}
	if !g.reposLoaded {
		t.Error("reposLoaded should be true")
	}
}

func TestLoadRepos_ErrorShowsInDetail(t *testing.T) {
	mc := &mockClient{err: fmt.Errorf("gh not found")}
	g := newTestGuiWithClient(mc)
	_ = g.loadRepos()
	if !strings.Contains(g.panels.Detail.Content, "gh not found") {
		t.Errorf("detail content %q should contain error", g.panels.Detail.Content)
	}
}

// --- loadItems ---

func TestLoadItems_PopulatesPanels(t *testing.T) {
	mc := &mockClient{
		prs:    []gh.PRItem{{Number: 1, Title: "Fix bug"}, {Number: 2, Title: "Add feat"}},
		issues: []gh.IssueItem{{Number: 10, Title: "Issue one"}},
	}
	g := newTestGuiWithClient(mc)
	g.panels.Repos.Items = repoItems("owner/repo")
	g.panels.Repos.Selected = 0

	if err := g.loadItems(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.panels.Issues.Items) != 1 {
		t.Fatalf("got %d issues, want 1", len(g.panels.Issues.Items))
	}
	if got := g.panels.Issues.Format(g.panels.Issues.Items[0]); got != "Issue #10 Issue one" {
		t.Errorf("issues[0] = %q, want %q", got, "Issue #10 Issue one")
	}
	if len(g.panels.PRs.Items) != 2 {
		t.Fatalf("got %d prs, want 2", len(g.panels.PRs.Items))
	}
	if got := g.panels.PRs.Format(g.panels.PRs.Items[0]); got != "PR #1 Fix bug" {
		t.Errorf("prs[0] = %q, want %q", got, "PR #1 Fix bug")
	}
	if g.panels.Issues.Selected != 0 {
		t.Errorf("Issues.Selected = %d, want 0", g.panels.Issues.Selected)
	}
	if g.panels.PRs.Selected != 0 {
		t.Errorf("PRs.Selected = %d, want 0", g.panels.PRs.Selected)
	}
}

func TestLoadItems_EmptyRepos(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	// repos が空なので何もしない
	if err := g.loadItems(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.panels.Issues.Items) != 0 {
		t.Errorf("issues should be empty")
	}
	if len(g.panels.PRs.Items) != 0 {
		t.Errorf("prs should be empty")
	}
}

// --- loadDetail ---

func TestLoadDetail(t *testing.T) {
	tests := []struct {
		name  string
		panel PanelType
		item  panels.Item
		mc    *mockClient
		want  string
	}{
		{
			name:  "Issue",
			panel: PanelIssues,
			item:  panels.Item{Number: 10, Title: "Bug"},
			mc:    &mockClient{issueView: "Issue detail content"},
			want:  "Issue detail content",
		},
		{
			name:  "PR",
			panel: PanelPRs,
			item:  panels.Item{Number: 1, Title: "Fix"},
			mc:    &mockClient{prView: "PR detail content"},
			want:  "PR detail content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithClient(tt.mc)
			g.state.ActivePanel = tt.panel
			g.panels.Repos.Items = repoItems("owner/repo")
			switch tt.panel {
			case PanelIssues:
				g.panels.Issues.Items = []panels.Item{tt.item}
			case PanelPRs:
				g.panels.PRs.Items = []panels.Item{tt.item}
			}
			if err := g.loadDetail(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if g.panels.Detail.Content != tt.want {
				t.Errorf("got %q, want %q", g.panels.Detail.Content, tt.want)
			}
		})
	}
}

func TestLoadDetail_NoopCases(t *testing.T) {
	tests := []struct {
		name  string
		panel PanelType
	}{
		{name: "EmptyItems", panel: PanelIssues},
		{name: "NotListPanel", panel: PanelRepos},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithClient(&mockClient{})
			g.state.ActivePanel = tt.panel
			g.panels.Repos.Items = repoItems("owner/repo")
			g.panels.Detail.SetContent("keep")
			if err := g.loadDetail(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if g.panels.Detail.Content != "keep" {
				t.Errorf("got %q, want %q", g.panels.Detail.Content, "keep")
			}
		})
	}
}

func TestNilClient(t *testing.T) {
	t.Run("LoadRepos", func(t *testing.T) {
		g := newTestGui()
		if err := g.loadRepos(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(g.panels.Repos.Items) != 0 {
			t.Errorf("repos should be empty, got %v", g.panels.Repos.Items)
		}
	})
	t.Run("LoadItems", func(t *testing.T) {
		g := newTestGui()
		g.panels.Repos.Items = repoItems("owner/repo")
		if err := g.loadItems(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("LoadDetail", func(t *testing.T) {
		g := newTestGui()
		g.state.ActivePanel = PanelPRs
		g.panels.Repos.Items = repoItems("owner/repo")
		g.panels.PRs.Items = []panels.Item{{Number: 1, Title: "x"}}
		if err := g.loadDetail(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// --- refreshDetailPreview ---

func TestRefreshDetailPreview_SetsContent(t *testing.T) {
	tests := []struct {
		name  string
		panel PanelType
		item  panels.Item
	}{
		{
			name:  "Issues",
			panel: PanelIssues,
			item:  panels.Item{Number: 7, Title: "My Issue"},
		},
		{
			name:  "PRs",
			panel: PanelPRs,
			item:  panels.Item{Number: 42, Title: "My PR"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGui()
			g.state.ActivePanel = tt.panel
			switch tt.panel {
			case PanelIssues:
				g.panels.Issues.Items = []panels.Item{tt.item}
			case PanelPRs:
				g.panels.PRs.Items = []panels.Item{tt.item}
			}
			g.refreshDetailPreview()

			want := panels.FormatIssueItem(tt.item)
			if tt.panel == PanelPRs {
				want = panels.FormatPRItem(tt.item)
			}
			if g.panels.Detail.Content != want {
				t.Errorf("got %q, want %q", g.panels.Detail.Content, want)
			}
		})
	}
}

func TestRefreshDetailPreview_EmptyItems(t *testing.T) {
	g := newTestGui()
	g.state.ActivePanel = PanelIssues
	// panic しないことを確認
	g.refreshDetailPreview()
}
