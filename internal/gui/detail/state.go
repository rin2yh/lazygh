package detail

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type State struct {
	viewport viewport.Model
	width    int
	height   int
	body     string
}

func NewState() State {
	return State{
		viewport: viewport.New(1, 1),
		width:    1,
		height:   1,
	}
}

func (s *State) Sync(width int, height int, body string) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if s.width != width || s.height != height {
		s.viewport.Width = width
		s.viewport.Height = height
		s.width = width
		s.height = height
	}
	if s.body != body {
		s.viewport.SetContent(body)
		s.body = body
		s.viewport.GotoTop()
	}
}

func (s *State) Height() int {
	return s.height
}

func (s *State) Update(msg tea.KeyMsg) bool {
	updated, _ := s.viewport.Update(msg)
	s.viewport = updated
	return true
}

func (s *State) ScrollDown(lines int) {
	s.viewport.ScrollDown(lines)
}

func (s *State) ScrollUp(lines int) {
	s.viewport.ScrollUp(lines)
}

func (s *State) View() string {
	return s.viewport.View()
}

func (s *State) YOffset() int {
	return s.viewport.YOffset
}
