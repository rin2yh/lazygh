package state

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/model"
	"github.com/rin2yh/lazygh/internal/pr"
)

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
	pr.ListState

	Detail DetailState

	Width  int
	Height int
}

func NewState() *State {
	return &State{
		ListState: pr.ListState{
			Items:  []model.Item{},
			Filter: model.PRFilterOpen,
		},
		Detail: DetailState{
			Mode: model.DetailModeOverview,
		},
	}
}

func (s *State) SetWindowSize(width int, height int) {
	s.Width = width
	s.Height = height
}

func (s *State) BeginLoadPRs() {
	s.Fetching = true
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

func (s *State) ApplyPRsResult(repo string, items []model.Item, err error) {
	s.Fetching = false
	s.Detail.Loading = model.LoadingNone
	if err != nil {
		s.showError("Error loading PRs", err)
		return
	}

	s.Repo = repo
	s.Items = items
	s.Selected = 0
	s.Detail.Mode = model.DetailModeOverview
	if len(items) == 0 {
		s.Detail.Content = "No pull requests"
		return
	}
	s.Detail.Content = s.SelectedOverview()
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
	changed := false
	if s.Selected < len(s.Items)-1 {
		s.Selected++
		changed = true
	}
	if changed && s.Detail.Mode == model.DetailModeOverview {
		s.refreshDetailPreview()
	}
	return changed
}

func (s *State) NavigateUp() bool {
	changed := false
	if s.Selected > 0 {
		s.Selected--
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
	s.refreshDetailPreview()
	return true
}

func (s *State) SwitchToDiff() bool {
	if s.Detail.Mode == model.DetailModeDiff {
		return false
	}
	s.Detail.Mode = model.DetailModeDiff
	s.Detail.Loading = model.LoadingNone
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
	if !hasClient || s.Fetching {
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
		return EnterAction{Kind: model.EnterLoadPRDiff, Repo: s.Repo, Number: item.Number}
	}
	return EnterAction{Kind: model.EnterLoadPRDetail, Repo: s.Repo, Number: item.Number}
}

func (s *State) refreshDetailPreview() {
	if overview := s.SelectedOverview(); overview != "" {
		s.Detail.Content = overview
	}
}

func (s *State) selectedPR() (model.Item, bool) {
	if len(s.Items) == 0 {
		return model.Item{}, false
	}
	if s.Selected < 0 || s.Selected >= len(s.Items) {
		return model.Item{}, false
	}
	return s.Items[s.Selected], true
}

func (s *State) showError(msg string, err error) {
	s.Detail.Loading = model.LoadingNone
	s.Detail.Content = model.SanitizeMultiline(fmt.Sprintf("%s: %v", msg, err))
}

func (s *State) SelectedPR() (model.Item, bool) {
	return s.selectedPR()
}

// ListRepo returns the current repository slug.
func (s *State) ListRepo() string {
	return s.Repo
}

func (s *State) OpenFilterSelect() {
	s.FilterOpen = true
	s.FilterCursor = 0
}

func (s *State) CloseFilterSelect() {
	s.FilterOpen = false
}

func (s *State) MoveFilterCursor(dir int) {
	n := len(model.PRFilterOptions)
	s.FilterCursor = (s.FilterCursor + dir + n) % n
}

// ToggleFilterAtCursor toggles the filter option under the cursor.
// It prevents deselecting all options (at least one must remain selected).
func (s *State) ToggleFilterAtCursor() {
	if s.FilterCursor < 0 || s.FilterCursor >= len(model.PRFilterOptions) {
		return
	}
	opt := model.PRFilterOptions[s.FilterCursor]
	next := s.Filter.Toggle(opt)
	if next == 0 {
		return // disallow empty selection
	}
	s.Filter = next
}
