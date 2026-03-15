package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/detail"
	guidiff "github.com/rin2yh/lazygh/internal/gui/diff"
	guireview "github.com/rin2yh/lazygh/internal/gui/review"
)

type PRClient interface {
	ResolveCurrentRepo() (string, error)
	ListPRs(repo string) ([]gh.PRItem, error)
	ViewPR(repo string, number int) (string, error)
	DiffPR(repo string, number int) (string, error)
}

const (
	ansiReset   = "\x1b[0m"
	ansiReverse = "\x1b[7m"
	ansiGreen   = "\x1b[32m"
	ansiRed     = "\x1b[31m"
	ansiYellow  = "\x1b[33m"
	ansiBlue    = "\x1b[34m"
	ansiCyan    = "\x1b[36m"
	ansiPurple  = "\x1b[35m"
	ansiGray    = "\x1b[90m"
)

type Gui struct {
	config *config.Config
	state  *core.State
	client PRClient

	focus panelFocus

	diff   guidiff.Selection
	detail detail.State

	review *guireview.Controller
}

func NewGui(cfg *config.Config, prClient PRClient, reviewClient guireview.PendingReviewClient) (*Gui, error) {
	gui := &Gui{
		config: cfg,
		state:  core.NewState(),
		client: prClient,
		focus:  panelPRs,
		detail: detail.NewState(),
	}
	gui.review = guireview.NewController(cfg, gui.state, reviewClient, &gui.diff, gui.setReviewFocus)
	return gui, nil
}

func (gui *Gui) Run() error {
	p := tea.NewProgram(&screen{gui: gui}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
