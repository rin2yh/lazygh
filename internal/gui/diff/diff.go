package diff

import (
	"strconv"
	"strings"

	"github.com/rin2yh/lazygh/internal/gh"
)

const (
	ansiReset  = "\x1b[0m"
	ansiGreen  = "\x1b[32m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiBlue   = "\x1b[34m"
	ansiCyan   = "\x1b[36m"
	ansiPurple = "\x1b[35m"
	ansiGray   = "\x1b[90m"
)

func RenderFileListLine(file gh.DiffFile) string {
	label := string(file.Status)
	if label == "" {
		label = string(gh.DiffFileStatusModified)
	}
	status := colorizeDiffStatus(label, file.Status)
	additions := colorize(ansiGreen, "+"+strconv.Itoa(file.Additions))
	deletions := colorize(ansiRed, "-"+strconv.Itoa(file.Deletions))
	return status + " " + file.Path + " " + additions + " " + deletions
}

func ColorizeContent(content string) string {
	if content == "" {
		return ""
	}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = ColorizeLine(line)
	}
	return strings.Join(lines, "\n")
}

func ColorizeLine(line string) string {
	switch {
	case strings.HasPrefix(line, "diff --git "):
		return colorize(ansiBlue, line)
	case strings.HasPrefix(line, "@@"):
		return colorize(ansiCyan, line)
	case strings.HasPrefix(line, "+++ "):
		return colorize(ansiGreen, line)
	case strings.HasPrefix(line, "--- "):
		return colorize(ansiRed, line)
	case strings.HasPrefix(line, "+"):
		return colorize(ansiGreen, line)
	case strings.HasPrefix(line, "-"):
		return colorize(ansiRed, line)
	case strings.HasPrefix(line, "new file mode "), strings.HasPrefix(line, "deleted file mode "):
		return colorize(ansiYellow, line)
	case strings.HasPrefix(line, "rename from "), strings.HasPrefix(line, "rename to "):
		return colorize(ansiPurple, line)
	case strings.HasPrefix(line, "index "), strings.HasPrefix(line, "similarity index "):
		return colorize(ansiGray, line)
	default:
		return line
	}
}

func colorizeDiffStatus(label string, status gh.DiffFileStatus) string {
	switch status {
	case gh.DiffFileStatusAdded:
		return colorize(ansiGreen, label)
	case gh.DiffFileStatusDeleted:
		return colorize(ansiRed, label)
	case gh.DiffFileStatusRenamed:
		return colorize(ansiCyan, label)
	case gh.DiffFileStatusCopied:
		return colorize(ansiBlue, label)
	case gh.DiffFileStatusType:
		return colorize(ansiPurple, label)
	default:
		return colorize(ansiYellow, label)
	}
}

func colorize(color string, text string) string {
	if text == "" {
		return ""
	}
	return color + text + ansiReset
}
