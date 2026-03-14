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
	GetReviewContext(repo string, number int) (ReviewContext, error)
	StartPendingReview(repo string, number int, ctx ReviewContext) (string, error)
	AddReviewComment(repo string, reviewID string, comment ReviewComment) error
	SubmitReview(repo string, reviewID string, body string) error
	DeletePendingReview(repo string, reviewID string) error
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

type ReviewContext struct {
	PullRequestID string
	CommitOID     string
}

type ReviewComment struct {
	Path      string
	Body      string
	Side      DiffSide
	Line      int
	StartSide DiffSide
	StartLine int
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

func (c *Client) runGraphQL(dst any, query string, variables ...string) error {
	args := []string{"api", "graphql", "-f", "query=" + query}
	args = append(args, variables...)
	return c.runJSON(dst, args...)
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

func (c *Client) GetReviewContext(repo string, number int) (ReviewContext, error) {
	owner, name, err := splitRepo(repo)
	if err != nil {
		return ReviewContext{}, err
	}
	var resp struct {
		Data struct {
			Repository struct {
				PullRequest struct {
					ID         string `json:"id"`
					HeadRefOID string `json:"headRefOid"`
				} `json:"pullRequest"`
			} `json:"repository"`
		} `json:"data"`
	}
	query := `query($owner:String!,$name:String!,$number:Int!){repository(owner:$owner,name:$name){pullRequest(number:$number){id headRefOid}}}`
	if err := c.runGraphQL(
		&resp,
		query,
		"-f", "owner="+owner,
		"-f", "name="+name,
		"-F", "number="+strconv.Itoa(number),
	); err != nil {
		return ReviewContext{}, err
	}
	ctx := ReviewContext{
		PullRequestID: strings.TrimSpace(resp.Data.Repository.PullRequest.ID),
		CommitOID:     strings.TrimSpace(resp.Data.Repository.PullRequest.HeadRefOID),
	}
	if ctx.PullRequestID == "" || ctx.CommitOID == "" {
		return ReviewContext{}, fmt.Errorf("review context is incomplete")
	}
	return ctx, nil
}

func (c *Client) StartPendingReview(_ string, _ int, ctx ReviewContext) (string, error) {
	var resp struct {
		Data struct {
			AddPullRequestReview struct {
				PullRequestReview struct {
					ID string `json:"id"`
				} `json:"pullRequestReview"`
			} `json:"addPullRequestReview"`
		} `json:"data"`
	}
	query := `mutation($pullRequestId:ID!,$commitOID:GitObjectID!){addPullRequestReview(input:{pullRequestId:$pullRequestId,commitOID:$commitOID}){pullRequestReview{id}}}`
	if err := c.runGraphQL(
		&resp,
		query,
		"-f", "pullRequestId="+ctx.PullRequestID,
		"-f", "commitOID="+ctx.CommitOID,
	); err != nil {
		return "", err
	}
	reviewID := strings.TrimSpace(resp.Data.AddPullRequestReview.PullRequestReview.ID)
	if reviewID == "" {
		return "", fmt.Errorf("pending review id is empty")
	}
	return reviewID, nil
}

func (c *Client) AddReviewComment(_ string, reviewID string, comment ReviewComment) error {
	if strings.TrimSpace(reviewID) == "" {
		return fmt.Errorf("review id is empty")
	}
	if strings.TrimSpace(comment.Path) == "" {
		return fmt.Errorf("comment path is empty")
	}
	if strings.TrimSpace(comment.Body) == "" {
		return fmt.Errorf("comment body is empty")
	}
	if comment.Line <= 0 {
		return fmt.Errorf("comment line is invalid")
	}

	query := `mutation($pullRequestReviewId:ID!,$body:String!,$path:String!,$line:Int!,$side:DiffSide!,$startLine:Int,$startSide:DiffSide){addPullRequestReviewThread(input:{pullRequestReviewId:$pullRequestReviewId,body:$body,path:$path,line:$line,side:$side,startLine:$startLine,startSide:$startSide}){thread{id}}}`
	args := []string{
		"-f", "pullRequestReviewId=" + reviewID,
		"-f", "body=" + comment.Body,
		"-f", "path=" + comment.Path,
		"-F", "line=" + strconv.Itoa(comment.Line),
		"-f", "side=" + string(comment.Side),
	}
	if comment.StartLine > 0 {
		args = append(args, "-F", "startLine="+strconv.Itoa(comment.StartLine))
		if comment.StartSide != "" {
			args = append(args, "-f", "startSide="+string(comment.StartSide))
		}
	}
	var resp struct {
		Data struct {
			AddPullRequestReviewThread struct {
				Thread struct {
					ID string `json:"id"`
				} `json:"thread"`
			} `json:"addPullRequestReviewThread"`
		} `json:"data"`
	}
	return c.runGraphQL(&resp, query, args...)
}

func (c *Client) SubmitReview(_ string, reviewID string, body string) error {
	var resp struct {
		Data struct {
			SubmitPullRequestReview struct {
				PullRequestReview struct {
					ID string `json:"id"`
				} `json:"pullRequestReview"`
			} `json:"submitPullRequestReview"`
		} `json:"data"`
	}
	query := `mutation($pullRequestReviewId:ID!,$body:String!){submitPullRequestReview(input:{pullRequestReviewId:$pullRequestReviewId,event:COMMENT,body:$body}){pullRequestReview{id}}}`
	return c.runGraphQL(
		&resp,
		query,
		"-f", "pullRequestReviewId="+reviewID,
		"-f", "body="+body,
	)
}

func (c *Client) DeletePendingReview(_ string, reviewID string) error {
	var resp struct {
		Data struct {
			DeletePullRequestReview struct {
				ClientMutationID string `json:"clientMutationId"`
			} `json:"deletePullRequestReview"`
		} `json:"data"`
	}
	query := `mutation($pullRequestReviewId:ID!){deletePullRequestReview(input:{pullRequestReviewId:$pullRequestReviewId}){clientMutationId}}`
	return c.runGraphQL(
		&resp,
		query,
		"-f", "pullRequestReviewId="+reviewID,
	)
}

func splitRepo(repo string) (string, string, error) {
	parts := strings.Split(strings.TrimSpace(repo), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo: %q", repo)
	}
	return parts[0], parts[1], nil
}
