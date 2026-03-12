package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
)

type Session struct {
	t        *testing.T
	logPath  string
	runCmd   *exec.Cmd
	ptmx     *os.File
	output   bytes.Buffer
	copyDone chan struct{}
}

func NewSession(t *testing.T, processBin string) *Session {
	t.Helper()
	tmpDir := t.TempDir()
	fakeBin := filepath.Join(tmpDir, "fake-bin")
	if err := os.MkdirAll(fakeBin, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	logPath := filepath.Join(tmpDir, "fake-gh.log")
	ghPath := filepath.Join(fakeBin, "gh")
	wrapper := fmt.Sprintf("#!/bin/sh\nexec '%s' -test.run '^TestFakeGHProcess$' -- \"$@\"\n", processBin)
	if err := os.WriteFile(ghPath, []byte(wrapper), 0o755); err != nil {
		t.Fatalf("write fake gh wrapper failed: %v", err)
	}

	binPath := filepath.Join(tmpDir, "lazygh-test-bin")
	buildCommand(t, "go", "build", "-o", binPath, ".")

	runCmd := exec.Command(binPath)
	runCmd.Env = append(
		os.Environ(),
		"PATH="+fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"FAKE_GH_LOG="+logPath,
		"GO_WANT_FAKE_GH_PROCESS=1",
		"TERM=dumb",
		"NO_COLOR=1",
	)

	ptmx, err := pty.Start(runCmd)
	if err != nil {
		t.Fatalf("pty start failed: %v", err)
	}

	s := &Session{
		t:        t,
		logPath:  logPath,
		runCmd:   runCmd,
		ptmx:     ptmx,
		copyDone: make(chan struct{}),
	}
	go func() {
		_, _ = s.output.ReadFrom(s.ptmx)
		close(s.copyDone)
	}()
	return s
}

func buildCommand(t *testing.T, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s failed: %v\n%s", strings.Join(append([]string{name}, args...), " "), err, string(out))
	}
}

func (s *Session) WaitLogContains(want string, timeout time.Duration) {
	s.t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if s.HasLogEntry(want) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	s.t.Fatalf("fake gh log did not contain %q in time. output:\n%s", want, s.output.String())
}

func (s *Session) HasLogEntry(want string) bool {
	b, err := os.ReadFile(s.logPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(b), want)
}

func (s *Session) WriteInput(in []byte) {
	s.t.Helper()
	if _, err := s.ptmx.Write(in); err != nil {
		s.t.Fatalf("write input failed: %v", err)
	}
}

func (s *Session) CloseAndWait() {
	s.t.Helper()
	_, _ = s.ptmx.Write([]byte("q"))
	time.Sleep(50 * time.Millisecond)
	_, _ = s.ptmx.Write([]byte{3})

	exitDone := make(chan error, 1)
	go func() {
		exitDone <- s.runCmd.Wait()
	}()

	select {
	case err := <-exitDone:
		if err != nil && !strings.Contains(err.Error(), "interrupt") {
			s.t.Fatalf("lazygh exited with error: %v", err)
		}
	case <-time.After(5 * time.Second):
		_ = s.runCmd.Process.Kill()
		s.t.Fatal("lazygh did not exit")
	}

	_ = s.ptmx.Close()
	<-s.copyDone
}

func (s *Session) AssertLogContainsAll(wants ...string) {
	s.t.Helper()
	logBytes, err := os.ReadFile(s.logPath)
	if err != nil {
		s.t.Fatalf("read fake gh log failed: %v", err)
	}
	logText := string(logBytes)
	for _, want := range wants {
		if !strings.Contains(logText, want) {
			s.t.Fatalf("fake gh should be called with %q, got:\n%s", want, logText)
		}
	}
}
