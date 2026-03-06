package panels

import "github.com/jesseduffield/gocui"

type ReposPanel struct {
	Repos    []string
	Selected int
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

func adjustScroll(v *gocui.View, selected int) {
	_, height := v.Size()
	_, originY := v.Origin()
	_ = v.SetOrigin(0, calcOriginY(selected, originY, height))
}

func (p *ReposPanel) Render(v *gocui.View) {
	adjustScroll(v, p.Selected)
	v.Clear()
	for i, repo := range p.Repos {
		prefix := "  "
		if i == p.Selected {
			prefix = "> "
		}
		_, _ = v.Write([]byte(prefix + repo + "\n"))
	}
}
