package gui

import "fmt"

func formatPanelTitle(base string, active bool) string {
	if active {
		return fmt.Sprintf("> %s <", base)
	}
	return fmt.Sprintf(" %s ", base)
}

func formatStatusLine(loading bool, diffMode bool, hasPR bool, focus panelFocus, hasFiles bool) string {
	var line string
	switch {
	case !diffMode && hasPR:
		line = "[PRs] [j/k/↑/↓]Move [enter]Reload [d]Diff [q]Quit"
	case !diffMode:
		line = "[d]Diff [q]Quit"
	case focus == panelPRs && hasPR:
		line = "[tab]Focus [PRs] [j/k/↑/↓]Move [l]Overview [enter]Reload [q]Quit"
	case focus == panelDiffFiles && hasFiles:
		line = "[tab]Focus [Files] [j/k/↑/↓]Move [l]Diff [o]Overview [q]Quit"
	case hasPR || hasFiles:
		line = "[tab]Focus [Diff] [j/k/↑/↓]Line [space/b]Page [g/G]Top/Bottom [h]Files [enter]Reload [o]Overview [q]Quit"
	default:
		line = "[o]Overview [q]Quit"
	}
	if loading {
		return fmt.Sprintf("Loading...  | %s", line)
	}
	return line
}

func formatRepoLine(repo string) string {
	return repo
}
