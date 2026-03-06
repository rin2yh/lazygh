package panels

type ListPanel struct {
	Selected int
}

func NewListPanel() ListPanel {
	return ListPanel{Selected: 0}
}

func formatListRow(row string, selected bool) string {
	prefix := "  "
	if selected {
		prefix = "> "
	}
	return prefix + row + "\n"
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
