package review

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
)

// Controller orchestrates the pending-review workflow and directly owns
// ReviewState (no *state.State reference).
type Controller struct {
	rs           *ReviewState
	isDiffMode   func() bool
	keys         config.KeyBindings
	comment      *comment
	summary      *summary
	rng          *rangeState
	pending      *pending
	view         *view
	setFocus     func(FocusTarget)
	threadClient ThreadClient
	app          AppState
	threadReply  *threadReply
}

// NewController creates a Controller. app provides list/detail context;
// client handles GitHub API calls.
func NewController(cfg *config.Config, app AppState, client PendingReviewClient, selection Selection, setFocus func(FocusTarget)) *Controller {
	rs := newReviewState()
	c := newComment(cfg, rs)
	s := newSummary(rs)
	rng := newRange(rs, selection)
	v := newView(rs, app, c, s)
	tr := newThreadReply(cfg, rs)
	return &Controller{
		rs:          rs,
		isDiffMode:  app.IsDiffMode,
		keys:        cfg.KeyBindings,
		comment:     c,
		summary:     s,
		rng:         rng,
		pending:     newPending(rs, app, client, selection, c, s),
		view:        v,
		setFocus:    setFocus,
		app:         app,
		threadReply: tr,
	}
}

// SetThreadClient injects the client used to fetch and reply to review threads.
func (c *Controller) SetThreadClient(tc ThreadClient) {
	c.threadClient = tc
}

// --- Reader interface ---

func (c *Controller) InputMode() InputMode   { return c.rs.InputMode }
func (c *Controller) IsInInputMode() bool    { return c.rs.InputMode != InputNone }
func (c *Controller) RangeStart() *Range     { return c.rs.RangeStart }
func (c *Controller) ShouldShowDrawer() bool { return c.view.ShouldShowDrawer() }

func (c *Controller) IsIndexWithinPendingRange(path string, commentable bool, idx int) bool {
	return c.rng.IsIndexWithinPendingRange(path, commentable, idx)
}

// BuildDrawerInput assembles an Input from the current review state.
// Returns nil when showDrawer is false.
func (c *Controller) BuildDrawerInput(showDrawer bool) *Input {
	if !showDrawer {
		return nil
	}
	inputMode := c.rs.InputMode
	summary := c.rs.Summary
	if inputMode == InputSummary {
		summary = c.summary.Text()
	}
	input := &Input{
		SummaryLines:     splitLines(summary),
		CommentModeLabel: CommentModeSingleLine,
		EventLabel:       c.rs.Event.Label(),
		Notice:           c.rs.Notice,
	}
	if rs := c.rs.RangeStart; rs != nil {
		input.CommentModeLabel = CommentModeRangeSelecting
		input.RangeStart = &DrawerRange{Path: rs.Path, Line: rs.Line}
		input.AnchorConflict = c.rng.HasConflict()
	}
	comments := c.rs.Comments
	input.Comments = make([]DrawerComment, 0, len(comments))
	for _, comment := range comments {
		input.Comments = append(input.Comments, DrawerComment{
			Path:      comment.Path,
			Line:      comment.Line,
			StartLine: comment.StartLine,
			Body:      comment.Body,
			Stale:     comment.Stale,
		})
	}
	input.SelectedCommentIdx = c.rs.SelectedCommentIdx
	if inputMode == InputComment {
		input.CommentInputLines = c.comment.InputLines()
	}
	if inputMode == InputSummary {
		input.SummaryInputLines = c.summary.Lines()
	}
	threads := c.rs.Threads
	input.Threads = make([]DrawerThread, 0, len(threads))
	for _, t := range threads {
		dt := DrawerThread{
			Path:       t.Path,
			Line:       t.Line,
			DiffSide:   t.DiffSide,
			IsResolved: t.IsResolved,
			IsOutdated: t.IsOutdated,
			Comments:   make([]DrawerThreadComment, 0, len(t.Comments)),
		}
		for _, tc := range t.Comments {
			dt.Comments = append(dt.Comments, DrawerThreadComment{
				Author: tc.Author,
				Body:   tc.Body,
			})
		}
		input.Threads = append(input.Threads, dt)
	}
	input.SelectedThreadIdx = c.rs.SelectedThreadIdx
	if inputMode == InputThreadReply {
		input.ThreadReplyLines = c.threadReply.InputLines()
	}
	return input
}

