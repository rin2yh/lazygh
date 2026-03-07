package app

import (
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui"
)

type App struct {
	Config *config.Config
	Gui    *gui.Gui
}

func NewApp(cfg *config.Config) (*App, error) {
	if err := gh.ValidateCLI(); err != nil {
		return nil, err
	}
	client := gh.NewClient()
	g, err := gui.NewGui(cfg, client)
	if err != nil {
		return nil, err
	}
	return &App{Config: cfg, Gui: g}, nil
}

func (a *App) Run() error {
	return a.Gui.Run()
}
