package review

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/pkg/gui/textarea"
)

type summary struct {
	rs *ReviewState
	textarea.State
}

func newSummary(rs *ReviewState) *summary {
	return &summary{
		rs:    rs,
		State: textarea.New("Review summary"),
	}
}

func (f *summary) BeginInput() {
	f.rs.BeginSummaryInput()
	f.Load(f.rs.Summary)
	f.Focus()
}

func (f *summary) StopInput() {
	f.Blur()
}

func (f *summary) HandleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	return f.Update(msg), true
}

func (f *summary) Save() {
	f.rs.SetSummary(f.Text())
}

func (f *summary) Clear() {
	f.State.Clear()
}
