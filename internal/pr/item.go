// Package pr defines shared pull request domain types.
package pr

// Item represents a pull request.
type Item struct {
	Number    int
	Title     string
	Status    string
	Assignees []string
}

const (
	PRStatusOpen   = "OPEN"
	PRStatusClosed = "CLOSED"
	PRStatusMerged = "MERGED"
	PRStatusDraft  = "DRAFT"
)
