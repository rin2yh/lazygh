package gui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

type panelFocus int

const (
	panelPRs panelFocus = iota
	panelDiffFiles
	panelDiffContent
)

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
	client gh.ClientInterface

	focus panelFocus

	diffFiles        []gh.DiffFile
	diffFileSelected int

	detailViewport       viewport.Model
	detailViewportWidth  int
	detailViewportHeight int
	detailViewportBody   string
}

func NewGui(cfg *config.Config, client gh.ClientInterface) (*Gui, error) {
	vp := viewport.New(1, 1)
	return &Gui{
		config:               cfg,
		state:                core.NewState(),
		client:               client,
		focus:                panelPRs,
		detailViewport:       vp,
		detailViewportWidth:  1,
		detailViewportHeight: 1,
	}, nil
}

func (gui *Gui) Run() error {
	p := tea.NewProgram(&model{gui: gui}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type prsLoadedMsg struct {
	repo string
	prs  []core.Item
	err  error
}

type detailLoadedMsg struct {
	mode    core.DetailMode
	number  int
	content string
	err     error
}

func (gui *Gui) applyPRsResult(msg prsLoadedMsg) {
	gui.state.ApplyPRsResult(msg.repo, msg.prs, msg.err)
	gui.focus = panelPRs
}

func (gui *Gui) applyDetailResult(msg detailLoadedMsg) {
	if !gui.state.ShouldApplyDetailResult(msg.mode, msg.number) {
		return
	}
	if msg.mode == core.DetailModeDiff {
		gui.state.ApplyDiffResult(msg.content, msg.err)
		if msg.err != nil {
			gui.diffFiles = nil
			gui.diffFileSelected = 0
			if gui.focus == panelDiffFiles {
				gui.focus = panelDiffContent
			}
			return
		}
		gui.updateDiffFiles(gui.state.DetailContent)
		return
	}
	gui.state.ApplyDetailResult(msg.content, msg.err)
}

func (gui *Gui) switchToOverview() bool {
	changed := gui.state.SwitchToOverview()
	if changed {
		gui.focus = panelPRs
	}
	return changed
}

func (gui *Gui) focusPRs() {
	gui.focus = panelPRs
}

func (gui *Gui) switchToDiff() bool {
	changed := gui.state.SwitchToDiff()
	if changed {
		gui.focus = panelDiffFiles
		gui.diffFiles = nil
		gui.diffFileSelected = 0
	}
	return changed
}

func (gui *Gui) cycleFocus() {
	if !gui.state.IsDiffMode() {
		gui.focus = panelPRs
		return
	}

	order := gui.focusOrder()
	if len(order) == 0 {
		gui.focus = panelPRs
		return
	}
	for i, focus := range order {
		if focus == gui.focus {
			gui.focus = order[(i+1)%len(order)]
			return
		}
	}
	gui.focus = order[0]
}

func (gui *Gui) focusOrder() []panelFocus {
	order := []panelFocus{panelPRs}
	if len(gui.diffFiles) > 0 {
		order = append(order, panelDiffFiles)
	}
	order = append(order, panelDiffContent)
	return order
}

func (gui *Gui) currentDiffContent() string {
	if len(gui.diffFiles) == 0 {
		return gui.state.DetailContent
	}
	if gui.diffFileSelected < 0 || gui.diffFileSelected >= len(gui.diffFiles) {
		return gui.state.DetailContent
	}
	return gui.diffFiles[gui.diffFileSelected].Content
}

func (gui *Gui) updateDiffFiles(content string) {
	files := gh.ParseUnifiedDiff(content)
	if len(files) == 0 {
		gui.diffFiles = nil
		gui.diffFileSelected = 0
		if gui.focus == panelDiffFiles {
			gui.focus = panelDiffContent
		}
		return
	}

	prevPath := ""
	if gui.diffFileSelected >= 0 && gui.diffFileSelected < len(gui.diffFiles) {
		prevPath = gui.diffFiles[gui.diffFileSelected].Path
	}

	gui.diffFiles = files
	gui.diffFileSelected = 0
	if prevPath != "" {
		for i, file := range files {
			if file.Path == prevPath {
				gui.diffFileSelected = i
				break
			}
		}
	}
}
