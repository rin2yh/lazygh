package review

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/model"
)

// FocusTarget identifies which UI panel should receive focus.
type FocusTarget int

const (
	FocusDiffContent FocusTarget = iota
	FocusReviewDrawer
)

// Selection provides the currently selected diff line/file to the review workflow.
type Selection interface {
	CurrentFile() (gh.DiffFile, bool)
	CurrentLine() (gh.DiffLine, bool)
	LineSelected() int
}

// AppState is the minimal interface the review package needs from the host
// application state (list/detail state).
type AppState interface {
	SelectedPR() (model.Item, bool)
	ListRepo() string
	BeginFetchReview()
	ClearFetching()
	IsDiffMode() bool
}

// Controller orchestrates the pending-review workflow and directly owns
// ReviewState (no *state.State reference).
type Controller struct {
	rs       *ReviewState
	keys     config.KeyBindings
	comment  *comment
	summary  *summary
	rng      *rangeState
	pending  *pending
	view     *view
	setFocus func(FocusTarget)
}

// NewController creates a Controller. host provides list/detail context;
// client handles GitHub API calls.
func NewController(cfg *config.Config, host AppState, client PendingReviewClient, selection Selection, setFocus func(FocusTarget)) *Controller {
	rs := newReviewState()
	c := newComment(cfg, rs)
	s := newSummary(rs)
	rng := newRange(rs, selection)
	v := newView(rs, host, c, s)
	return &Controller{
		rs:       rs,
		keys:     cfg.KeyBindings,
		comment:  c,
		summary:  s,
		rng:      rng,
		pending:  newPending(rs, host, client, selection, c, s),
		view:     v,
		setFocus: setFocus,
	}
}

// --- state accessors for the gui layer ---

func (c *Controller) InputMode() InputMode    { return c.rs.InputMode }
func (c *Controller) Summary() string         { return c.rs.Summary }
func (c *Controller) EventLabel() string      { return c.rs.Event.Label() }
func (c *Controller) Notice() string          { return c.rs.Notice }
func (c *Controller) RangeStart() *Range      { return c.rs.RangeStart }
func (c *Controller) Comments() []Comment     { return c.rs.Comments }
func (c *Controller) SelectedCommentIdx() int { return c.rs.SelectedCommentIdx }
func (c *Controller) HasRangeStart() bool     { return c.rs.RangeStart != nil }
func (c *Controller) IsInInputMode() bool     { return c.rs.InputMode != InputNone }
func (c *Controller) HasPendingReview() bool  { return c.rs.HasPendingReview() }
func (c *Controller) PRNumber() int           { return c.rs.PRNumber }
func (c *Controller) SetNotice(msg string)    { c.rs.SetNotice(msg) }
func (c *Controller) ClearRangeStart()        { c.rs.ClearRangeStart() }

// Reset clears review state (called when the PR list reloads).
func (c *Controller) Reset() { c.rs.Reset() }

// SetContext sets the pending review context (PR number, IDs).
func (c *Controller) SetContext(prNumber int, pullRequestID, commitOID, reviewID string) {
	c.rs.SetContext(prNumber, pullRequestID, commitOID, reviewID)
}

// OpenDrawer opens the review drawer.
func (c *Controller) OpenDrawer() { c.rs.OpenDrawer() }

// BeginCommentInput puts the drawer into comment input mode.
func (c *Controller) BeginCommentInput() { c.rs.BeginCommentInput() }

// --- view ---

func (c *Controller) ShouldShowDrawer() bool {
	return c.view.ShouldShowDrawer()
}

// --- comment editor ---

func (c *Controller) CommentValue() string {
	return c.comment.CurrentValue()
}

func (c *Controller) SetCommentValue(value string) {
	c.comment.SetValue(value)
}

func (c *Controller) SummaryValue() string {
	return c.summary.CurrentValue()
}

