package panels

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

type ItemKind int

const (
	ItemKindPR ItemKind = iota
	ItemKindIssue
)

type Item struct {
	Kind   ItemKind
	Number int
	Title  string
}

func (it Item) String() string {
	kind := "PR"
	if it.Kind == ItemKindIssue {
		kind = "Issue"
	}
	return fmt.Sprintf("%s #%d %s", kind, it.Number, it.Title)
}

type ItemsPanel struct {
	Items    []Item
	Selected int
}

func NewItemsPanel() *ItemsPanel {
	return &ItemsPanel{
		Items:    []Item{},
		Selected: 0,
	}
}

func (p *ItemsPanel) Render(v *gocui.View) {
	v.Clear()
	for i, item := range p.Items {
		prefix := "  "
		if i == p.Selected {
			prefix = "> "
		}
		_, _ = v.Write([]byte(prefix + item.String() + "\n"))
	}
}
