package gui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

func newTestGuiWithClient(client gh.ClientInterface) *Gui {
	g, _ := NewGui(config.Default(), client)
	return g
}

func newTestGuiWithPRs(client gh.ClientInterface, prs ...core.Item) *Gui {
	g := newTestGuiWithClient(client)
	g.state.ApplyPRsResult("owner/repo", prs, nil)
	return g
}

func TestNavigatePRList(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, []core.Item{{Number: 1, Title: "a"}, {Number: 2, Title: "b"}}...)

	g.navigateDown()
	if g.state.PRsSelected != 1 {
		t.Fatalf("got %d, want %d", g.state.PRsSelected, 1)
	}

	g.navigateUp()
	if g.state.PRsSelected != 0 {
		t.Fatalf("got %d, want %d", g.state.PRsSelected, 0)
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
				prs:  []core.Item{{Number: 1, Title: "Fix bug", Status: "OPEN", Assignees: []string{"alice"}}},
			},
			want: want{
				repo:   "owner/repo",
				prs:    []core.Item{{Number: 1, Title: "Fix bug", Status: "OPEN", Assignees: []string{"alice"}}},
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
			g := newTestGuiWithClient(&testmock.GHClient{})
			g.state.BeginLoadPRs()

			g.applyPRsResult(tt.msg)

			if g.state.PRsLoading {
				t.Fatal("expected PRsLoading=false")
			}
			if g.state.Loading != core.LoadingNone {
				t.Fatalf("got %v, want %v", g.state.Loading, core.LoadingNone)
			}
			if g.state.Repo != tt.want.repo {
				t.Fatalf("got %q, want %q", g.state.Repo, tt.want.repo)
			}
			if diff := cmp.Diff(tt.want.prs, g.state.PRs, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("prs mismatch (-want +got)\n%s", diff)
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
				mode:    core.DetailModeOverview,
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
				mode:   core.DetailModeOverview,
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
			g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "Fix bug"})
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

func TestApplyDetailResult_DiffUsesSanitizedContent(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "Fix bug"})
	g.switchToDiff()
	g.state.Loading = core.LoadingDetail

	raw := strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+ok\x1b[31mred",
	}, "\n")

	g.applyDetailResult(detailLoadedMsg{
		mode:    core.DetailModeDiff,
		number:  1,
		content: raw,
	})

	if strings.Contains(g.state.DetailContent, "\x1b") {
		t.Fatalf("detail content should be sanitized: %q", g.state.DetailContent)
	}
	if len(g.diffFiles) != 1 {
		t.Fatalf("got %d, want %d", len(g.diffFiles), 1)
	}
	if strings.Contains(g.diffFiles[0].Content, "\x1b") {
		t.Fatalf("diff file content should be sanitized: %q", g.diffFiles[0].Content)
	}
	if !strings.Contains(g.diffFiles[0].Content, "+ok[31mred") {
		t.Fatalf("unexpected diff content: %q", g.diffFiles[0].Content)
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

func TestModelHandleDetailLoad(t *testing.T) {
	tests := []struct {
		name         string
		client       *testmock.GHClient
		pr           core.Item
		switchToDiff bool
		wantMode     core.DetailMode
		wantContent  string
		wantNumber   int
	}{
		{
			name:        "overview",
			client:      &testmock.GHClient{PRView: "detail"},
			pr:          core.Item{Number: 1, Title: "x"},
			wantMode:    core.DetailModeOverview,
			wantContent: "detail",
			wantNumber:  1,
		},
		{
			name:         "diff",
			client:       &testmock.GHClient{PRDiff: "diff"},
			pr:           core.Item{Number: 2, Title: "x"},
			switchToDiff: true,
			wantMode:     core.DetailModeDiff,
			wantContent:  "diff",
			wantNumber:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithPRs(tt.client, tt.pr)
			if tt.switchToDiff {
				g.switchToDiff()
			}
			m := &model{gui: g}

			cmd := m.handleDetailLoad()
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
	})

	if len(items) != 2 {
		t.Fatalf("got %d, want %d", len(items), 2)
	}
	if items[0].Status != "OPEN" {
		t.Fatalf("got %q, want %q", items[0].Status, "OPEN")
	}
	if strings.Join(items[0].Assignees, ",") != "alice,bob" {
		t.Fatalf("got %q, want %q", strings.Join(items[0].Assignees, ","), "alice,bob")
	}
	if items[1].Status != "DRAFT" {
		t.Fatalf("got %q, want %q", items[1].Status, "DRAFT")
	}
}

func TestScrollDetailByKey(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*Gui)
		key            tea.KeyMsg
		wantHandled    bool
		wantOffsetMove bool
	}{
		{
			name: "diff mode page down",
			setup: func(g *Gui) {
				g.switchToDiff()
				g.focus = panelDiffContent
			},
			key:            tea.KeyMsg{Type: tea.KeyPgDown},
			wantHandled:    true,
			wantOffsetMove: true,
		},
		{
			name:           "overview mode page down",
			setup:          func(_ *Gui) {},
			key:            tea.KeyMsg{Type: tea.KeyPgDown},
			wantHandled:    false,
			wantOffsetMove: false,
		},
		{
			name: "diff mode d key",
			setup: func(g *Gui) {
				g.switchToDiff()
				g.focus = panelDiffContent
			},
			key:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantHandled:    false,
			wantOffsetMove: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
			tt.setup(g)

			g.syncDetailViewport(20, 4, strings.Repeat("line\n", 30))
			before := g.detailViewport.YOffset

			handled := g.scrollDetailByKey(tt.key)
			if handled != tt.wantHandled {
				t.Fatalf("got %v, want %v", handled, tt.wantHandled)
			}
			if tt.wantOffsetMove {
				if g.detailViewport.YOffset <= before {
					t.Fatalf("expected offset to increase, before=%d after=%d", before, g.detailViewport.YOffset)
				}
				return
			}
			if g.detailViewport.YOffset != before {
				t.Fatalf("got %d, want %d", g.detailViewport.YOffset, before)
			}
		})
	}
}

