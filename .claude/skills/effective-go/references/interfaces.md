# Interfaces

## Small interfaces

Prefer interfaces with few methods.

Example

type Reader interface {
    Read(p []byte) (int, error)
}

Avoid large interfaces.

## Consumer side interfaces

Define interfaces where they are used, not where they are implemented.

## Naming

Interfaces often end with -er.

Reader
Writer
Formatter

Avoid "Interface" suffix.

Bad

UserRepositoryInterface
