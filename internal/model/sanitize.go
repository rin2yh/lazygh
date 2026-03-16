package model

import "strings"

func SanitizeSingleLine(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == '\n' || r == '\r' || r == '\t':
			b.WriteByte(' ')
		case isDangerousControl(r):
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func SanitizeMultiline(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == '\r':
			continue
		case r == '\n':
			b.WriteByte('\n')
		case isDangerousControl(r):
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isDangerousControl(r rune) bool {
	return (r >= 0 && r < 0x20 && r != '\n' && r != '\t') || r == 0x7f
}
