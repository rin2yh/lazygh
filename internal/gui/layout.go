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
	case !diffMode && focus == panelRepo && hasPR:
		modeSpecific = "[Repo] [l]Next Panel [d]Diff"
	case !diffMode && focus == panelRepo:
		modeSpecific = "[Repo] [l]Next Panel [d]Diff"
	case !diffMode && focus == panelPRs && hasPR:
		modeSpecific = "[PRs] [h/l]Prev/Next Panel [j/k/↑/↓]Move [d]Diff"
	case !diffMode && focus == panelDiffContent && hasPR:
		modeSpecific = "[Overview] [h]Prev Panel [space/b]Page [enter]Reload [d]Diff"
	case !diffMode && focus == panelDiffContent:
		modeSpecific = "[Overview] [h]Prev Panel [d]Diff"
	case !diffMode:
		modeSpecific = "[l]Next Panel [d]Diff"
	case focus == panelRepo && hasPR:
		modeSpecific = "[tab]Focus [Repo] [l]Next Panel [d]Diff"
	case focus == panelRepo:
		modeSpecific = "[tab]Focus [Repo] [l]Next Panel [d]Diff"
	case focus == panelPRs && hasPR:
		modeSpecific = "[tab]Focus [PRs] [h/l]Prev/Next Panel [j/k/↑/↓]Move [c/R]Review"
	case focus == panelDiffFiles && hasFiles:
		modeSpecific = "[tab]Focus [Files] [j/k/↑/↓]Move [h/l]Prev/Next Panel [o]Overview [v]Range [c]Comment"
	case focus == panelReviewDrawer && inputMode == core.ReviewInputComment:
		modeSpecific = "[Ctrl+S]Save Comment [Esc]Cancel [S]Submit [X]Discard"
	case focus == panelReviewDrawer && inputMode == core.ReviewInputSummary:
		modeSpecific = "[Ctrl+S]Save Summary [Esc]Cancel [S]Submit [X]Discard"
	case focus == panelReviewDrawer:
		modeSpecific = "[h]Prev Panel [c]Comment [R]Summary [S]Submit [X]Discard [Esc]Diff"
	case hasPR || hasFiles || hasReviewDrawer:
		modeSpecific = "[tab]Focus [Diff] [j/k/↑/↓]Line [space/b]Page [g/G]Top/Bottom [h/l]Prev/Next Panel [v]Range [c]Comment [R]Summary [S]Submit [X]Discard [o]Overview"
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
