package factory

import (
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/pr"
)

func NewItem(number int, title string) pr.Item {
	return pr.Item{
		Number: number,
		Title:  title,
	}
}

func NewGHPRItem(number int, title string) gh.PRItem {
	return gh.PRItem{
		Number: number,
		Title:  title,
		State:  "OPEN",
	}
}
