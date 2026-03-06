package panels

import "testing"

func TestNewReposPanel(t *testing.T) {
	p := NewReposPanel()
	if p == nil {
		t.Fatal("NewReposPanel returned nil")
	}
	if p.Repos == nil {
		t.Error("Repos should not be nil")
	}
	if len(p.Repos) != 0 {
		t.Errorf("Repos length: got %d, want 0", len(p.Repos))
	}
	if p.Selected != 0 {
		t.Errorf("Selected: got %d, want 0", p.Selected)
	}
}
