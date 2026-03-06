package gh

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type ClientInterface interface {
	ListRepos() ([]string, error)
	ListPRs(repo string) ([]PRItem, error)
	ListIssues(repo string) ([]IssueItem, error)
	ViewPR(repo string, number int) (string, error)
	ViewIssue(repo string, number int) (string, error)
}

type PRItem struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

type IssueItem struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

type Client struct {
	execCommand func(name string, args ...string) *exec.Cmd
}

func NewClient() *Client {
	return &Client{execCommand: exec.Command}
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

func sanitizeOutput(out []byte) string {
	return strings.ToValidUTF8(string(out), "")
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

func (c *Client) ListRepos() ([]string, error) {
	type entry struct {
		NameWithOwner string `json:"nameWithOwner"`
	}
	out, err := c.runCommand("repo", "list", "--json", "nameWithOwner", "--limit", "100")
	if err != nil {
		return nil, err
	}
	var entries []entry
	if err := json.Unmarshal(out, &entries); err != nil {
		return nil, err
	}
	repos := make([]string, len(entries))
	for i, e := range entries {
		repos[i] = e.NameWithOwner
	}
	return repos, nil
}

func (c *Client) ListPRs(repo string) ([]PRItem, error) {
	var items []PRItem
	if err := c.runJSON(&items, "pr", "list", "--repo", repo, "--json", "number,title", "--limit", "100"); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Client) ListIssues(repo string) ([]IssueItem, error) {
	var items []IssueItem
	if err := c.runJSON(&items, "issue", "list", "--repo", repo, "--json", "number,title", "--limit", "100"); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Client) ViewPR(repo string, number int) (string, error) {
	out, err := c.runCommand(
		"pr", "view", strconv.Itoa(number),
		"--repo", repo,
		"--json", "title,body",
		"--template", "{{.title}}\n\n{{.body}}",
	)
	if err != nil {
		return "", err
	}
	return sanitizeOutput(out), nil
}

func (c *Client) ViewIssue(repo string, number int) (string, error) {
	out, err := c.runCommand(
		"issue", "view", strconv.Itoa(number),
		"--repo", repo,
		"--json", "title,body",
		"--template", "{{.title}}\n\n{{.body}}",
	)
	if err != nil {
		return "", err
	}
	return sanitizeOutput(out), nil
}
