package state

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/core"
)

type ReviewState struct {
	PRNumber           int
	PullRequestID      string
	CommitOID          string
	ReviewID           string
	DrawerOpen         bool
	InputMode          core.ReviewInputMode
	Event              core.ReviewEvent
	Summary            string
	Comments           []core.ReviewComment
	RangeStart         *core.ReviewRange
	Notice             string
	SelectedCommentIdx int
	EditingCommentIdx  int
}

// ListState holds PR list and selection state.
type ListState struct {
	Repo         string
	PRs          []core.Item
	PRsLoading   bool
	PRsSelected  int
	Filter       core.PRFilterMask
	FilterOpen   bool
	FilterCursor int
}

// DetailState holds detail panel display and loading state.
type DetailState struct {
	Mode    core.DetailMode
	Content string
	Loading core.LoadingKind
}

type EnterAction struct {
	Kind   core.EnterActionKind
	Repo   string
	Number int
}

type State struct {
	List   ListState
	Detail DetailState
	Review ReviewState

	Width  int
	Height int
}

func NewState() *State {
	return &State{
		List: ListState{
			PRs:    []core.Item{},
			Filter: core.PRFilterOpen,
		},
		Detail: DetailState{
			Mode: core.DetailModeOverview,
		},
		Review: ReviewState{
			Comments:          []core.ReviewComment{},
			EditingCommentIdx: core.NoEditingComment,
		},
	}
}

func (s *State) SetWindowSize(width int, height int) {
	s.Width = width
	s.Height = height
}

func (s *State) BeginLoadPRs() {
	s.List.PRsLoading = true
	s.Detail.Loading = core.LoadingPRs
}

// BeginReviewLoad marks a review operation as in-progress.
func (s *State) BeginReviewLoad() {
	s.Detail.Loading = core.LoadingReview
}

// ClearLoading clears any in-progress loading indicator.
func (s *State) ClearLoading() {
	s.Detail.Loading = core.LoadingNone
}

// StopReviewInput exits any active review input mode without closing the drawer.
func (s *State) StopReviewInput() {
	s.Review.InputMode = core.ReviewInputNone
}

func (s *State) ApplyPRsResult(repo string, prs []core.Item, err error) {
	s.List.PRsLoading = false
	s.Detail.Loading = core.LoadingNone
	if err != nil {
		s.showError("Error loading PRs", err)
		return
	}

	s.List.Repo = repo
	s.List.PRs = prs
	s.List.PRsSelected = 0
	s.Detail.Mode = core.DetailModeOverview
	s.resetReview()
	if len(prs) == 0 {
		s.Detail.Content = "No pull requests"
		return
	}
	s.Detail.Content = core.FormatPROverview(prs[s.List.PRsSelected])
}

func (s *State) ApplyDetailResult(content string, err error) {
	if err != nil {
		s.showError("Error loading detail", err)
		return
	}
	s.Detail.Loading = core.LoadingNone
	s.Detail.Content = core.SanitizeMultiline(content)
}

func (s *State) ApplyDiffResult(content string, err error) {
	if err != nil {
		s.showError("Error loading diff", err)
		return
	}
	s.Detail.Loading = core.LoadingNone
	s.Detail.Content = core.SanitizeMultiline(content)
}

func (s *State) NavigateDown() bool {
	if s.blocksPRSelectionChange() {
		s.Review.Notice = "Pending review exists. Submit with S or discard with X."
		return false
	}
	changed := false
	if s.List.PRsSelected < len(s.List.PRs)-1 {
		s.List.PRsSelected++
		changed = true
	}
	if changed && s.Detail.Mode == core.DetailModeOverview {
		s.refreshDetailPreview()
	}
	return changed
}

func (s *State) NavigateUp() bool {
	if s.blocksPRSelectionChange() {
		s.Review.Notice = "Pending review exists. Submit with S or discard with X."
		return false
	}
	changed := false
	if s.List.PRsSelected > 0 {
		s.List.PRsSelected--
		changed = true
	}
	if changed && s.Detail.Mode == core.DetailModeOverview {
		s.refreshDetailPreview()
	}
	return changed
}

func (s *State) SwitchToOverview() bool {
	if s.Detail.Mode == core.DetailModeOverview {
		return false
	}
	s.Detail.Mode = core.DetailModeOverview
	s.Detail.Loading = core.LoadingNone
	s.Review.InputMode = core.ReviewInputNone
	s.refreshDetailPreview()
	return true
}

