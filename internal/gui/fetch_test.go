package gui

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/model"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestModelInitLoadsPRs(t *testing.T) {
	mc := &testmock.GHClient{Repo: "owner/repo", PRs: []gh.PRItem{testfactory.NewGHPRItem(2, "p")}}
	g, err := NewGui(config.Default(), mc, mc)
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	m := &screen{gui: g}

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("expected init command")
	}
	msg := cmd().(prsLoadedMsg)
	if msg.err != nil {
		t.Fatalf("unexpected error: %v", msg.err)
	}
	if msg.repo != "owner/repo" {
		t.Fatalf("got %q, want %q", msg.repo, "owner/repo")
	}
	if len(msg.prs) != 1 {
		t.Fatalf("got %d, want %d", len(msg.prs), 1)
	}
}

func TestScreenOpenSelectedPR(t *testing.T) {
	tests := []struct {
		name         string
		client       *testmock.GHClient
		pr           model.Item
		switchToDiff bool
		wantMode     model.DetailMode
		wantContent  string
		wantNumber   int
	}{
		{
			name:        "overview",
			client:      &testmock.GHClient{PRView: "detail"},
			pr:          testfactory.NewItem(1, "x"),
			wantMode:    model.DetailModeOverview,
			wantContent: "detail",
			wantNumber:  1,
		},
		{
			name:         "diff",
			client:       &testmock.GHClient{PRDiff: "diff"},
			pr:           testfactory.NewItem(2, "x"),
			switchToDiff: true,
			wantMode:     model.DetailModeDiff,
			wantContent:  "diff",
			wantNumber:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGui(config.Default(), tt.client, tt.client)
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.state.ApplyPRsResult("owner/repo", []model.Item{tt.pr}, nil)
			if tt.switchToDiff {
				g.switchToDiff()
			}
			m := &screen{gui: g}

			cmd := m.openSelectedPR()
			if cmd == nil {
				t.Fatal("expected detail load command")
			}
			msg := cmd().(detailLoadedMsg)
			if msg.err != nil {
				t.Fatalf("unexpected error: %v", msg.err)
			}
			if msg.content != tt.wantContent {
				t.Fatalf("got %q, want %q", msg.content, tt.wantContent)
			}
			if msg.mode != tt.wantMode {
				t.Fatalf("got %v, want %v", msg.mode, tt.wantMode)
			}
			if msg.number != tt.wantNumber {
				t.Fatalf("got %d, want %d", msg.number, tt.wantNumber)
			}
		})
	}
}

func TestToCorePRsMapsStatusAndAssignees(t *testing.T) {
	items := toCorePRs([]gh.PRItem{
		{
			Number:  1,
			Title:   "open",
			State:   "OPEN",
			IsDraft: false,
			Assignees: []gh.GHUser{
				{Login: "alice"},
				{Login: "bob"},
			},
		},
		{
			Number:  2,
			Title:   "draft",
			State:   "OPEN",
			IsDraft: true,
		},
	}, model.PRFilterOpen)

	if len(items) != 2 {
		t.Fatalf("got %d, want %d", len(items), 2)
	}
	if items[0].Status != model.PRStatusOpen {
		t.Fatalf("got %q, want %q", items[0].Status, model.PRStatusOpen)
	}
	if strings.Join(items[0].Assignees, ",") != "alice,bob" {
		t.Fatalf("got %q, want %q", strings.Join(items[0].Assignees, ","), "alice,bob")
	}
	if items[1].Status != model.PRStatusDraft {
		t.Fatalf("got %q, want %q", items[1].Status, model.PRStatusDraft)
	}
}

