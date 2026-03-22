package gh

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// PRItem represents a pull request item from gh CLI.
type PRItem struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	State     string   `json:"state"`
	IsDraft   bool     `json:"isDraft"`
	Assignees []GHUser `json:"assignees"`
}

// GHUser represents a GitHub user.
type GHUser struct {
	Login string `json:"login"`
}

// Client wraps the gh CLI to interact with GitHub PRs and reviews.
type Client struct {
	runner *commandRunner
	api    *apiClient
}

// NewClient returns a Client backed by the real gh binary.
func NewClient() *Client {
	runner := &commandRunner{execCommand: exec.Command}
	return &Client{
		runner: runner,
		api:    &apiClient{runner: runner},
	}
}

// ValidateCLI returns an error if the gh binary is not found in PATH.
func ValidateCLI() error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("gh CLI is required but was not found in PATH: %w", err)
	}
	return nil
}

func (c *Client) ResolveCurrentRepo() (string, error) {
	type entry struct {
		NameWithOwner string `json:"nameWithOwner"`
	}
	var e entry
	if err := c.api.RunJSON(&e, "repo", "view", "--json", "nameWithOwner"); err != nil {
		return "", err
	}
	repo := strings.TrimSpace(e.NameWithOwner)
	if repo == "" {
		return "", fmt.Errorf("current repository is empty")
	}
	return repo, nil
}

func (c *Client) ListPRs(repo string, state string) ([]PRItem, error) {
	var items []PRItem
	if err := c.api.RunJSON(&items, "pr", "list", "--repo", repo, "--state", state, "--json", "number,title,state,isDraft,assignees", "--limit", "100"); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Client) ViewPR(repo string, number int) (string, error) {
	out, err := c.runner.Run(
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
	out, err := c.runner.Run(
		"pr", "diff", strconv.Itoa(number),
		"--repo", repo,
	)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func splitRepo(repo string) (string, string, error) {
	parts := strings.Split(strings.TrimSpace(repo), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo: %q", repo)
	}
	return parts[0], parts[1], nil
}
