package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

const statusBarHeight = 2

func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	leftWidth := maxX * 30 / 100
	contentHeight := maxY - statusBarHeight - 1
	leftTop := 1
	reposBottom := leftTop + (contentHeight / 3) - 1
	issuesTop := reposBottom + 1
	issuesBottom := leftTop + (contentHeight * 2 / 3) - 1
	prsTop := issuesBottom + 1

	// Repos panel (左上)
	if v, err := g.SetView("repos", 0, leftTop, leftWidth-1, reposBottom); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Title = formatPanelTitle("Repositories", gui.state.ActivePanel == PanelRepos)
		v.Wrap = false
		v.Highlight = shouldHighlightListPanel(gui.state.ActivePanel == PanelRepos, gui.panels.Repos.KeepSelectionOnBlur)
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		if err == gocui.ErrUnknownView {
			gui.panels.Repos.Render(v, gui.state.ActivePanel == PanelRepos)
		}
	}

	// Issues panel (左中)
	if v, err := g.SetView("issues", 0, issuesTop, leftWidth-1, issuesBottom); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Title = formatPanelTitle("Issues", gui.state.ActivePanel == PanelIssues)
		v.Wrap = false
		v.Highlight = shouldHighlightListPanel(gui.state.ActivePanel == PanelIssues, gui.panels.Issues.KeepSelectionOnBlur)
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		if err == gocui.ErrUnknownView {
			gui.panels.Issues.Render(v, gui.state.ActivePanel == PanelIssues)
		}
	}

	// PRs panel (左下)
	if v, err := g.SetView("prs", 0, prsTop, leftWidth-1, contentHeight); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Title = formatPanelTitle("PRs", gui.state.ActivePanel == PanelPRs)
		v.Wrap = false
		v.Highlight = shouldHighlightListPanel(gui.state.ActivePanel == PanelPRs, gui.panels.PRs.KeepSelectionOnBlur)
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		if err == gocui.ErrUnknownView {
			gui.panels.PRs.Render(v, gui.state.ActivePanel == PanelPRs)
		}
	}

	// Detail panel (右)
	if v, err := g.SetView("detail", leftWidth, 1, maxX-1, contentHeight); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Title = formatPanelTitle("Detail", gui.state.ActivePanel == PanelDetail)
		v.Wrap = true
		if err == gocui.ErrUnknownView {
			gui.panels.Detail.Render(v)
		}
	}

	// Status bar (Frame=false でも1行描画されるように高さを2確保)
	statusX0, statusY0, statusX1, statusY1 := statusViewBounds(maxX, maxY)
	if v, err := g.SetView("status", statusX0, statusY0, statusX1, statusY1); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Frame = false
		v.Clear()
		_, _ = v.Write([]byte(formatStatusLine(gui.state.ActivePanel)))
	}

	// フォーカス設定
	if err := g.SetCurrentView(gui.activeViewName()); err != nil {
		return err
	}

	return nil
}

func formatPanelTitle(base string, active bool) string {
	if active {
		return fmt.Sprintf("> %s <", base)
	}
	return fmt.Sprintf(" %s ", base)
}

func panelDisplayName(panel PanelType) string {
	switch panel {
	case PanelRepos:
		return "Repositories"
	case PanelIssues:
		return "Issues"
	case PanelPRs:
		return "PRs"
	case PanelDetail:
		return "Detail"
	default:
		return "Unknown"
	}
}

func formatStatusLine(activePanel PanelType) string {
	return fmt.Sprintf("Panel: %s  [q]Quit  [tab]Panel  [j/k]Navigate  [enter]Select", panelDisplayName(activePanel))
}

func statusViewBounds(maxX, maxY int) (int, int, int, int) {
	contentHeight := maxY - statusBarHeight - 1
	return 0, contentHeight + 1, maxX - 1, maxY
}

func shouldHighlightListPanel(active bool, keepSelectionOnBlur bool) bool {
	return active || keepSelectionOnBlur
}
