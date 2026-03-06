package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/rin2yh/lazygh/internal/gui/panels"
)

func (gui *Gui) setupKeybindings() error {
	// Quit
	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
		return err
	}
	for _, view := range []string{"repos", "issues", "prs", "detail", ""} {
		if err := gui.g.SetKeybinding(view, 'q', gocui.ModNone, gui.quit); err != nil {
			return err
		}
	}

	// Tab: 次パネル
	if err := gui.g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, gui.nextPanel); err != nil {
		return err
	}

	// Enter: repos -> loadItems, issues/prs -> loadDetail
	if err := gui.g.SetKeybinding("repos", gocui.KeyEnter, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		return gui.loadItems()
	}); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("issues", gocui.KeyEnter, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		return gui.loadDetail()
	}); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("prs", gocui.KeyEnter, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		return gui.loadDetail()
	}); err != nil {
		return err
	}

	// j/↓: 下へ, k/↑: 上へ
	for _, view := range []string{"repos", "issues", "prs", "detail"} {
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

func (gui *Gui) listPanelByViewName(viewName string) (PanelType, *panels.ItemsPanel, bool) {
	switch viewName {
	case "issues":
		return PanelIssues, gui.panels.Issues, true
	case "prs":
		return PanelPRs, gui.panels.PRs, true
	default:
		return 0, nil, false
	}
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
	default:
		panelType, p, ok := gui.listPanelByViewName(viewName)
		if !ok || gui.state.ActivePanel != panelType || len(p.Items) == 0 {
			return nil
		}
		if p.Selected < len(p.Items)-1 {
			p.Selected++
			gui.renderPanel(viewName)
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
	default:
		panelType, p, ok := gui.listPanelByViewName(viewName)
		if !ok || gui.state.ActivePanel != panelType || p.Selected <= 0 {
			return nil
		}
		p.Selected--
		gui.renderPanel(viewName)
		gui.refreshDetailPreview()
	}
	return nil
}
