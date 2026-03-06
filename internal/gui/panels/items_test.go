package panels

import "testing"

func TestItemString(t *testing.T) {
	tests := []struct {
		name string
		item Item
		want string
	}{
		{"PR", Item{Kind: ItemKindPR, Number: 42, Title: "fix bug"}, "PR #42 fix bug"},
		{"Issue", Item{Kind: ItemKindIssue, Number: 7, Title: "add feature"}, "Issue #7 add feature"},
		{"EmptyTitle", Item{Kind: ItemKindPR, Number: 1, Title: ""}, "PR #1 "},
	}
	for _, tt := range tests {
		got := tt.item.String()
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestNewItemsPanel(t *testing.T) {
	p := NewItemsPanel()
	if p == nil {
		t.Fatal("NewItemsPanel returned nil")
	}
	if p.Items == nil {
		t.Error("Items should not be nil")
	}
	if len(p.Items) != 0 {
		t.Errorf("Items length: got %d, want 0", len(p.Items))
	}
	if p.Selected != 0 {
		t.Errorf("Selected: got %d, want 0", p.Selected)
	}
}
