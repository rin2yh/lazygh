package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Viewport is the interface for scrollable content panels.
type Viewport interface {
	Sync(width, height int, body string)
	Height() int
	Update(msg tea.KeyMsg) (bool, tea.Cmd)
	ScrollDown(lines int)
	ScrollUp(lines int)
	GotoTop()
	GotoBottom()
	View() string
}

type State struct {
	vp     viewport.Model
	width  int
	height int
	body   string
}

func New() State {
	return State{
		vp:     viewport.New(1, 1),
		width:  1,
		height: 1,
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

func (s *State) Height() int {
	return s.height
}

func (s *State) Update(msg tea.KeyMsg) (bool, tea.Cmd) {
	updated, cmd := s.vp.Update(msg)
	s.vp = updated
	return true, cmd
}

func (s *State) ScrollDown(lines int) {
	s.vp.ScrollDown(lines)
}

func (s *State) ScrollUp(lines int) {
	s.vp.ScrollUp(lines)
}

func (s *State) GotoTop() {
	s.vp.GotoTop()
}

func (s *State) GotoBottom() {
	s.vp.GotoBottom()
}

func (s *State) View() string {
	return s.vp.View()
}
