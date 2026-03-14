package gui

import (
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/gh"
)

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
