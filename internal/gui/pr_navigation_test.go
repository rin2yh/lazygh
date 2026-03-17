package gui

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/model"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestNavigatePRList(t *testing.T) {
	g := mustNewGui(t, &testmock.GHClient{})
	g.state.ApplyPRsResult("owner/repo", []model.Item{testfactory.NewItem(1, "a"), testfactory.NewItem(2, "b")}, nil)

	g.navigateDown()
	if g.state.PRsSelected != 1 {
		t.Fatalf("got %d, want %d", g.state.PRsSelected, 1)
	}

	g.navigateUp()
	if g.state.PRsSelected != 0 {
		t.Fatalf("got %d, want %d", g.state.PRsSelected, 0)
	}
}
