package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/panels"
)

type PanelType int

const (
	PanelRepos PanelType = iota
	PanelItems
	PanelDetail
	panelCount
)

type State struct {
	ActivePanel PanelType
}

type Panels struct {
	Repos  *panels.ReposPanel
	Items  *panels.ItemsPanel
	Detail *panels.DetailPanel
}

type Gui struct {
	g           *gocui.Gui
	config      *config.Config
	state       *State
	panels      *Panels
	client      gh.ClientInterface
	reposLoaded bool
}

func NewGui(cfg *config.Config, client gh.ClientInterface) (*Gui, error) {
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		return nil, err
	}

	gui := &Gui{
		g:      g,
		config: cfg,
		state:  &State{ActivePanel: PanelRepos},
		panels: &Panels{
			Repos:  panels.NewReposPanel(),
			Items:  panels.NewItemsPanel(),
			Detail: panels.NewDetailPanel(),
		},
		client: client,
	}

	g.SetLayout(gui.layout)
	g.Mouse = false

	if err := gui.setupKeybindings(); err != nil {
		g.Close()
		return nil, err
	}

	return gui, nil
}

func (gui *Gui) Run() error {
	defer gui.g.Close()
	gui.g.Execute(func(_ *gocui.Gui) error {
		return gui.loadRepos()
	})
	return gui.g.MainLoop()
}

var panelViewNames = []string{"repos", "items", "detail"}

func (gui *Gui) activeViewName() string {
	return panelViewNames[gui.state.ActivePanel]
}

func (gui *Gui) showError(msg string, err error) {
	gui.panels.Detail.SetContent(fmt.Sprintf("%s: %v", msg, err))
	gui.renderPanel("detail")
}

func (gui *Gui) renderPanel(name string) {
	if gui.g == nil {
		return
	}
	v, err := gui.g.View(name)
	if err != nil {
		return
	}
	switch name {
	case "repos":
		gui.panels.Repos.Render(v)
	case "items":
		gui.panels.Items.Render(v)
	case "detail":
		gui.panels.Detail.Render(v)
	}
}

func (gui *Gui) loadRepos() error {
	if gui.client == nil {
		return nil
	}
	repos, err := gui.client.ListRepos()
	if err != nil {
		gui.showError("Error loading repos", err)
		return nil
	}
	gui.panels.Repos.Repos = repos
	gui.panels.Repos.Selected = 0
	gui.reposLoaded = true
	gui.renderPanel("repos")
	return nil
}

func (gui *Gui) loadItems() error {
	if gui.client == nil {
		return nil
	}
	if len(gui.panels.Repos.Repos) == 0 {
		return nil
	}
	repo := gui.panels.Repos.Repos[gui.panels.Repos.Selected]

	prs, err := gui.client.ListPRs(repo)
	if err != nil {
		gui.showError("Error loading PRs", err)
		return nil
	}

	issues, err := gui.client.ListIssues(repo)
	if err != nil {
		gui.showError("Error loading issues", err)
		return nil
	}

	items := make([]panels.Item, 0, len(prs)+len(issues))
	for _, pr := range prs {
		items = append(items, panels.Item{Kind: panels.ItemKindPR, Number: pr.Number, Title: pr.Title})
	}
	for _, issue := range issues {
		items = append(items, panels.Item{Kind: panels.ItemKindIssue, Number: issue.Number, Title: issue.Title})
	}

	gui.panels.Items.Items = items
	gui.panels.Items.Selected = 0
	gui.renderPanel("items")

	gui.panels.Detail.SetContent("")
	gui.renderPanel("detail")
	return nil
}

func (gui *Gui) loadDetail() error {
	if gui.client == nil {
		return nil
	}
	if len(gui.panels.Items.Items) == 0 {
		return nil
	}
	if len(gui.panels.Repos.Repos) == 0 {
		return nil
	}
	repo := gui.panels.Repos.Repos[gui.panels.Repos.Selected]
	item := gui.panels.Items.Items[gui.panels.Items.Selected]

	var content string
	var err error
	if item.Kind == panels.ItemKindPR {
		content, err = gui.client.ViewPR(repo, item.Number)
	} else {
		content, err = gui.client.ViewIssue(repo, item.Number)
	}
	if err != nil {
		gui.showError("Error loading detail", err)
		return nil
	}

	gui.panels.Detail.SetContent(content)
	gui.renderPanel("detail")
	return nil
}

func (gui *Gui) refreshDetailPreview() {
	if len(gui.panels.Items.Items) == 0 {
		return
	}
	item := gui.panels.Items.Items[gui.panels.Items.Selected]
	gui.panels.Detail.SetContent(item.String())
	gui.renderPanel("detail")
}
