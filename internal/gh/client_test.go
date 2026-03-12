package gh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/pkg/test/fixture"
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
	case args[1] == "repo" && len(args) > 2 && args[2] == "view":
		fmt.Print(`{"nameWithOwner":"owner/repo1"}`)
	case args[1] == "pr" && len(args) > 2 && args[2] == "list":
		fmt.Print(`[{"number":1,"title":"Fix bug"},{"number":2,"title":"Add feature"}]`)
	case args[1] == "pr" && len(args) > 2 && args[2] == "view":
		fmt.Print("PR view content")
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
		if prs[i].Number != fx.PRs[i].Number || prs[i].Title != fx.PRs[i].Title {
			t.Fatalf("unexpected PR[%d]: %+v", i, prs[i])
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

func TestResolveCurrentRepo_Error(t *testing.T) {
	c := &Client{execCommand: func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}}
	_, err := c.ResolveCurrentRepo()
	if err == nil {
		t.Fatal("expected error, got nil")
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
