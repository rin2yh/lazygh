package review

import (
	tea "github.com/charmbracelet/bubbletea"
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

// AppState is the minimal interface the review package needs from the host
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

// Reader はレビュー状態の読み取り専用インターフェース。ISP に従い Handler / Applier と分離している。
type Reader interface {
	ShouldShowDrawer() bool
	IsIndexWithinPendingRange(path string, commentable bool, idx int) bool
	SummaryValue() string
	CommentInputLines() []string
	SummaryInputLines() []string
	IsEditingComment() bool
	InputMode() InputMode
	Summary() string
	EventLabel() string
	Notice() string
	RangeStart() *Range
	Comments() []Comment
	SelectedCommentIdx() int
	HasRangeStart() bool
	IsInInputMode() bool
	HasPendingReview() bool
	PRNumber() int
	BuildDrawerInput(showDrawer bool) *DrawerInput
}

// Handler はユーザー入力によるレビュー操作を処理する。
type Handler interface {
	EditorKey(msg tea.KeyMsg) (tea.Cmd, bool)
	Submit() tea.Cmd
	Discard() tea.Cmd
	SaveComment() tea.Cmd
	SaveEditComment() tea.Cmd
	DeleteComment() tea.Cmd
	StopInput()
	ClearCommentInput()
	CycleReviewEvent()
	EditComment() bool
	SelectNextComment()
	SelectPrevComment()
	ToggleRangeSelection()
	BeginCommentFlow()
	BeginSummaryInput()
	Notify(msg string)
	ClearRangeStart()
}

// Applier は非同期操作の結果をレビュー状態に適用する。
type Applier interface {
	CommentResult(msg CommentSavedMsg)
	DeleteCommentResult(msg CommentDeletedMsg)
	EditCommentResult(msg CommentUpdatedMsg)
	SubmitResult(msg SubmittedMsg)
	DiscardResult(msg DiscardedMsg)
}
