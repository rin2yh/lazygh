package app

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/pr"
	"github.com/rin2yh/lazygh/internal/pr/list"
	"github.com/rin2yh/lazygh/internal/pr/overview"
	"github.com/rin2yh/lazygh/pkg/sanitize"
)

// EnterActionKind represents the type of action to take when entering a PR.
type EnterActionKind int

const (
	EnterNone EnterActionKind = iota
	EnterLoadPRDetail
	EnterLoadPRDiff
)

type EnterAction struct {
	Kind   EnterActionKind
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
	list.State

	Overview overview.State

	Width  int
	Height int

	review ReviewHook
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		State: list.State{
			Items:  []pr.Item{},
			Filter: list.PRFilterOpen,
		},
		Overview: overview.State{
			Mode: overview.DetailModeOverview,
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
	c.Overview.Fetching = overview.FetchingPRs
}

// ApplyPRsResult は PR 一覧結果を反映し、review をリセットする。
func (c *Coordinator) ApplyPRsResult(repo string, items []pr.Item, err error) {
	c.Fetching = false
	c.Overview.Fetching = overview.FetchNone
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
	c.Overview.Mode = overview.DetailModeOverview
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

func (c *Coordinator) SelectedPR() (pr.Item, bool) { return c.selectedPR() }
func (c *Coordinator) ListRepo() string            { return c.Repo }
func (c *Coordinator) BeginFetchReview()           { c.Overview.Fetching = overview.FetchingReview }
func (c *Coordinator) ClearFetching()              { c.Overview.Fetching = overview.FetchNone }
func (c *Coordinator) IsDiffMode() bool            { return c.Overview.Mode == overview.DetailModeDiff }

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
	c.Overview.Fetching = overview.FetchNone
	c.Overview.Content = sanitize.Multiline(content)
}

func (c *Coordinator) NavigateDown() bool {
	changed := c.State.NavigateDown()
	if changed && c.Overview.Mode == overview.DetailModeOverview {
		c.refreshOverviewPreview()
	}
	return changed
}

func (c *Coordinator) NavigateUp() bool {
	changed := c.State.NavigateUp()
	if changed && c.Overview.Mode == overview.DetailModeOverview {
		c.refreshOverviewPreview()
	}
	return changed
}

func (c *Coordinator) SwitchToOverview() bool {
	if c.Overview.Mode == overview.DetailModeOverview {
		return false
	}
	c.Overview.Mode = overview.DetailModeOverview
	c.Overview.Fetching = overview.FetchNone
	c.refreshOverviewPreview()
	return true
}

func (c *Coordinator) SwitchToDiff() bool {
	if c.Overview.Mode == overview.DetailModeDiff {
		return false
	}
	c.Overview.Mode = overview.DetailModeDiff
	c.Overview.Fetching = overview.FetchNone
	return true
}

func (c *Coordinator) ShouldApplyDetailResult(mode overview.DetailMode, number int) bool {
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
	c.Overview.Fetching = overview.FetchingDetail
	if c.Overview.Mode == overview.DetailModeDiff {
		return EnterAction{Kind: EnterLoadPRDiff, Repo: c.Repo, Number: item.Number}
	}
	return EnterAction{Kind: EnterLoadPRDetail, Repo: c.Repo, Number: item.Number}
}

func (c *Coordinator) refreshOverviewPreview() {
	if content, ok := c.SelectedOverview(); ok {
		c.Overview.Content = content
	}
}

func (c *Coordinator) selectedPR() (pr.Item, bool) {
	if len(c.Items) == 0 {
		return pr.Item{}, false
	}
	if c.Selected < 0 || c.Selected >= len(c.Items) {
		return pr.Item{}, false
	}
	return c.Items[c.Selected], true
}

func (c *Coordinator) showError(msg string, err error) {
	c.Overview.Fetching = overview.FetchNone
	c.Overview.Content = sanitize.Multiline(fmt.Sprintf("%s: %v", msg, err))
}
