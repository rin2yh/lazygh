package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/rin2yh/lazygh/pkg/test/e2e"
	"github.com/rin2yh/lazygh/pkg/test/fake"
)

func TestFakeGHProcess(t *testing.T) {
	if os.Getenv("GO_WANT_FAKE_GH_PROCESS") != "1" {
		return
	}

	table := map[string]fake.Response{
		"repo view": {
			Stdout:   `{"nameWithOwner":"owner/repo1"}`,
			ExitCode: 0,
		},
		"pr list": {
			Stdout:   `[{"number":1,"title":"Fix bug"}]`,
			ExitCode: 0,
		},
		"pr view": {
			Stdout:   "PR detail",
			ExitCode: 0,
		},
		"pr diff": {
			Stdout:   "diff --git a/main.go b/main.go\n--- a/main.go\n+++ b/main.go\n@@ -1 +1 @@\n-old line\n+new line\n",
			ExitCode: 0,
		},
		"api graphql headRefOid": {
			Stdout:   `{"data":{"repository":{"pullRequest":{"id":"PR_1","headRefOid":"abc123"}}}}`,
			ExitCode: 0,
		},
		"api graphql addPullRequestReview": {
			Stdout:   `{"data":{"addPullRequestReview":{"pullRequestReview":{"id":"PRR_1"}}}}`,
			ExitCode: 0,
		},
		"api graphql addPullRequestReviewThread": {
			Stdout:   `{"data":{"addPullRequestReviewThread":{"thread":{"id":"T_1"}}}}`,
			ExitCode: 0,
		},
		"api graphql submitPullRequestReview": {
			Stdout:   `{"data":{"submitPullRequestReview":{"pullRequestReview":{"id":"PRR_1"}}}}`,
			ExitCode: 0,
		},
	}

	gh := fake.Gh{
		Table:   table,
		LogPath: os.Getenv("FAKE_GH_LOG"),
	}

	ghArgs, err := gh.ParseArgs(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := gh.Log(ghArgs); err != nil {
		fmt.Fprintf(os.Stderr, "failed to append fake log: %v\n", err)
		os.Exit(1)
	}

	key, ok := gh.Key(ghArgs)
	if !ok {
		fmt.Fprintf(os.Stderr, "unexpected gh args: %s\n", strings.Join(ghArgs, " "))
		os.Exit(1)
	}
	resp, ok := gh.Find(key)
	if !ok {
		fmt.Fprintf(os.Stderr, "unexpected gh args: %s\n", strings.Join(ghArgs, " "))
		os.Exit(1)
	}
	gh.Write(resp)
	os.Exit(resp.ExitCode)
}

func skipIfNotE2E(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skip e2e in short mode")
	}
	if runtime.GOOS == "windows" {
		t.Skip("pty e2e is not supported on windows")
	}
}

func TestLazyghE2E_FakeGHViaPTY(t *testing.T) {
	skipIfNotE2E(t)

	s := e2e.NewSession(t, os.Args[0])
	defer s.CloseAndWait()

	s.WaitOutputContains("Fix bug", 5*time.Second)
	s.AssertLogContainsAll("repo view", "pr list")

	openPRDetailAndWait(t, s, 4*time.Second)
	s.AssertLogContainsAll("pr view")

	openPRDiffAndWait(t, s, 4*time.Second)
	s.AssertLogContainsAll("pr diff")
}

func TestLazyghE2E_ReviewFlow(t *testing.T) {
	skipIfNotE2E(t)

	s := e2e.NewSession(t, os.Args[0])
	defer s.CloseAndWait()

	s.WaitOutputContains("Fix bug", 5*time.Second)

	openPRDiffAndWait(t, s, 4*time.Second)

	// move focus from diff files panel to diff content panel;
	// wait for the Diff panel to show an active (green) border confirming focus change
	s.WriteInputAndWaitOutputContains([]byte("l"), "\x1b[32m Diff ", 3*time.Second)

	// navigate to the commentable DELETE line at index 4 (4 j-presses from index 0);
	// wait for the DELETE line (-1 location prefix with reverse-video highlight) to be selected
	s.WriteInputAndWaitOutputContains([]byte("j"), "\x1b[7m  \x1b[0m\x1b[31m--- a/main.go", 3*time.Second)
	s.WriteInputAndWaitOutputContains([]byte("j"), "\x1b[7m  \x1b[0m\x1b[32m+++ b/main.go", 3*time.Second)
	s.WriteInputAndWaitOutputContains([]byte("j"), "\x1b[7m  \x1b[0m\x1b[36m@@ -1 +1 @@", 3*time.Second)
	s.WriteInputAndWaitOutputContains([]byte("j"), "\x1b[7m-1", 3*time.Second)

	// begin comment input
	s.WriteInputAndWaitOutputContains([]byte{13}, "Comment Input", 3*time.Second)

	// type comment body
	s.WriteInput([]byte("e2e comment"))

	// save comment with Ctrl+S; waits for the addPullRequestReviewThread API call
	s.WriteInputAndWaitLogContains([]byte{19}, "addPullRequestReviewThread", 5*time.Second)

	// submit review with Ctrl+R; waits for the submitPullRequestReview API call
	s.WriteInputAndWaitLogContains([]byte{18}, "submitPullRequestReview", 5*time.Second)

	s.AssertLogContainsAll(
		"headRefOid",
		"addPullRequestReview(",
	)
}

func openPRDetailAndWait(t *testing.T, s *e2e.Session, timeout time.Duration) {
	t.Helper()
	s.WriteInputAndWaitOutputContains([]byte("r"), "PR detail", timeout)
}

func openPRDiffAndWait(t *testing.T, s *e2e.Session, timeout time.Duration) {
	t.Helper()
	s.WriteInputAndWaitOutputContains([]byte("d"), "main.go", timeout)
}
