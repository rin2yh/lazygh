package app

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

// ReviewHook は app/ が review 機能に要求する最小インターフェース（協調用）。
type ReviewHook interface {
	HasPendingReview() bool
	PRNumber() int
	Reset()
}

// Coordinator はアプリ全体の状態と機能間協調ロジックを保持する。
type Coordinator struct {
	list.ListState

	Overview overview.State

	Width  int
	Height int

	review ReviewHook
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		ListState: list.ListState{
			Items:  []model.Item{},
			Filter: model.PRFilterOpen,
		},
		Overview: overview.State{
			Mode: model.DetailModeOverview,
		},
	}
}

// SetReviewHook は review.Controller を注入する（gui.NewGui から呼ぶ）。
func (c *Coordinator) SetReviewHook(h ReviewHook) {
	c.review = h
}

// --- 機能間協調メソッド（app/ に集約）---

// BeginFetchPRs は PR 一覧ロード開始時にリスト・詳細の両状態を更新する。
func (c *Coordinator) BeginFetchPRs() {
	c.Fetching = true
	c.Overview.Fetching = model.FetchingPRs
}

// ApplyPRsResult は PR 一覧結果を反映し、review をリセットする。
func (c *Coordinator) ApplyPRsResult(repo string, items []model.Item, err error) {
	c.Fetching = false
	c.Overview.Fetching = model.FetchNone
	if err != nil {
		c.showError("Error fetching PRs", err)
		if c.review != nil {
			c.review.Reset()
		}
		return
	}

	c.Repo = repo
	c.Items = items
	c.Selected = 0
	c.Overview.Mode = model.DetailModeOverview
	if len(items) == 0 {
		c.Overview.Content = "No pull requests"
	} else if content, ok := c.SelectedOverview(); ok {
		c.Overview.Content = content
	}
	if c.review != nil {
		c.review.Reset()
	}
}

// BlocksPRSelectionChange は保留中レビューがある場合に PR 選択変更を禁止する。
func (c *Coordinator) BlocksPRSelectionChange() bool {
	item, ok := c.selectedPR()
	if !ok {
		return false
	}
	return c.review != nil && c.review.HasPendingReview() && c.review.PRNumber() == item.Number
}

// --- review.AppState インターフェースの実装 ---

func (c *Coordinator) SelectedPR() (model.Item, bool) { return c.selectedPR() }
func (c *Coordinator) ListRepo() string               { return c.Repo }
func (c *Coordinator) BeginFetchReview()              { c.Overview.Fetching = model.FetchingReview }
func (c *Coordinator) ClearFetching()                 { c.Overview.Fetching = model.FetchNone }
func (c *Coordinator) IsDiffMode() bool               { return c.Overview.Mode == model.DetailModeDiff }

// --- その他の state メソッド ---

func (c *Coordinator) SetWindowSize(width int, height int) {
	c.Width = width
	c.Height = height
}

func (c *Coordinator) ApplyDetailResult(content string, err error) {
	c.applyLoadedContent("Error fetching detail", content, err)
}

func (c *Coordinator) ApplyDiffResult(content string, err error) {
	c.applyLoadedContent("Error fetching diff", content, err)
}

func (c *Coordinator) applyLoadedContent(errPrefix, content string, err error) {
	if err != nil {
		c.showError(errPrefix, err)
		return
	}
	c.Overview.Fetching = model.FetchNone
	c.Overview.Content = model.SanitizeMultiline(content)
}

func (c *Coordinator) NavigateDown() bool {
	changed := false
	if c.Selected < len(c.Items)-1 {
		c.Selected++
		changed = true
	}
	if changed && c.Overview.Mode == model.DetailModeOverview {
		c.refreshOverviewPreview()
	}
	return changed
}

func (c *Coordinator) NavigateUp() bool {
	changed := false
	if c.Selected > 0 {
		c.Selected--
		changed = true
	}
	if changed && c.Overview.Mode == model.DetailModeOverview {
		c.refreshOverviewPreview()
	}
	return changed
}

func (c *Coordinator) SwitchToOverview() bool {
	if c.Overview.Mode == model.DetailModeOverview {
		return false
	}
	c.Overview.Mode = model.DetailModeOverview
	c.Overview.Fetching = model.FetchNone
	c.refreshOverviewPreview()
	return true
}

func (c *Coordinator) SwitchToDiff() bool {
	if c.Overview.Mode == model.DetailModeDiff {
		return false
	}
	c.Overview.Mode = model.DetailModeDiff
	c.Overview.Fetching = model.FetchNone
	return true
}

func (c *Coordinator) ShouldApplyDetailResult(mode model.DetailMode, number int) bool {
	if c.Overview.Mode != mode {
		return false
	}
	item, ok := c.selectedPR()
	if !ok {
		return false
	}
	return item.Number == number
}

func (c *Coordinator) PlanEnter(hasClient bool) EnterAction {
	if !hasClient || c.Fetching {
		return EnterAction{}
	}
	item, ok := c.selectedPR()
	if !ok {
		return EnterAction{}
	}
	c.Overview.Fetching = model.FetchingDetail
	if c.Overview.Mode == model.DetailModeDiff {
		return EnterAction{Kind: model.EnterLoadPRDiff, Repo: c.Repo, Number: item.Number}
	}
	return EnterAction{Kind: model.EnterLoadPRDetail, Repo: c.Repo, Number: item.Number}
}

func (c *Coordinator) OpenFilterSelect() {
	c.FilterOpen = true
	c.FilterCursor = 0
}

func (c *Coordinator) CloseFilterSelect() {
	c.FilterOpen = false
}

func (c *Coordinator) MoveFilterCursor(dir int) {
	n := len(model.PRFilterOptions)
	c.FilterCursor = (c.FilterCursor + dir + n) % n
}

// ToggleFilterAtCursor は選択中フィルタをトグルする（最低1つ必須）。
func (c *Coordinator) ToggleFilterAtCursor() {
	if c.FilterCursor < 0 || c.FilterCursor >= len(model.PRFilterOptions) {
		return
	}
	opt := model.PRFilterOptions[c.FilterCursor]
	next := c.Filter.Toggle(opt)
	if next == 0 {
		return
	}
	c.Filter = next
}

func (c *Coordinator) refreshOverviewPreview() {
	if content, ok := c.SelectedOverview(); ok {
		c.Overview.Content = content
	}
}

func (c *Coordinator) selectedPR() (model.Item, bool) {
	if len(c.Items) == 0 {
		return model.Item{}, false
	}
	if c.Selected < 0 || c.Selected >= len(c.Items) {
		return model.Item{}, false
	}
	return c.Items[c.Selected], true
}

func (c *Coordinator) showError(msg string, err error) {
	c.Overview.Fetching = model.FetchNone
	c.Overview.Content = model.SanitizeMultiline(fmt.Sprintf("%s: %v", msg, err))
}
