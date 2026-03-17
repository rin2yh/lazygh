package app

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/review"
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

// reviewCtrl はテスト用に Gui の ReviewController を具体型にキャストして返す。
// インターフェースに含めるほどでない初期化メソッド（SetContext 等）へのアクセスに使用する。
func reviewCtrl(g *Gui) *review.Controller {
	return g.review.(*review.Controller)
}
