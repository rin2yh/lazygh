package gui

import (
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
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
	ansiReset  = "\x1b[0m"
	ansiGreen  = "\x1b[32m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiBlue   = "\x1b[34m"
	ansiCyan   = "\x1b[36m"
	ansiPurple = "\x1b[35m"
	ansiGray   = "\x1b[90m"
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
		case "esc":
			m.gui.focusPRs()
			return m, nil
		case "tab":
			m.gui.cycleFocus()
			return m, nil
		case "j", "down":
			return m, m.handleDownKey()
		case "k", "up":
			return m, m.handleUpKey()
		case "pgdown", "f", " ", "pgup", "b", "home", "g", "end", "G":
			if m.gui.scrollDetailByKey(msg) {
				return m, nil
			}
			return m, nil
		case "h":
			return m, m.handleHKey()
		case "l":
			return m, m.handleLKey()
		case "o":
			m.gui.switchToOverview()
			return m, nil
		case "d":
			return m, m.handleDKey()
		case "enter":
			return m, m.handleDetailLoad()
		}
	}
	return m, nil
}

func (m *model) handleHKey() tea.Cmd {
	if !m.gui.state.IsDiffMode() {
		return nil
	}
	if len(m.gui.diffFiles) > 0 {
		m.gui.focus = panelDiffFiles
	}
	return nil
}

func (m *model) handleLKey() tea.Cmd {
	if m.gui.focus == panelPRs {
		if m.gui.state.IsDiffMode() {
			m.gui.switchToOverview()
		}
		return m.handleDetailLoad()
	}
	if m.gui.state.IsDiffMode() {
		m.gui.focus = panelDiffContent
		return nil
	}
	return nil
}

func (m *model) handleDKey() tea.Cmd {
	if m.gui.switchToDiff() {
		return m.handleDetailLoad()
	}
	return nil
}

func (m *model) handleDownKey() tea.Cmd {
	if m.gui.state.IsDiffMode() {
		switch m.gui.focus {
		case panelDiffFiles:
			m.gui.selectNextDiffFile()
			return nil
		case panelDiffContent:
			m.gui.scrollDetailDown()
			return nil
		default:
			changed := m.gui.navigateDown()
			if changed {
				return m.handleDetailLoad()
			}
			return nil
		}
	}

	m.gui.navigateDown()
	return nil
}

func (m *model) handleUpKey() tea.Cmd {
	if m.gui.state.IsDiffMode() {
		switch m.gui.focus {
		case panelDiffFiles:
			m.gui.selectPrevDiffFile()
			return nil
		case panelDiffContent:
			m.gui.scrollDetailUp()
			return nil
		default:
			changed := m.gui.navigateUp()
			if changed {
				return m.handleDetailLoad()
			}
			return nil
		}
	}

	m.gui.navigateUp()
	return nil
}

func (m *model) View() string {
	return m.gui.render()
}

