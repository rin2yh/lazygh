package review

import (
	tea "github.com/charmbracelet/bubbletea"
)

type summary struct {
	rs *ReviewState
	editorInput
}

func newSummary(rs *ReviewState) *summary {
	return &summary{
		rs:          rs,
		editorInput: newEditorInput("Review summary"),
	}
}

func (f *summary) BeginInput() {
	f.rs.BeginSummaryInput()
	f.load(f.rs.Summary)
	f.focus()
}

func (f *summary) StopInput() {
	f.blur()
}

func (f *summary) HandleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	return f.update(msg), true
}

func (f *summary) Save() {
	f.rs.SetSummary(f.text())
}

func (f *summary) Clear() {
	f.clear()
}
