package gui

import (
	"errors"
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func newTestGuiWithClient(client gh.ClientInterface) *Gui {
	g, _ := NewGui(config.Default(), client)
	return g
}

func TestNavigatePRList(t *testing.T) {
	g := newTestGuiWithClient(&testmock.GHClient{})
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
	type want struct {
		repo   string
		prs    []core.Item
		detail string
	}

	tests := []struct {
		name string
		msg  prsLoadedMsg
		want want
	}{
		{
			name: "success",
			msg: prsLoadedMsg{
				repo: "owner/repo",
				prs:  []core.Item{{Number: 1, Title: "Fix bug"}},
			},
			want: want{
				repo:   "owner/repo",
				prs:    []core.Item{{Number: 1, Title: "Fix bug"}},
				detail: "PR #1 Fix bug",
			},
		},
		{
			name: "empty",
			msg: prsLoadedMsg{
				repo: "owner/repo",
			},
			want: want{
				repo:   "owner/repo",
				prs:    nil,
				detail: "No pull requests",
			},
		},
		{
			name: "error",
			msg: prsLoadedMsg{
				err: errors.New("boom"),
			},
			want: want{
				repo:   "",
				prs:    nil,
				detail: "Error loading PRs: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithClient(&testmock.GHClient{})
			g.state.BeginLoadPRs()

			g.applyPRsResult(tt.msg)

			if g.state.PRsLoading {
				t.Fatal("loading should be false")
			}
			if g.state.Loading != core.LoadingNone {
				t.Fatalf("got %v, want %v", g.state.Loading, core.LoadingNone)
			}
			if g.state.Repo != tt.want.repo {
				t.Fatalf("got %q, want %q", g.state.Repo, tt.want.repo)
			}
			if len(g.state.PRs) != len(tt.want.prs) {
				t.Fatalf("got %d, want %d", len(g.state.PRs), len(tt.want.prs))
			}
			for i := range g.state.PRs {
				if g.state.PRs[i].Number != tt.want.prs[i].Number || g.state.PRs[i].Title != tt.want.prs[i].Title {
					t.Fatalf("unexpected PR[%d]: %+v", i, g.state.PRs[i])
				}
			}
			if g.state.DetailContent != tt.want.detail {
				t.Fatalf("got %q, want %q", g.state.DetailContent, tt.want.detail)
			}
		})
	}
}

func TestApplyDetailResult(t *testing.T) {
	type want struct {
		detail string
	}

	tests := []struct {
		name string
		msg  detailLoadedMsg
		want want
	}{
		{
			name: "success",
			msg: detailLoadedMsg{
				content: "hello",
			},
			want: want{
				detail: "hello",
			},
		},
		{
			name: "error",
			msg: detailLoadedMsg{
				err: errors.New("boom"),
			},
			want: want{
				detail: "Error loading detail: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithClient(&testmock.GHClient{})
			g.state.Loading = core.LoadingDetail

			g.applyDetailResult(tt.msg)

			if g.state.Loading != core.LoadingNone {
				t.Fatalf("got %v, want %v", g.state.Loading, core.LoadingNone)
			}
			if g.state.DetailContent != tt.want.detail {
				t.Fatalf("got %q, want %q", g.state.DetailContent, tt.want.detail)
			}
		})
	}
}

func TestModelInitLoadsPRs(t *testing.T) {
	mc := &testmock.GHClient{Repo: "owner/repo", PRs: []gh.PRItem{{Number: 2, Title: "p"}}}
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
	mc := &testmock.GHClient{PRView: "detail"}
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

func TestFramePanel(t *testing.T) {
	got := framePanel("Repo", false, []string{"body"}, 10, 3)
	want := []string{
		"┌ Repo ──┐",
		"│body    │",
		"└────────┘",
	}

	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFramePanelFallsBackWhenTooSmall(t *testing.T) {
	got := framePanel("Repo", false, []string{"x"}, 1, 2)
	want := []string{"x", ""}

	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestRenderLeftPanelsSeparated(t *testing.T) {
	g := newTestGuiWithClient(&testmock.GHClient{})
	g.state.Repo = "owner/repo"
	g.state.PRs = []core.Item{{Number: 1, Title: "Fix bug"}}

	lines := g.renderLeftPanels(20, 10)
	if len(lines) != 10 {
		t.Fatalf("got %d lines, want 10", len(lines))
	}
	if lines[0] != "┌ Repository ──────┐" {
		t.Fatalf("unexpected repository frame top: %q", lines[0])
	}
	if lines[3] != "└──────────────────┘" {
		t.Fatalf("unexpected repository frame bottom: %q", lines[3])
	}
	if !strings.HasPrefix(lines[4], "┌> PRs") || !strings.HasSuffix(lines[4], "┐") {
		t.Fatalf("unexpected pr frame top: %q", lines[4])
	}
}

func TestRenderPRPanel(t *testing.T) {
	type fixture struct {
		prsLoading bool
		prs        []core.Item
		selected   int
	}

	type want struct {
		line1 string
	}

	tests := []struct {
		name    string
		fixture fixture
		want    want
	}{
		{
			name:    "empty placeholder",
			fixture: fixture{},
			want: want{
				line1: "No pull requests",
			},
		},
		{
			name: "loading",
			fixture: fixture{
				prsLoading: true,
			},
			want: want{
				line1: "",
			},
		},
		{
			name: "with prs",
			fixture: fixture{
				prs:      []core.Item{{Number: 1, Title: "Fix bug"}},
				selected: 0,
			},
			want: want{
				line1: "> PR #1 Fix bug",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithClient(&testmock.GHClient{})
			g.state.PRsLoading = tt.fixture.prsLoading
			g.state.PRs = tt.fixture.prs
			g.state.PRsSelected = tt.fixture.selected
			lines := g.renderPRPanel(3)

			if len(lines) != 3 {
				t.Fatalf("got %d lines, want 3", len(lines))
			}
			if lines[0] != tt.want.line1 {
				t.Fatalf("got %q, want %q", lines[0], tt.want.line1)
			}
		})
	}
}

func TestRenderRepoPanel(t *testing.T) {
	type want struct {
		line1 string
	}

	tests := []struct {
		name string
		repo string
		want want
	}{
		{
			name: "show repo",
			repo: "owner/repo",
			want: want{
				line1: "owner/repo",
			},
		},
		{
			name: "empty repo",
			repo: "",
			want: want{
				line1: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithClient(&testmock.GHClient{})
			g.state.Repo = tt.repo
			lines := g.renderRepoPanel(2)

			if len(lines) != 2 {
				t.Fatalf("got %d lines, want 2", len(lines))
			}
			if lines[0] != tt.want.line1 {
				t.Fatalf("got %q, want %q", lines[0], tt.want.line1)
			}
		})
	}
}
