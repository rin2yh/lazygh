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
			Stdout:   "PR diff",
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

func TestLazyghE2E_FakeGHViaPTY(t *testing.T) {
	if testing.Short() {
		t.Skip("skip e2e in short mode")
	}
	if runtime.GOOS == "windows" {
		t.Skip("pty e2e is not supported on windows")
	}

	s := e2e.NewSession(t, os.Args[0])
	defer s.CloseAndWait()

	s.WaitLogContains("repo view", 3*time.Second)
	s.WaitLogContains("pr list", 3*time.Second)
	s.AssertLogContainsAll("repo view", "pr list")

	openPRDetailAndWait(t, s, 4*time.Second)
	s.AssertLogContainsAll("pr view")

	openPRDiffAndWait(t, s, 4*time.Second)
	s.AssertLogContainsAll("pr diff")
}

func openPRDetailAndWait(t *testing.T, s *e2e.Session, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		s.WriteInput([]byte("\r"))
		time.Sleep(80 * time.Millisecond)

		if s.HasLogEntry("pr view") {
			return
		}
	}
	t.Fatal("opening pr detail did not trigger pr view in time")
}

func openPRDiffAndWait(t *testing.T, s *e2e.Session, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		s.WriteInput([]byte("d"))
		time.Sleep(80 * time.Millisecond)

		if s.HasLogEntry("pr diff") {
			return
		}
	}
	t.Fatal("switching to diff did not trigger pr diff in time")
}