func (s *State) SwitchToDiff() bool {
	if s.Detail.Mode == core.DetailModeDiff {
		return false
	}
	s.Detail.Mode = core.DetailModeDiff
	s.Detail.Loading = core.LoadingNone
	s.Review.DrawerOpen = false
	s.Review.InputMode = core.ReviewInputNone
	return true
}

func (s *State) IsDiffMode() bool {
	return s.Detail.Mode == core.DetailModeDiff
}

func (s *State) ShouldApplyDetailResult(mode core.DetailMode, number int) bool {
	if s.Detail.Mode != mode {
		return false
	}
	item, ok := s.selectedPR()
	if !ok {
		return false
	}
	return item.Number == number
}

func (s *State) PlanEnter(hasClient bool, forcedDetailText string) EnterAction {
	if !hasClient || s.List.PRsLoading {
		return EnterAction{}
	}
	item, ok := s.selectedPR()
	if !ok {
		return EnterAction{}
	}
	if forcedDetailText != "" {
		s.Detail.Loading = core.LoadingNone
		s.Detail.Content = forcedDetailText
		return EnterAction{}
	}
	s.Detail.Loading = core.LoadingDetail
	if s.Detail.Mode == core.DetailModeDiff {
		return EnterAction{Kind: core.EnterLoadPRDiff, Repo: s.List.Repo, Number: item.Number}
	}
	return EnterAction{Kind: core.EnterLoadPRDetail, Repo: s.List.Repo, Number: item.Number}
}

func (s *State) refreshDetailPreview() {
	item, ok := s.selectedPR()
	if !ok {
		return
	}
	s.Detail.Content = core.FormatPROverview(item)
}

func (s *State) selectedPR() (core.Item, bool) {
	if len(s.List.PRs) == 0 {
		return core.Item{}, false
	}
	if s.List.PRsSelected < 0 || s.List.PRsSelected >= len(s.List.PRs) {
		return core.Item{}, false
	}
	return s.List.PRs[s.List.PRsSelected], true
}

func (s *State) showError(msg string, err error) {
	s.Detail.Loading = core.LoadingNone
	s.Detail.Content = core.SanitizeMultiline(fmt.Sprintf("%s: %v", msg, err))
}

func (s *State) SelectedPR() (core.Item, bool) {
	return s.selectedPR()
}

func (s *State) OpenReviewDrawer() {
	s.Review.DrawerOpen = true
}

func (s *State) CloseReviewDrawer() {
	s.Review.DrawerOpen = false
	s.Review.InputMode = core.ReviewInputNone
	s.Review.Notice = ""
}

func (s *State) BeginReviewCommentInput() {
	s.Review.DrawerOpen = true
	s.Review.InputMode = core.ReviewInputComment
	s.Review.Notice = ""
}

func (s *State) BeginReviewSummaryInput() {
	s.Review.DrawerOpen = true
	s.Review.InputMode = core.ReviewInputSummary
	s.Review.Notice = ""
}

func (s *State) SetReviewSummary(summary string) {
	s.Review.Summary = core.SanitizeMultiline(summary)
}

func (s *State) SetReviewContext(prNumber int, pullRequestID string, commitOID string, reviewID string) {
	s.Review.PRNumber = prNumber
	s.Review.PullRequestID = core.SanitizeSingleLine(pullRequestID)
	s.Review.CommitOID = core.SanitizeSingleLine(commitOID)
	s.Review.ReviewID = core.SanitizeSingleLine(reviewID)
}

func (s *State) AddReviewComment(comment core.ReviewComment) {
	s.Review.Comments = append(s.Review.Comments, core.ReviewComment{
		CommentID: comment.CommentID,
		Path:      core.SanitizeSingleLine(comment.Path),
		Body:      core.SanitizeMultiline(comment.Body),
		Side:      core.SanitizeSingleLine(comment.Side),
		Line:      comment.Line,
		StartSide: core.SanitizeSingleLine(comment.StartSide),
		StartLine: comment.StartLine,
	})
	s.Review.SelectedCommentIdx = len(s.Review.Comments) - 1
	s.Review.Notice = "Review comment added."
	s.Review.DrawerOpen = true
	s.Review.InputMode = core.ReviewInputNone
	s.Review.RangeStart = nil
}

func (s *State) SelectNextComment() {
	if len(s.Review.Comments) == 0 {
		return
	}
	if s.Review.SelectedCommentIdx < len(s.Review.Comments)-1 {
		s.Review.SelectedCommentIdx++
	}
}

