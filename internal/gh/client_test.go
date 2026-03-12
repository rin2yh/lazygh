package gh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	repo, err := c.ResolveCurrentRepo()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo != "owner/repo1" {
		t.Fatalf("got %q, want %q", repo, "owner/repo1")
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
	tests := []struct {
		name      string
		setupPath func(t *testing.T) string
		wantErr   bool
	}{
		{
			name: "ok",
			setupPath: func(t *testing.T) string {
				dir := t.TempDir()
				ghPath := filepath.Join(dir, "gh")
				if err := os.WriteFile(ghPath, []byte("#!/bin/sh\n"), 0o755); err != nil {
					t.Fatalf("write fake gh failed: %v", err)
				}
				return dir
			},
		},
		{
			name: "error",
			setupPath: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PATH", tt.setupPath(t))

			err := ValidateCLI()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "gh CLI is required") {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
