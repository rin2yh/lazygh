package gui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

type Gui struct {
	config *config.Config
	state  *core.State
	client gh.ClientInterface

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
	content string
	err     error
}

type model struct {
	gui *Gui
}

func (m *model) Init() tea.Cmd {
	if m.gui.client == nil {
		return nil
	}
	m.gui.state.BeginLoadPRs()
	return m.loadPRsCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.gui.state.SetWindowSize(msg.Width, msg.Height)
		return m, nil
	case prsLoadedMsg:
		m.gui.applyPRsResult(msg)
		return m, nil
	case detailLoadedMsg:
		m.gui.applyDetailResult(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			m.gui.navigateDown()
			return m, nil
		case "k", "up":
			m.gui.navigateUp()
			return m, nil
		case "enter":
			return m, m.handleEnter()
		}
	}
	return m, nil
}

func (m *model) View() string {
	return m.gui.render()
}

func toCorePRs(prs []gh.PRItem) []core.Item {
	items := make([]core.Item, 0, len(prs))
	for _, pr := range prs {
		items = append(items, core.Item{Number: pr.Number, Title: pr.Title})
	}
	return items
}

func (m *model) loadPRsCmd() tea.Cmd {
	return func() tea.Msg {
		repo, err := m.gui.client.ResolveCurrentRepo()
		if err != nil {
			return prsLoadedMsg{err: err}
		}
		prs, err := m.gui.client.ListPRs(repo)
		if err != nil {
			return prsLoadedMsg{repo: repo, err: err}
		}
		return prsLoadedMsg{repo: repo, prs: toCorePRs(prs)}
	}
}

func (m *model) loadDetailCmd(repo string, number int) tea.Cmd {
	return func() tea.Msg {
		content, err := m.gui.client.ViewPR(repo, number)
		return detailLoadedMsg{content: content, err: err}
	}
}

func (m *model) handleEnter() tea.Cmd {
	action := m.gui.state.PlanEnter(m.gui.client != nil, os.Getenv("LAZYGH_DEBUG_DETAIL_TEXT"))
	switch action.Kind {
	case core.EnterLoadPRDetail:
		return m.loadDetailCmd(action.Repo, action.Number)
	default:
		return nil
	}
}

func (gui *Gui) applyPRsResult(msg prsLoadedMsg) {
	gui.state.ApplyPRsResult(msg.repo, msg.prs, msg.err)
}

func (gui *Gui) applyDetailResult(msg detailLoadedMsg) {
	gui.state.ApplyDetailResult(msg.content, msg.err)
}

func (gui *Gui) navigateDown() {
	gui.state.NavigateDown()
}

func (gui *Gui) navigateUp() {
	gui.state.NavigateUp()
}

func (gui *Gui) render() string {
	w := gui.state.Width
	h := gui.state.Height
	if w <= 0 {
		w = 120
	}
	if h <= 0 {
		h = 40
	}

	leftWidth := w * 35 / 100
	if leftWidth < 1 {
		leftWidth = 1
	}
	if leftWidth > w-2 {
		leftWidth = w - 2
	}
	rightWidth := w - leftWidth - 1
	if rightWidth < 1 {
		rightWidth = 1
	}

	contentHeight := h - 1
	if contentHeight < 1 {
		contentHeight = 1
	}

	leftLines := gui.renderPRPanel("PRs (Open/Draft)", contentHeight)
	rightLines := gui.renderDetailPanel("Detail", rightWidth, contentHeight)

	var b strings.Builder
	for i := 0; i < contentHeight; i++ {
		left := ""
		if i < len(leftLines) {
			left = leftLines[i]
		}
		right := ""
		if i < len(rightLines) {
			right = rightLines[i]
		}
		b.WriteString(padOrTrim(left, leftWidth))
		b.WriteRune('│')
		b.WriteString(padOrTrim(right, rightWidth))
		b.WriteByte('\n')
	}
	b.WriteString(padOrTrim(formatStatusLine(gui.state.Repo), w))
	return b.String()
}

func (gui *Gui) renderPRPanel(title string, height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)
	lines = append(lines, formatPanelTitle(title, true))

	if gui.state.PRsLoading {
		for len(lines) < height {
			if len(lines) == 1 {
				lines = append(lines, "Loading...")
			} else {
				lines = append(lines, "")
			}
		}
		return lines
	}

	if len(gui.state.PRs) == 0 {
		for len(lines) < height {
			if len(lines) == 1 {
				lines = append(lines, "No pull requests")
			} else {
				lines = append(lines, "")
			}
		}
		return lines
	}

	for i := 0; len(lines) < height; i++ {
		if i >= len(gui.state.PRs) {
			lines = append(lines, "")
			continue
		}
		prefix := "  "
		if i == gui.state.PRsSelected {
			prefix = "> "
		}
		lines = append(lines, prefix+core.FormatPRItem(gui.state.PRs[i]))
	}
	return lines
}

func (gui *Gui) renderDetailPanel(title string, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	bodyHeight := height - 1
	if bodyHeight < 1 {
		bodyHeight = 1
	}
	gui.syncDetailViewport(width, bodyHeight, gui.state.DetailContent)

	lines := make([]string, 0, height)
	lines = append(lines, formatPanelTitle(title, false))
	for _, line := range strings.Split(gui.detailViewport.View(), "\n") {
		if len(lines) >= height {
			break
		}
		lines = append(lines, line)
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return lines
}

func (gui *Gui) syncDetailViewport(width int, height int, content string) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if gui.detailViewportWidth != width || gui.detailViewportHeight != height {
		gui.detailViewport.Width = width
		gui.detailViewport.Height = height
		gui.detailViewportWidth = width
		gui.detailViewportHeight = height
	}
	wrapped := wrapText(content, width)
	if gui.detailViewportBody != wrapped {
		gui.detailViewport.SetContent(wrapped)
		gui.detailViewportBody = wrapped
		gui.detailViewport.GotoTop()
	}
}

func wrapText(content string, width int) string {
	if width <= 0 || content == "" {
		return content
	}

	srcLines := strings.Split(content, "\n")
	dstLines := make([]string, 0, len(srcLines))
	for _, line := range srcLines {
		var b strings.Builder
		col := 0
		for _, r := range line {
			rw := runewidth.RuneWidth(r)
			if rw <= 0 {
				rw = 1
			}
			if col > 0 && col+rw > width {
				dstLines = append(dstLines, b.String())
				b.Reset()
				col = 0
			}
			b.WriteRune(r)
			col += rw
		}
		dstLines = append(dstLines, b.String())
	}
	return strings.Join(dstLines, "\n")
}

func padOrTrim(s string, width int) string {
	if width <= 0 {
		return ""
	}
	var b strings.Builder
	col := 0
	for _, r := range s {
		rw := runewidth.RuneWidth(r)
		if rw <= 0 {
			rw = 1
		}
		if col+rw > width {
			break
		}
		b.WriteRune(r)
		col += rw
	}
	if col < width {
		b.WriteString(strings.Repeat(" ", width-col))
	}
	return b.String()
}
