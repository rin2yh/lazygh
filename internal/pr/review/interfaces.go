package review

import (
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/model"
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
	SelectedPR() (model.Item, bool)
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
