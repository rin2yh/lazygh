package panels

import "testing"

func TestNewDetailPanel(t *testing.T) {
	p := NewDetailPanel()
	if p == nil {
		t.Fatal("NewDetailPanel returned nil")
	}
	if p.Content != "" {
		t.Errorf("Content: got %q, want empty", p.Content)
	}
	if p.ScrollY != 0 {
		t.Errorf("ScrollY: got %d, want 0", p.ScrollY)
	}
}

func TestDetailPanel_SetContent(t *testing.T) {
	tests := []struct {
		name        string
		initialScroll int
		content     string
		wantContent string
	}{
		{"WithContent", 5, "hello", "hello"},
		{"Empty", 0, "", ""},
	}
	for _, tt := range tests {
		p := NewDetailPanel()
		p.ScrollY = tt.initialScroll
		p.SetContent(tt.content)
		if p.Content != tt.wantContent {
			t.Errorf("%s Content: got %q, want %q", tt.name, p.Content, tt.wantContent)
		}
		if p.ScrollY != 0 {
			t.Errorf("%s ScrollY should be reset to 0, got %d", tt.name, p.ScrollY)
		}
	}
}
