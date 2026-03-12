package gui

import "fmt"

func formatPanelTitle(base string, active bool) string {
	if active {
		return fmt.Sprintf("> %s <", base)
	}
	return fmt.Sprintf(" %s ", base)
}

func formatStatusLine(repo string) string {
	if repo == "" {
		repo = "(resolving...)"
	}
	return fmt.Sprintf("Repo: %s  [q]Quit  [j/k]Move  [enter]Reload detail", repo)
}
