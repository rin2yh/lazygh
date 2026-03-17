package layout

const (
	defaultWidth  = 120
	defaultHeight = 40

	leftPanelRatioNormal = 26
	leftPanelRatioDiff   = 22
	percentDenom         = 100

	statusLineHeight      = 1
	reviewDrawerDivisor   = 3
	reviewDrawerMinHeight = 8
	repoPanelHeight       = 4

	diffSplitMinWidth = 20
	diffFilesRatio    = 30
	diffFilesMinWidth = 16
	diffMinDiffWidth  = 10
)

// DiffSplitWidths calculates the file list and diff content widths for split view.
// Returns filesWidth=0 when totalWidth is too narrow to split.
func DiffSplitWidths(totalWidth int) (filesWidth, diffWidth int) {
	if totalWidth < diffSplitMinWidth {
		return 0, totalWidth
	}
	filesWidth = totalWidth * diffFilesRatio / percentDenom
	if filesWidth < diffFilesMinWidth {
		filesWidth = diffFilesMinWidth
	}
	if filesWidth > totalWidth-diffMinDiffWidth {
		filesWidth = totalWidth - diffMinDiffWidth
	}
	diffWidth = totalWidth - filesWidth - 1
	if diffWidth < 1 {
		diffWidth = 1
	}
	return filesWidth, diffWidth
}

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
		w = defaultWidth
	}
	if h <= 0 {
		h = defaultHeight
	}

	leftRatio := leftPanelRatioNormal
	if diffMode {
		leftRatio = leftPanelRatioDiff
	}
	leftWidth := w * leftRatio / percentDenom
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

	contentHeight := h - statusLineHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	mainHeight := contentHeight
	drawerHeight := 0
	if hasReviewDrawer {
		drawerHeight = contentHeight / reviewDrawerDivisor
		if drawerHeight < reviewDrawerMinHeight {
			drawerHeight = reviewDrawerMinHeight
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

	repoHeight := repoPanelHeight
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
