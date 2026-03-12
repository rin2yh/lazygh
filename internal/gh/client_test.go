package gh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/pkg/test/fake"
	"github.com/rin2yh/lazygh/pkg/test/fixture"
	"github.com/rin2yh/lazygh/pkg/test/stub"
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
			Stdout:   `[{"number":1,"title":"Fix bug"},{"number":2,"title":"Add feature"}]`,
			ExitCode: 0,
		},
		"pr view": {
			Stdout:   "PR view content",
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
	return &Client{execCommand: stub.NewCommand(t, "TestFakeProcess", "GO_TEST_HELPER_PROCESS")}
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
