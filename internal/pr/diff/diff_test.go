package diff

import (
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/gh"
)

func TestRenderDiffFileListLineShowsColoredStatusAndCounts(t *testing.T) {
	line := RenderFileListLine(gh.DiffFile{
		Path:      "a.txt",
		Status:    gh.DiffFileStatusAdded,
		Additions: 3,
		Deletions: 1,
	})

	tests := []struct {
		substr string
	}{
		{ansiGreen + "A" + ansiReset},
		{ansiGreen + "+3" + ansiReset},
		{ansiRed + "-1" + ansiReset},
		{"a.txt"},
	}
	for _, tt := range tests {
		if !strings.Contains(line, tt.substr) {
			t.Errorf("line does not contain %q", tt.substr)
		}
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

	got := ColorizeContent(diff)

	tests := []struct {
		substr string
	}{
		{ansiBlue + "diff --git a/a.txt b/a.txt" + ansiReset},
		{ansiGray + "index 1111111..2222222 100644" + ansiReset},
		{ansiRed + "-old" + ansiReset},
		{ansiGreen + "+new" + ansiReset},
	}
	for _, tt := range tests {
		if !strings.Contains(got, tt.substr) {
			t.Errorf("diff content does not contain %q", tt.substr)
		}
	}
}
