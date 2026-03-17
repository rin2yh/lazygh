package gui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type viewportState struct {
	vp     viewport.Model
	width  int
	height int
	body   string
}

func newViewportState() viewportState {
	return viewportState{
		vp:     viewport.New(1, 1),
		width:  1,
		height: 1,
	}
}

func (s *viewportState) Sync(width int, height int, body string) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if s.width != width || s.height != height {
		s.vp.Width = width
		s.vp.Height = height
		s.width = width
		s.height = height
	}
	if s.body != body {
		s.vp.SetContent(body)
		s.body = body
		s.vp.GotoTop()
	}
}

func (s *viewportState) Height() int {
	return s.height
}

func (s *viewportState) Update(msg tea.KeyMsg) (bool, tea.Cmd) {
	updated, cmd := s.vp.Update(msg)
	s.vp = updated
	return true, cmd
}

func (s *viewportState) ScrollDown(lines int) {
	s.vp.ScrollDown(lines)
}

func (s *viewportState) ScrollUp(lines int) {
	s.vp.ScrollUp(lines)
}

func (s *viewportState) View() string {
	return s.vp.View()
}