func TestApplyPRsResult(t *testing.T) {
	type want struct {
		repo   string
		prs    []model.Item
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
				prs:  []model.Item{{Number: 1, Title: "Fix bug", Status: "OPEN", Assignees: []string{"alice"}}},
			},
			want: want{
				repo:   "owner/repo",
				prs:    []model.Item{{Number: 1, Title: "Fix bug", Status: "OPEN", Assignees: []string{"alice"}}},
				detail: "PR #1 Fix bug\nStatus: OPEN\nAssignee: alice",
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
			g, err := NewGui(config.Default(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.state.BeginLoadPRs()

			g.applyPRsResult(tt.msg)

			if g.state.Fetching {
				t.Fatal("expected PRsLoading=false")
			}
			if g.state.Overview.Loading != model.LoadingNone {
				t.Fatalf("got %v, want %v", g.state.Overview.Loading, model.LoadingNone)
			}
			if g.state.Repo != tt.want.repo {
				t.Fatalf("got %q, want %q", g.state.Repo, tt.want.repo)
			}
			if diff := cmp.Diff(tt.want.prs, g.state.Items, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("prs mismatch (-want +got)\n%s", diff)
			}
			if g.state.Overview.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", g.state.Overview.Content, tt.want.detail)
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
				mode:    model.DetailModeOverview,
				number:  1,
				content: "hello",
			},
			want: want{
				detail: "hello",
			},
		},
		{
			name: "error",
			msg: detailLoadedMsg{
				mode:   model.DetailModeOverview,
				number: 1,
				err:    errors.New("boom"),
			},
			want: want{
				detail: "Error loading detail: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGui(config.Default(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.state.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "Fix bug"}}, nil)
			g.state.Overview.Loading = model.LoadingDetail

			g.applyDetailResult(tt.msg)

			if g.state.Overview.Loading != model.LoadingNone {
				t.Fatalf("got %v, want %v", g.state.Overview.Loading, model.LoadingNone)
			}
			if g.state.Overview.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", g.state.Overview.Content, tt.want.detail)
			}
		})
	}
}

func TestApplyDetailResult_DiffUsesSanitizedContent(t *testing.T) {
	g, err := NewGui(config.Default(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.state.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "Fix bug"}}, nil)
	g.switchToDiff()
	g.state.Overview.Loading = model.LoadingDetail

	raw := strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+ok\x1b[31mred",
	}, "\n")

	g.applyDetailResult(detailLoadedMsg{
		mode:    model.DetailModeDiff,
		number:  1,
		content: raw,
	})

	if strings.Contains(g.state.Overview.Content, "\x1b") {
		t.Fatalf("detail content should be sanitized: %q", g.state.Overview.Content)
	}
	if len(g.diff.Files()) != 1 {
		t.Fatalf("got %d, want %d", len(g.diff.Files()), 1)
	}
	if strings.Contains(g.diff.Files()[0].Content, "\x1b") {
		t.Fatalf("diff file content should be sanitized: %q", g.diff.Files()[0].Content)
	}
	if !strings.Contains(g.diff.Files()[0].Content, "+ok[31mred") {
		t.Fatalf("unexpected diff content: %q", g.diff.Files()[0].Content)
	}
}

func TestUpdateDiffFiles(t *testing.T) {
	g, err := NewGui(config.Default(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	diff := strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+new",
		"diff --git a/b.txt b/b.txt",
		"--- a/b.txt",
		"+++ b/b.txt",
		"@@ -1 +1 @@",
		"-x",
		"+y",
	}, "\n")

	g.updateDiffFiles(diff)
	want := []gh.DiffFile{
		{Path: "a.txt", Content: strings.Join([]string{
			"diff --git a/a.txt b/a.txt",
			"--- a/a.txt",
			"+++ b/a.txt",
			"@@ -1 +1 @@",
			"-old",
			"+new",
		}, "\n"), Status: gh.DiffFileStatusModified, Additions: 1, Deletions: 1},
		{Path: "b.txt", Content: strings.Join([]string{
			"diff --git a/b.txt b/b.txt",
			"--- a/b.txt",
			"+++ b/b.txt",
			"@@ -1 +1 @@",
			"-x",
			"+y",
		}, "\n"), Status: gh.DiffFileStatusModified, Additions: 1, Deletions: 1},
	}
	if diff := cmp.Diff(want, g.diff.Files(), cmpopts.IgnoreFields(gh.DiffFile{}, "Lines")); diff != "" {
		t.Fatalf("diffFiles mismatch (-want +got)\n%s", diff)
	}

	g.diff.SetFileSelected(1)
	g.updateDiffFiles(diff)
	if g.diff.FileSelected() != 1 {
		t.Fatalf("got %d, want %d", g.diff.FileSelected(), 1)
	}
}
