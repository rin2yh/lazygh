package app

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func mustNewGui(t *testing.T, client *testmock.GHClient) *Gui {
	t.Helper()
	g, err := NewGui(config.Default(), NewCoordinator(), client, client)
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	return g
}
