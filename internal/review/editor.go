package review

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
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

// beginInput performs the shared steps to start a review input form:
// transition state, open drawer, populate editor, and focus it.
func beginInput(rs *ReviewState, editor *textarea.Model, transitionState func(), initialValue string) {
	transitionState()
	rs.OpenDrawer()
	editor.SetValue(initialValue)
	editor.Focus()
}
