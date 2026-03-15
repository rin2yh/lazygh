package review

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/rin2yh/lazygh/internal/core"
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
// transition state, open drawer, set focus, populate editor, and focus it.
func beginInput(state *core.State, setFocus func(FocusTarget), editor *textarea.Model, transitionState func(), initialValue string) {
	transitionState()
	state.OpenReviewDrawer()
	setFocus(FocusReviewDrawer)
	editor.SetValue(initialValue)
	editor.Focus()
}
