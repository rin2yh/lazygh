package state

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/model"
)

type ReviewState struct {
	PRNumber           int
	PullRequestID      string
	CommitOID          string
	ReviewID           string
	DrawerOpen         bool
	InputMode          model.ReviewInputMode
	Event              model.ReviewEvent
	Summary            string
	Comments           []model.ReviewComment
	RangeStart         *model.ReviewRange
	Notice             string
	SelectedCommentIdx int
	EditingCommentIdx  int
}

// ListState holds PR list and selection state.
type ListState struct {
	Repo         string
	PRs          []model.Item
	PRsLoading   bool
	PRsSelected  int
	Filter       model.PRFilterMask
	FilterOpen   bool
	FilterCursor int
}

// DetailState holds detail panel display and loading state.
type DetailState struct {
	Mode    model.DetailMode
	Content string
	Loading model.LoadingKind
}

type EnterAction struct {
	Kind   model.EnterActionKind
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
			PRs:    []model.Item{},
			Filter: model.PRFilterOpen,
		},
		Detail: DetailState{
			Mode: model.DetailModeOverview,
		},
		Review: ReviewState{
			Comments:          []model.ReviewComment{},
			EditingCommentIdx: model.NoEditingComment,
		},
	}
}

func (s *State) SetWindowSize(width int, height int) {
	s.Width = width
	s.Height = height
}

func (s *State) BeginLoadPRs() {
	s.List.PRsLoading = true
	s.Detail.Loading = model.LoadingPRs
}

// BeginReviewLoad marks a review operation as in-progress.
func (s *State) BeginReviewLoad() {
	s.Detail.Loading = model.LoadingReview
}

// ClearLoading clears any in-progress loading indicator.
func (s *State) ClearLoading() {
	s.Detail.Loading = model.LoadingNone
}

// StopReviewInput exits any active review input mode without closing the drawer.
func (s *State) StopReviewInput() {
	s.Review.InputMode = model.ReviewInputNone
}

func (s *State) ApplyPRsResult(repo string, prs []model.Item, err error) {
	s.List.PRsLoading = false
	s.Detail.Loading = model.LoadingNone
	if err != nil {
		s.showError("Error loading PRs", err)
		return
	}

	s.List.Repo = repo
	s.List.PRs = prs
	s.List.PRsSelected = 0
	s.Detail.Mode = model.DetailModeOverview
	s.resetReview()
	if len(prs) == 0 {
		s.Detail.Content = "No pull requests"
		return
	}
	s.Detail.Content = model.FormatPROverview(prs[s.List.PRsSelected])
}

func (s *State) ApplyDetailResult(content string, err error) {
	if err != nil {
		s.showError("Error loading detail", err)
		return
	}
	s.Detail.Loading = model.LoadingNone
	s.Detail.Content = model.SanitizeMultiline(content)
}

func (s *State) ApplyDiffResult(content string, err error) {
	if err != nil {
		s.showError("Error loading diff", err)
		return
	}
	s.Detail.Loading = model.LoadingNone
	s.Detail.Content = model.SanitizeMultiline(content)
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
	if changed && s.Detail.Mode == model.DetailModeOverview {
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
	if changed && s.Detail.Mode == model.DetailModeOverview {
		s.refreshDetailPreview()
	}
	return changed
}

func (s *State) SwitchToOverview() bool {
	if s.Detail.Mode == model.DetailModeOverview {
		return false
	}
	s.Detail.Mode = model.DetailModeOverview
	s.Detail.Loading = model.LoadingNone
	s.Review.InputMode = model.ReviewInputNone
	s.refreshDetailPreview()
	return true
}

func (s *State) SwitchToDiff() bool {
	if s.Detail.Mode == model.DetailModeDiff {
		return false
	}
	s.Detail.Mode = model.DetailModeDiff
	s.Detail.Loading = model.LoadingNone
	s.Review.DrawerOpen = false
	s.Review.InputMode = model.ReviewInputNone
	return true
}

func (s *State) IsDiffMode() bool {
	return s.Detail.Mode == model.DetailModeDiff
}

func (s *State) ShouldApplyDetailResult(mode model.DetailMode, number int) bool {
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
		s.Detail.Loading = model.LoadingNone
		s.Detail.Content = forcedDetailText
		return EnterAction{}
	}
	s.Detail.Loading = model.LoadingDetail
	if s.Detail.Mode == model.DetailModeDiff {
		return EnterAction{Kind: model.EnterLoadPRDiff, Repo: s.List.Repo, Number: item.Number}
	}
	return EnterAction{Kind: model.EnterLoadPRDetail, Repo: s.List.Repo, Number: item.Number}
}

func (s *State) refreshDetailPreview() {
	item, ok := s.selectedPR()
	if !ok {
		return
	}
	s.Detail.Content = model.FormatPROverview(item)
}

func (s *State) selectedPR() (model.Item, bool) {
	if len(s.List.PRs) == 0 {
		return model.Item{}, false
	}
	if s.List.PRsSelected < 0 || s.List.PRsSelected >= len(s.List.PRs) {
		return model.Item{}, false
	}
	return s.List.PRs[s.List.PRsSelected], true
}

func (s *State) showError(msg string, err error) {
	s.Detail.Loading = model.LoadingNone
	s.Detail.Content = model.SanitizeMultiline(fmt.Sprintf("%s: %v", msg, err))
}

