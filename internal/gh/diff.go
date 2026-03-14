package gh

import (
	"strings"
)

type DiffFile struct {
	Path      string
	Content   string
	Status    DiffFileStatus
	Additions int
	Deletions int
}

type DiffFileStatus string

const (
	DiffFileStatusModified DiffFileStatus = "M"
	DiffFileStatusAdded    DiffFileStatus = "A"
	DiffFileStatusDeleted  DiffFileStatus = "D"
	DiffFileStatusRenamed  DiffFileStatus = "R"
	DiffFileStatusCopied   DiffFileStatus = "C"
	DiffFileStatusType     DiffFileStatus = "T"
)

// UnifiedDiffParser parses unified diff text into per-file metadata and content.
type UnifiedDiffParser struct {
	content string
}

// NewUnifiedDiffParser creates a parser for unified diff content.
func NewUnifiedDiffParser(content string) *UnifiedDiffParser {
	return &UnifiedDiffParser{content: content}
}

// Parse parses unified diff content and returns file entries.
func (p *UnifiedDiffParser) Parse() []DiffFile {
	return newDiffFileParser(p.content).parse()
}

// ParseUnifiedDiff is a backward compatible helper around UnifiedDiffParser.
func ParseUnifiedDiff(content string) []DiffFile {
	return NewUnifiedDiffParser(content).Parse()
}

type diffFileParser struct {
	lines       []string
	files       []DiffFile
	start       int
	currentPath string
	pathParser  *diffPathParser
}

func newDiffFileParser(content string) *diffFileParser {
	return &diffFileParser{
		lines:       strings.Split(content, "\n"),
		files:       make([]DiffFile, 0),
		start:       -1,
		currentPath: "",
		pathParser:  newDiffPathParser(),
	}
}

func (p *diffFileParser) parse() []DiffFile {
	if len(p.lines) == 0 || (len(p.lines) == 1 && p.lines[0] == "") {
		return nil
	}

	for i, line := range p.lines {
		if !strings.HasPrefix(line, "diff --git ") {
			continue
		}
		p.appendFile(i)
		p.start = i
		p.currentPath = p.pathParser.parseFromDiffGitLine(line)
	}
	p.appendFile(len(p.lines))
	return p.files
}

func (p *diffFileParser) appendFile(end int) {
	if p.start < 0 || p.start >= end {
		return
	}

	segment := strings.Join(p.lines[p.start:end], "\n")
	if strings.TrimSpace(segment) == "" {
		return
	}

	path := p.currentPath
	if path == "" {
		path = "(unknown)"
	}

	status, additions, deletions, refinedPath := p.parseFileMetadata(p.lines[p.start:end])
	if refinedPath != "" {
		path = refinedPath
	}
	if status == "" {
		status = DiffFileStatusModified
	}

	p.files = append(p.files, DiffFile{
		Path:      path,
		Content:   segment,
		Status:    status,
		Additions: additions,
		Deletions: deletions,
	})
}

func (p *diffFileParser) parseFileMetadata(lines []string) (DiffFileStatus, int, int, string) {
	status := DiffFileStatusModified
	additions := 0
	deletions := 0
	path := ""
	pathPriority := 0
	inHunk := false

	setPath := func(raw string, priority int, stripGitPrefix bool) {
		if priority < pathPriority {
			return
		}
		clean := p.pathParser.normalizer.normalize(raw, stripGitPrefix)
		if clean == "" {
			return
		}
		path = clean
		pathPriority = priority
	}

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			inHunk = false
		case strings.HasPrefix(line, "@@"):
			inHunk = true
		case strings.HasPrefix(line, "new file mode "):
			status = DiffFileStatusAdded
			inHunk = false
		case strings.HasPrefix(line, "deleted file mode "):
			status = DiffFileStatusDeleted
			inHunk = false
		case strings.HasPrefix(line, "rename from "), strings.HasPrefix(line, "rename to "):
			status = DiffFileStatusRenamed
			inHunk = false
			if strings.HasPrefix(line, "rename to ") {
				setPath(strings.TrimSpace(strings.TrimPrefix(line, "rename to ")), 2, false)
			}
		case strings.HasPrefix(line, "copy from "), strings.HasPrefix(line, "copy to "):
			status = DiffFileStatusCopied
			inHunk = false
			if strings.HasPrefix(line, "copy to ") {
				setPath(strings.TrimSpace(strings.TrimPrefix(line, "copy to ")), 2, false)
			}
		case strings.HasPrefix(line, "old mode "), strings.HasPrefix(line, "new mode "):
			if status == DiffFileStatusModified {
				status = DiffFileStatusType
			}
			inHunk = false
		case !inHunk && strings.HasPrefix(line, "+++ "):
			setPath(strings.TrimSpace(strings.TrimPrefix(line, "+++ ")), 1, true)
			inHunk = false
		case !inHunk && strings.HasPrefix(line, "--- "):
			setPath(strings.TrimSpace(strings.TrimPrefix(line, "--- ")), 1, true)
			inHunk = false
		case inHunk && strings.HasPrefix(line, "+"):
			additions++
		case inHunk && strings.HasPrefix(line, "-"):
			deletions++
		}
	}

	return status, additions, deletions, path
}

