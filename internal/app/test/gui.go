package apptest

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/app"
	"github.com/rin2yh/lazygh/internal/config"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

// NewGui creates a *app.Gui for testing and fails the test on error.
func NewGui(t *testing.T, client *testmock.GHClient) *app.Gui {
	t.Helper()
	g, err := app.NewGui(config.Default(), app.NewCoordinator(), client, client)
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	return g
}
