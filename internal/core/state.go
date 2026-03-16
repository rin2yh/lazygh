package core

import (
	"fmt"
	"strings"
)

type Item struct {
	Number    int
	Title     string
	Status    string
	Assignees []string
}

type DetailMode int

const (
	DetailModeOverview DetailMode = iota
	DetailModeDiff
)

type EnterActionKind int

const (
	EnterNone EnterActionKind = iota
	EnterLoadPRDetail
	EnterLoadPRDiff
)

type LoadingKind int

const (
	LoadingNone LoadingKind = iota
	LoadingPRs
	LoadingDetail
	LoadingReview
)

type ReviewInputMode int

const (
	ReviewInputNone ReviewInputMode = iota
	ReviewInputComment
	ReviewInputSummary
)

type ReviewEvent int

const (
	ReviewEventComment ReviewEvent = iota
	ReviewEventApprove
	ReviewEventRequestChanges
)

func (e ReviewEvent) Label() string {
	switch e {
	case ReviewEventApprove:
		return "APPROVE"
	case ReviewEventRequestChanges:
		return "REQUEST CHANGES"
	default:
		return "COMMENT"
	}
}

type ReviewComment struct {
	Path      string
	Body      string
	Side      string
	Line      int
	StartSide string
	StartLine int
}

type ReviewRange struct {
	Path      string
	Index     int
	Side      string
	Line      int
	StartSide string
	StartLine int
}

type ReviewState struct {
	PRNumber      int
	PullRequestID string
	CommitOID     string
	ReviewID      string
	DrawerOpen    bool
	InputMode     ReviewInputMode
	Event         ReviewEvent
	Summary       string
	Comments      []ReviewComment
	RangeStart    *ReviewRange
	Notice        string
}

type PRFilterMask uint8

const (
	PRFilterOpen   PRFilterMask = 1 << iota // 1
	PRFilterClosed                          // 2
	PRFilterMerged                          // 4
)

// PRFilterOptions lists the filter options in display order.
var PRFilterOptions = []PRFilterMask{PRFilterOpen, PRFilterClosed, PRFilterMerged}

func (m PRFilterMask) Has(f PRFilterMask) bool { return m&f != 0 }

func (m PRFilterMask) Toggle(f PRFilterMask) PRFilterMask { return m ^ f }

func (m PRFilterMask) Label() string {
	if m == PRFilterOpen|PRFilterClosed|PRFilterMerged {
		return "All"
	}
	var parts []string
	if m.Has(PRFilterOpen) {
		parts = append(parts, "Open")
	}
	if m.Has(PRFilterClosed) {
		parts = append(parts, "Closed")
	}
	if m.Has(PRFilterMerged) {
		parts = append(parts, "Merged")
	}
	if len(parts) == 0 {
		return "None"
	}
	return strings.Join(parts, "+")
}

func (m PRFilterMask) StateArg() string {
	// single selection: use specific state arg for efficiency
	switch m {
	case PRFilterOpen:
		return "open"
	case PRFilterClosed:
		return "closed"
	case PRFilterMerged:
		return "merged"
	default:
		return "all"
	}
}

// Matches returns true if the gh state string matches this filter mask.
func (m PRFilterMask) Matches(state string) bool {
	switch state {
	case "OPEN":
		return m.Has(PRFilterOpen)
	case "CLOSED":
		return m.Has(PRFilterClosed)
	case "MERGED":
		return m.Has(PRFilterMerged)
	default:
		return false
	}
}

// ListState holds PR list and selection state.
type ListState struct {
	Repo         string
	PRs          []Item
	PRsLoading   bool
	PRsSelected  int
	Filter       PRFilterMask
	FilterOpen   bool
	FilterCursor int
}

// DetailState holds detail panel display and loading state.
type DetailState struct {
	Mode    DetailMode
	Content string
	Loading LoadingKind
}

type EnterAction struct {
	Kind   EnterActionKind
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
			PRs:    []Item{},
			Filter: PRFilterOpen,
		},
		Detail: DetailState{
			Mode: DetailModeOverview,
		},
		Review: ReviewState{
			Comments: []ReviewComment{},
		},
	}
}

func (s *State) SetWindowSize(width int, height int) {
	s.Width = width
	s.Height = height
}

func (s *State) BeginLoadPRs() {
	s.List.PRsLoading = true
	s.Detail.Loading = LoadingPRs
}

