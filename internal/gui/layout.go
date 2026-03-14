package gui

import "fmt"

func formatPanelTitle(base string) string {
	return fmt.Sprintf(" %s ", base)
}

func formatStatusLine(loading bool, diffMode bool, hasPR bool, focus panelFocus, hasFiles bool) string {
	base := "[q]Quit"
	if hasPR {
		base = fmt.Sprintf("%s [enter]Reload", base)
	}

	var modeSpecific string
	switch {
	case !diffMode && hasPR:
		modeSpecific = "[PRs] [j/k/↑/↓]Move [d]Diff"
	case !diffMode:
		modeSpecific = "[d]Diff"
	case focus == panelPRs && hasPR:
		modeSpecific = "[tab]Focus [PRs] [j/k/↑/↓]Move [l]Overview"
	case focus == panelDiffFiles && hasFiles:
		modeSpecific = "[tab]Focus [Files] [j/k/↑/↓]Move [l]Diff [o]Overview"
	case hasPR || hasFiles:
		modeSpecific = "[tab]Focus [Diff] [j/k/↑/↓]Line [space/b]Page [g/G]Top/Bottom [h]Files [o]Overview"
	default:
		modeSpecific = "[o]Overview"
	}

	line := base
	if modeSpecific != "" {
		line = fmt.Sprintf("%s | %s", base, modeSpecific)
	}
	if loading {
		return fmt.Sprintf("Loading...  | %s", line)
	}
	return line
}

func formatRepoLine(repo string) string {
	return repo
}
