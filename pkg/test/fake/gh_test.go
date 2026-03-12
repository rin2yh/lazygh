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

	t.Run("with gh prefix", func(t *testing.T) {
		key, ok := gh.Key([]string{"gh", "repo", "view", "--json", "nameWithOwner"})
		if !ok {
			t.Fatal("expected ok, got false")
		}
		if key != "repo view" {
			t.Fatalf("got %q, want %q", key, "repo view")
		}
	})

	t.Run("without gh prefix", func(t *testing.T) {
		key, ok := gh.Key([]string{"pr", "list"})
		if !ok {
			t.Fatal("expected ok, got false")
		}
		if key != "pr list" {
			t.Fatalf("got %q, want %q", key, "pr list")
		}
	})

	t.Run("too short", func(t *testing.T) {
		_, ok := gh.Key([]string{"gh"})
		if ok {
			t.Fatal("expected false, got true")
		}
	})
}

func TestFind(t *testing.T) {
	gh := Gh{Table: map[string]Response{"repo view": {Stdout: "ok", ExitCode: 0}}}

	t.Run("hit", func(t *testing.T) {
		resp, ok := gh.Find("repo view")
		if !ok {
			t.Fatal("expected hit, got miss")
		}
		if resp.Stdout != "ok" {
			t.Fatalf("got %q, want %q", resp.Stdout, "ok")
		}
	})

	t.Run("miss", func(t *testing.T) {
		_, ok := gh.Find("pr list")
		if ok {
			t.Fatal("expected miss, got hit")
		}
	})
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
