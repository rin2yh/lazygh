package gh

import "strings"

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
