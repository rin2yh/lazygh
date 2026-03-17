package list

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/model"
)

func formatItem(item model.Item) string {
	return fmt.Sprintf("#%d %s", item.Number, model.SanitizeSingleLine(item.Title))
}

func formatOverview(item model.Item) string {
	status := model.SanitizeSingleLine(item.Status)
	if status == "" {
		status = model.PRStatusOpen
	}

	assignee := "unassigned"
	first := ""
	extra := 0
	for _, name := range item.Assignees {
		n := model.SanitizeSingleLine(name)
		if n == "" {
			continue
		}
		if first == "" {
			first = n
		} else {
			extra++
		}
	}
	if first != "" {
		if extra > 0 {
			assignee = fmt.Sprintf("%s (+%d)", first, extra)
		} else {
			assignee = first
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
