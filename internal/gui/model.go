package gui

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

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
