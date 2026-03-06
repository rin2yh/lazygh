package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) setupKeybindings() error {
	// Quit
	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
		return err
	}
	for _, view := range []string{"repos", "items", "detail", ""} {
		if err := gui.g.SetKeybinding(view, 'q', gocui.ModNone, gui.quit); err != nil {
			return err
		}
	}

	// Tab: 次パネル
	if err := gui.g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, gui.nextPanel); err != nil {
		return err
	}

	// Enter: repos → loadItems, items → loadDetail
	if err := gui.g.SetKeybinding("repos", gocui.KeyEnter, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		return gui.loadItems()
	}); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("items", gocui.KeyEnter, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		return gui.loadDetail()
	}); err != nil {
		return err
	}

	// j/↓: 下へ, k/↑: 上へ
	for _, view := range []string{"repos", "items", "detail"} {
		v := view
		navDown := func(g *gocui.Gui, _ *gocui.View) error { return gui.navigateDown(g, v) }
		navUp := func(g *gocui.Gui, _ *gocui.View) error { return gui.navigateUp(g, v) }
		for _, kb := range []struct {
			key any
			fn  func(*gocui.Gui, *gocui.View) error
		}{
			{'j', navDown}, {gocui.KeyArrowDown, navDown},
			{'k', navUp}, {gocui.KeyArrowUp, navUp},
		} {
			if err := gui.g.SetKeybinding(v, kb.key, gocui.ModNone, kb.fn); err != nil {
				return err
			}
		}
	}

	return nil
}

func (gui *Gui) quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (gui *Gui) nextPanel(_ *gocui.Gui, _ *gocui.View) error {
	gui.state.ActivePanel = (gui.state.ActivePanel + 1) % panelCount
	return nil
}

func (gui *Gui) navigateDown(_ *gocui.Gui, viewName string) error {
	switch viewName {
	case "repos":
		p := gui.panels.Repos
		if gui.state.ActivePanel != PanelRepos || p.Selected >= len(p.Repos)-1 {
			return nil
		}
		p.Selected++
		gui.renderPanel("repos")
	case "items":
		if gui.state.ActivePanel != PanelItems || len(gui.panels.Items.Items) == 0 {
			return nil
		}
		p := gui.panels.Items
		if p.Selected < len(p.Items)-1 {
			p.Selected++
			gui.renderPanel("items")
		}
		gui.refreshDetailPreview()
	}
	return nil
}

func (gui *Gui) navigateUp(_ *gocui.Gui, viewName string) error {
	switch viewName {
	case "repos":
		p := gui.panels.Repos
		if gui.state.ActivePanel == PanelRepos && p.Selected > 0 {
			p.Selected--
			gui.renderPanel("repos")
		}
	case "items":
		p := gui.panels.Items
		if gui.state.ActivePanel == PanelItems && p.Selected > 0 {
			p.Selected--
			gui.renderPanel("items")
			gui.refreshDetailPreview()
		}
	}
	return nil
}
