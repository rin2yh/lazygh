package gui

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/core"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
)

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

func TestRenderRightPanels_OverviewShowsPanelTitle(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
	g.state.DetailContent = "detail"

	lines := g.renderRightPanels(40, 6)
	if len(lines) != 6 {
		t.Fatalf("got %d, want %d", len(lines), 6)
	}
	if !strings.Contains(lines[0], "Overview") {
		t.Fatalf("line does not contain %q: %q", "Overview", lines[0])
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
	if xansi.Strip(lines[0]) != "┌ Repository ──────┐" {
		t.Fatalf("got %q, want %q", xansi.Strip(lines[0]), "┌ Repository ──────┐")
	}
	if xansi.Strip(lines[3]) != "└──────────────────┘" {
		t.Fatalf("got %q, want %q", xansi.Strip(lines[3]), "└──────────────────┘")
	}
	if !strings.Contains(xansi.Strip(lines[4]), "PRs (Open/Draft)") {
		t.Fatalf("line does not contain expected title: %q", xansi.Strip(lines[4]))
	}
	if !strings.HasSuffix(xansi.Strip(lines[4]), "┐") {
		t.Fatalf("line does not have expected suffix: %q", xansi.Strip(lines[4]))
	}
}

func TestRenderPRPanel(t *testing.T) {
	type fixture struct {
		prsLoading bool
		prs        []core.Item
		selected   int
	}

	type want struct {
		line1          string
		line1Highlight bool
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
				line1:          "No pull requests",
				line1Highlight: false,
			},
		},
		{
			name: "loading",
			fixture: fixture{
				prsLoading: true,
			},
			want: want{
				line1:          "",
				line1Highlight: false,
			},
		},
		{
			name: "with prs",
			fixture: fixture{
				prs:      []core.Item{{Number: 1, Title: "Fix bug"}},
				selected: 0,
			},
			want: want{
				line1:          "> PR #1 Fix bug",
				line1Highlight: true,
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
			if xansi.Strip(lines[0]) != tt.want.line1 {
				t.Fatalf("got %q, want %q", xansi.Strip(lines[0]), tt.want.line1)
			}
			if tt.want.line1Highlight && !strings.Contains(lines[0], ansiReverse) {
				t.Fatalf("selected line should be highlighted: %q", lines[0])
			}
			if !tt.want.line1Highlight && strings.Contains(lines[0], ansiReverse) {
				t.Fatalf("line should not be highlighted: %q", lines[0])
			}
		})
	}
}

func TestRenderDiffFilesPanel_HighlightsSelectedLine(t *testing.T) {
	g := newTestGuiWithPRs(&testmock.GHClient{}, core.Item{Number: 1, Title: "x"})
	g.switchToDiff()
	g.updateDiffFiles(strings.Join([]string{
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
		"-one",
		"+two",
	}, "\n"))

	lines := g.renderDiffFilesPanel(40, 6)
	if len(lines) != 6 {
		t.Fatalf("got %d, want %d", len(lines), 6)
	}
	if !strings.Contains(lines[1], ansiReverse) {
		t.Fatalf("selected file line should be highlighted: %q", lines[1])
	}
	if !strings.Contains(xansi.Strip(lines[1]), "> M a.txt +1 -1") {
		t.Fatalf("unexpected selected file row: %q", xansi.Strip(lines[1]))
	}
	if strings.Contains(lines[2], ansiReverse) {
		t.Fatalf("non-selected line should not be highlighted: %q", lines[2])
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
