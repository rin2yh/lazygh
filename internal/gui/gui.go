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

type PanelType = core.PanelType

const (
	PanelRepos  = core.PanelRepos
	PanelIssues = core.PanelIssues
	PanelPRs    = core.PanelPRs
	PanelDetail = core.PanelDetail
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

type reposLoadedMsg struct {
	repos []string
	err   error
}

type itemsLoadedMsg struct {
	repo   string
	issues []core.Item
	prs    []core.Item
	err    error
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
	m.gui.state.BeginLoadRepos()
	return m.loadReposCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.gui.state.SetWindowSize(msg.Width, msg.Height)
		return m, nil
	case reposLoadedMsg:
		m.gui.applyReposResult(msg.repos, msg.err)
		return m, nil
	case itemsLoadedMsg:
		m.gui.applyItemsResult(msg)
		return m, nil
	case detailLoadedMsg:
		m.gui.applyDetailResult(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.gui.nextPanel()
			return m, nil
		case "shift+tab":
			m.gui.prevPanel()
			return m, nil
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

func (m *model) loadReposCmd() tea.Cmd {
	return func() tea.Msg {
		repos, err := m.gui.client.ListRepos()
		return reposLoadedMsg{repos: repos, err: err}
	}
}

func toCoreIssues(issues []gh.IssueItem) []core.Item {
	items := make([]core.Item, 0, len(issues))
	for _, issue := range issues {
		items = append(items, core.Item{Number: issue.Number, Title: issue.Title})
	}
	return items
}

func toCorePRs(prs []gh.PRItem) []core.Item {
	items := make([]core.Item, 0, len(prs))
	for _, pr := range prs {
		items = append(items, core.Item{Number: pr.Number, Title: pr.Title})
	}
	return items
}

func (m *model) loadItemsCmd(repo string) tea.Cmd {
	return func() tea.Msg {
		issues, err := m.gui.client.ListIssues(repo)
		if err != nil {
			return itemsLoadedMsg{repo: repo, err: err}
		}
		prs, err := m.gui.client.ListPRs(repo)
		if err != nil {
			return itemsLoadedMsg{repo: repo, err: err}
		}
		return itemsLoadedMsg{repo: repo, issues: toCoreIssues(issues), prs: toCorePRs(prs)}
	}
}

type detailLoader func(repo string, number int) (string, error)

func (m *model) loadDetailCmd(repo string, number int, loader detailLoader) tea.Cmd {
	return func() tea.Msg {
		content, err := loader(repo, number)
		return detailLoadedMsg{content: content, err: err}
	}
}

func (m *model) handleEnter() tea.Cmd {
	action := m.gui.state.PlanEnter(m.gui.client != nil, os.Getenv("LAZYGH_DEBUG_DETAIL_TEXT"))
	switch action.Kind {
	case core.EnterLoadItems:
		return m.loadItemsCmd(action.Repo)
	case core.EnterLoadIssueDetail:
		return m.loadDetailCmd(action.Repo, action.Number, m.gui.client.ViewIssue)
	case core.EnterLoadPRDetail:
		return m.loadDetailCmd(action.Repo, action.Number, m.gui.client.ViewPR)
	default:
		return nil
	}
}

func (gui *Gui) applyReposResult(repos []string, err error) {
	gui.state.ApplyReposResult(repos, err)
}

func (gui *Gui) applyItemsResult(msg itemsLoadedMsg) {
	gui.state.ApplyItemsResult(msg.repo, msg.issues, msg.prs, msg.err)
}

func (gui *Gui) applyDetailResult(msg detailLoadedMsg) {
	gui.state.ApplyDetailResult(msg.content, msg.err)
}

func (gui *Gui) nextPanel() {
	gui.state.NextPanel()
}

func (gui *Gui) prevPanel() {
	gui.state.PrevPanel()
}

func (gui *Gui) navigateDown() {
	if gui.state.ActivePanel == PanelDetail {
		gui.detailViewport.LineDown(1)
		return
	}
	gui.state.NavigateDown()
}

func (gui *Gui) navigateUp() {
	if gui.state.ActivePanel == PanelDetail {
		gui.detailViewport.LineUp(1)
		return
	}
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

	leftWidth := w * 30 / 100
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

	reposH := contentHeight / 3
	issuesH := contentHeight / 3
	prsH := contentHeight - reposH - issuesH

	leftLines := make([]string, 0, contentHeight)
	leftLines = append(leftLines, gui.renderItemsPanel("Repositories", gui.state.Repos, gui.state.ReposSelected, gui.state.ReposLoading, core.FormatRepoItem, gui.state.ActivePanel == PanelRepos, reposH)...)
	leftLines = append(leftLines, gui.renderItemsPanel("Issues", gui.state.Issues, gui.state.IssuesSelected, gui.state.IssuesLoading, core.FormatIssueItem, gui.state.ActivePanel == PanelIssues, issuesH)...)
	leftLines = append(leftLines, gui.renderItemsPanel("PRs", gui.state.PRs, gui.state.PRsSelected, gui.state.PRsLoading, core.FormatPRItem, gui.state.ActivePanel == PanelPRs, prsH)...)

	rightLines := gui.renderDetailPanel("Detail", gui.state.ActivePanel == PanelDetail, rightWidth, contentHeight)

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
	b.WriteString(padOrTrim(formatStatusLine(gui.state.ActivePanel), w))
	return b.String()
}

func (gui *Gui) renderItemsPanel(title string, items []core.Item, selected int, loading bool, formatter func(core.Item) string, active bool, height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)
	lines = append(lines, formatPanelTitle(title, active))
	if loading {
		for len(lines) < height {
			if len(lines) == 1 {
				lines = append(lines, "Loading...")
			} else {
				lines = append(lines, "")
			}
		}
		return lines
	}

	for i := 0; len(lines) < height; i++ {
		if i >= len(items) {
			lines = append(lines, "")
			continue
		}
		prefix := "  "
		if i == selected {
			prefix = "> "
		}
		lines = append(lines, prefix+formatter(items[i]))
	}
	return lines
}

func (gui *Gui) renderDetailPanel(title string, active bool, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	bodyHeight := height - 1
	if bodyHeight < 1 {
		bodyHeight = 1
	}
	gui.syncDetailViewport(width, bodyHeight, gui.state.DetailContent)

	lines := make([]string, 0, height)
	lines = append(lines, formatPanelTitle(title, active))
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