func toCorePRs(prs []gh.PRItem) []core.Item {
	items := make([]core.Item, 0, len(prs))
	for _, pr := range prs {
		status := pr.State
		if pr.IsDraft {
			status = "DRAFT"
		}
		assignees := make([]string, 0, len(pr.Assignees))
		for _, user := range pr.Assignees {
			name := strings.TrimSpace(user.Login)
			if name != "" {
				assignees = append(assignees, name)
			}
		}
		items = append(items, core.Item{
			Number:    pr.Number,
			Title:     pr.Title,
			Status:    status,
			Assignees: assignees,
		})
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

func (m *model) loadDetailCmd(repo string, number int, mode core.DetailMode) tea.Cmd {
	return func() tea.Msg {
		var (
			content string
			err     error
		)
		switch mode {
		case core.DetailModeDiff:
			content, err = m.gui.client.DiffPR(repo, number)
		default:
			content, err = m.gui.client.ViewPR(repo, number)
		}
		return detailLoadedMsg{mode: mode, number: number, content: content, err: err}
	}
}

func (m *model) handleDetailLoad() tea.Cmd {
	action := m.gui.state.PlanEnter(m.gui.client != nil, os.Getenv("LAZYGH_DEBUG_DETAIL_TEXT"))
	switch action.Kind {
	case core.EnterLoadPRDiff:
		return m.loadDetailCmd(action.Repo, action.Number, core.DetailModeDiff)
	case core.EnterLoadPRDetail:
		return m.loadDetailCmd(action.Repo, action.Number, core.DetailModeOverview)
	default:
		return nil
	}
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

func (gui *Gui) navigateDown() bool {
	return gui.state.NavigateDown()
}

func (gui *Gui) navigateUp() bool {
	return gui.state.NavigateUp()
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

func (gui *Gui) scrollDetailByKey(msg tea.KeyMsg) bool {
	if !gui.state.IsDiffMode() || gui.focus != panelDiffContent {
		return false
	}

	switch msg.String() {
	case "pgdown", "f", " ", "pgup", "b":
		updated, _ := gui.detailViewport.Update(msg)
		gui.detailViewport = updated
		return true
	case "home", "g":
		gui.detailViewport.GotoTop()
		return true
	case "end", "G":
		gui.detailViewport.GotoBottom()
		return true
	default:
		return false
	}
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

func (gui *Gui) selectNextDiffFile() bool {
	if len(gui.diffFiles) == 0 || gui.diffFileSelected >= len(gui.diffFiles)-1 {
		return false
	}
	gui.diffFileSelected++
	return true
}

func (gui *Gui) selectPrevDiffFile() bool {
	if len(gui.diffFiles) == 0 || gui.diffFileSelected <= 0 {
		return false
	}
	gui.diffFileSelected--
	return true
}

func (gui *Gui) scrollDetailDown() {
	gui.detailViewport.ScrollDown(1)
}

func (gui *Gui) scrollDetailUp() {
	gui.detailViewport.ScrollUp(1)
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

	leftRatio := 26
	if gui.state.IsDiffMode() {
		leftRatio = 22
	}
	leftWidth := w * leftRatio / 100
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

	leftLines := gui.renderLeftPanels(leftWidth, contentHeight)
	rightLines := gui.renderRightPanels(rightWidth, contentHeight)

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
		b.WriteRune(' ')
		b.WriteString(padOrTrim(right, rightWidth))
		b.WriteByte('\n')
	}
	b.WriteString(padOrTrim(
		formatStatusLine(
			gui.state.Loading != core.LoadingNone,
			gui.state.IsDiffMode(),
			len(gui.state.PRs) > 0,
			gui.focus,
			len(gui.diffFiles) > 0,
		),
		w,
	))
	return b.String()
}

func (gui *Gui) renderRightPanels(width int, height int) []string {
	if !gui.state.IsDiffMode() {
		return gui.renderDetailPanel("", false, width, height, gui.state.DetailContent)
	}
	coloredDiff := colorizeDiffContent(gui.currentDiffContent())

	if width < 20 {
		return gui.renderDetailPanel("Diff", gui.focus == panelDiffContent, width, height, coloredDiff)
	}

	filesWidth := width * 30 / 100
	if filesWidth < 16 {
		filesWidth = 16
	}
	if filesWidth > width-10 {
		filesWidth = width - 10
	}
	diffWidth := width - filesWidth - 1
	if diffWidth < 1 {
		diffWidth = 1
	}

	filesLines := gui.renderDiffFilesPanel(filesWidth, height)
	diffLines := gui.renderDetailPanel("Diff", gui.focus == panelDiffContent, diffWidth, height, coloredDiff)

	lines := make([]string, 0, height)
	for i := 0; i < height; i++ {
		left := ""
		if i < len(filesLines) {
			left = filesLines[i]
		}
		right := ""
		if i < len(diffLines) {
			right = diffLines[i]
		}
		lines = append(lines, padOrTrim(left, filesWidth)+" "+padOrTrim(right, diffWidth))
	}
	return lines
}

func (gui *Gui) renderPRPanel(height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)

	if gui.state.PRsLoading {
		for len(lines) < height {
			lines = append(lines, "")
		}
		return lines
	}

	if len(gui.state.PRs) == 0 {
		for len(lines) < height {
			if len(lines) == 0 {
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

func (gui *Gui) renderRepoPanel(height int) []string {
	if height <= 0 {
		return nil
	}
	lines := make([]string, 0, height)
	for len(lines) < height {
		if len(lines) == 0 {
			lines = append(lines, formatRepoLine(gui.state.Repo))
		} else {
			lines = append(lines, "")
		}
	}
	return lines
}

func (gui *Gui) renderLeftPanels(width int, height int) []string {
	if height <= 0 {
		return nil
	}

	repoPanelHeight := 4
	if height < repoPanelHeight+1 {
		repoPanelHeight = height / 2
	}
	if repoPanelHeight < 1 {
		repoPanelHeight = 1
	}
	prPanelHeight := height - repoPanelHeight
	if prPanelHeight < 1 {
		prPanelHeight = 1
		repoPanelHeight = height - prPanelHeight
	}

	repoInnerHeight := repoPanelHeight
	if repoPanelHeight > 2 {
		repoInnerHeight = repoPanelHeight - 2
	}
	prInnerHeight := prPanelHeight
	if prPanelHeight > 2 {
		prInnerHeight = prPanelHeight - 2
	}

	repoLines := framePanel("Repository", false, gui.renderRepoPanel(repoInnerHeight), width, repoPanelHeight)
	prLines := framePanel("PRs (Open/Draft)", gui.focus == panelPRs, gui.renderPRPanel(prInnerHeight), width, prPanelHeight)

	lines := make([]string, 0, height)
	lines = append(lines, repoLines...)
	lines = append(lines, prLines...)
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return lines
}

func (gui *Gui) renderDiffFilesPanel(width int, height int) []string {
	if height <= 0 {
		return nil
	}

	innerHeight := height
	if height > 2 {
		innerHeight = height - 2
	}
	lines := make([]string, 0, innerHeight)

	if len(gui.diffFiles) == 0 {
		for len(lines) < innerHeight {
			if len(lines) == 0 {
				lines = append(lines, "No changed files")
			} else {
				lines = append(lines, "")
			}
		}
		return framePanel("Files", gui.focus == panelDiffFiles, lines, width, height)
	}

	start := 0
	if gui.diffFileSelected >= innerHeight {
		start = gui.diffFileSelected - innerHeight + 1
	}
	for i := 0; len(lines) < innerHeight; i++ {
		idx := start + i
		if idx >= len(gui.diffFiles) {
			lines = append(lines, "")
			continue
		}
		prefix := "  "
		if idx == gui.diffFileSelected {
			prefix = "> "
		}
		lines = append(lines, prefix+renderDiffFileListLine(gui.diffFiles[idx]))
	}
	return framePanel("Files", gui.focus == panelDiffFiles, lines, width, height)
}

func (gui *Gui) renderDetailPanel(title string, active bool, width int, height int, content string) []string {
	if height <= 0 {
		return nil
	}

	innerWidth := width
	if width > 2 {
		innerWidth = width - 2
	}
	innerHeight := height
	if height > 2 {
		innerHeight = height - 2
	}
	bodyHeight := innerHeight
	if bodyHeight < 1 {
		bodyHeight = 1
	}
	gui.syncDetailViewport(innerWidth, bodyHeight, content)

	lines := make([]string, 0, innerHeight)
	for _, line := range strings.Split(gui.detailViewport.View(), "\n") {
		if len(lines) >= innerHeight {
			break
		}
		lines = append(lines, line)
	}
	for len(lines) < innerHeight {
		lines = append(lines, "")
	}
	return framePanel(title, active, lines, width, height)
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
func renderDiffFileListLine(file gh.DiffFile) string {
	label := string(file.Status)
	if label == "" {
		label = string(gh.DiffFileStatusModified)
	}
	status := colorizeDiffStatus(label, file.Status)
	additions := colorize(ansiGreen, "+"+strconv.Itoa(file.Additions))
	deletions := colorize(ansiRed, "-"+strconv.Itoa(file.Deletions))
	return status + " " + file.Path + " " + additions + " " + deletions
}

func colorizeDiffStatus(label string, status gh.DiffFileStatus) string {
	switch status {
	case gh.DiffFileStatusAdded:
		return colorize(ansiGreen, label)
	case gh.DiffFileStatusDeleted:
		return colorize(ansiRed, label)
	case gh.DiffFileStatusRenamed:
		return colorize(ansiCyan, label)
	case gh.DiffFileStatusCopied:
		return colorize(ansiBlue, label)
	case gh.DiffFileStatusType:
		return colorize(ansiPurple, label)
	default:
		return colorize(ansiYellow, label)
	}
}

func colorize(color string, text string) string {
	if text == "" {
		return ""
	}
	return color + text + ansiReset
}

func colorizeDiffContent(content string) string {
	if content == "" {
		return ""
	}

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			lines[i] = colorize(ansiBlue, line)
		case strings.HasPrefix(line, "@@"):
			lines[i] = colorize(ansiCyan, line)
		case strings.HasPrefix(line, "+++ "):
			lines[i] = colorize(ansiGreen, line)
		case strings.HasPrefix(line, "--- "):
			lines[i] = colorize(ansiRed, line)
		case strings.HasPrefix(line, "+"):
			lines[i] = colorize(ansiGreen, line)
		case strings.HasPrefix(line, "-"):
			lines[i] = colorize(ansiRed, line)
		case strings.HasPrefix(line, "new file mode "), strings.HasPrefix(line, "deleted file mode "):
			lines[i] = colorize(ansiYellow, line)
		case strings.HasPrefix(line, "rename from "), strings.HasPrefix(line, "rename to "):
			lines[i] = colorize(ansiPurple, line)
		case strings.HasPrefix(line, "index "), strings.HasPrefix(line, "similarity index "):
			lines[i] = colorize(ansiGray, line)
		}
	}

	return strings.Join(lines, "\n")
}

func wrapText(content string, width int) string {
	if width <= 0 || content == "" {
		return content
	}

	srcLines := strings.Split(content, "\n")
	dstLines := make([]string, 0, len(srcLines))
	for _, line := range srcLines {
		lineWidth := xansi.StringWidth(line)
		if lineWidth <= width {
			dstLines = append(dstLines, line)
			continue
		}
		for left := 0; left < lineWidth; left += width {
			right := left + width
			dstLines = append(dstLines, xansi.Cut(line, left, right))
		}
	}
	return strings.Join(dstLines, "\n")
}

func padOrTrim(s string, width int) string {
	if width <= 0 {
		return ""
	}
	out := xansi.Truncate(s, width, "")
	col := xansi.StringWidth(out)
	if col < width {
		out += strings.Repeat(" ", width-col)
	}
	return out
}

func framePanel(title string, active bool, content []string, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	if width < 2 || height < 3 {
		lines := make([]string, 0, height)
		for i := 0; i < height; i++ {
			if i < len(content) {
				lines = append(lines, content[i])
			} else {
				lines = append(lines, "")
			}
		}
		return lines
	}

	innerWidth := width - 2
	innerHeight := height - 2
	lines := make([]string, 0, height)
	top := strings.Repeat("─", innerWidth)
	if strings.TrimSpace(title) != "" {
		topLabel := formatPanelTitle(title, active)
		labelWidth := runewidth.StringWidth(topLabel)
		if labelWidth > 0 {
			if labelWidth >= innerWidth {
				top = padOrTrim(topLabel, innerWidth)
			} else {
				top = topLabel + strings.Repeat("─", innerWidth-labelWidth)
			}
		}
	}
	lines = append(lines, "┌"+top+"┐")
	for i := 0; i < innerHeight; i++ {
		row := ""
		if i < len(content) {
			row = content[i]
		}
		lines = append(lines, "│"+padOrTrim(row, innerWidth)+"│")
	}
	lines = append(lines, "└"+strings.Repeat("─", innerWidth)+"┘")
	return lines
}
