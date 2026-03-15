package gh

import "strings"

type diffFileMetadata struct {
	status    DiffFileStatus
	additions int
	deletions int
	path      string
}

type diffMetadataParser struct {
	pathParser *diffPathParser
}

func (p *diffMetadataParser) parse(lines []string) diffFileMetadata {
	metadata := diffFileMetadata{
		status: DiffFileStatusModified,
	}
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
		metadata.path = clean
		pathPriority = priority
	}

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			setPath(p.pathParser.parseFromDiffGitLine(line), 0, false)
			inHunk = false
		case strings.HasPrefix(line, "@@"):
			inHunk = true
		case strings.HasPrefix(line, "new file mode "):
			metadata.status = DiffFileStatusAdded
			inHunk = false
		case strings.HasPrefix(line, "deleted file mode "):
			metadata.status = DiffFileStatusDeleted
			inHunk = false
		case strings.HasPrefix(line, "rename from "), strings.HasPrefix(line, "rename to "):
			metadata.status = DiffFileStatusRenamed
			inHunk = false
			if strings.HasPrefix(line, "rename to ") {
				setPath(strings.TrimSpace(strings.TrimPrefix(line, "rename to ")), 2, false)
			}
		case strings.HasPrefix(line, "copy from "), strings.HasPrefix(line, "copy to "):
			metadata.status = DiffFileStatusCopied
			inHunk = false
			if strings.HasPrefix(line, "copy to ") {
				setPath(strings.TrimSpace(strings.TrimPrefix(line, "copy to ")), 2, false)
			}
		case strings.HasPrefix(line, "old mode "), strings.HasPrefix(line, "new mode "):
			if metadata.status == DiffFileStatusModified {
				metadata.status = DiffFileStatusType
			}
			inHunk = false
		case !inHunk && strings.HasPrefix(line, "+++ "):
			setPath(strings.TrimSpace(strings.TrimPrefix(line, "+++ ")), 1, true)
			inHunk = false
		case !inHunk && strings.HasPrefix(line, "--- "):
			setPath(strings.TrimSpace(strings.TrimPrefix(line, "--- ")), 1, true)
			inHunk = false
		case inHunk && strings.HasPrefix(line, "+"):
			metadata.additions++
		case inHunk && strings.HasPrefix(line, "-"):
			metadata.deletions++
		}
	}

	return metadata
}
