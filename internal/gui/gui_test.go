package gui

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gui/panels"
)

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
	g := newTestGui()

	g.state.ActivePanel = PanelRepos
	_ = g.nextPanel(nil, nil)
	if g.state.ActivePanel != PanelItems {
		t.Errorf("after Repos: got %d, want PanelItems", g.state.ActivePanel)
	}

	_ = g.nextPanel(nil, nil)
	if g.state.ActivePanel != PanelDetail {
		t.Errorf("after Items: got %d, want PanelDetail", g.state.ActivePanel)
	}

	_ = g.nextPanel(nil, nil)
	if g.state.ActivePanel != PanelRepos {
		t.Errorf("after Detail: got %d, want PanelRepos", g.state.ActivePanel)
	}
}

func TestNavigateDown_Repos(t *testing.T) {
	twoItems := []panels.Item{
		{Kind: panels.ItemKindPR, Number: 1, Title: "a"},
		{Kind: panels.ItemKindPR, Number: 2, Title: "b"},
	}
	tests := []struct {
		name        string
		active      PanelType
		repos       []string
		repoSel     int
		wantRepoSel int
	}{
		{"Empty", PanelRepos, nil, 0, 0},
		{"UpperBound", PanelRepos, []string{"a", "b", "c"}, 2, 2},
		{"Normal", PanelRepos, []string{"a", "b", "c"}, 0, 1},
		{"InactivePanel", PanelItems, []string{"a", "b"}, 0, 0},
	}
	for _, tt := range tests {
		g := newTestGui()
		g.state.ActivePanel = tt.active
		g.panels.Repos.Repos = tt.repos
		g.panels.Repos.Selected = tt.repoSel
		_ = g.navigateDown(nil, "repos")
		if g.panels.Repos.Selected != tt.wantRepoSel {
			t.Errorf("%s: got %d, want %d", tt.name, g.panels.Repos.Selected, tt.wantRepoSel)
		}
	}

	// items panel navigate down
	g := newTestGui()
	g.state.ActivePanel = PanelItems
	g.panels.Items.Items = twoItems
	g.panels.Items.Selected = 0
	_ = g.navigateDown(nil, "items")
	if g.panels.Items.Selected != 1 {
		t.Errorf("Items navigateDown: got %d, want 1", g.panels.Items.Selected)
	}
}

func TestNavigateUp_Repos(t *testing.T) {
	twoItems := []panels.Item{
		{Kind: panels.ItemKindPR, Number: 1, Title: "a"},
		{Kind: panels.ItemKindPR, Number: 2, Title: "b"},
	}
	tests := []struct {
		name        string
		repos       []string
		repoSel     int
		wantRepoSel int
	}{
		{"LowerBound", []string{"a", "b"}, 0, 0},
		{"Normal", []string{"a", "b", "c"}, 2, 1},
	}
	for _, tt := range tests {
		g := newTestGui()
		g.state.ActivePanel = PanelRepos
		g.panels.Repos.Repos = tt.repos
		g.panels.Repos.Selected = tt.repoSel
		_ = g.navigateUp(nil, "repos")
		if g.panels.Repos.Selected != tt.wantRepoSel {
			t.Errorf("%s: got %d, want %d", tt.name, g.panels.Repos.Selected, tt.wantRepoSel)
		}
	}

	// items panel navigate up
	g := newTestGui()
	g.state.ActivePanel = PanelItems
	g.panels.Items.Items = twoItems
	g.panels.Items.Selected = 1
	_ = g.navigateUp(nil, "items")
	if g.panels.Items.Selected != 0 {
		t.Errorf("Items navigateUp: got %d, want 0", g.panels.Items.Selected)
	}
}
