package review

// InputMode represents the active text-input state within the review drawer.
type InputMode int

const (
	InputNone InputMode = iota
	InputComment
	InputSummary
	InputThreadReply
)
