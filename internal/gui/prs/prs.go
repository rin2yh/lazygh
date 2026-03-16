package prs

import (
	"github.com/rin2yh/lazygh/internal/gui/layout"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

type Input struct {
	Repo       string
	PRsLoading bool
	PRs        []string
	PRSelected int
	Filter     string
}

func RenderLeft(input Input, repoHeight, prHeight int, active func(layout.Focus) bool, style func(bool) widget.PanelStyle, width int) []string {
	height := repoHeight + prHeight
	repoLines := widget.FramePanel("Repository", renderRepo(input), width, repoHeight, style(active(layout.FocusRepo)))
	prTitle := "PRs [" + input.Filter + "]"
	prLines := widget.FramePanel(prTitle, renderPRs(input), width, prHeight, style(active(layout.FocusPRs)))
	lines := make([]string, 0, height)
	lines = append(lines, repoLines...)
	lines = append(lines, prLines...)
	return lines
}

func renderRepo(input Input) []string {
	return []string{input.Repo}
}

func renderPRs(input Input) []string {
	if input.PRsLoading {
		return nil
	}
	if len(input.PRs) == 0 {
		return []string{"No pull requests"}
	}
	lines := make([]string, 0, len(input.PRs))
	for i, pr := range input.PRs {
		line := widget.ListItem(pr, i == input.PRSelected)
		lines = append(lines, line)
	}
	return lines
}