func (s *State) SelectedPR() (model.Item, bool) {
	return s.selectedPR()
}

func (s *State) OpenReviewDrawer() {
	s.Review.DrawerOpen = true
}

func (s *State) CloseReviewDrawer() {
	s.Review.DrawerOpen = false
	s.Review.InputMode = model.ReviewInputNone
	s.Review.Notice = ""
}

func (s *State) BeginReviewCommentInput() {
	s.Review.DrawerOpen = true
	s.Review.InputMode = model.ReviewInputComment
	s.Review.Notice = ""
}

func (s *State) BeginReviewSummaryInput() {
	s.Review.DrawerOpen = true
	s.Review.InputMode = model.ReviewInputSummary
	s.Review.Notice = ""
}

func (s *State) SetReviewSummary(summary string) {
	s.Review.Summary = model.SanitizeMultiline(summary)
}

func (s *State) SetReviewContext(prNumber int, pullRequestID string, commitOID string, reviewID string) {
	s.Review.PRNumber = prNumber
	s.Review.PullRequestID = model.SanitizeSingleLine(pullRequestID)
	s.Review.CommitOID = model.SanitizeSingleLine(commitOID)
	s.Review.ReviewID = model.SanitizeSingleLine(reviewID)
}

func (s *State) AddReviewComment(comment model.ReviewComment) {
	s.Review.Comments = append(s.Review.Comments, model.ReviewComment{
		CommentID: comment.CommentID,
		Path:      model.SanitizeSingleLine(comment.Path),
		Body:      model.SanitizeMultiline(comment.Body),
		Side:      model.SanitizeSingleLine(comment.Side),
		Line:      comment.Line,
		StartSide: model.SanitizeSingleLine(comment.StartSide),
		StartLine: comment.StartLine,
	})
	s.Review.SelectedCommentIdx = len(s.Review.Comments) - 1
	s.Review.Notice = "Review comment added."
	s.Review.DrawerOpen = true
	s.Review.InputMode = model.ReviewInputNone
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

func (s *State) DeleteSelectedComment() (model.ReviewComment, bool) {
	idx := s.Review.SelectedCommentIdx
	if idx < 0 || idx >= len(s.Review.Comments) {
		return model.ReviewComment{}, false
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

func (s *State) SelectedComment() (model.ReviewComment, bool) {
	idx := s.Review.SelectedCommentIdx
	if idx < 0 || idx >= len(s.Review.Comments) {
		return model.ReviewComment{}, false
	}
	return s.Review.Comments[idx], true
}

func (s *State) BeginEditComment() {
	s.Review.EditingCommentIdx = s.Review.SelectedCommentIdx
	s.Review.InputMode = model.ReviewInputComment
	s.Review.DrawerOpen = true
	s.Review.Notice = ""
}

func (s *State) ApplyEditComment(newBody string) {
	idx := s.Review.EditingCommentIdx
	if idx < 0 || idx >= len(s.Review.Comments) {
		return
	}
	s.Review.Comments[idx].Body = model.SanitizeMultiline(newBody)
	s.Review.EditingCommentIdx = model.NoEditingComment
	s.Review.InputMode = model.ReviewInputNone
	s.Review.Notice = "Comment updated."
}

func (s *State) ClearEditingComment() {
	s.Review.EditingCommentIdx = model.NoEditingComment
}

func (s *State) SetReviewNotice(msg string) {
	s.Review.Notice = model.SanitizeMultiline(msg)
}

func (s *State) ClearReviewNotice() {
	s.Review.Notice = ""
}

func (s *State) HasPendingReview() bool {
	return s.Review.ReviewID != ""
}

func (s *State) MarkReviewRangeStart(anchor model.ReviewRange) {
	copied := anchor
	s.Review.RangeStart = &copied
	s.Review.DrawerOpen = true
	s.Review.Notice = "Range start selected."
}

func (s *State) CycleReviewEvent() {
	switch s.Review.Event {
	case model.ReviewEventComment:
		s.Review.Event = model.ReviewEventApprove
	case model.ReviewEventApprove:
		s.Review.Event = model.ReviewEventRequestChanges
	default:
		s.Review.Event = model.ReviewEventComment
	}
}

func (s *State) ClearReviewRangeStart() {
	s.Review.RangeStart = nil
}

func (s *State) ResetReviewAfterSubmit(notice string) {
	s.resetReview()
	s.Review.Notice = model.SanitizeMultiline(notice)
}

func (s *State) ResetReviewAfterDiscard(notice string) {
	s.resetReview()
	s.Review.Notice = model.SanitizeMultiline(notice)
}

func (s *State) OpenFilterSelect() {
	s.List.FilterOpen = true
	s.List.FilterCursor = 0
}

func (s *State) CloseFilterSelect() {
	s.List.FilterOpen = false
}

func (s *State) MoveFilterCursor(dir int) {
	n := len(model.PRFilterOptions)
	s.List.FilterCursor = (s.List.FilterCursor + dir + n) % n
}

// ToggleFilterAtCursor toggles the filter option under the cursor.
// It prevents deselecting all options (at least one must remain selected).
func (s *State) ToggleFilterAtCursor() {
	if s.List.FilterCursor < 0 || s.List.FilterCursor >= len(model.PRFilterOptions) {
		return
	}
	opt := model.PRFilterOptions[s.List.FilterCursor]
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
		EditingCommentIdx: model.NoEditingComment,
	}
}
