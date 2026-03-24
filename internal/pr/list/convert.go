package list

import (
	"strings"

	"github.com/rin2yh/lazygh/internal/gh"
	pr "github.com/rin2yh/lazygh/internal/pr"
)

// Convert transforms a slice of gh.PRItem into pr.Item, filtering by the given mask.
func Convert(ghprs []gh.PRItem, filter PRFilterMask) []pr.Item {
	items := make([]pr.Item, 0, len(ghprs))
	for _, ghpr := range ghprs {
		if !filter.Matches(ghpr.State) {
			continue
		}
		status := ghpr.State
		if ghpr.IsDraft {
			status = pr.PRStatusDraft
		}
		assignees := make([]string, 0, len(ghpr.Assignees))
		for _, user := range ghpr.Assignees {
			name := strings.TrimSpace(user.Login)
			if name != "" {
				assignees = append(assignees, name)
			}
		}
		items = append(items, pr.Item{
			Number:    ghpr.Number,
			Title:     ghpr.Title,
			Status:    status,
			Assignees: assignees,
		})
	}
	return items
}