func (c *Controller) CommentInputLines() []string {
	return c.comment.InputLines()
}

func (c *Controller) SummaryInputLines() []string {
	return c.summary.InputLines()
}

func (c *Controller) BeginSummaryInput() {
	c.summary.BeginInput()
	c.setFocus(FocusReviewDrawer)
}

func (c *Controller) StopInput() {
	if t, ok := c.view.StopInput(); ok {
		c.setFocus(t)
	}
}

func (c *Controller) ClearCommentInput() {
	c.comment.Clear()
}

func (c *Controller) EditorKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.Type {
	case tea.KeyEsc:
		handled := c.view.HandleEsc()
		if handled {
			c.setFocus(FocusDiffContent)
		}
		return nil, handled
	}
	if c.keys.Matches(msg, config.ActionReviewSave) && c.view.InputMode() == InputSummary {
		if t, ok := c.view.HandleSummarySave(); ok {
			c.setFocus(t)
		}
		return nil, true
	}

	switch c.view.InputMode() {
	case InputComment:
		return c.comment.HandleKey(msg)
	case InputSummary:
		return c.summary.HandleKey(msg)
	default:
		return nil, false
	}
}

// --- range ---

func (c *Controller) ToggleRangeSelection() {
	if c.rng.ToggleSelection() {
		c.setFocus(FocusDiffContent)
	}
}

func (c *Controller) IsIndexWithinPendingRange(path string, commentable bool, idx int) bool {
	return c.rng.IsIndexWithinPendingRange(path, commentable, idx)
}

// --- pending review actions ---

func (c *Controller) BuildCommentDraft(body string) (gh.ReviewComment, error) {
	return c.comment.BuildDraft(body, c.rng.RangeStart())
}

func (c *Controller) SaveComment() tea.Cmd {
	return c.pending.HandleCommentSave()
}

func (c *Controller) Submit() tea.Cmd {
	return c.pending.HandleSubmit()
}

func (c *Controller) Discard() tea.Cmd {
	return c.pending.HandleDiscard()
}

func (c *Controller) CommentResult(msg CommentSavedMsg) {
	if c.pending.ApplyCommentResult(msg) {
		c.setFocus(FocusReviewDrawer)
	}
}

func (c *Controller) SubmitResult(msg SubmittedMsg) {
	c.pending.ApplySubmitResult(msg)
	if msg.Err == nil {
		c.setFocus(FocusDiffContent)
	}
}

func (c *Controller) DiscardResult(msg DiscardedMsg) {
	c.pending.ApplyDiscardResult(msg)
	if msg.Err == nil {
		c.setFocus(FocusDiffContent)
	}
}

func (c *Controller) DeleteComment() tea.Cmd {
	return c.pending.HandleDeleteComment()
}

func (c *Controller) EditComment() bool {
	if !c.pending.BeginEditComment() {
		return false
	}
	c.setFocus(FocusDiffContent)
	return true
}

func (c *Controller) SaveEditComment() tea.Cmd {
	return c.pending.HandleEditCommentSave()
}

func (c *Controller) DeleteCommentResult(msg CommentDeletedMsg) {
	c.pending.ApplyDeleteCommentResult(msg)
}

func (c *Controller) EditCommentResult(msg CommentUpdatedMsg) {
	c.pending.ApplyEditCommentResult(msg)
	if msg.Err == nil {
		c.setFocus(FocusReviewDrawer)
	}
}

func (c *Controller) SelectNextComment() {
	c.pending.SelectNextComment()
}

func (c *Controller) SelectPrevComment() {
	c.pending.SelectPrevComment()
}

func (c *Controller) IsEditingComment() bool {
	return c.pending.IsEditingComment()
}

func (c *Controller) CycleReviewEvent() {
	c.rs.CycleEvent()
}

func (c *Controller) BeginCommentFlow() {
	c.comment.BeginInput()
	c.setFocus(FocusReviewDrawer)
}
