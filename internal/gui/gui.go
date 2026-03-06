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
	PanelIssues
	PanelPRs
	PanelDetail
	panelCount
)

type State struct {
	ActivePanel PanelType
}

type Panels struct {
	Repos  *panels.ItemsPanel
	Issues *panels.ItemsPanel
	PRs    *panels.ItemsPanel
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
			Repos:  panels.NewItemsPanel(panels.FormatRepoItem, true),
			Issues: panels.NewItemsPanel(panels.FormatIssueItem, false),
			PRs:    panels.NewItemsPanel(panels.FormatPRItem, false),
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
		if gui.client == nil {
			return nil
		}
		gui.panels.Repos.Loading = true
		gui.renderPanel("repos")
		go func() {
			repos, err := gui.client.ListRepos()
			gui.g.Execute(func(_ *gocui.Gui) error {
				return gui.applyReposResult(repos, err)
			})
		}()
		return nil
	})
	return gui.g.MainLoop()
}

var panelViewNames = []string{"repos", "issues", "prs", "detail"}

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
		gui.panels.Repos.Render(v, gui.state.ActivePanel == PanelRepos)
	case "issues":
		gui.panels.Issues.Render(v, gui.state.ActivePanel == PanelIssues)
	case "prs":
		gui.panels.PRs.Render(v, gui.state.ActivePanel == PanelPRs)
	case "detail":
		gui.panels.Detail.Render(v)
	}
}

func (gui *Gui) applyReposResult(repos []string, err error) error {
	gui.panels.Repos.Loading = false
	if err != nil {
		gui.showError("Error loading repos", err)
		return nil
	}
	gui.setItemsPanel(gui.panels.Repos, "repos", toRepoItems(repos))
	gui.reposLoaded = true
	return nil
}

func (gui *Gui) loadRepos() error {
	if gui.client == nil {
		return nil
	}
	repos, err := gui.client.ListRepos()
	return gui.applyReposResult(repos, err)
}

func (gui *Gui) loadItems() error {
	if gui.client == nil {
		return nil
	}
	repo, ok := gui.selectedRepo()
	if !ok {
		return nil
	}

	issues, err := gui.client.ListIssues(repo)
	if err != nil {
		gui.showError("Error loading issues", err)
		return nil
	}

	prs, err := gui.client.ListPRs(repo)
	if err != nil {
		gui.showError("Error loading PRs", err)
		return nil
	}

	issueItems := make([]panels.Item, 0, len(issues))
	for _, issue := range issues {
		issueItems = append(issueItems, panels.Item{Number: issue.Number, Title: issue.Title})
	}
	prItems := make([]panels.Item, 0, len(prs))
	for _, pr := range prs {
		prItems = append(prItems, panels.Item{Number: pr.Number, Title: pr.Title})
	}

	gui.setItemsPanel(gui.panels.Issues, "issues", issueItems)
	gui.setItemsPanel(gui.panels.PRs, "prs", prItems)

	gui.panels.Detail.SetContent("")
	gui.renderPanel("detail")
	return nil
}

func (gui *Gui) setItemsPanel(panel *panels.ItemsPanel, viewName string, items []panels.Item) {
	panel.Items = items
	panel.Selected = 0
	gui.renderPanel(viewName)
}

type detailLoader func(repo string, number int) (string, error)

func (gui *Gui) activeItemsPanel() (*panels.ItemsPanel, bool) {
	switch gui.state.ActivePanel {
	case PanelIssues:
		return gui.panels.Issues, true
	case PanelPRs:
		return gui.panels.PRs, true
	default:
		return nil, false
	}
}

func (gui *Gui) activeDetailLoader() (detailLoader, bool) {
	switch gui.state.ActivePanel {
	case PanelIssues:
		return gui.client.ViewIssue, true
	case PanelPRs:
		return gui.client.ViewPR, true
	default:
		return nil, false
	}
}

func (gui *Gui) loadDetail() error {
	if gui.client == nil {
		return nil
	}
	repo, ok := gui.selectedRepo()
	if !ok {
		return nil
	}
	itemsPanel, ok := gui.activeItemsPanel()
	if !ok || len(itemsPanel.Items) == 0 {
		return nil
	}
	loader, ok := gui.activeDetailLoader()
	if !ok {
		return nil
	}

	item := itemsPanel.Items[itemsPanel.Selected]
	content, err := loader(repo, item.Number)
	if err != nil {
		gui.showError("Error loading detail", err)
		return nil
	}

	gui.panels.Detail.SetContent(content)
	gui.renderPanel("detail")
	return nil
}

func (gui *Gui) refreshDetailPreview() {
	itemsPanel, ok := gui.activeItemsPanel()
	if !ok || len(itemsPanel.Items) == 0 {
		return
	}
	item := itemsPanel.Items[itemsPanel.Selected]
	gui.panels.Detail.SetContent(itemsPanel.Format(item))
	gui.renderPanel("detail")
}

func toRepoItems(repos []string) []panels.Item {
	items := make([]panels.Item, 0, len(repos))
	for _, repo := range repos {
		items = append(items, panels.Item{Title: repo})
	}
	return items
}

func (gui *Gui) selectedRepo() (string, bool) {
	if len(gui.panels.Repos.Items) == 0 {
		return "", false
	}
	return gui.panels.Repos.Format(gui.panels.Repos.Items[gui.panels.Repos.Selected]), true
}
