package layout

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/config"
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
	Keys            config.KeyBindings
}

func (s Status) String() string {
	base := fmt.Sprintf("[%s]Quit", s.Keys.QuitLabel())
	if s.HasPR {
		base = fmt.Sprintf("%s [%s]Reload", base, s.Keys.ReloadLabel())
	}

	var modeSpecific string
	switch {
	case !s.DiffMode && s.Focus == FocusRepo:
		modeSpecific = fmt.Sprintf("[Repo] [%s]Next Panel [%s]Diff", s.Keys.Label(config.ActionPanelNext), s.Keys.DiffLabel())
	case !s.DiffMode && s.Focus == FocusPRs && s.HasPR:
		modeSpecific = fmt.Sprintf("[PRs] [%s]Prev/Next Panel [%s]Move [%s]Diff", s.Keys.PanelLabel(), s.Keys.MoveLabel(), s.Keys.DiffLabel())
	case !s.DiffMode && s.Focus == FocusDiffContent && s.HasPR:
		modeSpecific = fmt.Sprintf("[Overview] [%s]Prev Panel [%s]Page [%s]Reload [%s]Diff", s.Keys.Label(config.ActionPanelPrev), s.Keys.PageLabel(), s.Keys.ReloadLabel(), s.Keys.DiffLabel())
	case !s.DiffMode && s.Focus == FocusDiffContent:
		modeSpecific = fmt.Sprintf("[Overview] [%s]Prev Panel [%s]Diff", s.Keys.Label(config.ActionPanelPrev), s.Keys.DiffLabel())
	case !s.DiffMode:
		modeSpecific = fmt.Sprintf("[%s]Next Panel [%s]Diff", s.Keys.Label(config.ActionPanelNext), s.Keys.DiffLabel())
	case s.Focus == FocusRepo:
		modeSpecific = fmt.Sprintf("[%s]Focus [Repo] [%s]Next Panel [%s]Diff", s.Keys.FocusLabel(), s.Keys.Label(config.ActionPanelNext), s.Keys.DiffLabel())
	case s.Focus == FocusPRs && s.HasPR:
		modeSpecific = fmt.Sprintf("[%s]Focus [PRs] [%s]Prev/Next Panel [%s]Move [%s]Review", s.Keys.FocusLabel(), s.Keys.PanelLabel(), s.Keys.MoveLabel(), s.Keys.ReviewModeLabel())
	case s.Focus == FocusDiffFiles && s.HasFiles:
		modeSpecific = fmt.Sprintf("[%s]Focus [Files] [%s]Move [%s]Prev/Next Panel [%s]Overview [%s]Range [%s]Comment", s.Keys.FocusLabel(), s.Keys.MoveLabel(), s.Keys.PanelLabel(), s.Keys.OverviewLabel(), s.Keys.RangeLabel(), s.Keys.CommentLabel())
	case s.Focus == FocusReviewDrawer && s.InputMode == core.ReviewInputComment:
		modeSpecific = fmt.Sprintf("[%s]Save Comment [%s]Cancel [%s]Submit [%s]Discard", s.Keys.SaveLabel(), s.Keys.CancelLabel(), s.Keys.SubmitLabel(), s.Keys.DiscardLabel())
	case s.Focus == FocusReviewDrawer && s.InputMode == core.ReviewInputSummary:
		modeSpecific = fmt.Sprintf("[%s]Save Summary [%s]Cancel [%s]Submit [%s]Discard", s.Keys.SaveLabel(), s.Keys.CancelLabel(), s.Keys.SubmitLabel(), s.Keys.DiscardLabel())
	case s.Focus == FocusReviewDrawer:
		modeSpecific = fmt.Sprintf("[%s]Prev Panel [%s]Comment [%s]Summary [%s]Submit [%s]Discard [%s]Diff", s.Keys.Label(config.ActionPanelPrev), s.Keys.CommentLabel(), s.Keys.SummaryLabel(), s.Keys.SubmitLabel(), s.Keys.DiscardLabel(), s.Keys.CancelLabel())
	case s.HasPR || s.HasFiles || s.HasReviewDrawer:
		modeSpecific = fmt.Sprintf("[%s]Focus [Diff] [%s]Line [%s]Page [%s]Top/Bottom [%s]Prev/Next Panel [%s]Range [%s]Comment [%s]Summary [%s]Submit [%s]Discard [%s]Overview", s.Keys.FocusLabel(), s.Keys.MoveLabel(), s.Keys.PageLabel(), s.Keys.TopBottomLabel(), s.Keys.PanelLabel(), s.Keys.RangeLabel(), s.Keys.CommentLabel(), s.Keys.SummaryLabel(), s.Keys.SubmitLabel(), s.Keys.DiscardLabel(), s.Keys.OverviewLabel())
	default:
		modeSpecific = fmt.Sprintf("[%s]Overview", s.Keys.OverviewLabel())
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
