package review

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func newEditor(placeholder string) textarea.Model {
	editor := textarea.New()
	editor.Placeholder = placeholder
	editor.ShowLineNumbers = false
	editor.SetHeight(4)
	editor.Prompt = ""
	editor.CharLimit = 0
	return editor
}

func editorLines(e textarea.Model) []string {
	return strings.Split(e.View(), "\n")
}

// editorInput wraps textarea.Model with common editor operations shared by
// comment and summary input forms.
type editorInput struct {
	editor textarea.Model
}

func newEditorInput(placeholder string) editorInput {
	return editorInput{editor: newEditor(placeholder)}
}

func (e *editorInput) text() string    { return e.editor.Value() }
func (e *editorInput) lines() []string { return editorLines(e.editor) }
func (e *editorInput) blur()           { e.editor.Blur() }
func (e *editorInput) focus()          { e.editor.Focus() }

func (e *editorInput) clear()              { e.editor.SetValue("") }
func (e *editorInput) load(content string) { e.editor.SetValue(content) }

func (e *editorInput) update(msg tea.KeyMsg) tea.Cmd {
	updated, cmd := e.editor.Update(msg)
	e.editor = updated
	return cmd
}
