// Package sanitize provides string sanitization utilities for safe display.
package sanitize

import "strings"

// SingleLine replaces newlines and tabs with spaces and removes dangerous control characters.
func SingleLine(s string) string {
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

// Multiline strips carriage returns and removes dangerous control characters while preserving newlines.
func Multiline(s string) string {
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
	return (r < 0x20 && r != '\n' && r != '\t') || r == 0x7f
}
