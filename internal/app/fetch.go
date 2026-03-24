package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/pr"
	"github.com/rin2yh/lazygh/internal/pr/diff"
	"github.com/rin2yh/lazygh/internal/pr/list"
	"github.com/rin2yh/lazygh/internal/pr/overview"
)

type prsLoadedMsg struct {
	repo string
	prs  []pr.Item
	err  error
}

type detailLoadedMsg struct {
	mode    overview.DetailMode
	number  int
	content string
	err     error
}

func (s *screen) loadPRsCmd() tea.Cmd {
	filter := s.gui.coord.Filter
	return func() tea.Msg {
		repo, err := s.gui.client.ResolveCurrentRepo()
		if err != nil {
			return prsLoadedMsg{err: err}
		}
		prs, err := s.gui.client.ListPRs(repo, filter.StateArg())
		if err != nil {
			return prsLoadedMsg{repo: repo, err: err}
		}
		return prsLoadedMsg{repo: repo, prs: list.Convert(prs, filter)}
	}
}

func (s *screen) loadDetailCmd(repo string, number int, mode overview.DetailMode) tea.Cmd {
	return func() tea.Msg {
		var (
			content string
			err     error
		)
		switch mode {
		case overview.DetailModeDiff:
			content, err = s.gui.client.DiffPR(repo, number)
		default:
			content, err = s.gui.client.ViewPR(repo, number)
		}
		return detailLoadedMsg{mode: mode, number: number, content: content, err: err}
	}
}

func (gui *Gui) applyPRsResult(msg prsLoadedMsg) {
	gui.coord.ApplyPRsResult(msg.repo, msg.prs, msg.err)
	gui.focus = layout.FocusPRs
}

func (gui *Gui) applyDetailResult(msg detailLoadedMsg) {
	if !gui.coord.ShouldApplyDetailResult(msg.mode, msg.number) {
		return
	}
	if msg.mode == overview.DetailModeDiff {
		gui.coord.ApplyDiffResult(msg.content, msg.err)
		if msg.err != nil {
			gui.diff.Reset()
			gui.resetDiffFocusIfOnFiles()
			return
		}
		gui.updateDiffFiles(gui.coord.Overview.Content)
		return
	}
	gui.coord.ApplyDetailResult(msg.content, msg.err)
}

func (gui *Gui) currentDiffContent() string {
	files := gui.diff.Files()
	selected := gui.diff.FileSelected()
	if len(files) == 0 {
		return gui.coord.Overview.Content
	}
	if selected < 0 || selected >= len(files) {
		return gui.coord.Overview.Content
	}
	return files[selected].Content
}

func (gui *Gui) updateDiffFiles(content string) {
	files, selected, lineSelected := diff.ParseFiles(gui.diff.Files(), gui.diff.FileSelected(), content)
	if len(files) == 0 {
		gui.diff.Reset()
		gui.resetDiffFocusIfOnFiles()
		return
	}

	gui.diff.SetFiles(files)
	gui.diff.SetFileSelected(selected)
	gui.diff.SetLineSelected(lineSelected)
	gui.diff.EnsureLineSelection()
}
