package gui

import "fmt"

func formatPanelTitle(base string, active bool) string {
	if active {
		return fmt.Sprintf("> %s <", base)
	}
	return fmt.Sprintf(" %s ", base)
}

func formatStatusLine(_ string) string {
	return fmt.Sprintf("[q]Quit  [j/k]Move  [enter]Reload detail")
}

func formatRepoLine(repo string) string {
	if repo == "" {
		repo = "(resolving...)"
	}
	return repo
}
