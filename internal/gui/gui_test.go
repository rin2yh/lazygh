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
			Repos:  panels.NewReposPanel(),
			Items:  panels.NewItemsPanel(),
			Detail: panels.NewDetailPanel(),
		},
	}
}

func newTestGuiWithClient(client gh.ClientInterface) *Gui {
	g := newTestGui()
	g.client = client
	return g
}

func TestActiveViewName(t *testing.T) {
	tests := []struct {
		panel PanelType
		want  string
	}{
		{PanelRepos, "repos"},
		{PanelItems, "items"},
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
		{PanelRepos, PanelItems},
		{PanelItems, PanelDetail},
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
	tests := []struct {
		name        string
		active      PanelType
		viewName    string
		repos       []string
		repoSel     int
		wantRepoSel int
		items       []panels.Item
		itemSel     int
		wantItemSel int
	}{
		{"Empty", PanelRepos, "repos", nil, 0, 0, nil, 0, 0},
		{"UpperBound", PanelRepos, "repos", []string{"a", "b", "c"}, 2, 2, nil, 0, 0},
		{"Normal", PanelRepos, "repos", []string{"a", "b", "c"}, 0, 1, nil, 0, 0},
		{"InactivePanel", PanelItems, "repos", []string{"a", "b"}, 0, 0, nil, 0, 0},
		{
			name:        "ItemsNormal",
			active:      PanelItems,
			viewName:    "items",
			items:       []panels.Item{{Kind: panels.ItemKindPR, Number: 1, Title: "a"}, {Kind: panels.ItemKindPR, Number: 2, Title: "b"}},
			wantItemSel: 1,
		},
	}
	for _, tt := range tests {
		g := newTestGui()
		g.state.ActivePanel = tt.active
		g.panels.Repos.Repos = tt.repos
		g.panels.Repos.Selected = tt.repoSel
		g.panels.Items.Items = tt.items
		g.panels.Items.Selected = tt.itemSel
		_ = g.navigateDown(nil, tt.viewName)
		if g.panels.Repos.Selected != tt.wantRepoSel {
			t.Errorf("%s: repos: got %d, want %d", tt.name, g.panels.Repos.Selected, tt.wantRepoSel)
		}
		if g.panels.Items.Selected != tt.wantItemSel {
			t.Errorf("%s: items: got %d, want %d", tt.name, g.panels.Items.Selected, tt.wantItemSel)
		}
	}
}

func TestNavigateUp_Repos(t *testing.T) {
	tests := []struct {
		name        string
		active      PanelType
		viewName    string
		repos       []string
		repoSel     int
		wantRepoSel int
		items       []panels.Item
		itemSel     int
		wantItemSel int
	}{
		{"LowerBound", PanelRepos, "repos", []string{"a", "b"}, 0, 0, nil, 0, 0},
		{"Normal", PanelRepos, "repos", []string{"a", "b", "c"}, 2, 1, nil, 0, 0},
		{
			name:        "ItemsNormal",
			active:      PanelItems,
			viewName:    "items",
			items:       []panels.Item{{Kind: panels.ItemKindPR, Number: 1, Title: "a"}, {Kind: panels.ItemKindPR, Number: 2, Title: "b"}},
			itemSel:     1,
			wantItemSel: 0,
		},
	}
	for _, tt := range tests {
		g := newTestGui()
		g.state.ActivePanel = tt.active
		g.panels.Repos.Repos = tt.repos
		g.panels.Repos.Selected = tt.repoSel
		g.panels.Items.Items = tt.items
		g.panels.Items.Selected = tt.itemSel
		_ = g.navigateUp(nil, tt.viewName)
		if g.panels.Repos.Selected != tt.wantRepoSel {
			t.Errorf("%s: repos: got %d, want %d", tt.name, g.panels.Repos.Selected, tt.wantRepoSel)
		}
		if g.panels.Items.Selected != tt.wantItemSel {
			t.Errorf("%s: items: got %d, want %d", tt.name, g.panels.Items.Selected, tt.wantItemSel)
		}
	}
}

// --- renderPanel ---

func TestRenderPanel_NilGui(t *testing.T) {
	g := newTestGui() // g.g == nil
	// panic しないことを確認
	g.renderPanel("repos")
	g.renderPanel("items")
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
	if len(g.panels.Repos.Repos) != 2 {
		t.Fatalf("got %d repos, want 2", len(g.panels.Repos.Repos))
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
	if len(g.panels.Repos.Repos) != 2 {
		t.Fatalf("got %d repos, want 2", len(g.panels.Repos.Repos))
	}
	if g.panels.Repos.Repos[0] != "owner/repo1" {
		t.Errorf("repos[0] = %q, want %q", g.panels.Repos.Repos[0], "owner/repo1")
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

func TestLoadItems_PopulatesPanel(t *testing.T) {
	mc := &mockClient{
		prs:    []gh.PRItem{{Number: 1, Title: "Fix bug"}, {Number: 2, Title: "Add feat"}},
		issues: []gh.IssueItem{{Number: 10, Title: "Issue one"}},
	}
	g := newTestGuiWithClient(mc)
	g.panels.Repos.Repos = []string{"owner/repo"}
	g.panels.Repos.Selected = 0

	if err := g.loadItems(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.panels.Items.Items) != 3 {
		t.Fatalf("got %d items, want 3", len(g.panels.Items.Items))
	}
	// PR が先、Issue が後
	if g.panels.Items.Items[0].Kind != panels.ItemKindPR {
		t.Error("items[0] should be PR")
	}
	if g.panels.Items.Items[2].Kind != panels.ItemKindIssue {
		t.Error("items[2] should be Issue")
	}
	if g.panels.Items.Selected != 0 {
		t.Errorf("Selected = %d, want 0", g.panels.Items.Selected)
	}
}

func TestLoadItems_EmptyRepos(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	// repos が空なので何もしない
	if err := g.loadItems(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.panels.Items.Items) != 0 {
		t.Errorf("items should be empty")
	}
}

// --- loadDetail ---

func TestLoadDetail(t *testing.T) {
	tests := []struct {
		name string
		item panels.Item
		mc   *mockClient
		want string
	}{
		{
			name: "PR",
			item: panels.Item{Kind: panels.ItemKindPR, Number: 1, Title: "Fix"},
			mc:   &mockClient{prView: "PR detail content"},
			want: "PR detail content",
		},
		{
			name: "Issue",
			item: panels.Item{Kind: panels.ItemKindIssue, Number: 10, Title: "Bug"},
			mc:   &mockClient{issueView: "Issue detail content"},
			want: "Issue detail content",
		},
	}
	for _, tt := range tests {
		g := newTestGuiWithClient(tt.mc)
		g.panels.Repos.Repos = []string{"owner/repo"}
		g.panels.Items.Items = []panels.Item{tt.item}
		g.panels.Items.Selected = 0
		if err := g.loadDetail(); err != nil {
			t.Fatalf("%s: unexpected error: %v", tt.name, err)
		}
		if g.panels.Detail.Content != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, g.panels.Detail.Content, tt.want)
		}
	}
}

func TestLoadDetail_EmptyItems(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.panels.Repos.Repos = []string{"owner/repo"}
	if err := g.loadDetail(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNilClient(t *testing.T) {
	t.Run("LoadRepos", func(t *testing.T) {
		g := newTestGui()
		if err := g.loadRepos(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(g.panels.Repos.Repos) != 0 {
			t.Errorf("repos should be empty, got %v", g.panels.Repos.Repos)
		}
	})
	t.Run("LoadItems", func(t *testing.T) {
		g := newTestGui()
		g.panels.Repos.Repos = []string{"owner/repo"}
		if err := g.loadItems(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("LoadDetail", func(t *testing.T) {
		g := newTestGui()
		g.panels.Repos.Repos = []string{"owner/repo"}
		g.panels.Items.Items = []panels.Item{{Kind: panels.ItemKindPR, Number: 1, Title: "x"}}
		if err := g.loadDetail(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// --- refreshDetailPreview ---

func TestRefreshDetailPreview_SetsContent(t *testing.T) {
	g := newTestGui()
	g.panels.Items.Items = []panels.Item{
		{Kind: panels.ItemKindPR, Number: 42, Title: "My PR"},
		{Kind: panels.ItemKindIssue, Number: 7, Title: "My Issue"},
	}
	g.panels.Items.Selected = 0
	g.refreshDetailPreview()
	want := g.panels.Items.Items[0].String()
	if g.panels.Detail.Content != want {
		t.Errorf("got %q, want %q", g.panels.Detail.Content, want)
	}
}

func TestRefreshDetailPreview_EmptyItems(t *testing.T) {
	g := newTestGui()
	// panic しないことを確認
	g.refreshDetailPreview()
}
