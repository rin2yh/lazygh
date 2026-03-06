package main

import (
	"fmt"
	"os"

	"github.com/rin2yh/lazygh/internal/app"
	"github.com/rin2yh/lazygh/internal/config"
)

func main() {
	cfg := config.Default()
	a, err := app.NewApp(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if err := a.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
