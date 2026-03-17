package gui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/gh"
	guidiff "github.com/rin2yh/lazygh/internal/gui/diff"
	"github.com/rin2yh/lazygh/internal/model"
)

type prsLoadedMsg struct {
	repo string
	prs  []model.Item
	err  error
}

type detailLoadedMsg struct {
	mode    model.DetailMode
	number  int
	content string
	err     error
}

func toCorePRs(prs []gh.PRItem, filter model.PRFilterMask) []model.Item {
	items := make([]model.Item, 0, len(prs))
	for _, pr := range prs {
		if !filter.Matches(pr.State) {
			continue
		}
		status := pr.State
		if pr.IsDraft {
			status = model.PRStatusDraft
		}
		assignees := make([]string, 0, len(pr.Assignees))
		for _, user := range pr.Assignees {
			name := strings.TrimSpace(user.Login)
			if name != "" {
				assignees = append(assignees, name)
			}
		}
		items = append(items, model.Item{
			Number:    pr.Number,
			Title:     pr.Title,
			Status:    status,
			Assignees: assignees,
		})
	}
	return items
}

func (s *screen) loadPRsCmd() tea.Cmd {
	filter := s.gui.state.Filter
	return func() tea.Msg {
		repo, err := s.gui.client.ResolveCurrentRepo()
		if err != nil {
			return prsLoadedMsg{err: err}
		}
		prs, err := s.gui.client.ListPRs(repo, filter.StateArg())
		if err != nil {
			return prsLoadedMsg{repo: repo, err: err}
		}
		return prsLoadedMsg{repo: repo, prs: toCorePRs(prs, filter)}
	}
}

func (s *screen) loadDetailCmd(repo string, number int, mode model.DetailMode) tea.Cmd {
	return func() tea.Msg {
		var (
			content string
			err     error
		)
		switch mode {
		case model.DetailModeDiff:
			content, err = s.gui.client.DiffPR(repo, number)
		default:
			content, err = s.gui.client.ViewPR(repo, number)
		}
		return detailLoadedMsg{mode: mode, number: number, content: content, err: err}
	}
}

func (gui *Gui) applyPRsResult(msg prsLoadedMsg) {
	gui.state.ApplyPRsResult(msg.repo, msg.prs, msg.err)
	gui.review.Reset()
	gui.focus = panelPRs
}

func (gui *Gui) applyDetailResult(msg detailLoadedMsg) {
	if !gui.state.ShouldApplyDetailResult(msg.mode, msg.number) {
		return
	}
	if msg.mode == model.DetailModeDiff {
		gui.state.ApplyDiffResult(msg.content, msg.err)
		if msg.err != nil {
			gui.diff.Reset()
			if gui.focus == panelDiffFiles {
				gui.focus = panelDiffContent
			}
			return
		}
		gui.updateDiffFiles(gui.state.Detail.Content)
		return
	}
	gui.state.ApplyDetailResult(msg.content, msg.err)
}

func (gui *Gui) currentDiffContent() string {
	files := gui.diff.Files()
	selected := gui.diff.FileSelected()
	if len(files) == 0 {
		return gui.state.Detail.Content
	}
	if selected < 0 || selected >= len(files) {
		return gui.state.Detail.Content
	}
	return files[selected].Content
}

func (gui *Gui) updateDiffFiles(content string) {
	files, selected, lineSelected := guidiff.ParseFiles(gui.diff.Files(), gui.diff.FileSelected(), content)
	if len(files) == 0 {
		gui.diff.Reset()
		if gui.focus == panelDiffFiles {
			gui.focus = panelDiffContent
		}
		return
	}

	gui.diff.SetFiles(files)
	gui.diff.SetFileSelected(selected)
	gui.diff.SetLineSelected(lineSelected)
	gui.diff.EnsureLineSelection()
}
