package gui

import (
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
	if v, err := g.SetView("repos", 0, leftTop, leftWidth-1, reposBottom); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Repositories "
		v.Wrap = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		gui.panels.Repos.Render(v)
	}

	// Issues panel (左中)
	if v, err := g.SetView("issues", 0, issuesTop, leftWidth-1, issuesBottom); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Issues "
		v.Wrap = false
		gui.panels.Issues.Render(v)
	}

	// PRs panel (左下)
	if v, err := g.SetView("prs", 0, prsTop, leftWidth-1, contentHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " PRs "
		v.Wrap = false
		gui.panels.PRs.Render(v)
	}

	// Detail panel (右)
	if v, err := g.SetView("detail", leftWidth, 1, maxX-1, contentHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Detail "
		v.Wrap = true
		gui.panels.Detail.Render(v)
	}

	// Status bar
	if v, err := g.SetView("status", 0, contentHeight+1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		_, _ = v.Write([]byte("[q]Quit  [tab]Panel  [j/k]Navigate  [enter]Select"))
	}

	// フォーカス設定
	if err := g.SetCurrentView(gui.activeViewName()); err != nil {
		return err
	}

	// アクティブパネルの枠色
	gui.updateBorderColors(g)

	return nil
}

func (gui *Gui) updateBorderColors(g *gocui.Gui) {
	views := []string{"repos", "issues", "prs", "detail"}
	active := gui.activeViewName()
	for _, name := range views {
		v, err := g.View(name)
		if err != nil {
			continue
		}
		if name == active {
			v.FgColor = gocui.ColorGreen
		} else {
			v.FgColor = gocui.ColorDefault
		}
	}
}
