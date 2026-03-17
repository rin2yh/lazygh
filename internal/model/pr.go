package model

import (
	"strings"
)

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

type PRFilterMask uint8

const (
	PRFilterOpen   PRFilterMask = 1 << iota // 1
	PRFilterClosed                          // 2
	PRFilterMerged                          // 4
)

// PRFilterOptions lists the filter options in display order.
var PRFilterOptions = []PRFilterMask{PRFilterOpen, PRFilterClosed, PRFilterMerged}

func (m PRFilterMask) Has(f PRFilterMask) bool { return m&f != 0 }

func (m PRFilterMask) Toggle(f PRFilterMask) PRFilterMask { return m ^ f }

func (m PRFilterMask) Label() string {
	if m == PRFilterOpen|PRFilterClosed|PRFilterMerged {
		return "All"
	}
	var parts []string
	if m.Has(PRFilterOpen) {
		parts = append(parts, "Open")
	}
	if m.Has(PRFilterClosed) {
		parts = append(parts, "Closed")
	}
	if m.Has(PRFilterMerged) {
		parts = append(parts, "Merged")
	}
	if len(parts) == 0 {
		return "None"
	}
	return strings.Join(parts, ",")
}

func (m PRFilterMask) StateArg() string {
	// single selection: use specific state arg for efficiency
	switch m {
	case PRFilterOpen:
		return "open"
	case PRFilterClosed:
		return "closed"
	case PRFilterMerged:
		return "merged"
	default:
		return "all"
	}
}

// Matches returns true if the gh state string matches this filter mask.
func (m PRFilterMask) Matches(state string) bool {
	switch state {
	case PRStatusOpen:
		return m.Has(PRFilterOpen)
	case PRStatusClosed:
		return m.Has(PRFilterClosed)
	case PRStatusMerged:
		return m.Has(PRFilterMerged)
	default:
		return false
	}
}
