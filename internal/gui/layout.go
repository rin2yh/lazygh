package gui

import "fmt"

func formatPanelTitle(base string, active bool) string {
	if active {
		return fmt.Sprintf("> %s <", base)
	}
	return fmt.Sprintf(" %s ", base)
}

func formatStatusLine(loading bool) string {
	line := "[q]Quit  [j/k]Move  [enter]Reload detail"
	if loading {
		return fmt.Sprintf("Loading...  | %s", line)
	}
	return line
}

func formatRepoLine(repo string) string {
	return repo
}
