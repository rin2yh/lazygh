# lazygh

A TUI for fast PR review on the current GitHub repository, with lazygit-style keybindings.
Accesses the GitHub API through the `gh` CLI.

## Prerequisites

- Go 1.25+
- [gh CLI](https://cli.github.com/) installed
- Authenticated via `gh auth login`
- Run from within a Git-managed repository

## Installation

```sh
go install github.com/rin2yh/lazygh@latest
```

Build from source:

```sh
git clone https://github.com/rin2yh/lazygh
cd lazygh
go build -o lazygh .
```

## Usage

```sh
lazygh
```

## CI

GitHub Actions automatically runs `go fmt`, `go vet`, and `go test` on every PR and push, serving as a quality gate.

## Task Management

Managed via GitHub Milestones.
Details are described in the Issues within each milestone.
