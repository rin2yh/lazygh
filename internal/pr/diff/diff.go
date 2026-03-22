package diff

import (
	"strconv"
	"strings"

	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/pkg/gui/ansi"
)

func RenderFileListLine(file gh.DiffFile) string {
	label := string(file.Status)
	if label == "" {
		label = string(gh.DiffFileStatusModified)
	}
	status := colorizeDiffStatus(label, file.Status)
	additions := colorize(ansi.Green, "+"+strconv.Itoa(file.Additions))
	deletions := colorize(ansi.Red, "-"+strconv.Itoa(file.Deletions))
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
		return colorize(ansi.Blue, line)
	case strings.HasPrefix(line, "@@"):
		return colorize(ansi.Cyan, line)
	case strings.HasPrefix(line, "+++ "):
		return colorize(ansi.Green, line)
	case strings.HasPrefix(line, "--- "):
		return colorize(ansi.Red, line)
	case strings.HasPrefix(line, "+"):
		return colorize(ansi.Green, line)
	case strings.HasPrefix(line, "-"):
		return colorize(ansi.Red, line)
	case strings.HasPrefix(line, "new file mode "), strings.HasPrefix(line, "deleted file mode "):
		return colorize(ansi.Yellow, line)
	case strings.HasPrefix(line, "rename from "), strings.HasPrefix(line, "rename to "):
		return colorize(ansi.Purple, line)
	case strings.HasPrefix(line, "index "), strings.HasPrefix(line, "similarity index "):
		return colorize(ansi.Gray, line)
	default:
		return line
	}
}

func colorizeDiffStatus(label string, status gh.DiffFileStatus) string {
	switch status {
	case gh.DiffFileStatusAdded:
		return colorize(ansi.Green, label)
	case gh.DiffFileStatusDeleted:
		return colorize(ansi.Red, label)
	case gh.DiffFileStatusRenamed:
		return colorize(ansi.Cyan, label)
	case gh.DiffFileStatusCopied:
		return colorize(ansi.Blue, label)
	case gh.DiffFileStatusType:
		return colorize(ansi.Purple, label)
	default:
		return colorize(ansi.Yellow, label)
	}
}

func colorize(color string, text string) string {
	if text == "" {
		return ""
	}
	return color + text + ansi.Reset
}
