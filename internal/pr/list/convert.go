package list

import (
	"strings"

	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/model"
)

// Convert transforms a slice of gh.PRItem into model.Item, filtering by the given mask.
func Convert(prs []gh.PRItem, filter PRFilterMask) []model.Item {
	items := make([]model.Item, 0, len(prs))
	for _, pr := range prs {
		if !filter.Matches(pr.State) {
			continue
		}
		status := pr.State
		if pr.IsDraft {
			status = model.PRStatusDraft
		}
		assignees := make([]string, 0, len(pr.Assignees))
		for _, user := range pr.Assignees {
			name := strings.TrimSpace(user.Login)
			if name != "" {
				assignees = append(assignees, name)
			}
		}
		items = append(items, model.Item{
			Number:    pr.Number,
			Title:     pr.Title,
			Status:    status,
			Assignees: assignees,
		})
	}
	return items
}
