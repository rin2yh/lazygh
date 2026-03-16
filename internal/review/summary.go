package review

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type summary struct {
	rs       *ReviewState
	setFocus func(FocusTarget)
	editor   textarea.Model
}

func newSummary(rs *ReviewState, setFocus func(FocusTarget)) *summary {
	return &summary{
		rs:       rs,
		setFocus: setFocus,
		editor:   newEditor("Review summary"),
	}
}

func (f *summary) CurrentValue() string {
	return f.editor.Value()
}

func (f *summary) InputLines() []string {
	return editorLines(f.editor)
}

func (f *summary) BeginInput() {
	beginInput(f.rs, f.setFocus, &f.editor, f.rs.BeginSummaryInput, f.rs.Summary)
}

func (f *summary) StopInput() {
	f.editor.Blur()
}

func (f *summary) HandleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	updated, cmd := f.editor.Update(msg)
	f.editor = updated
	return cmd, true
}

func (f *summary) Save() {
	f.rs.SetSummary(f.editor.Value())
}

func (f *summary) Clear() {
	f.editor.SetValue("")
}
