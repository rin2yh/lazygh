package gui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
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

	diffFiles        []gh.DiffFile
	diffFileSelected int
	diffLineSelected int

	detailViewport       viewport.Model
	detailViewportWidth  int
	detailViewportHeight int
	detailViewportBody   string

	review *guireview.Controller
}

func NewGui(cfg *config.Config, prClient PRClient, reviewClient guireview.PendingReviewClient) (*Gui, error) {
	vp := viewport.New(1, 1)
	gui := &Gui{
		config:               cfg,
		state:                core.NewState(),
		client:               prClient,
		focus:                panelPRs,
		detailViewport:       vp,
		detailViewportWidth:  1,
		detailViewportHeight: 1,
	}
	gui.review = guireview.NewController(gui.state, reviewClient, gui, gui.setReviewFocus)
	return gui, nil
}

func (gui *Gui) Run() error {
	p := tea.NewProgram(&screen{gui: gui}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
