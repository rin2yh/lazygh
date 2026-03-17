package diff

import (
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/gh"
)

func TestSelectionSelectNextFileResetsLineAndKeepsValidSelection(t *testing.T) {
	selection := Selection{}
	selection.SetFiles([]gh.DiffFile{
		{
			Path: "a.txt",
			Lines: []gh.DiffLine{
				{Text: "a1"},
			},
		},
		{
			Path: "b.txt",
			Lines: []gh.DiffLine{
				{Text: "b1"},
				{Text: "b2"},
			},
		},
	})
	selection.SetLineSelected(10)

	if !selection.SelectNextFile() {
		t.Fatal("expected file selection to move")
	}
	if selection.FileSelected() != 1 {
		t.Fatalf("got %d, want %d", selection.FileSelected(), 1)
	}
	if selection.LineSelected() != 0 {
		t.Fatalf("got %d, want %d", selection.LineSelected(), 0)
	}
}

func TestSelectionLineNavigation(t *testing.T) {
	selection := Selection{}
	selection.SetFiles(ParseFilesMust(strings.Join([]string{
		"diff --git a/a.txt b/a.txt",
		"--- a/a.txt",
		"+++ b/a.txt",
		"@@ -1,6 +1,6 @@",
		" line1",
		" line2",
		"-line3",
		"+line3x",
		" line4",
		" line5",
		" line6",
	}, "\n")))

	if !selection.SelectNextLine(4) {
		t.Fatal("expected line selection to move")
	}
	if selection.LineSelected() == 0 {
		t.Fatal("expected selected line to advance")
	}
	if !selection.GotoLastLine() {
		t.Fatal("expected goto last line to move")
	}

	file, ok := selection.CurrentFile()
	if !ok {
		t.Fatal("expected current file")
	}
	if selection.LineSelected() != len(file.Lines)-1 {
		t.Fatalf("got %d, want %d", selection.LineSelected(), len(file.Lines)-1)
	}
	if !selection.GotoFirstLine() {
		t.Fatal("expected goto first line to move")
	}
	if selection.LineSelected() != 0 {
		t.Fatalf("got %d, want %d", selection.LineSelected(), 0)
	}
}

func TestSelectionCurrentDiffAccessors(t *testing.T) {
	selection := Selection{}
	selection.SetFiles([]gh.DiffFile{
		{
			Path: "a.txt",
			Lines: []gh.DiffLine{
				{Text: "x"},
				{Text: "y"},
			},
		},
	})
	selection.SetLineSelected(1)

	file, ok := selection.CurrentFile()
	if !ok {
		t.Fatal("expected current file")
	}
	if file.Path != "a.txt" {
		t.Fatalf("got %q, want %q", file.Path, "a.txt")
	}

	line, ok := selection.CurrentLine()
	if !ok {
		t.Fatal("expected current line")
	}
	if line.Text != "y" {
		t.Fatalf("got %q, want %q", line.Text, "y")
	}
	if selection.LineSelected() != 1 {
		t.Fatalf("got %d, want %d", selection.LineSelected(), 1)
	}
}

func ParseFilesMust(content string) []gh.DiffFile {
	files, _, _ := ParseFiles(nil, 0, content)
	if len(files) == 0 {
		panic("expected parsed files")
	}
	return files
}
