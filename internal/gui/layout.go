package gui

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/core"
)

func formatPanelTitle(base string) string {
	return fmt.Sprintf(" %s ", base)
}

func formatStatusLine(loading bool, diffMode bool, hasPR bool, focus panelFocus, hasFiles bool, hasReviewDrawer bool, inputMode core.ReviewInputMode) string {
	base := "[q]Quit"
	if hasPR {
		base = fmt.Sprintf("%s [enter]Reload", base)
	}

	var modeSpecific string
	switch {
	case !diffMode && hasPR:
		modeSpecific = "[PRs] [j/k/↑/↓]Move [d]Diff"
	case !diffMode:
		modeSpecific = "[d]Diff"
	case focus == panelPRs && hasPR:
		modeSpecific = "[tab]Focus [PRs] [j/k/↑/↓]Move [l]Overview [c/R]Review"
	case focus == panelDiffFiles && hasFiles:
		modeSpecific = "[tab]Focus [Files] [j/k/↑/↓]Move [l]Diff [o]Overview [v]Range [c]Comment"
	case focus == panelReviewDrawer && inputMode == core.ReviewInputComment:
		modeSpecific = "[Ctrl+S]Save Comment [Esc]Cancel [S]Submit [X]Discard"
	case focus == panelReviewDrawer && inputMode == core.ReviewInputSummary:
		modeSpecific = "[Ctrl+S]Save Summary [Esc]Cancel [S]Submit [X]Discard"
	case focus == panelReviewDrawer:
		modeSpecific = "[c]Comment [R]Summary [S]Submit [X]Discard [Esc]Diff"
	case hasPR || hasFiles || hasReviewDrawer:
		modeSpecific = "[tab]Focus [Diff] [j/k/↑/↓]Line [space/b]Page [g/G]Top/Bottom [h]Files [v]Range [c]Comment [R]Summary [S]Submit [X]Discard [o]Overview"
	default:
		modeSpecific = "[o]Overview"
	}

	line := base
	if modeSpecific != "" {
		line = fmt.Sprintf("%s | %s", base, modeSpecific)
	}
	if loading {
		return fmt.Sprintf("Loading...  | %s", line)
	}
	return line
}

func formatRepoLine(repo string) string {
	return repo
}
