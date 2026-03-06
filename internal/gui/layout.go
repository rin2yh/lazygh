package gui

import "fmt"

const statusBarHeight = 2

func formatPanelTitle(base string, active bool) string {
	if active {
		return fmt.Sprintf("> %s <", base)
	}
	return fmt.Sprintf(" %s ", base)
}

func panelDisplayName(panel PanelType) string {
	switch panel {
	case PanelRepos:
		return "Repositories"
	case PanelIssues:
		return "Issues"
	case PanelPRs:
		return "PRs"
	case PanelDetail:
		return "Detail"
	default:
		return "Unknown"
	}
}

func formatStatusLine(activePanel PanelType) string {
	return fmt.Sprintf("Panel: %s  [q]Quit  [tab]Panel  [j/k]Navigate  [enter]Select", panelDisplayName(activePanel))
}

func shouldHighlightListPanel(active bool, keepSelectionOnBlur bool) bool {
	return active || keepSelectionOnBlur
}

func statusViewBounds(maxX, maxY int) (int, int, int, int) {
	contentHeight := maxY - statusBarHeight - 1
	return 0, contentHeight + 1, maxX - 1, maxY
}