func TestUpdateDiffFiles(t *testing.T) {
	g := newTestGuiWithClient(&testmock.GHClient{})
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
	if diff := cmp.Diff(want, g.diffFiles); diff != "" {
		t.Fatalf("diffFiles mismatch (-want +got)\n%s", diff)
	}

	g.diffFileSelected = 1
	g.updateDiffFiles(diff)
	if g.diffFileSelected != 1 {
		t.Fatalf("got %d, want %d", g.diffFileSelected, 1)
	}
}

func TestRenderDiffFileListLineShowsColoredStatusAndCounts(t *testing.T) {
	line := renderDiffFileListLine(gh.DiffFile{
		Path:      "a.txt",
		Status:    gh.DiffFileStatusAdded,
		Additions: 3,
		Deletions: 1,
	})

	if !strings.Contains(line, ansiGreen+"A"+ansiReset) {
		t.Fatalf("line does not contain %q", ansiGreen+"A"+ansiReset)
	}
	if !strings.Contains(line, ansiGreen+"+3"+ansiReset) {
		t.Fatalf("line does not contain %q", ansiGreen+"+3"+ansiReset)
	}
	if !strings.Contains(line, ansiRed+"-1"+ansiReset) {
		t.Fatalf("line does not contain %q", ansiRed+"-1"+ansiReset)
	}
	if !strings.Contains(line, "a.txt") {
		t.Fatalf("line does not contain %q", "a.txt")
	}
}

func TestColorizeDiffContent(t *testing.T) {
	diff := strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"index 1111111..2222222 100644",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+new",
		" context",
	}, "\n")

	got := colorizeDiffContent(diff)
	if !strings.Contains(got, ansiBlue+"diff --git a/a.txt b/a.txt"+ansiReset) {
		t.Fatalf("diff content does not contain expected header")
	}
	if !strings.Contains(got, ansiGray+"index 1111111..2222222 100644"+ansiReset) {
		t.Fatalf("diff content does not contain expected index line")
	}
	if !strings.Contains(got, ansiRed+"-old"+ansiReset) {
		t.Fatalf("diff content does not contain expected removal line")
	}
	if !strings.Contains(got, ansiGreen+"+new"+ansiReset) {
		t.Fatalf("diff content does not contain expected addition line")
	}
}

func TestCycleFocus_DiffMode(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
	g.switchToDiff()
	g.diffFiles = []gh.DiffFile{{Path: "a.txt", Content: "x"}}

	if g.focus != panelDiffFiles {
		t.Fatalf("got %v, want %v", g.focus, panelDiffFiles)
	}

	g.cycleFocus()
	if g.focus != panelDiffContent {
		t.Fatalf("got %v, want %v", g.focus, panelDiffContent)
	}
	g.cycleFocus()
	if g.focus != panelPRs {
		t.Fatalf("got %v, want %v", g.focus, panelPRs)
	}
	g.cycleFocus()
	if g.focus != panelDiffFiles {
		t.Fatalf("got %v, want %v", g.focus, panelDiffFiles)
	}
}

func TestModelHandleLKeyShowsOverviewFromPRsInDiffMode(t *testing.T) {
	mc := &testmock.GHClient{PRView: "overview"}
	g := newTestGuiWithPRs(mc, core.Item{Number: 1, Title: "x"})
	g.switchToDiff()
	g.focus = panelPRs

	m := &model{gui: g}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if g.state.IsDiffMode() {
		t.Fatal("expected overview mode")
	}
	if cmd == nil {
		t.Fatal("expected detail load command")
	}

	msg := cmd().(detailLoadedMsg)
	if msg.err != nil {
		t.Fatalf("unexpected error: %v", msg.err)
	}
	if msg.mode != core.DetailModeOverview {
		t.Fatalf("got %v, want %v", msg.mode, core.DetailModeOverview)
	}
	if msg.number != 1 {
		t.Fatalf("got %d, want %d", msg.number, 1)
	}
	if msg.content != "overview" {
		t.Fatalf("got %q, want %q", msg.content, "overview")
	}
}

