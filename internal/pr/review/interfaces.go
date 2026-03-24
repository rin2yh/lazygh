package review

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/pr"
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

// AppState is the minimal interface the review package needs from the app
// application state (list/detail state).
type AppState interface {
	SelectedPR() (pr.Item, bool)
	ListRepo() string
	BeginFetchReview()
	ClearFetching()
	IsDiffMode() bool
}

// PendingReviewClient handles GitHub API calls for the pending review workflow.
type PendingReviewClient interface {
	GetReviewContext(repo string, number int) (gh.ReviewContext, error)
	StartPendingReview(repo string, number int, ctx gh.ReviewContext) (string, error)
	AddReviewComment(repo string, reviewID string, comment gh.ReviewComment) (string, error)
	SubmitReview(repo string, reviewID string, event gh.ReviewEvent, body string) error
	DeletePendingReview(repo string, reviewID string) error
	DeletePendingReviewComment(commentID string) error
	UpdatePendingReviewComment(commentID string, body string) error
}

// Reader はGUIレイヤーが参照するレビュー状態の読み取りインターフェース。
// BuildDrawerInput が描画用DTOを一括提供するため、個々の状態フィールドは含めない。
type Reader interface {
	ShouldShowDrawer() bool
	IsIndexWithinPendingRange(path string, commentable bool, idx int) bool
	InputMode() InputMode
	IsInInputMode() bool
	RangeStart() *Range
	BuildDrawerInput(showDrawer bool) *Input
}

// Handler はユーザー入力によるレビュー操作を処理する。
type Handler interface {
	HandleInputKey(msg tea.KeyMsg) (tea.Cmd, bool)
	HandleAction(action config.Action) tea.Cmd
	// HandleCancel はレビュー固有のキャンセル処理を行う。
	// range選択中またはinput mode中であれば処理してtrueを返す。
	HandleCancel() bool
	SelectNextComment()
	SelectPrevComment()
	Notify(msg string)
}

// Applier は外部から得た結果をレビュー状態に適用する。
type Applier interface {
	CommentResult(msg CommentSavedMsg)
	DeleteCommentResult(msg CommentDeletedMsg)
	EditCommentResult(msg CommentUpdatedMsg)
	SubmitResult(msg SubmittedMsg)
	DiscardResult(msg DiscardedMsg)
	MarkStaleComments(files []gh.DiffFile)
}
