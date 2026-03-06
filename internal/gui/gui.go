package gui

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/panels"
)

type PanelType int

const (
	PanelRepos PanelType = iota
	PanelIssues
	PanelPRs
	PanelDetail
	panelCount
)

type State struct {
	ActivePanel PanelType
}

type Panels struct {
	Repos  *panels.ItemsPanel
	Issues *panels.ItemsPanel
	PRs    *panels.ItemsPanel
	Detail *panels.DetailPanel
}

type Gui struct {
	config      *config.Config
	state       *State
	panels      *Panels
	client      gh.ClientInterface
	reposLoaded bool
	width       int
	height      int
}

func NewGui(cfg *config.Config, client gh.ClientInterface) (*Gui, error) {
	return &Gui{
		config: cfg,
		state:  &State{ActivePanel: PanelRepos},
		panels: &Panels{
			Repos:  panels.NewItemsPanel(panels.FormatRepoItem, true),
			Issues: panels.NewItemsPanel(panels.FormatIssueItem, false),
			PRs:    panels.NewItemsPanel(panels.FormatPRItem, false),
			Detail: panels.NewDetailPanel(),
		},
		client: client,
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
	issues []gh.IssueItem
	prs    []gh.PRItem
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
	m.gui.panels.Repos.Loading = true
	return m.loadReposCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.gui.width = msg.Width
		m.gui.height = msg.Height
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
		return itemsLoadedMsg{repo: repo, issues: issues, prs: prs}
	}
}

func (m *model) loadDetailCmd(repo string, number int, loader detailLoader) tea.Cmd {
	return func() tea.Msg {
		content, err := loader(repo, number)
		return detailLoadedMsg{content: content, err: err}
	}
}

func (m *model) handleEnter() tea.Cmd {
	switch m.gui.state.ActivePanel {
	case PanelRepos:
		repo, ok := m.gui.selectedRepo()
		if !ok || m.gui.client == nil {
			return nil
		}
		m.gui.panels.Issues.Loading = true
		m.gui.panels.PRs.Loading = true
		m.gui.panels.Detail.SetContent("Loading items...")
		return m.loadItemsCmd(repo)
	case PanelIssues, PanelPRs:
		if m.gui.client == nil {
			return nil
		}
		repo, ok := m.gui.selectedRepo()
		if !ok {
			return nil
		}
		itemsPanel, ok := m.gui.activeItemsPanel()
		if !ok || len(itemsPanel.Items) == 0 {
			return nil
		}
		loader, ok := m.gui.activeDetailLoader()
		if !ok {
			return nil
		}
		if forced := os.Getenv("LAZYGH_DEBUG_DETAIL_TEXT"); forced != "" {
			m.gui.panels.Detail.SetContent(forced)
			return nil
		}
		item := itemsPanel.Items[itemsPanel.Selected]
		m.gui.panels.Detail.SetContent("Loading detail...")
		return m.loadDetailCmd(repo, item.Number, loader)
	default:
		return nil
	}
}

func (gui *Gui) showError(msg string, err error) {
	gui.panels.Detail.SetContent(fmt.Sprintf("%s: %v", msg, err))
}

func (gui *Gui) applyReposResult(repos []string, err error) {
	gui.panels.Repos.Loading = false
	if err != nil {
		gui.showError("Error loading repos", err)
		return
	}
	gui.panels.Repos.Items = toRepoItems(repos)
	gui.panels.Repos.Selected = 0
	gui.reposLoaded = true
}

func (gui *Gui) applyItemsResult(msg itemsLoadedMsg) {
	gui.panels.Issues.Loading = false
	gui.panels.PRs.Loading = false
	if msg.err != nil {
		gui.showError("Error loading items", msg.err)
		return
	}
	currentRepo, ok := gui.selectedRepo()
	if !ok || currentRepo != msg.repo {
		return
	}

	issueItems := make([]panels.Item, 0, len(msg.issues))
	for _, issue := range msg.issues {
		issueItems = append(issueItems, panels.Item{Number: issue.Number, Title: issue.Title})
	}
	prItems := make([]panels.Item, 0, len(msg.prs))
	for _, pr := range msg.prs {
		prItems = append(prItems, panels.Item{Number: pr.Number, Title: pr.Title})
	}
	gui.panels.Issues.Items = issueItems
	gui.panels.PRs.Items = prItems
	gui.panels.Issues.Selected = 0
	gui.panels.PRs.Selected = 0
	gui.panels.Detail.SetContent("")
}

func (gui *Gui) applyDetailResult(msg detailLoadedMsg) {
	if msg.err != nil {
		gui.showError("Error loading detail", msg.err)
		return
	}
	gui.panels.Detail.SetContent(msg.content)
}

func (gui *Gui) nextPanel() {
	gui.state.ActivePanel = (gui.state.ActivePanel + 1) % panelCount
}

func (gui *Gui) prevPanel() {
	if gui.state.ActivePanel == PanelRepos {
		gui.state.ActivePanel = PanelDetail
		return
	}
	gui.state.ActivePanel--
}

