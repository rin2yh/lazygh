package gui

import (
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
)

func wrapText(content string, width int) string {
	if width <= 0 || content == "" {
		return content
	}

	srcLines := strings.Split(content, "\n")
	dstLines := make([]string, 0, len(srcLines))
	for _, line := range srcLines {
		lineWidth := xansi.StringWidth(line)
		if lineWidth <= width {
			dstLines = append(dstLines, line)
			continue
		}
		for left := 0; left < lineWidth; left += width {
			right := left + width
			dstLines = append(dstLines, xansi.Cut(line, left, right))
		}
	}
	return strings.Join(dstLines, "\n")
}