// BeginReviewLoad marks a review operation as in-progress.
func (s *State) BeginReviewLoad() {
	s.Detail.Loading = LoadingReview
}

// ClearLoading clears any in-progress loading indicator.
func (s *State) ClearLoading() {
	s.Detail.Loading = LoadingNone
}

// StopReviewInput exits any active review input mode without closing the drawer.
func (s *State) StopReviewInput() {
	s.Review.InputMode = ReviewInputNone
}

func (s *State) ApplyPRsResult(repo string, prs []Item, err error) {
	s.List.PRsLoading = false
	s.Detail.Loading = LoadingNone
	if err != nil {
		s.showError("Error loading PRs", err)
		return
	}

	s.List.Repo = repo
	s.List.PRs = prs
	s.List.PRsSelected = 0
	s.Detail.Mode = DetailModeOverview
	s.resetReview()
	if len(prs) == 0 {
		s.Detail.Content = "No pull requests"
		return
	}
	s.Detail.Content = FormatPROverview(prs[s.List.PRsSelected])
}

func (s *State) ApplyDetailResult(content string, err error) {
	if err != nil {
		s.showError("Error loading detail", err)
		return
	}
	s.Detail.Loading = LoadingNone
	s.Detail.Content = sanitizeMultiline(content)
}

func (s *State) ApplyDiffResult(content string, err error) {
	if err != nil {
		s.showError("Error loading diff", err)
		return
	}
	s.Detail.Loading = LoadingNone
	s.Detail.Content = sanitizeMultiline(content)
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
	if changed && s.Detail.Mode == DetailModeOverview {
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
	if changed && s.Detail.Mode == DetailModeOverview {
		s.refreshDetailPreview()
	}
	return changed
}

func (s *State) SwitchToOverview() bool {
	if s.Detail.Mode == DetailModeOverview {
		return false
	}
	s.Detail.Mode = DetailModeOverview
	s.Detail.Loading = LoadingNone
	s.Review.InputMode = ReviewInputNone
	s.refreshDetailPreview()
	return true
}

func (s *State) SwitchToDiff() bool {
	if s.Detail.Mode == DetailModeDiff {
		return false
	}
	s.Detail.Mode = DetailModeDiff
	s.Detail.Loading = LoadingNone
	s.Review.DrawerOpen = false
	s.Review.InputMode = ReviewInputNone
	return true
}

func (s *State) IsDiffMode() bool {
	return s.Detail.Mode == DetailModeDiff
}

func (s *State) ShouldApplyDetailResult(mode DetailMode, number int) bool {
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
		s.Detail.Loading = LoadingNone
		s.Detail.Content = forcedDetailText
		return EnterAction{}
	}
	s.Detail.Loading = LoadingDetail
	if s.Detail.Mode == DetailModeDiff {
		return EnterAction{Kind: EnterLoadPRDiff, Repo: s.List.Repo, Number: item.Number}
	}
	return EnterAction{Kind: EnterLoadPRDetail, Repo: s.List.Repo, Number: item.Number}
}

func (s *State) refreshDetailPreview() {
	item, ok := s.selectedPR()
	if !ok {
		return
	}
	s.Detail.Content = FormatPROverview(item)
}

func (s *State) selectedPR() (Item, bool) {
	if len(s.List.PRs) == 0 {
		return Item{}, false
	}
	if s.List.PRsSelected < 0 || s.List.PRsSelected >= len(s.List.PRs) {
		return Item{}, false
	}
	return s.List.PRs[s.List.PRsSelected], true
}

func (s *State) showError(msg string, err error) {
	s.Detail.Loading = LoadingNone
	s.Detail.Content = sanitizeMultiline(fmt.Sprintf("%s: %v", msg, err))
}

func (s *State) SelectedPR() (Item, bool) {
	return s.selectedPR()
}

func (s *State) OpenReviewDrawer() {
	s.Review.DrawerOpen = true
}

func (s *State) CloseReviewDrawer() {
	s.Review.DrawerOpen = false
	s.Review.InputMode = ReviewInputNone
	s.Review.Notice = ""
}

func (s *State) BeginReviewCommentInput() {
	s.Review.DrawerOpen = true
	s.Review.InputMode = ReviewInputComment
	s.Review.Notice = ""
}

