package review

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/pkg/gui/textarea"
)

type comment struct {
	keys      config.KeyBindings
	rs        *ReviewState
	selection Selection
	textarea.State
}

func newComment(cfg *config.Config, rs *ReviewState) *comment {
	return &comment{
		keys:  cfg.KeyBindings,
		rs:    rs,
		State: textarea.New("Add review comment"),
	}
}

func (f *comment) bindSelection(selection Selection) {
	f.selection = selection
}

func (f *comment) CurrentValue() string {
	return f.Text()
}

func (f *comment) SetValue(v string) {
	f.Load(v)
}

func (f *comment) InputLines() []string {
	return f.Lines()
}

func (f *comment) BeginInput() {
	f.rs.BeginCommentInput()
	f.Clear()
	f.Focus()
}

func (f *comment) Clear() {
	f.State.Clear()
	f.rs.Notify("Comment input cleared.")
}

func (f *comment) StartEdit(body string) {
	f.Load(body)
	f.Focus()
}

func (f *comment) StopInput() {
	f.Blur()
	f.State.Clear()
}

func (f *comment) HandleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if f.keys.Matches(msg, config.ActionReviewSave) {
		return nil, true
	}
	return f.Update(msg), true
}

func (f *comment) BuildDraft(body string, start *Range) (gh.ReviewComment, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return gh.ReviewComment{}, fmt.Errorf("comment body is empty")
	}
	line, ok := f.selection.CurrentLine()
	if !ok || !line.Commentable {
		return gh.ReviewComment{}, fmt.Errorf("current line is not commentable")
	}
	c := gh.ReviewComment{
		Path: line.Path,
		Body: body,
		Side: line.Side,
	}
	if line.NewLine > 0 && line.Side != gh.DiffSideLeft {
		c.Line = line.NewLine
	} else {
		c.Line = line.OldLine
	}
	if c.Line <= 0 {
		return gh.ReviewComment{}, fmt.Errorf("comment line is invalid")
	}
	if start == nil {
		return c, nil
	}
	if start.Path != c.Path {
		return gh.ReviewComment{}, fmt.Errorf("range must stay within one file")
	}
	if start.Index != f.selection.LineSelected() {
		c.StartLine = start.Line
		c.StartSide = gh.DiffSide(start.Side)
		if start.Index > f.selection.LineSelected() {
			c.StartLine, c.Line = c.Line, c.StartLine
			c.StartSide, c.Side = c.Side, c.StartSide
		}
	}
	return c, nil
}

func (f *comment) ApplySaved() {
	f.State.Clear()
	f.Blur()
}
