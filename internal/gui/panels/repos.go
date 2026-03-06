package panels

import "github.com/jesseduffield/gocui"

type ReposPanel struct {
	Repos    []string
	Selected int
	Loading  bool
}

func NewReposPanel() *ReposPanel {
	return &ReposPanel{
		Repos:    []string{},
		Selected: 0,
	}
}

func calcOriginY(selected, originY, height int) int {
	if selected < originY {
		return selected
	}
	if selected >= originY+height {
		return selected - height + 1
	}
	return originY
}

func calcCursorY(selected, originY, height int) int {
	if height <= 0 {
		return 0
	}

	cursorY := selected - originY
	if cursorY < 0 {
		return 0
	}
	if cursorY >= height {
		return height - 1
	}
	return cursorY
}

func adjustScroll(v *gocui.View, selected int) {
	_, height := v.Size()
	_, originY := v.Origin()
	_ = v.SetOrigin(0, calcOriginY(selected, originY, height))
}

func (p *ReposPanel) Render(v *gocui.View) {
	v.Clear()
	if p.Loading {
		_ = v.SetCursor(0, 0)
		_, _ = v.Write([]byte("Loading...\n"))
		return
	}
	adjustScroll(v, p.Selected)
	_, originY := v.Origin()
	_, height := v.Size()
	_ = v.SetCursor(0, calcCursorY(p.Selected, originY, height))
	for i, repo := range p.Repos {
		prefix := "  "
		if i == p.Selected {
			prefix = "> "
		}
		_, _ = v.Write([]byte(prefix + repo + "\n"))
	}
}
