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
		if gui.state.ActivePanel == PanelRepos && len(gui.panels.Repos.Repos) > 0 {
			if gui.panels.Repos.Selected < len(gui.panels.Repos.Repos)-1 {
				gui.panels.Repos.Selected++
			}
		}
	case "items":
		if gui.state.ActivePanel == PanelItems && len(gui.panels.Items.Items) > 0 {
			if gui.panels.Items.Selected < len(gui.panels.Items.Items)-1 {
				gui.panels.Items.Selected++
			}
		}
	}
	return nil
}

func (gui *Gui) navigateUp(_ *gocui.Gui, viewName string) error {
	switch viewName {
	case "repos":
		if gui.state.ActivePanel == PanelRepos && gui.panels.Repos.Selected > 0 {
			gui.panels.Repos.Selected--
		}
	case "items":
		if gui.state.ActivePanel == PanelItems && gui.panels.Items.Selected > 0 {
			gui.panels.Items.Selected--
		}
	}
	return nil
}
