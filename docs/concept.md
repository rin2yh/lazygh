# lazygh Concept

## Why lazygh

A TUI for performing PR reviews on the current repository via the shortest possible workflow, using lazygit-style keybindings.

### Differentiation from Existing Tools

| Aspect | gh CLI | lazygit | lazygh |
|--------|--------|---------|--------|
| Usability | Requires memorizing commands | Single-key focused | Single-key focused (for GitHub PRs) |
| Scope | GitHub in general | Local Git operations | PR review for the current repo |
| Authentication | `gh auth` | Not required | Delegated to `gh auth` |
| Diff viewing | Command output | Fast | PR unified diff displayed in center panel |

## Architecture Overview

```
lazygh (Go TUI)
    ↕ exec
gh CLI
    ↕
GitHub API
```

### Layer Responsibilities

- **lazygh (Go TUI)**: Screen rendering, key input, state transitions
- **gh CLI**: API access and authentication
- **GitHub API**: PR / comment / diff data

## Design Principles

### Fixed to Current Repository
Resolves the current repository at startup and handles only PRs for that repo.
Multi-repo switching is planned for v2 and beyond.

### Focused on PR Review
v1 is scoped to the minimum features required for PR review.
Issue workflows are not provided in v1.

### lazygit-style Keybindings
Prioritizes `h/j/k/l`-centered navigation for fast selection and view switching.

- Left: PR list
- Center: PR details (Overview / Diff)
- Bottom: Status / key guide
