package gui

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestNavigatePRList(t *testing.T) {
	g, err := NewGui(config.Default(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []core.Item{testfactory.CoreItem(1, "a"), testfactory.CoreItem(2, "b")}, nil)

	g.navigateDown()
	if g.state.List.PRsSelected != 1 {
		t.Fatalf("got %d, want %d", g.state.List.PRsSelected, 1)
	}

	g.navigateUp()
	if g.state.List.PRsSelected != 0 {
		t.Fatalf("got %d, want %d", g.state.List.PRsSelected, 0)
	}
}
