package fixture

import (
	"os"
	"path/filepath"
	"testing"
)

type PR struct {
	Number    int
	Title     string
	State     string
	IsDraft   bool
	Assignees []string
}

type GHSuccess struct {
	Repo    string
	PRs     []PR
	Content string
	Diff    string
}

func NewGHSuccess() GHSuccess {
	return GHSuccess{
		Repo: "owner/repo1",
		PRs: []PR{
			{Number: 1, Title: "Fix bug", State: "OPEN", IsDraft: false, Assignees: []string{"alice"}},
			{Number: 2, Title: "Add feature", State: "OPEN", IsDraft: true, Assignees: nil},
		},
		Content: "PR view content",
		Diff:    "PR diff content",
	}
}

func NewPathWithGH(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	ghPath := filepath.Join(dir, "gh")
	if err := os.WriteFile(ghPath, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write fake gh failed: %v", err)
	}
	return dir
}

func NewEmptyPath(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}