// --- ReviewHook (coordinator.ReviewHook) ---

func (c *Controller) HasPendingReview() bool { return c.rs.HasPendingReview() }
func (c *Controller) PRNumber() int          { return c.rs.PRNumber }
func (c *Controller) Reset()                 { c.rs.Reset() }

// --- Handler interface ---

// Notify shows a notice message in the review drawer.
func (c *Controller) Notify(msg string) { c.rs.Notify(msg) }

// HandleCancel handles review-specific cancel logic, returning true if consumed.
// Clears an active range selection or stops an active input mode.
// Focus is moved to FocusDiffContent via the injected setFocus callback.
func (c *Controller) HandleCancel() bool {
	if c.rs.InputMode == InputNone && c.rs.RangeStart != nil {
		c.rs.ClearRangeStart()
		c.rs.Notify("Range selection cleared.")
		c.setFocus(FocusDiffContent)
		return true
	}
	if c.rs.InputMode == InputThreadReply {
		c.threadReply.StopInput()
		c.rs.StopInput()
		c.setFocus(FocusReviewDrawer)
		return true
	}
	if c.rs.InputMode != InputNone {
		c.view.StopInput()
		c.setFocus(FocusDiffContent)
		return true
	}
	return false
}

func (c *Controller) HandleInputKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	action, ok := c.keys.ActionFor(msg)
	if ok {
		switch action {
		case config.ActionReviewSubmit:
			return c.pending.HandleSubmit(), true
		case config.ActionReviewDiscard:
			return c.pending.HandleDiscard(), true
		case config.ActionReviewSave:
			if c.rs.InputMode == InputComment {
				if c.pending.IsEditingComment() {
					return c.pending.HandleEditCommentSave(), true
				}
				return c.pending.HandleCommentSave(), true
			}
			if c.rs.InputMode == InputThreadReply {
				return c.handleThreadReplySave(), true
			}
		}
	}
	if cmd, handled := c.editorKey(msg); handled {
		return cmd, true
	}
	return nil, false
}

func (c *Controller) HandleAction(action config.Action) tea.Cmd {
	switch action {
	case config.ActionReviewRange:
		return c.requireDiffMode("Review range selection is only available in diff view.", c.toggleRangeSelection)
	case config.ActionReviewComment:
		return c.requireDiffMode("Review comments are only available in diff view.", c.beginCommentFlow)
	case config.ActionReviewSummary:
		return c.requireDiffMode("Review summary is only available in diff view.", c.beginSummaryInput)
	case config.ActionReviewSubmit:
		return c.pending.HandleSubmit()
	case config.ActionReviewDiscard:
		return c.pending.HandleDiscard()
	case config.ActionReviewClearComment:
		if c.rs.InputMode == InputComment {
			c.comment.Clear()
		}
	case config.ActionReviewEvent:
		if c.isDiffMode() {
			c.rs.CycleEvent()
		}
	case config.ActionReviewDeleteComment:
		return c.pending.HandleDeleteComment()
	case config.ActionReviewEditComment:
		c.editComment()
	case config.ActionReviewReplyThread:
		c.beginThreadReply()
	}
	return nil
}

func (c *Controller) SelectNextComment() { c.pending.SelectNextComment() }
func (c *Controller) SelectPrevComment() { c.pending.SelectPrevComment() }
func (c *Controller) SelectNextThread()  { c.rs.SelectNextThread() }
func (c *Controller) SelectPrevThread()  { c.rs.SelectPrevThread() }

// --- Applier interface ---

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

func (c *Controller) DeleteCommentResult(msg CommentDeletedMsg) {
	c.pending.ApplyDeleteCommentResult(msg)
}

func (c *Controller) EditCommentResult(msg CommentUpdatedMsg) {
	c.pending.ApplyEditCommentResult(msg)
	if msg.Err == nil {
		c.setFocus(FocusReviewDrawer)
	}
}

func (c *Controller) ThreadsResult(msg ThreadsLoadedMsg) {
	c.app.ClearFetching()
	if msg.Err != nil {
		c.rs.Notify(msg.Err.Error())
		return
	}
	c.rs.LoadThreads(msg.Threads)
}

func (c *Controller) ThreadReplyResult(msg ThreadReplyMsg) {
	c.app.ClearFetching()
	c.threadReply.StopInput()
	c.rs.StopInput()
	if msg.Err != nil {
		c.rs.Notify(msg.Err.Error())
		return
	}
	c.rs.Notify("Reply posted.")
	c.setFocus(FocusReviewDrawer)
}

