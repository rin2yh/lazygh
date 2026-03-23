// Package textarea wraps charmbracelet/bubbles textarea as a reusable TUI component.
package textarea

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// State is a stateful textarea component.
type State struct {
	m textarea.Model
}

// New returns a State with sensible defaults and the given placeholder.
func New(placeholder string) State {
	m := textarea.New()
	m.Placeholder = placeholder
	m.ShowLineNumbers = false
	m.SetHeight(4)
	m.Prompt = ""
	m.CharLimit = 0
	return State{m: m}
}

func (s *State) Text() string        { return s.m.Value() }
func (s *State) View() string        { return s.m.View() }
func (s *State) Lines() []string     { return strings.Split(s.m.View(), "\n") }
func (s *State) Focus()              { s.m.Focus() }
func (s *State) Blur()               { s.m.Blur() }
func (s *State) Clear()              { s.m.SetValue("") }
func (s *State) Load(content string) { s.m.SetValue(content) }
func (s *State) Update(msg tea.KeyMsg) tea.Cmd {
	updated, cmd := s.m.Update(msg)
	s.m = updated
	return cmd
}
