package gh

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func withGHCommandEnv(base []string) []string {
	env := os.Environ()
	if len(base) > 0 {
		env = append(env, base...)
	}
	return append(env,
		"NO_COLOR=1",
		"CLICOLOR=0",
		"GH_PAGER=cat",
	)
}

// CommandError is returned when the gh CLI exits with a non-zero status.
type CommandError struct {
	Command []string
	Stderr  string
	Err     error
}

func (e *CommandError) Error() string {
	command := "gh"
	if len(e.Command) > 0 {
		command += " " + strings.Join(e.Command, " ")
	}
	if stderr := strings.TrimSpace(e.Stderr); stderr != "" {
		return fmt.Sprintf("%s failed: %s: %v", command, stderr, e.Err)
	}
	return fmt.Sprintf("%s failed: %v", command, e.Err)
}

func (e *CommandError) Unwrap() error {
	return e.Err
}

// commandRunner executes gh CLI subprocesses.
type commandRunner struct {
	execCommand func(name string, args ...string) *exec.Cmd
}

func (r *commandRunner) Run(args ...string) ([]byte, error) {
	cmd := r.execCommand("gh", args...)
	cmd.Env = withGHCommandEnv(cmd.Env)
	out, err := cmd.Output()
	if err == nil {
		return out, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return nil, &CommandError{
			Command: append([]string(nil), args...),
			Stderr:  string(exitErr.Stderr),
			Err:     err,
		}
	}

	return nil, &CommandError{
		Command: append([]string(nil), args...),
		Err:     err,
	}
}

// apiClient wraps commandRunner with JSON and GraphQL helpers.
type apiClient struct {
	runner *commandRunner
}

func (a *apiClient) RunJSON(dst any, args ...string) error {
	out, err := a.runner.Run(args...)
	if err != nil {
		return err
	}
	return json.Unmarshal(out, dst)
}

func (a *apiClient) RunGraphQL(dst any, query string, variables ...string) error {
	args := []string{"api", "graphql", "-f", "query=" + query}
	args = append(args, variables...)
	return a.RunJSON(dst, args...)
}
