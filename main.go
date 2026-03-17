package main

import (
	"fmt"
	"os"

	"github.com/rin2yh/lazygh/internal/app"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
)

func main() {
	cfg := config.Default()
	a, err := app.NewApp(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	client := gh.NewClient()
	g, err := app.NewGui(cfg, a.Coordinator, client, client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if err := g.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
