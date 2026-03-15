package gh

import "strings"

type diffFileSegment struct {
	lines       []string
	content     string
	initialPath string
}

type diffFileParser struct {
	content         string
	segmentSplitter *diffSegmentSplitter
	metadataParser  *diffMetadataParser
	lineParser      *diffLineParser
}

func newDiffFileParser(content string) *diffFileParser {
	pathParser := newDiffPathParser()
	return &diffFileParser{
		content:         content,
		segmentSplitter: &diffSegmentSplitter{pathParser: pathParser},
		metadataParser:  &diffMetadataParser{pathParser: pathParser},
		lineParser:      &diffLineParser{},
	}
}

func (p *diffFileParser) parse() []DiffFile {
	segments := p.segmentSplitter.split(p.content)
	if len(segments) == 0 {
		return nil
	}

	files := make([]DiffFile, 0, len(segments))
	for _, segment := range segments {
		metadata := p.metadataParser.parse(segment.lines)
		path := segment.initialPath
		if metadata.path != "" {
			path = metadata.path
		}
		if path == "" {
			path = "(unknown)"
		}

		files = append(files, DiffFile{
			Path:      path,
			Content:   segment.content,
			Status:    metadata.status,
			Additions: metadata.additions,
			Deletions: metadata.deletions,
			Lines:     p.lineParser.parse(path, segment.lines),
		})
	}
	return files
}

type diffSegmentSplitter struct {
	pathParser *diffPathParser
}

func (s *diffSegmentSplitter) split(content string) []diffFileSegment {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return nil
	}

	segments := make([]diffFileSegment, 0)
	start := -1
	currentPath := ""

	appendSegment := func(end int) {
		if start < 0 || start >= end {
			return
		}
		part := lines[start:end]
		content := strings.Join(part, "\n")
		if strings.TrimSpace(content) == "" {
			return
		}
		segments = append(segments, diffFileSegment{
			lines:       part,
			content:     content,
			initialPath: currentPath,
		})
	}

	for i, line := range lines {
		if !strings.HasPrefix(line, "diff --git ") {
			continue
		}
		appendSegment(i)
		start = i
		currentPath = s.pathParser.parseFromDiffGitLine(line)
	}
	appendSegment(len(lines))
	return segments
}
