// Package app is the BubbleTea adapter layer for lazygh.
//
// # Responsibilities
//
// app owns:
//   - [Gui]: the top-level BubbleTea Model (Init / Update / View)
//   - [Coordinator]: application-wide state (selected PR, fetch progress, focus)
//   - Key-input routing and navigation (input.go, navigation.go, focus.go)
//   - Asynchronous data fetching (fetch.go)
//   - TUI rendering (render.go, screen.go)
//
// # Interface placement
//
// Interfaces in this package follow Go's "consumer-defined" convention:
// each interface is declared in the file where it is consumed, not in the
// package that implements it.  This keeps interfaces small and avoids
// import cycles.
//
// Interfaces declared here and their implementing types:
//
//   - [PRClient]          consumed by Gui / Coordinator; implemented by [github.com/rin2yh/lazygh/internal/gh.Client]
//   - [ReviewController]  consumed by Gui / input handlers; implemented by [github.com/rin2yh/lazygh/internal/pr/review.Controller]
//   - [ReviewReader]      read-only sub-set of ReviewController, used by render.go
//   - [ReviewHandler]     mutation sub-set of ReviewController, used by input.go
//   - [ReviewApplier]     async-result sub-set of ReviewController, used by coordinator.go
//   - [DetailViewport]    consumed by Gui; implemented by [github.com/rin2yh/lazygh/pkg/gui/viewport.State]
//
// Interfaces consumed by [github.com/rin2yh/lazygh/internal/pr/review] are defined in that
// package (review.Selection, review.AppState, review.PendingReviewClient) to
// avoid an import from review → app.
package app
