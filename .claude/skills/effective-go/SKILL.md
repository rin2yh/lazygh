---
name: effective-go
description: Apply idiomatic Go conventions from Effective Go when writing or reviewing Go code.
---

# Role

Write Go code following Effective Go conventions.

# When to use

Use this skill when:

- writing Go code
- reviewing Go code
- refactoring Go code
- designing Go packages

# Core Rules

## Formatting

- Always produce gofmt-compatible code
- Avoid unnecessary parentheses
- Prefer short variable declarations when appropriate

See: references/formatting.md

## Naming

- Package names are short and lowercase
- Avoid redundant names (httpserver.Server)
- Interfaces often end with -er

See: references/naming.md

## Interfaces

- Prefer small interfaces
- Define interfaces on the consumer side
- Avoid large interface types

See: references/interfaces.md

## Errors

- Return errors explicitly
- Do not use panic for normal errors
- Add context when wrapping errors

See: references/errors.md

## Packages

- Keep packages cohesive
- Avoid circular dependencies
- Package names should describe their responsibility

See: references/packages.md

## Concurrency

- Use goroutines for concurrency
- Prefer channels for communication
- Avoid shared mutable state when possible

See: references/concurrency.md
