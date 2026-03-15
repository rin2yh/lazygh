package gui

func (gui *Gui) navigateDown() bool {
	return gui.state.NavigateDown()
}

func (gui *Gui) navigateUp() bool {
	return gui.state.NavigateUp()
}
