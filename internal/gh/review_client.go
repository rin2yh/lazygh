package gh

import (
	"fmt"
	"strconv"
	"strings"
)

// ReviewContext holds the identifiers needed to start or interact with a GitHub PR review.
type ReviewContext struct {
	PullRequestID string
	CommitOID     string
}

// ReviewEvent represents the type of review submission event.
type ReviewEvent string

const (
	ReviewEventComment        ReviewEvent = "COMMENT"
	ReviewEventApprove        ReviewEvent = "APPROVE"
	ReviewEventRequestChanges ReviewEvent = "REQUEST_CHANGES"
)

// ReviewComment holds the parameters for adding a review thread comment.
type ReviewComment struct {
	Path      string
	Body      string
	Side      DiffSide
	Line      int
	StartSide DiffSide
	StartLine int
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
	if err := c.api.RunGraphQL(
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
	if err := c.api.RunGraphQL(
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

func (c *Client) AddReviewComment(_ string, reviewID string, comment ReviewComment) (string, error) {
	if strings.TrimSpace(reviewID) == "" {
		return "", fmt.Errorf("review id is empty")
	}
	if strings.TrimSpace(comment.Path) == "" {
		return "", fmt.Errorf("comment path is empty")
	}
	if strings.TrimSpace(comment.Body) == "" {
		return "", fmt.Errorf("comment body is empty")
	}
	if comment.Line <= 0 {
		return "", fmt.Errorf("comment line is invalid")
	}

	query := `mutation($pullRequestReviewId:ID!,$body:String!,$path:String!,$line:Int!,$side:DiffSide!,$startLine:Int,$startSide:DiffSide){addPullRequestReviewThread(input:{pullRequestReviewId:$pullRequestReviewId,body:$body,path:$path,line:$line,side:$side,startLine:$startLine,startSide:$startSide}){thread{comments(first:1){nodes{id}}}}}`
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
					Comments struct {
						Nodes []struct {
							ID string `json:"id"`
						} `json:"nodes"`
					} `json:"comments"`
				} `json:"thread"`
			} `json:"addPullRequestReviewThread"`
		} `json:"data"`
	}
	if err := c.api.RunGraphQL(&resp, query, args...); err != nil {
		return "", err
	}
	nodes := resp.Data.AddPullRequestReviewThread.Thread.Comments.Nodes
	if len(nodes) == 0 {
		return "", nil
	}
	return strings.TrimSpace(nodes[0].ID), nil
}

func (c *Client) DeletePendingReviewComment(commentID string) error {
	if strings.TrimSpace(commentID) == "" {
		return fmt.Errorf("comment id is empty")
	}
	var resp struct {
		Data struct {
			DeletePullRequestReviewComment struct {
				ClientMutationID string `json:"clientMutationId"`
			} `json:"deletePullRequestReviewComment"`
		} `json:"data"`
	}
	query := `mutation($id:ID!){deletePullRequestReviewComment(input:{id:$id}){clientMutationId}}`
	return c.api.RunGraphQL(&resp, query, "-f", "id="+commentID)
}

func (c *Client) UpdatePendingReviewComment(commentID string, body string) error {
	if strings.TrimSpace(commentID) == "" {
		return fmt.Errorf("comment id is empty")
	}
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("comment body is empty")
	}
	var resp struct {
		Data struct {
			UpdatePullRequestReviewComment struct {
				PullRequestReviewComment struct {
					ID string `json:"id"`
				} `json:"pullRequestReviewComment"`
			} `json:"updatePullRequestReviewComment"`
		} `json:"data"`
	}
	query := `mutation($id:ID!,$body:String!){updatePullRequestReviewComment(input:{pullRequestReviewCommentId:$id,body:$body}){pullRequestReviewComment{id}}}`
	return c.api.RunGraphQL(&resp, query, "-f", "id="+commentID, "-f", "body="+body)
}

func (c *Client) SubmitReview(_ string, reviewID string, event ReviewEvent, body string) error {
	var resp struct {
		Data struct {
			SubmitPullRequestReview struct {
				PullRequestReview struct {
					ID string `json:"id"`
				} `json:"pullRequestReview"`
			} `json:"submitPullRequestReview"`
		} `json:"data"`
	}
	query := `mutation($pullRequestReviewId:ID!,$event:PullRequestReviewEvent!,$body:String!){submitPullRequestReview(input:{pullRequestReviewId:$pullRequestReviewId,event:$event,body:$body}){pullRequestReview{id}}}`
	return c.api.RunGraphQL(
		&resp,
		query,
		"-f", "pullRequestReviewId="+reviewID,
		"-f", "event="+string(event),
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
	return c.api.RunGraphQL(
		&resp,
		query,
		"-f", "pullRequestReviewId="+reviewID,
	)
}
