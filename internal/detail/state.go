// Package detail は PR 詳細パネルのドメイン状態を所有するパッケージ。
package detail

import "github.com/rin2yh/lazygh/internal/model"

// State holds detail panel display and loading state.
type State struct {
	Mode    model.DetailMode
	Content string
	Loading model.LoadingKind
}
