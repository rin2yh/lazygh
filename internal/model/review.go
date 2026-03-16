package model

// ReviewInputMode and ReviewEvent are used by both internal/state and multiple
// internal/gui subpackages (input routing, layout/status). If state were split
// into feature packages, these would still cross package boundaries, so they
// should stay here unless input routing is also refactored to delegate to the
// review package.

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

// ReviewComment, ReviewRange, and NoEditingComment are used exclusively within
// the review domain (internal/state.ReviewState + internal/gui/review/).
// If ReviewState were moved into internal/gui/review/, these could move to
// internal/gui/review/model.go.

type ReviewComment struct {
	CommentID string
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

const NoEditingComment = -1
