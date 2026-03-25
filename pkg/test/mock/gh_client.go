package mock

import (
	"sync"

	"github.com/rin2yh/lazygh/internal/gh"
)

type GHClient struct {
	Repo             string
	PRs              []gh.PRItem
	PRView           string
	PRDiff           string
	ReviewContext    gh.ReviewContext
	PendingReviewID  string
	PendingCommentID string
	ReviewComments   []gh.ReviewComment
	ReviewThreads    []gh.ReviewThread
	SubmittedReviews []string
	DeletedReviews   []string
	DeletedComments  []string
	UpdatedComments  []string
	ThreadReplies    []string
	Err              error
}

func (m *GHClient) ResolveCurrentRepo() (string, error) {
	return m.Repo, m.Err
}

func (m *GHClient) ListPRs(_ string, _ string) ([]gh.PRItem, error) {
	return m.PRs, m.Err
}

func (m *GHClient) ViewPR(_ string, _ int) (string, error) {
	return m.PRView, m.Err
}

func (m *GHClient) DiffPR(_ string, _ int) (string, error) {
	return m.PRDiff, m.Err
}

func (m *GHClient) GetReviewContext(_ string, _ int) (gh.ReviewContext, error) {
	return m.ReviewContext, m.Err
}

func (m *GHClient) StartPendingReview(_ string, _ int, _ gh.ReviewContext) (string, error) {
	return m.PendingReviewID, m.Err
}

func (m *GHClient) AddReviewComment(_ string, _ string, comment gh.ReviewComment) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	m.ReviewComments = append(m.ReviewComments, comment)
	return m.PendingCommentID, nil
}

func (m *GHClient) DeletePendingReviewComment(commentID string) error {
	if m.Err != nil {
		return m.Err
	}
	m.DeletedComments = append(m.DeletedComments, commentID)
	return nil
}

func (m *GHClient) UpdatePendingReviewComment(commentID string, body string) error {
	if m.Err != nil {
		return m.Err
	}
	m.UpdatedComments = append(m.UpdatedComments, commentID+":"+body)
	return nil
}

func (m *GHClient) SubmitReview(_ string, reviewID string, event gh.ReviewEvent, body string) error {
	if m.Err != nil {
		return m.Err
	}
	m.SubmittedReviews = append(m.SubmittedReviews, reviewID+":"+string(event)+":"+body)
	return nil
}

func (m *GHClient) DeletePendingReview(_ string, reviewID string) error {
	if m.Err != nil {
		return m.Err
	}
	m.DeletedReviews = append(m.DeletedReviews, reviewID)
	return nil
}

func (m *GHClient) GetReviewThreads(_ string, _ int) ([]gh.ReviewThread, error) {
	return m.ReviewThreads, m.Err
}

func (m *GHClient) AddReplyToReviewThread(threadID string, body string) error {
	if m.Err != nil {
		return m.Err
	}
	m.ThreadReplies = append(m.ThreadReplies, threadID+":"+body)
	return nil
}

type ControlledGHClient struct {
	Repo            string
	PRs             []gh.PRItem
	PRView          string
	PRDiff          string
	ReviewContext   gh.ReviewContext
	PendingReviewID string
	Err             error

	ResolveCalled chan struct{}
	PRsCalled     chan struct{}
	DetailCalled  chan struct{}
	DiffCalled    chan struct{}
	ReviewCalled  chan struct{}

	ReleaseResolve <-chan struct{}
	ReleasePRs     <-chan struct{}
	ReleaseDetail  <-chan struct{}
	ReleaseDiff    <-chan struct{}

	resolveOnce sync.Once
	prsOnce     sync.Once
	detailOnce  sync.Once
	diffOnce    sync.Once
	reviewOnce  sync.Once
}

func (c *ControlledGHClient) ResolveCurrentRepo() (string, error) {
	if c.ResolveCalled != nil {
		c.resolveOnce.Do(func() { close(c.ResolveCalled) })
	}
	if c.ReleaseResolve != nil {
		<-c.ReleaseResolve
	}
	if c.Err != nil {
		return "", c.Err
	}
	return c.Repo, nil
}

func (c *ControlledGHClient) ListPRs(_ string, _ string) ([]gh.PRItem, error) {
	if c.PRsCalled != nil {
		c.prsOnce.Do(func() { close(c.PRsCalled) })
	}
	if c.ReleasePRs != nil {
		<-c.ReleasePRs
	}
	if c.Err != nil {
		return nil, c.Err
	}
	return c.PRs, nil
}

func (c *ControlledGHClient) ViewPR(_ string, _ int) (string, error) {
	if c.DetailCalled != nil {
		c.detailOnce.Do(func() { close(c.DetailCalled) })
	}
	if c.ReleaseDetail != nil {
		<-c.ReleaseDetail
	}
	if c.Err != nil {
		return "", c.Err
	}
	return c.PRView, nil
}

func (c *ControlledGHClient) DiffPR(_ string, _ int) (string, error) {
	if c.DiffCalled != nil {
		c.diffOnce.Do(func() { close(c.DiffCalled) })
	}
	if c.ReleaseDiff != nil {
		<-c.ReleaseDiff
	}
	if c.Err != nil {
		return "", c.Err
	}
	return c.PRDiff, nil
}

func (c *ControlledGHClient) GetReviewContext(_ string, _ int) (gh.ReviewContext, error) {
	if c.ReviewCalled != nil {
		c.reviewOnce.Do(func() { close(c.ReviewCalled) })
	}
	if c.Err != nil {
		return gh.ReviewContext{}, c.Err
	}
	return c.ReviewContext, nil
}

func (c *ControlledGHClient) StartPendingReview(_ string, _ int, _ gh.ReviewContext) (string, error) {
	if c.Err != nil {
		return "", c.Err
	}
	return c.PendingReviewID, nil
}

func (c *ControlledGHClient) AddReviewComment(_ string, _ string, _ gh.ReviewComment) (string, error) {
	return "", c.Err
}

func (c *ControlledGHClient) DeletePendingReviewComment(_ string) error {
	return c.Err
}

func (c *ControlledGHClient) UpdatePendingReviewComment(_ string, _ string) error {
	return c.Err
}

func (c *ControlledGHClient) SubmitReview(_ string, _ string, _ gh.ReviewEvent, _ string) error {
	return c.Err
}

func (c *ControlledGHClient) DeletePendingReview(_ string, _ string) error {
	return c.Err
}

func (c *ControlledGHClient) GetReviewThreads(_ string, _ int) ([]gh.ReviewThread, error) {
	return nil, c.Err
}

func (c *ControlledGHClient) AddReplyToReviewThread(_ string, _ string) error {
	return c.Err
}
