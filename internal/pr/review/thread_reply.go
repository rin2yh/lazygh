package review

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/pkg/gui/textarea"
)

type threadReply struct {
	keys config.KeyBindings
	rs   *ReviewState
	textarea.State
}

func newThreadReply(cfg *config.Config, rs *ReviewState) *threadReply {
	return &threadReply{
		keys:  cfg.KeyBindings,
		rs:    rs,
		State: textarea.New("Reply to thread"),
	}
}

func (r *threadReply) BeginInput() {
	r.rs.BeginThreadReplyInput()
	r.Clear()
	r.Focus()
}

func (r *threadReply) CurrentValue() string { return r.Text() }

func (r *threadReply) InputLines() []string { return r.Lines() }

func (r *threadReply) StopInput() {
	r.Blur()
	r.State.Clear()
}

func (r *threadReply) HandleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if r.keys.Matches(msg, config.ActionReviewSave) {
		return nil, true
	}
	return r.Update(msg), true
}
