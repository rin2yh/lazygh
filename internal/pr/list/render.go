package list

import (
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/pr"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

// Input holds the data needed to render the left PR list panels.
type Input struct {
	Repo     string
	Fetching bool
	Items    []pr.Item
	Selected int
	Filter   string
}

const repoPanelHeight = 4

var (
	prefixOpen   = widget.Colorize("O", "green")
	prefixDraft  = widget.Colorize("D", "gray")
	prefixClosed = widget.Colorize("C", "red")
	prefixMerged = widget.Colorize("M", "purple")
)

func statusPrefix(status string) string {
	switch status {
	case pr.PRStatusDraft:
		return prefixDraft
	case pr.PRStatusClosed:
		return prefixClosed
	case pr.PRStatusMerged:
		return prefixMerged
	default:
		return prefixOpen
	}
}

func splitHeight(total int) (repoH, prH int) {
	repoH = repoPanelHeight
	if total < repoH+1 {
		repoH = total / 2
	}
	if repoH < 1 {
		repoH = 1
	}
	prH = total - repoH
	if prH < 1 {
		prH = 1
		repoH = total - prH
	}
	return repoH, prH
}

// RenderLeft renders the Repository and PR panels on the left side.
func RenderLeft(input Input, style func(layout.Focus) widget.PanelStyle, width, height int) []string {
	repoHeight, prHeight := splitHeight(height)
	repoLines := widget.FramePanel("Repository", renderRepo(input), width, repoHeight, style(layout.FocusRepo))
	prTitle := "PR [" + input.Filter + "]"
	prLines := widget.FramePanel(prTitle, renderPRs(input), width, prHeight, style(layout.FocusPRs))
	lines := make([]string, 0, height)
	lines = append(lines, repoLines...)
	lines = append(lines, prLines...)
	return lines
}

func renderRepo(input Input) []string {
	return []string{input.Repo}
}

func renderPRs(input Input) []string {
	if input.Fetching {
		return nil
	}
	if len(input.Items) == 0 {
		return []string{"No pull requests"}
	}
	lines := make([]string, 0, len(input.Items))
	for i, item := range input.Items {
		line := widget.ListItem(statusPrefix(item.Status)+" "+formatItem(item), i == input.Selected)
		lines = append(lines, line)
	}
	return lines
}
