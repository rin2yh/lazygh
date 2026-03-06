package panels

import (
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

const wideRunePad = '\x00'

func normalizeDisplayText(text string) string {
	if text == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(text))

	for _, r := range text {
		b.WriteRune(r)
		if unicode.IsControl(r) {
			continue
		}
		for i := 1; i < runeDisplayWidth(r); i++ {
			b.WriteRune(wideRunePad)
		}
	}

	return b.String()
}

func runeDisplayWidth(r rune) int {
	w := runewidth.RuneWidth(r)
	if w <= 0 {
		return 1
	}
	if w == 2 && runewidth.IsAmbiguousWidth(r) {
		return 1
	}
	return w
}
