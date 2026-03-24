package review

// InputMode represents the active text-input state within the review drawer.
type InputMode int

const (
	InputNone InputMode = iota
	InputComment
	InputSummary
)

// Event represents the type of review action to submit.
type Event int

const (
	EventComment Event = iota
	EventApprove
	EventRequestChanges
)

// Label returns the GitHub API string for the review event.
func (e Event) Label() string {
	switch e {
	case EventApprove:
		return "APPROVE"
	case EventRequestChanges:
		return "REQUEST CHANGES"
	default:
		return "COMMENT"
	}
}

// Comment holds a single pending review comment.
type Comment struct {
	CommentID string
	Path      string
	Body      string
	Side      string
	Line      int
	StartSide string
	StartLine int
	// Stale is set when the comment's anchor position no longer exists in the
	// current diff (e.g. after the PR is updated with new commits).
	Stale bool
}

// Range identifies a diff line position for range-based comments.
type Range struct {
	Path      string
	Index     int
	Side      string
	Line      int
	StartSide string
	StartLine int
}

const noEditingComment = -1
