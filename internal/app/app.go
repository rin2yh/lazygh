package app

import (
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
)

type App struct {
	Config      *config.Config
	Coordinator *Coordinator
}

func NewApp(cfg *config.Config) (*App, error) {
	if err := gh.ValidateCLI(); err != nil {
		return nil, err
	}
	return &App{Config: cfg, Coordinator: NewCoordinator()}, nil
}
