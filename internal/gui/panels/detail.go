package panels

type DetailPanel struct {
	Content string
	ScrollY int
}

func NewDetailPanel() *DetailPanel {
	return &DetailPanel{}
}

func (p *DetailPanel) SetContent(content string) {
	p.Content = content
	p.ScrollY = 0
}
