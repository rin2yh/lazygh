package layout

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/model"
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
	Fetching  bool
	DiffMode  bool
	Focus     Focus
	InputMode model.ReviewInputMode
	Keys      config.KeyBindings
}

func (s Status) String() string {
	base := fmt.Sprintf("[%s]Quit [%s]Help", s.Keys.PrimaryLabel(config.ActionQuit), s.Keys.PrimaryLabel(config.ActionShowHelp))

	var ctx string
	switch {
	case s.InputMode == model.ReviewInputComment:
		ctx = fmt.Sprintf("[%s]Save Comment [%s]Cancel", s.Keys.PrimaryLabel(config.ActionReviewSave), s.Keys.PrimaryLabel(config.ActionCancel))
	case s.InputMode == model.ReviewInputSummary:
		ctx = fmt.Sprintf("[%s]Save Summary [%s]Cancel", s.Keys.PrimaryLabel(config.ActionReviewSave), s.Keys.PrimaryLabel(config.ActionCancel))
	case s.Focus == FocusReviewDrawer:
		ctx = fmt.Sprintf("[Review] [%s]Submit [%s]Discard [%s]Cancel", s.Keys.PrimaryLabel(config.ActionReviewSubmit), s.Keys.PrimaryLabel(config.ActionReviewDiscard), s.Keys.PrimaryLabel(config.ActionCancel))
	case !s.DiffMode:
		ctx = fmt.Sprintf("[%s]Panels [%s]Diff", s.Keys.PanelLabel(), s.Keys.PrimaryLabel(config.ActionShowDiff))
	default:
		ctx = fmt.Sprintf("[%s]Panels [%s]Overview", s.Keys.PanelLabel(), s.Keys.PrimaryLabel(config.ActionShowOverview))
	}

	line := fmt.Sprintf("%s | %s", base, ctx)
	if s.Fetching {
		return fmt.Sprintf("Fetching... | %s", line)
	}
	return line
}
