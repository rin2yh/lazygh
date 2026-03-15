package review

import "github.com/charmbracelet/bubbles/textarea"

func newEditor(placeholder string) textarea.Model {
	editor := textarea.New()
	editor.Placeholder = placeholder
	editor.ShowLineNumbers = false
	editor.SetHeight(4)
	editor.Prompt = ""
	editor.CharLimit = 0
	return editor
}
