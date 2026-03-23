package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/pr/diff"
	"github.com/rin2yh/lazygh/internal/pr/review"
	"github.com/rin2yh/lazygh/pkg/gui/viewport"
)

type Gui struct {
	config *config.Config
	coord  *Coordinator
	client PRClient

	focus    layout.Focus
	showHelp bool

	diff   diff.Selection
	detail DetailViewport

	review ReviewController
}

func NewGui(cfg *config.Config, coord *Coordinator, prClient PRClient, reviewClient review.PendingReviewClient) (*Gui, error) {
	vp := viewport.New()
	g := &Gui{
		config: cfg,
		coord:  coord,
		client: prClient,
		focus:  layout.FocusPRs,
		detail: &vp,
	}
	revCtrl := review.NewController(cfg, coord, reviewClient, &g.diff, g.setReviewFocus)
	g.review = revCtrl
	coord.SetReviewHook(revCtrl)
	return g, nil
}

func (gui *Gui) Run() error {
	p := tea.NewProgram(&screen{gui: gui}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
