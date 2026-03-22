// Package ansi provides ANSI escape sequence constants for terminal coloring.
package ansi

const (
	Reset   = "\x1b[0m"
	Black   = "\x1b[30m"
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Purple  = Magenta // alias for Magenta
	Cyan    = "\x1b[36m"
	White   = "\x1b[37m"
	Gray    = "\x1b[90m"
	Reverse = "\x1b[7m"
)
