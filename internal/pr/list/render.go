package list

import (
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/model"
	"github.com/rin2yh/lazygh/pkg/gui/widget"
)

// PanelInput holds the data needed to render the left PR list panels.
type PanelInput struct {
	Repo     string
	Fetching bool
	Items    []model.Item
	Selected int
	Filter   string
}

var (
	prefixOpen   = widget.Colorize("O", "green")
	prefixDraft  = widget.Colorize("D", "gray")
	prefixClosed = widget.Colorize("C", "red")
	prefixMerged = widget.Colorize("M", "purple")
)

func statusPrefix(status string) string {
	switch status {
	case model.PRStatusDraft:
		return prefixDraft
	case model.PRStatusClosed:
		return prefixClosed
	case model.PRStatusMerged:
		return prefixMerged
	default:
		return prefixOpen
	}
}

// RenderLeft renders the Repository and PR panels on the left side.
func RenderLeft(input PanelInput, repoHeight, prHeight int, active func(layout.Focus) bool, style func(bool) widget.PanelStyle, width int) []string {
	height := repoHeight + prHeight
	repoLines := widget.FramePanel("Repository", renderRepo(input), width, repoHeight, style(active(layout.FocusRepo)))
	prTitle := "PR [" + input.Filter + "]"
	prLines := widget.FramePanel(prTitle, renderPRs(input), width, prHeight, style(active(layout.FocusPRs)))
	lines := make([]string, 0, height)
	lines = append(lines, repoLines...)
	lines = append(lines, prLines...)
	return lines
}

func renderRepo(input PanelInput) []string {
	return []string{input.Repo}
}

func renderPRs(input PanelInput) []string {
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
