package state

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/model"
	"github.com/rin2yh/lazygh/internal/pr/list"
	"github.com/rin2yh/lazygh/internal/pr/overview"
)

type EnterAction struct {
	Kind   model.EnterActionKind
	Repo   string
	Number int
}

type State struct {
	list.ListState

	Overview overview.State

	Width  int
	Height int
}

func NewState() *State {
	return &State{
		ListState: list.ListState{
			Items:  []model.Item{},
			Filter: model.PRFilterOpen,
		},
		Overview: overview.State{
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
	s.Overview.Loading = model.LoadingPRs
}

// BeginReviewLoad marks a review operation as in-progress.
func (s *State) BeginReviewLoad() {
	s.Overview.Loading = model.LoadingReview
}

// ClearLoading clears any in-progress loading indicator.
func (s *State) ClearLoading() {
	s.Overview.Loading = model.LoadingNone
}

func (s *State) ApplyPRsResult(repo string, items []model.Item, err error) {
	s.Fetching = false
	s.Overview.Loading = model.LoadingNone
	if err != nil {
		s.showError("Error loading PRs", err)
		return
	}

	s.Repo = repo
	s.Items = items
	s.Selected = 0
	s.Overview.Mode = model.DetailModeOverview
	if len(items) == 0 {
		s.Overview.Content = "No pull requests"
		return
	}
	if content, ok := s.SelectedOverview(); ok {
		s.Overview.Content = content
	}
}

func (s *State) ApplyDetailResult(content string, err error) {
	if err != nil {
		s.showError("Error loading detail", err)
		return
	}
	s.Overview.Loading = model.LoadingNone
	s.Overview.Content = model.SanitizeMultiline(content)
}

func (s *State) ApplyDiffResult(content string, err error) {
	if err != nil {
		s.showError("Error loading diff", err)
		return
	}
	s.Overview.Loading = model.LoadingNone
	s.Overview.Content = model.SanitizeMultiline(content)
}

func (s *State) NavigateDown() bool {
	changed := false
	if s.Selected < len(s.Items)-1 {
		s.Selected++
		changed = true
	}
	if changed && s.Overview.Mode == model.DetailModeOverview {
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
	if changed && s.Overview.Mode == model.DetailModeOverview {
		s.refreshDetailPreview()
	}
	return changed
}

func (s *State) SwitchToOverview() bool {
	if s.Overview.Mode == model.DetailModeOverview {
		return false
	}
	s.Overview.Mode = model.DetailModeOverview
	s.Overview.Loading = model.LoadingNone
	s.refreshDetailPreview()
	return true
}

func (s *State) SwitchToDiff() bool {
	if s.Overview.Mode == model.DetailModeDiff {
		return false
	}
	s.Overview.Mode = model.DetailModeDiff
	s.Overview.Loading = model.LoadingNone
	return true
}

func (s *State) IsDiffMode() bool {
	return s.Overview.Mode == model.DetailModeDiff
}

func (s *State) ShouldApplyDetailResult(mode model.DetailMode, number int) bool {
	if s.Overview.Mode != mode {
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
		s.Overview.Loading = model.LoadingNone
		s.Overview.Content = forcedDetailText
		return EnterAction{}
	}
	s.Overview.Loading = model.LoadingDetail
	if s.Overview.Mode == model.DetailModeDiff {
		return EnterAction{Kind: model.EnterLoadPRDiff, Repo: s.Repo, Number: item.Number}
	}
	return EnterAction{Kind: model.EnterLoadPRDetail, Repo: s.Repo, Number: item.Number}
}

func (s *State) refreshDetailPreview() {
	if content, ok := s.SelectedOverview(); ok {
		s.Overview.Content = content
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
	s.Overview.Loading = model.LoadingNone
	s.Overview.Content = model.SanitizeMultiline(fmt.Sprintf("%s: %v", msg, err))
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
