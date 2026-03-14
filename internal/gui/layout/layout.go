package layout

type Screen struct {
	Width         int
	Height        int
	LeftWidth     int
	RightWidth    int
	ContentHeight int
	MainHeight    int
	DrawerHeight  int
	RepoHeight    int
	PRHeight      int
}

func New(width int, height int, diffMode bool, hasReviewDrawer bool) Screen {
	w := width
	h := height
	if w <= 0 {
		w = 120
	}
	if h <= 0 {
		h = 40
	}

	leftRatio := 26
	if diffMode {
		leftRatio = 22
	}
	leftWidth := w * leftRatio / 100
	if leftWidth < 1 {
		leftWidth = 1
	}
	if leftWidth > w-2 {
		leftWidth = w - 2
	}
	rightWidth := w - leftWidth - 1
	if rightWidth < 1 {
		rightWidth = 1
	}

	contentHeight := h - 1
	if contentHeight < 1 {
		contentHeight = 1
	}

	mainHeight := contentHeight
	drawerHeight := 0
	if hasReviewDrawer {
		drawerHeight = contentHeight / 3
		if drawerHeight < 8 {
			drawerHeight = 8
		}
		if drawerHeight >= contentHeight {
			drawerHeight = contentHeight - 1
		}
		if drawerHeight < 0 {
			drawerHeight = 0
		}
		mainHeight = contentHeight - drawerHeight
	}
	if mainHeight < 1 {
		mainHeight = 1
	}

	repoHeight := 4
	if mainHeight < repoHeight+1 {
		repoHeight = mainHeight / 2
	}
	if repoHeight < 1 {
		repoHeight = 1
	}
	prHeight := mainHeight - repoHeight
	if prHeight < 1 {
		prHeight = 1
		repoHeight = mainHeight - prHeight
	}

	return Screen{
		Width:         w,
		Height:        h,
		LeftWidth:     leftWidth,
		RightWidth:    rightWidth,
		ContentHeight: contentHeight,
		MainHeight:    mainHeight,
		DrawerHeight:  drawerHeight,
		RepoHeight:    repoHeight,
		PRHeight:      prHeight,
	}
}

func (s Screen) InnerHeight(height int) int {
	if height > 2 {
		return height - 2
	}
	return height
}

func (s Screen) InnerWidth(width int) int {
	if width > 2 {
		return width - 2
	}
	return width
}
