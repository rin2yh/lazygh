package review

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

type FocusTarget int

const (
	FocusDiffContent FocusTarget = iota
	FocusReviewDrawer
)

type Selection interface {
	CurrentDiffFile() (gh.DiffFile, bool)
	CurrentDiffLine() (gh.DiffLine, bool)
	CurrentLineIndex() int
}

type Controller struct {
	keys    config.KeyBindings
	comment *comment
	summary *summary
	rng     *rangeState
	pending *pending
	view    *view
}

func NewController(cfg *config.Config, state *core.State, client PendingReviewClient, selection Selection, setFocus func(FocusTarget)) *Controller {
	comment := newComment(cfg, state, setFocus)
	summary := newSummary(state, setFocus)
	rng := newRange(state, selection, setFocus)
	view := newView(state, setFocus, comment, summary)
	return &Controller{
		keys:    cfg.KeyBindings,
		comment: comment,
		summary: summary,
		rng:     rng,
		pending: newPending(state, client, selection, setFocus, comment, summary),
		view:    view,
	}
}

func (c *Controller) ShouldShowDrawer() bool {
	return c.view.ShouldShowDrawer()
}

func (c *Controller) CurrentCommentValue() string {
	return c.comment.CurrentValue()
}

func (c *Controller) SetCommentValue(value string) {
	c.comment.SetValue(value)
}

func (c *Controller) CurrentSummaryValue() string {
	return c.summary.CurrentValue()
}

func (c *Controller) CommentInputLines() []string {
	return c.comment.InputLines()
}

func (c *Controller) SummaryInputLines() []string {
	return c.summary.InputLines()
}

func (c *Controller) BeginCommentInput() {
	c.comment.BeginInput()
}

func (c *Controller) BeginSummaryInput() {
	c.summary.BeginInput()
}

func (c *Controller) StopInput() {
	c.view.StopInput()
}

func (c *Controller) HandleEditorKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.Type {
	case tea.KeyEsc:
		return nil, c.view.HandleEsc()
	}
	if c.keys.Matches(msg, config.ActionReviewSave) && c.view.InputMode() == core.ReviewInputSummary {
		return nil, c.view.HandleSummarySave()
	}

	switch c.view.InputMode() {
	case core.ReviewInputComment:
		handled, cmd := c.comment.HandleKey(msg)
		return cmd, handled
	case core.ReviewInputSummary:
		handled, cmd := c.summary.HandleKey(msg)
		return cmd, handled
	default:
		return nil, false
	}
}

func (c *Controller) ToggleRangeSelection() {
	c.rng.ToggleSelection()
}

func (c *Controller) BeginCommentFlow() {
	c.comment.BeginInput()
}

func (c *Controller) ClearCommentInput() {
	c.comment.Clear()
}

func (c *Controller) BuildCommentDraft(body string) (gh.ReviewComment, error) {
	return c.comment.BuildDraft(body, c.rng.RangeStart())
}

func (c *Controller) HandleCommentSave() tea.Cmd {
	return c.pending.HandleCommentSave()
}

func (c *Controller) HandleSubmit() tea.Cmd {
	return c.pending.HandleSubmit()
}

func (c *Controller) HandleDiscard() tea.Cmd {
	return c.pending.HandleDiscard()
}

func (c *Controller) ApplyCommentResult(msg CommentSavedMsg) {
	c.pending.ApplyCommentResult(msg)
}

func (c *Controller) ApplySubmitResult(msg SubmittedMsg) {
	c.pending.ApplySubmitResult(msg)
}

func (c *Controller) ApplyDiscardResult(msg DiscardedMsg) {
	c.pending.ApplyDiscardResult(msg)
}

func (c *Controller) IsLineWithinPendingRange(line gh.DiffLine) bool {
	return c.rng.IsLineWithinPendingRange(line)
}
