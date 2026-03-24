package list

import (
	"fmt"

	"github.com/rin2yh/lazygh/internal/pr"
	"github.com/rin2yh/lazygh/pkg/sanitize"
)

func formatItem(item pr.Item) string {
	return fmt.Sprintf("#%d %s", item.Number, sanitize.SingleLine(item.Title))
}

func formatOverview(item pr.Item) string {
	status := sanitize.SingleLine(item.Status)
	if status == "" {
		status = pr.PRStatusOpen
	}

	assignee := "unassigned"
	first := ""
	extra := 0
	for _, name := range item.Assignees {
		n := sanitize.SingleLine(name)
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
		sanitize.SingleLine(item.Title),
		status,
		assignee,
	)
}
