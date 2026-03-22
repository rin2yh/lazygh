package gh

import "fmt"

type DiffFile struct {
	Path      string
	Content   string
	Status    DiffFileStatus
	Additions int
	Deletions int
	Lines     []DiffLine
}

type DiffFileStatus string

type DiffSide string

const (
	DiffSideLeft  DiffSide = "LEFT"
	DiffSideRight DiffSide = "RIGHT"
)

type DiffLineKind string

const (
	DiffLineKindMeta    DiffLineKind = "meta"
	DiffLineKindHunk    DiffLineKind = "hunk"
	DiffLineKindContext DiffLineKind = "context"
	DiffLineKindAdd     DiffLineKind = "add"
	DiffLineKindDelete  DiffLineKind = "delete"
)

type DiffLine struct {
	Text        string
	Kind        DiffLineKind
	Path        string
	OldLine     int
	NewLine     int
	Side        DiffSide
	Commentable bool
}

const (
	DiffFileStatusModified DiffFileStatus = "M"
	DiffFileStatusAdded    DiffFileStatus = "A"
	DiffFileStatusDeleted  DiffFileStatus = "D"
	DiffFileStatusRenamed  DiffFileStatus = "R"
	DiffFileStatusCopied   DiffFileStatus = "C"
	DiffFileStatusType     DiffFileStatus = "T"
)

// ParseUnifiedDiff parses unified diff text into per-file metadata and content.
func ParseUnifiedDiff(content string) []DiffFile {
	return newDiffFileParser(content).parse()
}

func FormatDiffLineLocation(line DiffLine) string {
	switch {
	case line.OldLine > 0 && line.NewLine > 0:
		return fmt.Sprintf("%d/%d", line.OldLine, line.NewLine)
	case line.NewLine > 0:
		return fmt.Sprintf("+%d", line.NewLine)
	case line.OldLine > 0:
		return fmt.Sprintf("-%d", line.OldLine)
	default:
		return ""
	}
}
