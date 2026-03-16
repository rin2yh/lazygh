package fake

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	gh := Gh{}

	t.Run("ok", func(t *testing.T) {
		args, err := gh.ParseArgs([]string{"test-bin", "-test.run=TestFakeProcess", "--", "gh", "repo", "view"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got, want := strings.Join(args, " "), "gh repo view"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("missing separator", func(t *testing.T) {
		_, err := gh.ParseArgs([]string{"test-bin", "gh", "repo", "view"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestLog(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "fake-gh.log")
	gh := Gh{LogPath: logPath}

	if err := gh.Log([]string{"repo", "view"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := gh.Log([]string{"pr", "list"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	got := string(b)
	if !strings.Contains(got, "repo view\n") {
		t.Fatalf("expected repo view log, got %q", got)
	}
	if !strings.Contains(got, "pr list\n") {
		t.Fatalf("expected pr list log, got %q", got)
	}

	t.Run("without log path", func(t *testing.T) {
		if err := (Gh{}).Log([]string{"repo", "view"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestKey(t *testing.T) {
	gh := Gh{}

	tests := []struct {
		name    string
		args    []string
		wantKey string
		wantOK  bool
	}{
		{"with gh prefix", []string{"gh", "repo", "view", "--json", "nameWithOwner"}, "repo view", true},
		{"without gh prefix", []string{"pr", "list"}, "pr list", true},
		{"too short", []string{"gh"}, "", false},
		{"graphql headRefOid", []string{"api", "graphql", "-f", "query=query($owner:String!,$name:String!,$number:Int!){repository(owner:$owner,name:$name){pullRequest(number:$number){id headRefOid}}}"}, "api graphql headRefOid", true},
		{"graphql addPullRequestReview", []string{"api", "graphql", "-f", "query=mutation($pullRequestId:ID!,$commitOID:GitObjectID!){addPullRequestReview(input:{pullRequestId:$pullRequestId,commitOID:$commitOID}){pullRequestReview{id}}}"}, "api graphql addPullRequestReview", true},
		{"graphql addPullRequestReviewThread", []string{"api", "graphql", "-f", "query=mutation($pullRequestReviewId:ID!,...){addPullRequestReviewThread(...)...}"}, "api graphql addPullRequestReviewThread", true},
		{"graphql submitPullRequestReview", []string{"api", "graphql", "-f", "query=mutation($pullRequestReviewId:ID!){submitPullRequestReview(...)...}"}, "api graphql submitPullRequestReview", true},
		{"graphql without query", []string{"api", "graphql", "-f", "owner=foo"}, "api graphql", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, ok := gh.Key(tt.args)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && key != tt.wantKey {
				t.Fatalf("got %q, want %q", key, tt.wantKey)
			}
		})
	}
}

func TestFind(t *testing.T) {
	gh := Gh{Table: map[string]Response{"repo view": {Stdout: "ok", ExitCode: 0}}}

	tests := []struct {
		name       string
		key        string
		wantOK     bool
		wantStdout string
	}{
		{"hit", "repo view", true, "ok"},
		{"miss", "pr list", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, ok := gh.Find(tt.key)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && resp.Stdout != tt.wantStdout {
				t.Fatalf("got %q, want %q", resp.Stdout, tt.wantStdout)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	gh := Gh{}

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("stdout pipe failed: %v", err)
	}
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		t.Fatalf("stderr pipe failed: %v", err)
	}

	os.Stdout = stdoutW
	os.Stderr = stderrW

	gh.Write(Response{Stdout: "out", Stderr: "err"})

	_ = stdoutW.Close()
	_ = stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	outBytes, err := io.ReadAll(stdoutR)
	if err != nil {
		t.Fatalf("read stdout failed: %v", err)
	}
	errBytes, err := io.ReadAll(stderrR)
	if err != nil {
		t.Fatalf("read stderr failed: %v", err)
	}

	if got := string(outBytes); got != "out" {
		t.Fatalf("stdout got %q, want %q", got, "out")
	}
	if got := string(errBytes); got != "err" {
		t.Fatalf("stderr got %q, want %q", got, "err")
	}
}