func (gui *Gui) navigateDown() {
	panelType, p, ok := gui.listPanelByPanel(gui.state.ActivePanel)
	if !ok || len(p.Items) == 0 {
		return
	}
	if p.Selected < len(p.Items)-1 {
		p.Selected++
	}
	if panelType != PanelRepos {
		gui.refreshDetailPreview()
	}
}

func (gui *Gui) navigateUp() {
	panelType, p, ok := gui.listPanelByPanel(gui.state.ActivePanel)
	if !ok || p.Selected <= 0 {
		return
	}
	p.Selected--
	if panelType != PanelRepos {
		gui.refreshDetailPreview()
	}
}

func (gui *Gui) listPanelByPanel(panel PanelType) (PanelType, *panels.ItemsPanel, bool) {
	switch panel {
	case PanelRepos:
		return PanelRepos, gui.panels.Repos, true
	case PanelIssues:
		return PanelIssues, gui.panels.Issues, true
	case PanelPRs:
		return PanelPRs, gui.panels.PRs, true
	default:
		return 0, nil, false
	}
}

type detailLoader func(repo string, number int) (string, error)

func (gui *Gui) activeItemsPanel() (*panels.ItemsPanel, bool) {
	switch gui.state.ActivePanel {
	case PanelIssues:
		return gui.panels.Issues, true
	case PanelPRs:
		return gui.panels.PRs, true
	default:
		return nil, false
	}
}

func (gui *Gui) activeDetailLoader() (detailLoader, bool) {
	switch gui.state.ActivePanel {
	case PanelIssues:
		return gui.client.ViewIssue, true
	case PanelPRs:
		return gui.client.ViewPR, true
	default:
		return nil, false
	}
}

func (gui *Gui) refreshDetailPreview() {
	itemsPanel, ok := gui.activeItemsPanel()
	if !ok || len(itemsPanel.Items) == 0 {
		return
	}
	item := itemsPanel.Items[itemsPanel.Selected]
	gui.panels.Detail.SetContent(itemsPanel.Format(item))
}

func toRepoItems(repos []string) []panels.Item {
	items := make([]panels.Item, 0, len(repos))
	for _, repo := range repos {
		items = append(items, panels.Item{Title: repo})
	}
	return items
}

func (gui *Gui) selectedRepo() (string, bool) {
	if len(gui.panels.Repos.Items) == 0 {
		return "", false
	}
	return gui.panels.Repos.Format(gui.panels.Repos.Items[gui.panels.Repos.Selected]), true
}

func (gui *Gui) render() string {
	w := gui.width
	h := gui.height
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
	leftLines = append(leftLines, gui.renderItemsPanel("Repositories", gui.panels.Repos, gui.state.ActivePanel == PanelRepos, reposH)...)
	leftLines = append(leftLines, gui.renderItemsPanel("Issues", gui.panels.Issues, gui.state.ActivePanel == PanelIssues, issuesH)...)
	leftLines = append(leftLines, gui.renderItemsPanel("PRs", gui.panels.PRs, gui.state.ActivePanel == PanelPRs, prsH)...)

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

func (gui *Gui) renderItemsPanel(title string, panel *panels.ItemsPanel, active bool, height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)
	lines = append(lines, formatPanelTitle(title, active))
	if panel.Loading {
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
		if i >= len(panel.Items) {
			lines = append(lines, "")
			continue
		}
		prefix := "  "
		if i == panel.Selected {
			prefix = "> "
		}
		lines = append(lines, prefix+panel.Format(panel.Items[i]))
	}
	return lines
}

func (gui *Gui) renderDetailPanel(title string, active bool, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)
	lines = append(lines, formatPanelTitle(title, active))
	for _, line := range strings.Split(gui.panels.Detail.Content, "\n") {
		if len(lines) >= height {
			break
		}
		// lines = append(lines, sanitizeRenderLine(line))
		lines = append(lines, line)
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	_ = width
	return lines
}

func padOrTrim(s string, width int) string {
	if width <= 0 {
		return ""
	}
	// s = sanitizeRenderLine(s)
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

func sanitizeRenderLine(s string) string {
	// Normalize tabs and remove control/escape sequences to keep column alignment stable.
	s = strings.ReplaceAll(s, "\r", "")
	var b strings.Builder
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '\x1b' {
			// CSI
			if i+1 < len(runes) && runes[i+1] == '[' {
				i += 2
				for i < len(runes) && !(runes[i] >= 0x40 && runes[i] <= 0x7e) {
					i++
				}
				continue
			}
			// OSC
			if i+1 < len(runes) && runes[i+1] == ']' {
				i += 2
				for i < len(runes) {
					if runes[i] == '\a' {
						break
					}
					if runes[i] == '\x1b' && i+1 < len(runes) && runes[i+1] == '\\' {
						i++
						break
					}
					i++
				}
				continue
			}
			continue
		}
		if r == '\t' {
			b.WriteString("    ")
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
