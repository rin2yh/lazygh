package layout

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/core"
)

type Focus int

const (
	FocusRepo Focus = iota
	FocusPRs
	FocusDiffFiles
	FocusDiffContent
	FocusReviewDrawer
)

type Status struct {
	Loading         bool
	DiffMode        bool
	HasPR           bool
	Focus           Focus
	HasFiles        bool
	HasReviewDrawer bool
	InputMode       core.ReviewInputMode
}

func (s Status) String() string {
	base := "[q]Quit"
	if s.HasPR {
		base = fmt.Sprintf("%s [enter]Reload", base)
	}

	var modeSpecific string
	switch {
	case !s.DiffMode && s.Focus == FocusRepo && s.HasPR:
		modeSpecific = "[Repo] [l]Next Panel [d]Diff"
	case !s.DiffMode && s.Focus == FocusRepo:
		modeSpecific = "[Repo] [l]Next Panel [d]Diff"
	case !s.DiffMode && s.Focus == FocusPRs && s.HasPR:
		modeSpecific = "[PRs] [h/l]Prev/Next Panel [j/k/↑/↓]Move [d]Diff"
	case !s.DiffMode && s.Focus == FocusDiffContent && s.HasPR:
		modeSpecific = "[Overview] [h]Prev Panel [space/b]Page [enter]Reload [d]Diff"
	case !s.DiffMode && s.Focus == FocusDiffContent:
		modeSpecific = "[Overview] [h]Prev Panel [d]Diff"
	case !s.DiffMode:
		modeSpecific = "[l]Next Panel [d]Diff"
	case s.Focus == FocusRepo && s.HasPR:
		modeSpecific = "[tab]Focus [Repo] [l]Next Panel [d]Diff"
	case s.Focus == FocusRepo:
		modeSpecific = "[tab]Focus [Repo] [l]Next Panel [d]Diff"
	case s.Focus == FocusPRs && s.HasPR:
		modeSpecific = "[tab]Focus [PRs] [h/l]Prev/Next Panel [j/k/↑/↓]Move [c/R]Review"
	case s.Focus == FocusDiffFiles && s.HasFiles:
		modeSpecific = "[tab]Focus [Files] [j/k/↑/↓]Move [h/l]Prev/Next Panel [o]Overview [v]Range [c]Comment"
	case s.Focus == FocusReviewDrawer && s.InputMode == core.ReviewInputComment:
		modeSpecific = "[Ctrl+S]Save Comment [Esc]Cancel [S]Submit [X]Discard"
	case s.Focus == FocusReviewDrawer && s.InputMode == core.ReviewInputSummary:
		modeSpecific = "[Ctrl+S]Save Summary [Esc]Cancel [S]Submit [X]Discard"
	case s.Focus == FocusReviewDrawer:
		modeSpecific = "[h]Prev Panel [c]Comment [R]Summary [S]Submit [X]Discard [Esc]Diff"
	case s.HasPR || s.HasFiles || s.HasReviewDrawer:
		modeSpecific = "[tab]Focus [Diff] [j/k/↑/↓]Line [space/b]Page [g/G]Top/Bottom [h/l]Prev/Next Panel [v]Range [c]Comment [R]Summary [S]Submit [X]Discard [o]Overview"
	default:
		modeSpecific = "[o]Overview"
	}

	line := base
	if modeSpecific != "" {
		line = fmt.Sprintf("%s | %s", base, modeSpecific)
	}
	if s.Loading {
		return fmt.Sprintf("Loading...  | %s", line)
	}
	return line
}
