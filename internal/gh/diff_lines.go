package gh

import (
	"strconv"
	"strings"
)

type diffLineParser struct{}

func (p *diffLineParser) parse(path string, lines []string) []DiffLine {
	result := make([]DiffLine, 0, len(lines))
	inHunk := false
	oldLine := 0
	newLine := 0

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "@@"):
			result = append(result, DiffLine{
				Text: line,
				Kind: DiffLineKindHunk,
				Path: path,
			})
			parsedOld, parsedNew, ok := parseHunkHeader(line)
			if ok {
				oldLine = parsedOld
				newLine = parsedNew
				inHunk = true
			} else {
				inHunk = false
			}
		case !inHunk:
			result = append(result, DiffLine{
				Text: line,
				Kind: DiffLineKindMeta,
				Path: path,
			})
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++ "):
			result = append(result, DiffLine{
				Text:        line,
				Kind:        DiffLineKindAdd,
				Path:        path,
				NewLine:     newLine,
				Side:        DiffSideRight,
				Commentable: true,
			})
			newLine++
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "--- "):
			result = append(result, DiffLine{
				Text:        line,
				Kind:        DiffLineKindDelete,
				Path:        path,
				OldLine:     oldLine,
				Side:        DiffSideLeft,
				Commentable: true,
			})
			oldLine++
		case strings.HasPrefix(line, " "):
			result = append(result, DiffLine{
				Text:        line,
				Kind:        DiffLineKindContext,
				Path:        path,
				OldLine:     oldLine,
				NewLine:     newLine,
				Side:        DiffSideRight,
				Commentable: true,
			})
			oldLine++
			newLine++
		case strings.HasPrefix(line, "\\"):
			result = append(result, DiffLine{
				Text: line,
				Kind: DiffLineKindMeta,
				Path: path,
			})
		default:
			result = append(result, DiffLine{
				Text: line,
				Kind: DiffLineKindMeta,
				Path: path,
			})
		}
	}

	return result
}

func parseHunkHeader(line string) (int, int, bool) {
	if !strings.HasPrefix(line, "@@") {
		return 0, 0, false
	}
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return 0, 0, false
	}
	oldLine, ok := parseHunkRange(parts[1], '-')
	if !ok {
		return 0, 0, false
	}
	newLine, ok := parseHunkRange(parts[2], '+')
	if !ok {
		return 0, 0, false
	}
	return oldLine, newLine, true
}

func parseHunkRange(token string, prefix byte) (int, bool) {
	if len(token) < 2 || token[0] != prefix {
		return 0, false
	}
	token = token[1:]
	if idx := strings.IndexByte(token, ','); idx >= 0 {
		token = token[:idx]
	}
	value, err := strconv.Atoi(token)
	if err != nil {
		return 0, false
	}
	if value < 0 {
		return 0, false
	}
	return value, true
}