type diffPathParser struct {
	normalizer *diffPathNormalizer
}

func newDiffPathParser() *diffPathParser {
	return &diffPathParser{
		normalizer: &diffPathNormalizer{},
	}
}

func (p *diffPathParser) parseFromDiffGitLine(line string) string {
	payload := strings.TrimSpace(strings.TrimPrefix(line, "diff --git "))
	if payload == "" {
		return ""
	}

	oldPath, newPath, ok := p.splitGitPaths(payload)
	if !ok {
		return ""
	}

	path := p.normalizer.normalize(newPath, true)
	if path != "" {
		return path
	}
	return p.normalizer.normalize(oldPath, true)
}

func (p *diffPathParser) splitGitPaths(line string) (string, string, bool) {
	line = strings.TrimLeft(line, " \t")
	if line == "" {
		return "", "", false
	}

	if line[0] != '"' {
		return p.splitUnquotedGitPaths(line)
	}

	first, rest, ok := p.consumePathToken(line)
	if !ok {
		return "", "", false
	}
	rest = strings.TrimLeft(rest, " \t")
	if rest == "" {
		return "", "", false
	}
	second, _, ok := p.consumePathToken(rest)
	if !ok {
		return "", "", false
	}
	return first, second, true
}

func (p *diffPathParser) splitUnquotedGitPaths(line string) (string, string, bool) {
	type candidate struct {
		old string
		new string
	}
	candidates := make([]candidate, 0)

	for i := 0; i < len(line)-2; i++ {
		if line[i] != ' ' && line[i] != '\t' {
			continue
		}
		if !strings.HasPrefix(line[i+1:], "b/") {
			continue
		}

		oldPath := strings.TrimSpace(line[:i])
		newPath := strings.TrimSpace(line[i+1:])
		if oldPath == "" || newPath == "" {
			continue
		}
		if !strings.HasPrefix(oldPath, "a/") || !strings.HasPrefix(newPath, "b/") {
			continue
		}
		candidates = append(candidates, candidate{old: oldPath, new: newPath})
	}

	if len(candidates) == 0 {
		return "", "", false
	}

	for _, c := range candidates {
		if p.normalizer.normalize(c.old, true) == p.normalizer.normalize(c.new, true) {
			return c.old, c.new, true
		}
	}

	return candidates[0].old, candidates[0].new, true
}

func (p *diffPathParser) consumePathToken(s string) (string, string, bool) {
	if s == "" {
		return "", "", false
	}

	if s[0] == '"' {
		escaped := false
		for i := 1; i < len(s); i++ {
			switch {
			case escaped:
				escaped = false
			case s[i] == '\\':
				escaped = true
			case s[i] == '"':
				return s[:i+1], s[i+1:], true
			}
		}
		return "", "", false
	}

	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' {
			return s[:i], s[i+1:], true
		}
	}
	return s, "", true
}

type diffPathNormalizer struct{}

func (n *diffPathNormalizer) normalize(path string, stripGitPrefix bool) string {
	path = strings.TrimSpace(path)
	if path == "" || path == "/dev/null" {
		return ""
	}

	path = n.trimQuoted(path)
	path = strings.TrimSpace(path)
	if path == "" || path == "/dev/null" {
		return ""
	}
	if stripGitPrefix && (strings.HasPrefix(path, "a/") || strings.HasPrefix(path, "b/")) {
		path = path[2:]
	}
	path = n.sanitize(path)
	path = strings.TrimSpace(path)
	if path == "" || path == "/dev/null" {
		return ""
	}
	return path
}

func (n *diffPathNormalizer) trimQuoted(path string) string {
	if len(path) < 2 || path[0] != '"' || path[len(path)-1] != '"' {
		return path
	}
	return path[1 : len(path)-1]
}

func (n *diffPathNormalizer) sanitize(path string) string {
	var b strings.Builder
	for _, r := range path {
		switch {
		case r == '\n' || r == '\r' || r == '\t':
			b.WriteByte(' ')
		case (r >= 0 && r < 0x20) || r == 0x7f:
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
