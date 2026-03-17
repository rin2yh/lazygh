// Package overview は PR 概要パネルのドメイン状態を所有するパッケージ。
package overview

import "github.com/rin2yh/lazygh/internal/model"

// State holds overview panel display and fetching state.
type State struct {
	Mode     model.DetailMode
	Content  string
	Fetching model.FetchKind
}
