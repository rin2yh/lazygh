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
	BeginReviewLoad()
	ClearLoading()
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
	c := newComment(cfg, rs, setFocus)
	s := newSummary(rs, setFocus)
	rng := newRange(rs, selection, setFocus)
	v := newView(rs, host, setFocus, c, s)
	return &Controller{
		rs:       rs,
		keys:     cfg.KeyBindings,
		comment:  c,
		summary:  s,
		rng:      rng,
		pending:  newPending(rs, host, client, selection, setFocus, c, s),
		view:     v,
		setFocus: setFocus,
	}
}

// --- state accessors for the gui layer ---

func (c *Controller) InputMode() model.ReviewInputMode { return c.rs.InputMode }
func (c *Controller) Summary() string                  { return c.rs.Summary }
func (c *Controller) EventLabel() string               { return c.rs.Event.Label() }
func (c *Controller) Notice() string                   { return c.rs.Notice }
func (c *Controller) RangeStart() *model.ReviewRange   { return c.rs.RangeStart }
func (c *Controller) Comments() []model.ReviewComment  { return c.rs.Comments }
func (c *Controller) SelectedCommentIdx() int          { return c.rs.SelectedCommentIdx }
func (c *Controller) HasRangeStart() bool              { return c.rs.RangeStart != nil }
func (c *Controller) IsInInputMode() bool              { return c.rs.InputMode != model.ReviewInputNone }
func (c *Controller) HasPendingReview() bool           { return c.rs.HasPendingReview() }
func (c *Controller) PRNumber() int                    { return c.rs.PRNumber }
func (c *Controller) SetNotice(msg string)             { c.rs.SetNotice(msg) }
func (c *Controller) ClearRangeStart()                 { c.rs.ClearRangeStart() }

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

func (c *Controller) BeginSummaryInput() {
	c.summary.BeginInput()
}

func (c *Controller) StopInput() {
	c.view.StopInput()
}

func (c *Controller) ClearCommentInput() {
	c.comment.Clear()
}

func (c *Controller) HandleEditorKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.Type {
	case tea.KeyEsc:
		return nil, c.view.HandleEsc()
	}
	if c.keys.Matches(msg, config.ActionReviewSave) && c.view.InputMode() == model.ReviewInputSummary {
		return nil, c.view.HandleSummarySave()
	}

	switch c.view.InputMode() {
	case model.ReviewInputComment:
		return c.comment.HandleKey(msg)
	case model.ReviewInputSummary:
		return c.summary.HandleKey(msg)
	default:
		return nil, false
	}
}

// --- range ---

func (c *Controller) ToggleRangeSelection() {
	c.rng.ToggleSelection()
}

func (c *Controller) IsIndexWithinPendingRange(path string, commentable bool, idx int) bool {
	return c.rng.IsIndexWithinPendingRange(path, commentable, idx)
}

// --- pending review actions ---

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

func (c *Controller) HandleDeleteComment() tea.Cmd {
	return c.pending.HandleDeleteComment()
}

func (c *Controller) BeginEditComment() bool {
	if !c.pending.BeginEditComment() {
		return false
	}
	c.setFocus(FocusDiffContent)
	return true
}

func (c *Controller) HandleEditCommentSave() tea.Cmd {
	return c.pending.HandleEditCommentSave()
}

func (c *Controller) ApplyDeleteCommentResult(msg CommentDeletedMsg) {
	c.pending.ApplyDeleteCommentResult(msg)
}

func (c *Controller) ApplyEditCommentResult(msg CommentUpdatedMsg) {
	c.pending.ApplyEditCommentResult(msg)
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
}
