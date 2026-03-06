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

func (p *ReposPanel) Render(v *gocui.View) {
	v.Clear()
	for i, repo := range p.Repos {
		prefix := "  "
		if i == p.Selected {
			prefix = "> "
		}
		_, _ = v.Write([]byte(prefix + repo + "\n"))
	}
}
