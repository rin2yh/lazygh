package core

import "fmt"

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
	Summary       string
	Comments      []ReviewComment
	RangeStart    *ReviewRange
	Notice        string
}

type EnterAction struct {
	Kind   EnterActionKind
	Repo   string
	Number int
}

type State struct {
	Repo string

	PRs         []Item
	PRsLoading  bool
	PRsSelected int

	DetailMode    DetailMode
	DetailContent string
	Loading       LoadingKind
	Review        ReviewState

	Width  int
	Height int
}

func NewState() *State {
	return &State{
		PRs:        []Item{},
		DetailMode: DetailModeOverview,
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
	s.PRsLoading = true
	s.Loading = LoadingPRs
}

func (s *State) ApplyPRsResult(repo string, prs []Item, err error) {
	s.PRsLoading = false
	s.Loading = LoadingNone
	if err != nil {
		s.showError("Error loading PRs", err)
		return
	}

	s.Repo = repo
	s.PRs = prs
	s.PRsSelected = 0
	s.DetailMode = DetailModeOverview
	s.resetReview()
	if len(prs) == 0 {
		s.DetailContent = "No pull requests"
		return
	}
	s.DetailContent = FormatPROverview(prs[s.PRsSelected])
}

func (s *State) ApplyDetailResult(content string, err error) {
	if err != nil {
		s.showError("Error loading detail", err)
		return
	}
	s.Loading = LoadingNone
	s.DetailContent = sanitizeMultiline(content)
}

func (s *State) ApplyDiffResult(content string, err error) {
	if err != nil {
		s.showError("Error loading diff", err)
		return
	}
	s.Loading = LoadingNone
	s.DetailContent = sanitizeMultiline(content)
}

func (s *State) NavigateDown() bool {
	if s.blocksPRSelectionChange() {
		s.Review.Notice = "Pending review exists. Submit with S or discard with X."
		return false
	}
	changed := false
	if s.PRsSelected < len(s.PRs)-1 {
		s.PRsSelected++
		changed = true
	}
	if changed && s.DetailMode == DetailModeOverview {
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
	if s.PRsSelected > 0 {
		s.PRsSelected--
		changed = true
	}
	if changed && s.DetailMode == DetailModeOverview {
		s.refreshDetailPreview()
	}
	return changed
}

func (s *State) SwitchToOverview() bool {
	if s.DetailMode == DetailModeOverview {
		return false
	}
	s.DetailMode = DetailModeOverview
	s.Loading = LoadingNone
	s.Review.InputMode = ReviewInputNone
	s.refreshDetailPreview()
	return true
}

func (s *State) SwitchToDiff() bool {
	if s.DetailMode == DetailModeDiff {
		return false
	}
	s.DetailMode = DetailModeDiff
	s.Loading = LoadingNone
	s.Review.DrawerOpen = false
	s.Review.InputMode = ReviewInputNone
	return true
}

func (s *State) IsDiffMode() bool {
	return s.DetailMode == DetailModeDiff
}

func (s *State) ShouldApplyDetailResult(mode DetailMode, number int) bool {
	if s.DetailMode != mode {
		return false
	}
	item, ok := s.selectedPR()
	if !ok {
		return false
	}
	return item.Number == number
}

func (s *State) PlanEnter(hasClient bool, forcedDetailText string) EnterAction {
	if !hasClient || s.PRsLoading {
		return EnterAction{}
	}
	item, ok := s.selectedPR()
	if !ok {
		return EnterAction{}
	}
	if forcedDetailText != "" {
		s.Loading = LoadingNone
		s.DetailContent = forcedDetailText
		return EnterAction{}
	}
	s.Loading = LoadingDetail
	if s.DetailMode == DetailModeDiff {
		return EnterAction{Kind: EnterLoadPRDiff, Repo: s.Repo, Number: item.Number}
	}
	return EnterAction{Kind: EnterLoadPRDetail, Repo: s.Repo, Number: item.Number}
}

func (s *State) refreshDetailPreview() {
	item, ok := s.selectedPR()
	if !ok {
		return
	}
	s.DetailContent = FormatPROverview(item)
}

func (s *State) selectedPR() (Item, bool) {
	if len(s.PRs) == 0 {
		return Item{}, false
	}
	if s.PRsSelected < 0 || s.PRsSelected >= len(s.PRs) {
		return Item{}, false
	}
	return s.PRs[s.PRsSelected], true
}

func (s *State) showError(msg string, err error) {
	s.Loading = LoadingNone
	s.DetailContent = sanitizeMultiline(fmt.Sprintf("%s: %v", msg, err))
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
