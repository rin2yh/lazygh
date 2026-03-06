package gh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestHelperProcess は実際のテストではなく、fake gh コマンドとして動作する。
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_TEST_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	sep := -1
	for i, a := range args {
		if a == "--" {
			sep = i
			break
		}
	}
	if sep < 0 || sep+1 >= len(args) {
		os.Exit(1)
	}
	args = args[sep+1:]

	if len(args) < 2 {
		os.Exit(1)
	}

	switch {
	case args[1] == "repo" && len(args) > 2 && args[2] == "list":
		fmt.Print(`[{"nameWithOwner":"owner/repo1"},{"nameWithOwner":"owner/repo2"}]`)
	case args[1] == "pr" && len(args) > 2 && args[2] == "list":
		fmt.Print(`[{"number":1,"title":"Fix bug"},{"number":2,"title":"Add feature"}]`)
	case args[1] == "issue" && len(args) > 2 && args[2] == "list":
		fmt.Print(`[{"number":10,"title":"Issue one"}]`)
	case args[1] == "pr" && len(args) > 2 && args[2] == "view":
		fmt.Print("PR view content")
	case args[1] == "issue" && len(args) > 2 && args[2] == "view":
		fmt.Print("Issue view content")
	default:
		fmt.Fprintf(os.Stderr, "unknown: %s\n", strings.Join(args, " "))
		os.Exit(1)
	}
	os.Exit(0)
}

func helperCmd(t *testing.T) func(string, ...string) *exec.Cmd {
	t.Helper()
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_TEST_HELPER_PROCESS=1"}
		return cmd
	}
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	return &Client{execCommand: helperCmd(t)}
}

func TestListRepos(t *testing.T) {
	c := newTestClient(t)
	repos, err := c.ListRepos()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 2 {
		t.Fatalf("got %d repos, want 2", len(repos))
	}
	if repos[0] != "owner/repo1" || repos[1] != "owner/repo2" {
		t.Errorf("unexpected repos: %v", repos)
	}
}

func TestListPRs(t *testing.T) {
	c := newTestClient(t)
	prs, err := c.ListPRs("owner/repo1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prs) != 2 {
		t.Fatalf("got %d PRs, want 2", len(prs))
	}
	if prs[0].Number != 1 || prs[0].Title != "Fix bug" {
		t.Errorf("unexpected PR[0]: %+v", prs[0])
	}
	if prs[1].Number != 2 || prs[1].Title != "Add feature" {
		t.Errorf("unexpected PR[1]: %+v", prs[1])
	}
}

func TestListIssues(t *testing.T) {
	c := newTestClient(t)
	issues, err := c.ListIssues("owner/repo1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("got %d issues, want 1", len(issues))
	}
	if issues[0].Number != 10 || issues[0].Title != "Issue one" {
		t.Errorf("unexpected issue: %+v", issues[0])
	}
}

func TestViewPR(t *testing.T) {
	c := newTestClient(t)
	content, err := c.ViewPR("owner/repo1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != "PR view content" {
		t.Errorf("got %q, want %q", content, "PR view content")
	}
}

func TestViewIssue(t *testing.T) {
	c := newTestClient(t)
	content, err := c.ViewIssue("owner/repo1", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != "Issue view content" {
		t.Errorf("got %q, want %q", content, "Issue view content")
	}
}

func TestListRepos_Error(t *testing.T) {
	c := &Client{execCommand: func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}}
	_, err := c.ListRepos()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