func (s *State) SelectPrevComment() {
	if s.Review.SelectedCommentIdx > 0 {
		s.Review.SelectedCommentIdx--
	}
}

func (s *State) DeleteSelectedComment() (core.ReviewComment, bool) {
	idx := s.Review.SelectedCommentIdx
	if idx < 0 || idx >= len(s.Review.Comments) {
		return core.ReviewComment{}, false
	}
	deleted := s.Review.Comments[idx]
	s.Review.Comments = append(s.Review.Comments[:idx], s.Review.Comments[idx+1:]...)
	if len(s.Review.Comments) == 0 {
		s.Review.SelectedCommentIdx = 0
	} else if s.Review.SelectedCommentIdx >= len(s.Review.Comments) {
		s.Review.SelectedCommentIdx = len(s.Review.Comments) - 1
	}
	return deleted, true
}

func (s *State) SelectedComment() (core.ReviewComment, bool) {
	idx := s.Review.SelectedCommentIdx
	if idx < 0 || idx >= len(s.Review.Comments) {
		return core.ReviewComment{}, false
	}
	return s.Review.Comments[idx], true
}

func (s *State) BeginEditComment() {
	s.Review.EditingCommentIdx = s.Review.SelectedCommentIdx
	s.Review.InputMode = core.ReviewInputComment
	s.Review.DrawerOpen = true
	s.Review.Notice = ""
}

func (s *State) ApplyEditComment(newBody string) {
	idx := s.Review.EditingCommentIdx
	if idx < 0 || idx >= len(s.Review.Comments) {
		return
	}
	s.Review.Comments[idx].Body = core.SanitizeMultiline(newBody)
	s.Review.EditingCommentIdx = core.NoEditingComment
	s.Review.InputMode = core.ReviewInputNone
	s.Review.Notice = "Comment updated."
}

func (s *State) ClearEditingComment() {
	s.Review.EditingCommentIdx = core.NoEditingComment
}

func (s *State) SetReviewNotice(msg string) {
	s.Review.Notice = core.SanitizeMultiline(msg)
}

func (s *State) ClearReviewNotice() {
	s.Review.Notice = ""
}

func (s *State) HasPendingReview() bool {
	return s.Review.ReviewID != ""
}

func (s *State) MarkReviewRangeStart(anchor core.ReviewRange) {
	copied := anchor
	s.Review.RangeStart = &copied
	s.Review.DrawerOpen = true
	s.Review.Notice = "Range start selected."
}

func (s *State) CycleReviewEvent() {
	switch s.Review.Event {
	case core.ReviewEventComment:
		s.Review.Event = core.ReviewEventApprove
	case core.ReviewEventApprove:
		s.Review.Event = core.ReviewEventRequestChanges
	default:
		s.Review.Event = core.ReviewEventComment
	}
}

func (s *State) ClearReviewRangeStart() {
	s.Review.RangeStart = nil
}

func (s *State) ResetReviewAfterSubmit(notice string) {
	s.resetReview()
	s.Review.Notice = core.SanitizeMultiline(notice)
}

func (s *State) ResetReviewAfterDiscard(notice string) {
	s.resetReview()
	s.Review.Notice = core.SanitizeMultiline(notice)
}

func (s *State) OpenFilterSelect() {
	s.List.FilterOpen = true
	s.List.FilterCursor = 0
}

func (s *State) CloseFilterSelect() {
	s.List.FilterOpen = false
}

func (s *State) MoveFilterCursor(dir int) {
	n := len(core.PRFilterOptions)
	s.List.FilterCursor = (s.List.FilterCursor + dir + n) % n
}

// ToggleFilterAtCursor toggles the filter option under the cursor.
// It prevents deselecting all options (at least one must remain selected).
func (s *State) ToggleFilterAtCursor() {
	if s.List.FilterCursor < 0 || s.List.FilterCursor >= len(core.PRFilterOptions) {
		return
	}
	opt := core.PRFilterOptions[s.List.FilterCursor]
	next := s.List.Filter.Toggle(opt)
	if next == 0 {
		return // disallow empty selection
	}
	s.List.Filter = next
}

func (s *State) blocksPRSelectionChange() bool {
	item, ok := s.selectedPR()
	if !ok {
		return false
	}
	return s.HasPendingReview() && s.Review.PRNumber == item.Number
}

func (s *State) resetReview() {
	s.Review = ReviewState{
		Notice:            s.Review.Notice,
		EditingCommentIdx: core.NoEditingComment,
	}
}
