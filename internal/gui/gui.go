package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/rin2yh/lazygh/internal/config"
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
	g      *gocui.Gui
	config *config.Config
	state  *State
	panels *Panels
}

func NewGui(cfg *config.Config) (*Gui, error) {
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
	return gui.g.MainLoop()
}

var panelViewNames = []string{"repos", "items", "detail"}

func (gui *Gui) activeViewName() string {
	return panelViewNames[gui.state.ActivePanel]
}
