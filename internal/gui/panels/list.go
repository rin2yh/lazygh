package panels

import "github.com/jesseduffield/gocui"

type ListPanel struct {
	Selected int
}

func NewListPanel() ListPanel {
	return ListPanel{Selected: 0}
}

type RowRenderer func(index int) string

func (p *ListPanel) Render(v *gocui.View, count int, renderRow RowRenderer, showSelection bool) {
	adjustScroll(v, p.Selected)
	v.Clear()
	if count == 0 || renderRow == nil {
		_ = v.SetCursor(0, 0)
		return
	}
	_, originY := v.Origin()
	_, height := v.Size()
	_ = v.SetCursor(0, calcCursorY(p.Selected, originY, height))
	for i := 0; i < count; i++ {
		_, _ = v.Write([]byte(formatListRow(renderRow(i), showSelection && i == p.Selected)))
	}
}

func formatListRow(row string, selected bool) string {
	prefix := "  "
	if selected {
		prefix = "> "
	}
	return prefix + normalizeDisplayText(row) + "\n"
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
