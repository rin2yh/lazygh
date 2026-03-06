package gh

import (
	"encoding/json"
	"os/exec"
	"strconv"
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

func (c *Client) runJSON(dst any, args ...string) error {
	out, err := c.execCommand("gh", args...).Output()
	if err != nil {
		return err
	}
	return json.Unmarshal(out, dst)
}

func (c *Client) ListRepos() ([]string, error) {
	type entry struct {
		NameWithOwner string `json:"nameWithOwner"`
	}
	out, err := c.execCommand("gh", "repo", "list", "--json", "nameWithOwner", "--limit", "100").Output()
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
	out, err := c.execCommand("gh", "pr", "view", strconv.Itoa(number), "--repo", repo).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (c *Client) ViewIssue(repo string, number int) (string, error) {
	out, err := c.execCommand("gh", "issue", "view", strconv.Itoa(number), "--repo", repo).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
