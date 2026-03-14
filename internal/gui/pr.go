package gui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

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

func (s *screen) loadPRsCmd() tea.Cmd {
	return func() tea.Msg {
		repo, err := s.gui.client.ResolveCurrentRepo()
		if err != nil {
			return prsLoadedMsg{err: err}
		}
		prs, err := s.gui.client.ListPRs(repo)
		if err != nil {
			return prsLoadedMsg{repo: repo, err: err}
		}
		return prsLoadedMsg{repo: repo, prs: toCorePRs(prs)}
	}
}

func (s *screen) loadDetailCmd(repo string, number int, mode core.DetailMode) tea.Cmd {
	return func() tea.Msg {
		var (
			content string
			err     error
		)
		switch mode {
		case core.DetailModeDiff:
			content, err = s.gui.client.DiffPR(repo, number)
		default:
			content, err = s.gui.client.ViewPR(repo, number)
		}
		return detailLoadedMsg{mode: mode, number: number, content: content, err: err}
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
		gui.diffLineSelected = 0
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
	gui.diffLineSelected = 0
	if prevPath != "" {
		for i, file := range files {
			if file.Path == prevPath {
				gui.diffFileSelected = i
				break
			}
		}
	}
	gui.ensureDiffLineSelection()
}
