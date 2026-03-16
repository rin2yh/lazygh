package review

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	appstate "github.com/rin2yh/lazygh/internal/state"
)

type comment struct {
	keys      config.KeyBindings
	state     *appstate.State
	selection Selection
	setFocus  func(FocusTarget)
	editor    textarea.Model
}

func newComment(cfg *config.Config, state *appstate.State, setFocus func(FocusTarget)) *comment {
	return &comment{
		keys:     cfg.KeyBindings,
		state:    state,
		setFocus: setFocus,
		editor:   newEditor("Add review comment"),
	}
}

func (f *comment) bindSelection(selection Selection) {
	f.selection = selection
}

func (f *comment) CurrentValue() string {
	return f.editor.Value()
}

func (f *comment) SetValue(value string) {
	f.editor.SetValue(value)
}

func (f *comment) InputLines() []string {
	return editorLines(f.editor)
}

func (f *comment) BeginInput() {
	beginInput(f.state, f.setFocus, &f.editor, f.state.BeginReviewCommentInput, "")
}

func (f *comment) Clear() {
	f.editor.SetValue("")
	f.state.SetReviewNotice("Comment input cleared.")
}

func (f *comment) StartEdit(body string) {
	f.editor.SetValue(body)
	f.editor.Focus()
}

func (f *comment) StopInput() {
	f.editor.Blur()
	f.editor.SetValue("")
}

func (f *comment) HandleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch {
	case f.keys.Matches(msg, config.ActionReviewSave):
		return nil, true
	}

	updated, cmd := f.editor.Update(msg)
	f.editor = updated
	return cmd, true
}

func (f *comment) BuildDraft(body string, start *core.ReviewRange) (gh.ReviewComment, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return gh.ReviewComment{}, fmt.Errorf("comment body is empty")
	}
	line, ok := f.selection.CurrentDiffLine()
	if !ok || !line.Commentable {
		return gh.ReviewComment{}, fmt.Errorf("current line is not commentable")
	}
	comment := gh.ReviewComment{
		Path: line.Path,
		Body: body,
		Side: line.Side,
	}
	if line.NewLine > 0 && line.Side != gh.DiffSideLeft {
		comment.Line = line.NewLine
	} else {
		comment.Line = line.OldLine
	}
	if comment.Line <= 0 {
		return gh.ReviewComment{}, fmt.Errorf("comment line is invalid")
	}
	if start == nil {
		return comment, nil
	}
	if start.Path != comment.Path {
		return gh.ReviewComment{}, fmt.Errorf("range must stay within one file")
	}
	if start.Index != f.selection.CurrentLineIndex() {
		comment.StartLine = start.Line
		comment.StartSide = gh.DiffSide(start.Side)
		if start.Index > f.selection.CurrentLineIndex() {
			comment.StartLine, comment.Line = comment.Line, comment.StartLine
			comment.StartSide, comment.Side = comment.Side, comment.StartSide
		}
	}
	return comment, nil
}

func (f *comment) ApplySaved() {
	f.editor.SetValue("")
	f.editor.Blur()
	f.setFocus(FocusReviewDrawer)
}
