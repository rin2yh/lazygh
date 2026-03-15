package review

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
)

type summary struct {
	state    *core.State
	setFocus func(FocusTarget)
	editor   textarea.Model
}

func newSummary(state *core.State, setFocus func(FocusTarget)) *summary {
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
	return strings.Split(f.editor.View(), "\n")
}

func (f *summary) BeginInput() {
	beginInput(f.state, f.setFocus, &f.editor, f.state.BeginReviewSummaryInput, f.state.Review.Summary)
}

func (f *summary) StopInput() {
	f.editor.Blur()
}

func (f *summary) HandleKey(msg tea.KeyMsg) (bool, tea.Cmd) {
	updated, cmd := f.editor.Update(msg)
	f.editor = updated
	return true, cmd
}

func (f *summary) Save() {
	f.state.SetReviewSummary(f.editor.Value())
}

func (f *summary) Clear() {
	f.editor.SetValue("")
}
