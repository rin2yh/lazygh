package prs

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/model"
)

// FormatPRItem formats a PR for display in the list panel.
func FormatPRItem(item model.Item) string {
	return fmt.Sprintf("#%d %s", item.Number, model.SanitizeSingleLine(item.Title))
}

// FormatPROverview formats a PR for display in the overview panel.
func FormatPROverview(item model.Item) string {
	status := model.SanitizeSingleLine(item.Status)
	if status == "" {
		status = model.PRStatusOpen
	}

	assignee := "unassigned"
	if len(item.Assignees) > 0 {
		list := make([]string, 0, len(item.Assignees))
		for _, name := range item.Assignees {
			n := model.SanitizeSingleLine(name)
			if n != "" {
				list = append(list, n)
			}
		}
		if len(list) > 0 {
			assignee = list[0]
			if len(list) > 1 {
				assignee = fmt.Sprintf("%s (+%d)", list[0], len(list)-1)
			}
		}
	}

	return fmt.Sprintf(
		"PR #%d %s\nStatus: %s\nAssignee: %s",
		item.Number,
		model.SanitizeSingleLine(item.Title),
		status,
		assignee,
	)
}
