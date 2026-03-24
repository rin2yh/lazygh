package list

import (
	"strings"

	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/pr"
)

// Convert transforms a slice of gh.PRItem into pr.Item, filtering by the given mask.
func Convert(prs []gh.PRItem, filter PRFilterMask) []pr.Item {
	items := make([]pr.Item, 0, len(prs))
	for _, p := range prs {
		if !filter.Matches(p.State) {
			continue
		}
		status := p.State
		if p.IsDraft {
			status = pr.PRStatusDraft
		}
		assignees := make([]string, 0, len(p.Assignees))
		for _, user := range p.Assignees {
			name := strings.TrimSpace(user.Login)
			if name != "" {
				assignees = append(assignees, name)
			}
		}
		items = append(items, pr.Item{
			Number:    p.Number,
			Title:     p.Title,
			Status:    status,
			Assignees: assignees,
		})
	}
	return items
}
