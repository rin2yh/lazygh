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
	Loading   bool
	DiffMode  bool
	Focus     Focus
	InputMode core.ReviewInputMode
	Keys      config.KeyBindings
}

func (s Status) String() string {
	base := fmt.Sprintf("[%s]Quit [%s]Help", s.Keys.QuitLabel(), s.Keys.HelpLabel())

	var ctx string
	switch {
	case s.InputMode == core.ReviewInputComment:
		ctx = fmt.Sprintf("[%s]Save Comment [%s]Cancel", s.Keys.SaveLabel(), s.Keys.CancelLabel())
	case s.InputMode == core.ReviewInputSummary:
		ctx = fmt.Sprintf("[%s]Save Summary [%s]Cancel", s.Keys.SaveLabel(), s.Keys.CancelLabel())
	case s.Focus == FocusReviewDrawer:
		ctx = fmt.Sprintf("[Review] [%s]Submit [%s]Discard [%s]Cancel", s.Keys.SubmitLabel(), s.Keys.DiscardLabel(), s.Keys.CancelLabel())
	case !s.DiffMode:
		ctx = fmt.Sprintf("[%s]Panels [%s]Diff", s.Keys.PanelLabel(), s.Keys.DiffLabel())
	default:
		ctx = fmt.Sprintf("[%s]Panels [%s]Overview", s.Keys.PanelLabel(), s.Keys.OverviewLabel())
	}

	line := fmt.Sprintf("%s | %s", base, ctx)
	if s.Loading {
		return fmt.Sprintf("Loading... | %s", line)
	}
	return line
}
