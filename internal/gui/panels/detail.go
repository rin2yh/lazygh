package panels

import "github.com/jesseduffield/gocui"

type DetailPanel struct {
	Content string
	ScrollY int
}

func NewDetailPanel() *DetailPanel {
	return &DetailPanel{}
}

func (p *DetailPanel) Render(v *gocui.View) {
	v.Clear()
	_, _ = v.Write([]byte(normalizeDisplayText(p.Content)))
}

func (p *DetailPanel) SetContent(content string) {
	p.Content = content
	p.ScrollY = 0
}