func (s *State) BeginReviewSummaryInput() {
	s.Review.DrawerOpen = true
	s.Review.InputMode = ReviewInputSummary
	s.Review.Notice = ""
}

func (s *State) SetReviewSummary(summary string) {
	s.Review.Summary = sanitizeMultiline(summary)
}

func (s *State) SetReviewContext(prNumber int, pullRequestID string, commitOID string, reviewID string) {
	s.Review.PRNumber = prNumber
	s.Review.PullRequestID = sanitizeSingleLine(pullRequestID)
	s.Review.CommitOID = sanitizeSingleLine(commitOID)
	s.Review.ReviewID = sanitizeSingleLine(reviewID)
}

func (s *State) AddReviewComment(comment ReviewComment) {
	s.Review.Comments = append(s.Review.Comments, ReviewComment{
		Path:      sanitizeSingleLine(comment.Path),
		Body:      sanitizeMultiline(comment.Body),
		Side:      sanitizeSingleLine(comment.Side),
		Line:      comment.Line,
		StartSide: sanitizeSingleLine(comment.StartSide),
		StartLine: comment.StartLine,
	})
	s.Review.Notice = "Review comment added."
	s.Review.DrawerOpen = true
	s.Review.InputMode = ReviewInputNone
	s.Review.RangeStart = nil
}

func (s *State) SetReviewNotice(msg string) {
	s.Review.Notice = sanitizeMultiline(msg)
}

func (s *State) ClearReviewNotice() {
	s.Review.Notice = ""
}

func (s *State) HasPendingReview() bool {
	return s.Review.ReviewID != ""
}

func (s *State) MarkReviewRangeStart(anchor ReviewRange) {
	copied := anchor
	s.Review.RangeStart = &copied
	s.Review.DrawerOpen = true
	s.Review.Notice = "Range start selected."
}

func (s *State) CycleReviewEvent() {
	switch s.Review.Event {
	case ReviewEventComment:
		s.Review.Event = ReviewEventApprove
	case ReviewEventApprove:
		s.Review.Event = ReviewEventRequestChanges
	default:
		s.Review.Event = ReviewEventComment
	}
}

func (s *State) ClearReviewRangeStart() {
	s.Review.RangeStart = nil
}

func (s *State) ResetReviewAfterSubmit(notice string) {
	s.resetReview()
	s.Review.Notice = sanitizeMultiline(notice)
}

func (s *State) ResetReviewAfterDiscard(notice string) {
	s.resetReview()
	s.Review.Notice = sanitizeMultiline(notice)
}

func (s *State) OpenFilterSelect() {
	s.List.FilterOpen = true
	s.List.FilterCursor = 0
}

func (s *State) CloseFilterSelect() {
	s.List.FilterOpen = false
}

func (s *State) MoveFilterCursor(dir int) {
	n := len(PRFilterOptions)
	s.List.FilterCursor = (s.List.FilterCursor + dir + n) % n
}

// ToggleFilterAtCursor toggles the filter option under the cursor.
// It prevents deselecting all options (at least one must remain selected).
func (s *State) ToggleFilterAtCursor() {
	if s.List.FilterCursor < 0 || s.List.FilterCursor >= len(PRFilterOptions) {
		return
	}
	opt := PRFilterOptions[s.List.FilterCursor]
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
		Comments: []ReviewComment{},
		Notice:   s.Review.Notice,
	}
}

func FormatPRItem(item Item) string {
	return fmt.Sprintf("PR #%d %s", item.Number, sanitizeSingleLine(item.Title))
}

func FormatPROverview(item Item) string {
	status := sanitizeSingleLine(item.Status)
	if status == "" {
		status = "OPEN"
	}

	assignee := "unassigned"
	if len(item.Assignees) > 0 {
		list := make([]string, 0, len(item.Assignees))
		for _, name := range item.Assignees {
			n := sanitizeSingleLine(name)
			if n != "" {
				list = append(list, n)
			}
		}
		if len(list) > 0 {
			assignee = list[0]
			if len(list) > 1 {
				assignee = fmt.Sprintf("%s (+%d)", list[0], len(list)-1)
			}
		}
	}

	return fmt.Sprintf(
		"PR #%d %s\nStatus: %s\nAssignee: %s",
		item.Number,
		sanitizeSingleLine(item.Title),
		status,
		assignee,
	)
}
