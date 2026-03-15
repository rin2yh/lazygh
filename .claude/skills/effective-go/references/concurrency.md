# Concurrency

## Goroutines

Use goroutines to perform concurrent work.

go process()

## Channels

Prefer channels for communication between goroutines.

Avoid sharing mutable state.

## Channel ownership

The sender should usually close the channel.

Receivers should not close channels they did not create.
