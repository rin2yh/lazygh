package review

// Event represents the type of review action to submit.
type Event int

const (
	EventComment Event = iota
	EventApprove
	EventRequestChanges
)

// Label returns the display label for the review event.
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
