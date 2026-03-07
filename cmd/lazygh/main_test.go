package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestFakeGHHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_FAKE_GH_HELPER") != "1" {
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
		fmt.Fprintln(os.Stderr, "missing -- separator")
		os.Exit(1)
	}
	ghArgs := args[sep+1:]
	if logPath := os.Getenv("FAKE_GH_LOG"); logPath != "" {
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err == nil {
			_, _ = f.WriteString(strings.Join(ghArgs, " ") + "\n")
			_ = f.Close()
		}
	}

	switch {
	case len(ghArgs) >= 2 && ghArgs[0] == "repo" && ghArgs[1] == "list":
		fmt.Print(`[{"nameWithOwner":"owner/repo1"}]`)
		os.Exit(0)
	case len(ghArgs) >= 2 && ghArgs[0] == "issue" && ghArgs[1] == "list":
		fmt.Print(`[{"number":10,"title":"Issue one"}]`)
		os.Exit(0)
	case len(ghArgs) >= 2 && ghArgs[0] == "pr" && ghArgs[1] == "list":
		fmt.Print(`[{"number":1,"title":"Fix bug"}]`)
		os.Exit(0)
	case len(ghArgs) >= 2 && ghArgs[0] == "issue" && ghArgs[1] == "view":
		fmt.Print("Issue detail")
		os.Exit(0)
	case len(ghArgs) >= 2 && ghArgs[0] == "pr" && ghArgs[1] == "view":
		fmt.Print("PR detail")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unexpected gh args: %s\n", strings.Join(ghArgs, " "))
		os.Exit(1)
	}
}

func TestLazyghE2E_FakeGHViaPTY(t *testing.T) {
	if testing.Short() {
		t.Skip("skip e2e in short mode")
	}
	if runtime.GOOS == "windows" {
		t.Skip("pty e2e is not supported on windows")
	}

	s := newE2ESession(t, os.Args[0])
	defer s.closeAndWait()

	s.waitLogContains("repo list", 3*time.Second)
	s.writeInput([]byte("\r"))
	s.waitLogContains("issue list", 3*time.Second)
	s.waitLogContains("pr list", 3*time.Second)
	s.writeInput([]byte{3})
	s.assertLogContainsAll("repo list", "issue list", "pr list")
}
