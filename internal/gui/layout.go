package gui

import (
	"github.com/jesseduffield/gocui"
)

const statusBarHeight = 2

func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	leftWidth := maxX * 30 / 100
	contentHeight := maxY - statusBarHeight - 1
	reposHeight := contentHeight * 40 / 100
	itemsTop := reposHeight + 1

	// Repos panel (左上)
	if v, err := g.SetView("repos", 0, 1, leftWidth-1, reposHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Repositories "
		v.Wrap = false
		gui.panels.Repos.Render(v)
	}

	// Items panel (左下)
	if v, err := g.SetView("items", 0, itemsTop, leftWidth-1, contentHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Items "
		v.Wrap = false
		gui.panels.Items.Render(v)
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
	views := []string{"repos", "items", "detail"}
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