func TestModelUpdateFocusKeysInDiffMode(t *testing.T) {
	tests := []struct {
		name      string
		key       tea.KeyMsg
		start     panelFocus
		wantFocus panelFocus
	}{
		{
			name:      "l moves files to diff",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			start:     panelDiffFiles,
			wantFocus: panelDiffContent,
		},
		{
			name:      "h moves diff to files",
			key:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			start:     panelDiffContent,
			wantFocus: panelDiffFiles,
		},
		{
			name:      "esc moves to prs",
			key:       tea.KeyMsg{Type: tea.KeyEsc},
			start:     panelDiffContent,
			wantFocus: panelPRs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
			g.switchToDiff()
			g.diffFiles = []gh.DiffFile{{Path: "a.txt", Content: "x"}}
			g.focus = tt.start
			m := &model{gui: g}

			_, cmd := m.Update(tt.key)
			if cmd != nil {
				t.Fatal("did not expect command")
			}
			if g.focus != tt.wantFocus {
				t.Fatalf("got %v, want %v", g.focus, tt.wantFocus)
			}
		})
	}
}

func TestRenderRightPanels_DiffModeHasFilesPanel(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
	g.switchToDiff()
	g.updateDiffFiles(strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1 +1 @@",
		"-old",
		"+new",
	}, "\n"))

	lines := g.renderRightPanels(60, 10)
	if len(lines) != 10 {
		t.Fatalf("got %d, want %d", len(lines), 10)
	}
	if !strings.Contains(lines[0], "Files") {
		t.Fatalf("line does not contain %q: %q", "Files", lines[0])
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

func TestWrapTextWithANSI(t *testing.T) {
	got := wrapText(ansiGreen+"abcdef"+ansiReset, 3)
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d, want %d", len(lines), 2)
	}
	if xansi.Strip(lines[0]) != "abc" {
		t.Fatalf("got %q, want %q", xansi.Strip(lines[0]), "abc")
	}
	if xansi.StringWidth(lines[0]) != 3 {
		t.Fatalf("got %d, want %d", xansi.StringWidth(lines[0]), 3)
	}
	if xansi.Strip(lines[1]) != "def" {
		t.Fatalf("got %q, want %q", xansi.Strip(lines[1]), "def")
	}
	if xansi.StringWidth(lines[1]) != 3 {
		t.Fatalf("got %d, want %d", xansi.StringWidth(lines[1]), 3)
	}
}

func TestFramePanel(t *testing.T) {
	got := framePanel("Repo", false, []string{"body"}, 10, 3)
	want := []string{
		"┌ Repo ──┐",
		"│body    │",
		"└────────┘",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("frame mismatch (-want +got)\n%s", diff)
	}
}

func TestFramePanelFallsBackWhenTooSmall(t *testing.T) {
	got := framePanel("Repo", false, []string{"x"}, 1, 2)
	want := []string{"x", ""}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("frame mismatch (-want +got)\n%s", diff)
	}
}

func TestPadOrTrimHandlesANSI(t *testing.T) {
	colored := ansiGreen + "+10" + ansiReset
	got := padOrTrim(colored, 4)
	if !strings.Contains(got, colored) {
		t.Fatalf("result does not contain colored text: %q", got)
	}
}

func TestRenderLeftPanelsSeparated(t *testing.T) {
	g := newTestGuiWithClient(&testmock.GHClient{})
	g.state.Repo = "owner/repo"
	g.state.PRs = []core.Item{{Number: 1, Title: "Fix bug"}}

	lines := g.renderLeftPanels(20, 10)
	if len(lines) != 10 {
		t.Fatalf("got %d, want %d", len(lines), 10)
	}
	if lines[0] != "┌ Repository ──────┐" {
		t.Fatalf("got %q, want %q", lines[0], "┌ Repository ──────┐")
	}
	if lines[3] != "└──────────────────┘" {
		t.Fatalf("got %q, want %q", lines[3], "└──────────────────┘")
	}
	if !strings.HasPrefix(lines[4], "┌> PRs") {
		t.Fatalf("line does not have expected prefix: %q", lines[4])
	}
	if !strings.HasSuffix(lines[4], "┐") {
		t.Fatalf("line does not have expected suffix: %q", lines[4])
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
				t.Fatalf("got %d, want %d", len(lines), 3)
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
				t.Fatalf("got %d, want %d", len(lines), 2)
			}
			if lines[0] != tt.want.line1 {
				t.Fatalf("got %q, want %q", lines[0], tt.want.line1)
			}
		})
	}
}
