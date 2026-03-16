package review

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	appstate "github.com/rin2yh/lazygh/internal/state"
)

type summary struct {
	state    *appstate.State
	setFocus func(FocusTarget)
	editor   textarea.Model
}

func newSummary(state *appstate.State, setFocus func(FocusTarget)) *summary {
	return &summary{
		state:    state,
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
	beginInput(f.state, f.setFocus, &f.editor, f.state.BeginReviewSummaryInput, f.state.Review.Summary)
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
	f.state.SetReviewSummary(f.editor.Value())
}

func (f *summary) Clear() {
	f.editor.SetValue("")
}
