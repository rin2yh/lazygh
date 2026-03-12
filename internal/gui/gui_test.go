package gui

import (
	"errors"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

type mockClient struct {
	repo   string
	prs    []gh.PRItem
	prView string
	err    error
}

func (m *mockClient) ResolveCurrentRepo() (string, error)    { return m.repo, m.err }
func (m *mockClient) ListPRs(_ string) ([]gh.PRItem, error)  { return m.prs, m.err }
func (m *mockClient) ViewPR(_ string, _ int) (string, error) { return m.prView, m.err }

func newTestGuiWithClient(client gh.ClientInterface) *Gui {
	g, _ := NewGui(config.Default(), client)
	return g
}

func TestNavigatePRList(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.ApplyPRsResult("owner/repo", []core.Item{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}}, nil)

	g.navigateDown()
	if g.state.PRsSelected != 1 {
		t.Fatalf("got %d, want 1", g.state.PRsSelected)
	}

	g.navigateUp()
	if g.state.PRsSelected != 0 {
		t.Fatalf("got %d, want 0", g.state.PRsSelected)
	}
}

func TestApplyPRsResult(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.BeginLoadPRs()
	g.applyPRsResult(prsLoadedMsg{
		repo: "owner/repo",
		prs:  []core.Item{{Number: 1, Title: "Fix bug"}},
	})
	if g.state.PRsLoading {
		t.Fatal("loading should be false")
	}
	if g.state.Repo != "owner/repo" {
		t.Fatalf("got %q, want owner/repo", g.state.Repo)
	}
	if got := core.FormatPRItem(g.state.PRs[0]); got != "PR #1 Fix bug" {
		t.Fatalf("unexpected pr row: %q", got)
	}
	if g.state.Loading != core.LoadingNone {
		t.Fatalf("got %v, want %v", g.state.Loading, core.LoadingNone)
	}
}

func TestApplyPRsResult_Error(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.BeginLoadPRs()
	g.applyPRsResult(prsLoadedMsg{err: errors.New("boom")})
	if g.state.DetailContent == "" {
		t.Fatal("error message should be shown")
	}
}

func TestApplyDetailResult(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.applyDetailResult(detailLoadedMsg{content: "hello"})
	if g.state.DetailContent != "hello" {
		t.Fatalf("got %q, want hello", g.state.DetailContent)
	}
}

func TestApplyDetailResult_Error(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.applyDetailResult(detailLoadedMsg{err: errors.New("boom")})
	if g.state.DetailContent == "" {
		t.Fatal("error message should be shown")
	}
}

func TestModelInitLoadsPRs(t *testing.T) {
	mc := &mockClient{repo: "owner/repo", prs: []gh.PRItem{{Number: 2, Title: "p"}}}
	g := newTestGuiWithClient(mc)
	m := &model{gui: g}

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("expected init command")
	}
	msg := cmd().(prsLoadedMsg)
	if msg.err != nil || msg.repo != "owner/repo" || len(msg.prs) != 1 {
		t.Fatalf("unexpected prsLoadedMsg: %+v", msg)
	}
}

func TestModelHandleEnterDetail(t *testing.T) {
	mc := &mockClient{prView: "detail"}
	g := newTestGuiWithClient(mc)
	g.state.ApplyPRsResult("owner/repo", []core.Item{{Number: 1, Title: "x"}}, nil)
	m := &model{gui: g}

	cmd := m.handleEnter()
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd().(detailLoadedMsg)
	if msg.err != nil || msg.content != "detail" {
		t.Fatalf("unexpected msg: %+v", msg)
	}
}

func TestWrapText(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
		want    string
	}{
		{
			name:    "wrap long line",
			content: "abcdefghij",
			width:   4,
			want:    "abcd\nefgh\nij",
		},
		{
			name:    "keep existing line breaks",
			content: "abcde\nfghij",
			width:   3,
			want:    "abc\nde\nfgh\nij",
		},
		{
			name:    "no wrap when width is enough",
			content: "abc",
			width:   10,
			want:    "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wrapText(tt.content, tt.width); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderPRPanel_EmptyPlaceholder(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	lines := g.renderPRPanel("PRs", 3)

	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	if lines[1] != "No pull requests" {
		t.Fatalf("got %q, want %q", lines[1], "No pull requests")
	}
}

func TestRenderPRPanel_LoadingHidesLoadingText(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.BeginLoadPRs()
	lines := g.renderPRPanel("PRs", 3)

	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	if lines[1] != "" {
		t.Fatalf("got %q, want empty line", lines[1])
	}
}

func TestRenderRepoPanel_ShowsRepoName(t *testing.T) {
	g := newTestGuiWithClient(&mockClient{})
	g.state.Repo = "owner/repo"
	lines := g.renderRepoPanel("Repository", 2)

	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}
	if lines[0] != " Repository " {
		t.Fatalf("got %q, want %q", lines[0], " Repository ")
	}
	if lines[1] != "owner/repo" {
		t.Fatalf("got %q, want %q", lines[1], "owner/repo")
	}
}