// MarkStaleComments marks pending comments whose anchor position no longer
// exists in files as stale. Call this whenever the diff is refreshed.
func (c *Controller) MarkStaleComments(files []gh.DiffFile) {
	if len(c.rs.Comments) == 0 || len(files) == 0 {
		return
	}
	type pos struct {
		side gh.DiffSide
		line int
	}
	valid := make(map[string]map[pos]bool, len(files))
	for _, file := range files {
		m := make(map[pos]bool)
		for _, l := range file.Lines {
			if !l.Commentable {
				continue
			}
			if l.NewLine > 0 {
				m[pos{gh.DiffSideRight, l.NewLine}] = true
			}
			if l.OldLine > 0 {
				m[pos{gh.DiffSideLeft, l.OldLine}] = true
			}
		}
		valid[file.Path] = m
	}
	for i := range c.rs.Comments {
		comment := &c.rs.Comments[i]
		m, ok := valid[comment.Path]
		if !ok {
			comment.Stale = true
			continue
		}
		comment.Stale = !m[pos{gh.DiffSide(comment.Side), comment.Line}]
	}
}

// --- Lifecycle (not in interfaces; accessed via concrete type in tests) ---

// SetContext sets the pending review context (PR number, IDs).
func (c *Controller) SetContext(prNumber int, pullRequestID, commitOID, reviewID string) {
	c.rs.SetContext(prNumber, pullRequestID, commitOID, reviewID)
}

// OpenDrawer opens the review drawer.
func (c *Controller) OpenDrawer() { c.rs.OpenDrawer() }

// BeginCommentInput puts the drawer into comment input mode.
func (c *Controller) BeginCommentInput() { c.rs.BeginCommentInput() }

// --- Test support (not in interfaces) ---

func (c *Controller) CommentValue() string         { return c.comment.CurrentValue() }
func (c *Controller) SetCommentValue(value string) { c.comment.SetValue(value) }
func (c *Controller) SummaryValue() string         { return c.summary.Text() }

// --- internal ---

func (c *Controller) editorKey(msg tea.KeyMsg) (tea.Cmd, bool) {
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
	case InputThreadReply:
		return c.threadReply.HandleKey(msg)
	default:
		return nil, false
	}
}

func (c *Controller) toggleRangeSelection() {
	if c.rng.ToggleSelection() {
		c.setFocus(FocusDiffContent)
	}
}

func (c *Controller) beginCommentFlow() {
	c.comment.BeginInput()
	c.setFocus(FocusReviewDrawer)
}

func (c *Controller) beginSummaryInput() {
	c.summary.BeginInput()
	c.setFocus(FocusReviewDrawer)
}

func (c *Controller) editComment() bool {
	if !c.pending.BeginEditComment() {
		return false
	}
	c.setFocus(FocusDiffContent)
	return true
}

func (c *Controller) saveComment() tea.Cmd { return c.pending.HandleCommentSave() }
func (c *Controller) submit() tea.Cmd      { return c.pending.HandleSubmit() }

func (c *Controller) beginThreadReply() {
	if _, ok := c.rs.SelectedThread(); !ok {
		c.rs.Notify("No thread selected.")
		return
	}
	c.threadReply.BeginInput()
	c.setFocus(FocusReviewDrawer)
}

func (c *Controller) handleThreadReplySave() tea.Cmd {
	thread, ok := c.rs.SelectedThread()
	if !ok {
		c.rs.Notify("No thread selected.")
		return nil
	}
	body := strings.TrimSpace(c.threadReply.CurrentValue())
	if body == "" {
		c.rs.Notify("Reply body is empty.")
		return nil
	}
	if c.threadClient == nil {
		c.rs.Notify("Thread client not available.")
		return nil
	}
	threadID := thread.ID
	c.app.BeginFetchReview()
	return func() tea.Msg {
		err := c.threadClient.AddReplyToReviewThread(threadID, body)
		return ThreadReplyMsg{Err: err}
	}
}

func (c *Controller) requireDiffMode(notice string, fn func()) tea.Cmd {
	if !c.isDiffMode() {
		c.rs.Notify(notice)
		return nil
	}
	fn()
	return nil
}

func splitLines(content string) []string {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	return strings.Split(content, "\n")
}
