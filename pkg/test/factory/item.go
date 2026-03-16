package factory

import (
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/model"
)

func NewItem(number int, title string) model.Item {
	return model.Item{
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
