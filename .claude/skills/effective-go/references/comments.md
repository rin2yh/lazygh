# Comments

## Package comments

Every package should have a package comment.

Example

// Package http provides HTTP client and server implementations.

## Exported identifiers

All exported types, functions, and variables must have comments.

Example

// Serve starts the HTTP server.
func Serve() {}

## Avoid obvious comments

Bad

// increment i
i++

Prefer comments that explain intent rather than code.
