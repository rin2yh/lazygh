package model

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
