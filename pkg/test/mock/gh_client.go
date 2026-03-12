package mock

import (
	"sync"

	"github.com/rin2yh/lazygh/internal/gh"
)

type GHClient struct {
	Repo   string
	PRs    []gh.PRItem
	PRView string
	PRDiff string
	Err    error
}

func (m *GHClient) ResolveCurrentRepo() (string, error) {
	return m.Repo, m.Err
}

func (m *GHClient) ListPRs(_ string) ([]gh.PRItem, error) {
	return m.PRs, m.Err
}

func (m *GHClient) ViewPR(_ string, _ int) (string, error) {
	return m.PRView, m.Err
}

func (m *GHClient) DiffPR(_ string, _ int) (string, error) {
	return m.PRDiff, m.Err
}

type ControlledGHClient struct {
	Repo   string
	PRs    []gh.PRItem
	PRView string
	PRDiff string
	Err    error

	ResolveCalled chan struct{}
	PRsCalled     chan struct{}
	DetailCalled  chan struct{}
	DiffCalled    chan struct{}

	ReleaseResolve <-chan struct{}
	ReleasePRs     <-chan struct{}
	ReleaseDetail  <-chan struct{}
	ReleaseDiff    <-chan struct{}

	resolveOnce sync.Once
	prsOnce     sync.Once
	detailOnce  sync.Once
	diffOnce    sync.Once
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

func (c *ControlledGHClient) ListPRs(_ string) ([]gh.PRItem, error) {
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
