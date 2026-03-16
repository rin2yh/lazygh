package gh

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/pkg/test/fake"
	"github.com/rin2yh/lazygh/pkg/test/fixture"
	ghstub "github.com/rin2yh/lazygh/pkg/test/stub/gh"
)

// TestFakeProcess は実際のテストではなく、fake gh コマンドとして動作する。
func TestFakeProcess(t *testing.T) {
	if os.Getenv("GO_TEST_HELPER_PROCESS") != "1" {
		return
	}

	table := map[string]fake.Response{
		"repo view": {
			Stdout:   `{"nameWithOwner":"owner/repo1"}`,
			ExitCode: 0,
		},
		"pr list": {
			Stdout:   `[{"number":1,"title":"Fix bug","state":"OPEN","isDraft":false,"assignees":[{"login":"alice"}]},{"number":2,"title":"Add feature","state":"OPEN","isDraft":true,"assignees":[]}]`,
			ExitCode: 0,
		},
		"pr view": {
			Stdout:   "PR view content",
			ExitCode: 0,
		},
		"pr diff": {
			Stdout:   "PR diff content",
			ExitCode: 0,
		},
		"api graphql headRefOid": {
			Stdout:   `{"data":{"repository":{"pullRequest":{"id":"PR_kwDOAA","headRefOid":"deadbeef"}}}}`,
			ExitCode: 0,
		},
	}

	gh := fake.Gh{Table: table}

	ghArgs, err := gh.ParseArgs(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	key, ok := gh.Key(ghArgs)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown: %s\n", strings.Join(ghArgs, " "))
		os.Exit(1)
	}
	resp, ok := gh.Find(key)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown: %s\n", strings.Join(ghArgs, " "))
		os.Exit(1)
	}
	gh.Write(resp)
	os.Exit(resp.ExitCode)
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	runner := &commandRunner{execCommand: ghstub.NewCommand(t, "TestFakeProcess", "GO_TEST_HELPER_PROCESS")}
	return &Client{runner: runner, api: &apiClient{runner: runner}}
}

func TestResolveCurrentRepo(t *testing.T) {
	c := newTestClient(t)
	fx := fixture.NewGHSuccess()

	repo, err := c.ResolveCurrentRepo()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo != fx.Repo {
		t.Fatalf("got %q, want %q", repo, fx.Repo)
	}
}

func TestListPRs(t *testing.T) {
	c := newTestClient(t)
	fx := fixture.NewGHSuccess()

	prs, err := c.ListPRs(fx.Repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prs) != len(fx.PRs) {
		t.Fatalf("got %d PRs, want %d", len(prs), len(fx.PRs))
	}
	for i := range prs {
		if prs[i].Number != fx.PRs[i].Number ||
			prs[i].Title != fx.PRs[i].Title ||
			prs[i].State != fx.PRs[i].State ||
			prs[i].IsDraft != fx.PRs[i].IsDraft {
			t.Fatalf("unexpected PR[%d]: %+v", i, prs[i])
		}
		if len(prs[i].Assignees) != len(fx.PRs[i].Assignees) {
			t.Fatalf("unexpected assignee length at PR[%d]: %+v", i, prs[i])
		}
		for j := range prs[i].Assignees {
			if prs[i].Assignees[j].Login != fx.PRs[i].Assignees[j] {
				t.Fatalf("unexpected assignee at PR[%d][%d]: %+v", i, j, prs[i].Assignees[j])
			}
		}
	}
}

func TestViewPR(t *testing.T) {
	c := newTestClient(t)
	fx := fixture.NewGHSuccess()

	content, err := c.ViewPR(fx.Repo, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != fx.Content {
		t.Fatalf("got %q, want %q", content, fx.Content)
	}
}

func TestDiffPR(t *testing.T) {
	c := newTestClient(t)
	fx := fixture.NewGHSuccess()

	content, err := c.DiffPR(fx.Repo, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != fx.Diff {
		t.Fatalf("got %q, want %q", content, fx.Diff)
	}
}

func TestResolveCurrentRepo_Error(t *testing.T) {
	runner := &commandRunner{execCommand: func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}}
	c := &Client{runner: runner, api: &apiClient{runner: runner}}
	_, err := c.ResolveCurrentRepo()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolveCurrentRepo_CommandErrorIncludesContext(t *testing.T) {
	runner := &commandRunner{execCommand: func(name string, args ...string) *exec.Cmd {
		return exec.Command("sh", "-c", "echo permission denied >&2; exit 1")
	}}
	c := &Client{runner: runner, api: &apiClient{runner: runner}}

	_, err := c.ResolveCurrentRepo()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var cmdErr *CommandError
	if !errors.As(err, &cmdErr) {
		t.Fatalf("expected CommandError, got %T", err)
	}
	if got := strings.Join(cmdErr.Command, " "); got != "repo view --json nameWithOwner" {
		t.Fatalf("unexpected command: %q", got)
	}
	if !strings.Contains(cmdErr.Stderr, "permission denied") {
		t.Fatalf("stderr was not captured: %q", cmdErr.Stderr)
	}
	if !strings.Contains(err.Error(), "gh repo view --json nameWithOwner failed") {
		t.Fatalf("missing command context: %v", err)
	}
}

func TestListPRs_InvalidJSON(t *testing.T) {
	runner := &commandRunner{execCommand: func(name string, args ...string) *exec.Cmd {
		return exec.Command("sh", "-c", "printf 'not-json'")
	}}
	c := &Client{runner: runner, api: &apiClient{runner: runner}}

	_, err := c.ListPRs("owner/repo1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid character") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetReviewContext(t *testing.T) {
	c := newTestClient(t)

	ctx, err := c.GetReviewContext("owner/repo1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.PullRequestID != "PR_kwDOAA" {
		t.Fatalf("got %q, want %q", ctx.PullRequestID, "PR_kwDOAA")
	}
	if ctx.CommitOID != "deadbeef" {
		t.Fatalf("got %q, want %q", ctx.CommitOID, "deadbeef")
	}
}

func TestValidateCLI(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		t.Setenv("PATH", fixture.NewPathWithGH(t))

		err := ValidateCLI()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Setenv("PATH", fixture.NewEmptyPath(t))

		err := ValidateCLI()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "gh CLI is required") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
