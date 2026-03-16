package factory

import (
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

func CoreItem(number int, title string) core.Item {
	return core.Item{
		Number: number,
		Title:  title,
	}
}

func GHPRItem(number int, title string) gh.PRItem {
	return gh.PRItem{
		Number: number,
		Title:  title,
		State:  "OPEN",
	}
}
