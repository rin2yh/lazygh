package app

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/pr"
	"github.com/rin2yh/lazygh/internal/pr/overview"
	testfactory "github.com/rin2yh/lazygh/pkg/test/factory"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func TestModelInitLoadsPRs(t *testing.T) {
	mc := &testmock.GHClient{Repo: "owner/repo", PRs: []gh.PRItem{testfactory.NewGHPRItem(2, "p")}}
	g, err := NewGui(config.Default(), NewCoordinator(), mc, mc)
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
		pr           pr.Item
		switchToDiff bool
		wantMode     overview.DetailMode
		wantContent  string
		wantNumber   int
	}{
		{
			name:        "overview",
			client:      &testmock.GHClient{PRView: "detail"},
			pr:          testfactory.NewItem(1, "x"),
			wantMode:    overview.DetailModeOverview,
			wantContent: "detail",
			wantNumber:  1,
		},
		{
			name:         "diff",
			client:       &testmock.GHClient{PRDiff: "diff"},
			pr:           testfactory.NewItem(2, "x"),
			switchToDiff: true,
			wantMode:     overview.DetailModeDiff,
			wantContent:  "diff",
			wantNumber:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGui(config.Default(), NewCoordinator(), tt.client, tt.client)
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.coord.ApplyPRsResult("owner/repo", []pr.Item{tt.pr}, nil)
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

func TestGuiApplyPRsResult(t *testing.T) {
	type want struct {
		repo   string
		prs    []pr.Item
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
				prs:  []pr.Item{{Number: 1, Title: "Fix bug", Status: "OPEN", Assignees: []string{"alice"}}},
			},
			want: want{
				repo:   "owner/repo",
				prs:    []pr.Item{{Number: 1, Title: "Fix bug", Status: "OPEN", Assignees: []string{"alice"}}},
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
				detail: "Error fetching PRs: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.coord.BeginFetchPRs()

			g.applyPRsResult(tt.msg)

			if g.coord.IsFetching() {
				t.Fatal("expected PRsLoading=false")
			}
			if k := g.coord.Overview.FetchKind(); k != overview.FetchNone {
				t.Fatalf("got %v, want %v", k, overview.FetchNone)
			}
			if g.coord.Repo() != tt.want.repo {
				t.Fatalf("got %q, want %q", g.coord.Repo(), tt.want.repo)
			}
			if diff := cmp.Diff(tt.want.prs, g.coord.Items(), cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("prs mismatch (-want +got)\n%s", diff)
			}
			if g.coord.Overview.Content() != tt.want.detail {
				t.Fatalf("got %q, want %q", g.coord.Overview.Content(), tt.want.detail)
			}
		})
	}
}

func TestGuiApplyDetailResult(t *testing.T) {
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
				mode:    overview.DetailModeOverview,
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
				mode:   overview.DetailModeOverview,
				number: 1,
				err:    errors.New("boom"),
			},
			want: want{
				detail: "Error fetching detail: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
			if err != nil {
				t.Fatalf("NewGui failed: %v", err)
			}
			g.coord.ApplyPRsResult("owner/repo", []pr.Item{{Number: 1, Title: "Fix bug"}}, nil)
			g.coord.Overview.StartFetching(overview.FetchingDetail)

			g.applyDetailResult(tt.msg)

			if k := g.coord.Overview.FetchKind(); k != overview.FetchNone {
				t.Fatalf("got %v, want %v", k, overview.FetchNone)
			}
			if g.coord.Overview.Content() != tt.want.detail {
				t.Fatalf("got %q, want %q", g.coord.Overview.Content(), tt.want.detail)
			}
		})
	}
}

func TestApplyDetailResult_DiffUsesSanitizedContent(t *testing.T) {
	g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	g.coord.ApplyPRsResult("owner/repo", []pr.Item{{Number: 1, Title: "Fix bug"}}, nil)
	g.switchToDiff()
	g.coord.Overview.StartFetching(overview.FetchingDetail)

	raw := strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+ok\x1b[31mred",
	}, "\n")

	g.applyDetailResult(detailLoadedMsg{
		mode:    overview.DetailModeDiff,
		number:  1,
		content: raw,
	})

	if strings.Contains(g.coord.Overview.Content(), "\x1b") {
		t.Fatalf("detail content should be sanitized: %q", g.coord.Overview.Content())
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
	g, err := NewGui(config.Default(), NewCoordinator(), &testmock.GHClient{}, &testmock.GHClient{})
	if err != nil {
		t.Fatalf("NewGui failed: %v", err)
	}
	diffContent := strings.Join([]string{
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

	g.updateDiffFiles(diffContent)
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
	if d := cmp.Diff(want, g.diff.Files(), cmpopts.IgnoreFields(gh.DiffFile{}, "Lines")); d != "" {
		t.Fatalf("diffFiles mismatch (-want +got)\n%s", d)
	}

	g.diff.SetFileSelected(1)
	g.updateDiffFiles(diffContent)
	if g.diff.FileSelected() != 1 {
		t.Fatalf("got %d, want %d", g.diff.FileSelected(), 1)
	}
}
