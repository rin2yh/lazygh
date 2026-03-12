package gh

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type ClientInterface interface {
	ResolveCurrentRepo() (string, error)
	ListPRs(repo string) ([]PRItem, error)
	ViewPR(repo string, number int) (string, error)
	DiffPR(repo string, number int) (string, error)
}

type PRItem struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	State     string   `json:"state"`
	IsDraft   bool     `json:"isDraft"`
	Assignees []GHUser `json:"assignees"`
}

type GHUser struct {
	Login string `json:"login"`
}

type Client struct {
	execCommand func(name string, args ...string) *exec.Cmd
}

func NewClient() *Client {
	return &Client{execCommand: exec.Command}
}

func ValidateCLI() error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("gh CLI is required but was not found in PATH: %w", err)
	}
	return nil
}

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

func (c *Client) runCommand(args ...string) ([]byte, error) {
	cmd := c.execCommand("gh", args...)
	cmd.Env = withGHCommandEnv(cmd.Env)
	return cmd.Output()
}

func (c *Client) runJSON(dst any, args ...string) error {
	out, err := c.runCommand(args...)
	if err != nil {
		return err
	}
	return json.Unmarshal(out, dst)
}

func (c *Client) ResolveCurrentRepo() (string, error) {
	type entry struct {
		NameWithOwner string `json:"nameWithOwner"`
	}
	var e entry
	if err := c.runJSON(&e, "repo", "view", "--json", "nameWithOwner"); err != nil {
		return "", err
	}
	repo := strings.TrimSpace(e.NameWithOwner)
	if repo == "" {
		return "", fmt.Errorf("current repository is empty")
	}
	return repo, nil
}

func (c *Client) ListPRs(repo string) ([]PRItem, error) {
	var items []PRItem
	if err := c.runJSON(&items, "pr", "list", "--repo", repo, "--state", "open", "--json", "number,title,state,isDraft,assignees", "--limit", "100"); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Client) ViewPR(repo string, number int) (string, error) {
	out, err := c.runCommand(
		"pr", "view", strconv.Itoa(number),
		"--repo", repo,
		"--json", "title,body,state,isDraft,assignees",
		"--template", "{{.title}}\nStatus: {{if .isDraft}}DRAFT{{else}}{{.state}}{{end}}\nAssignee: {{if .assignees}}{{range $i, $a := .assignees}}{{if $i}}, {{end}}{{$a.login}}{{end}}{{else}}unassigned{{end}}\n\n{{.body}}",
	)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (c *Client) DiffPR(repo string, number int) (string, error) {
	out, err := c.runCommand(
		"pr", "diff", strconv.Itoa(number),
		"--repo", repo,
	)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
