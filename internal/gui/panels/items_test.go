package panels

import "testing"

func TestFormatters(t *testing.T) {
	tests := []struct {
		name string
		fn   ItemFormatter
		item Item
		want string
	}{
		{name: "Repo", fn: FormatRepoItem, item: Item{Title: "owner/repo"}, want: "owner/repo"},
		{name: "PR", fn: FormatPRItem, item: Item{Number: 42, Title: "fix bug"}, want: "PR #42 fix bug"},
		{name: "Issue", fn: FormatIssueItem, item: Item{Number: 7, Title: "add feature"}, want: "Issue #7 add feature"},
		{name: "EmptyTitle", fn: FormatPRItem, item: Item{Number: 1, Title: ""}, want: "PR #1 "},
	}
	for _, tt := range tests {
		got := tt.fn(tt.item)
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestNewItemsPanel(t *testing.T) {
	p := NewItemsPanel(FormatIssueItem, false)
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
	if p.Loading {
		t.Error("Loading should be false")
	}
	if p.KeepSelectionOnBlur {
		t.Error("KeepSelectionOnBlur should be false")
	}
}

func TestItemsPanelFormatRows(t *testing.T) {
	p := NewItemsPanel(FormatIssueItem, false)
	p.Items = []Item{
		{Number: 7, Title: "issue"},
		{Number: 42, Title: "pr"},
	}

	tests := []struct {
		index int
		want  string
	}{
		{index: 0, want: "Issue #7 issue"},
		{index: 1, want: "Issue #42 pr"},
	}

	for _, tt := range tests {
		got := p.Format(p.Items[tt.index])
		if got != tt.want {
			t.Errorf("Format(%d)=%q, want %q", tt.index, got, tt.want)
		}
	}
}

func TestItemsPanelFormat_DefaultFormatter(t *testing.T) {
	p := NewItemsPanel(nil, false)
	got := p.Format(Item{Title: "owner/repo"})
	if got != "owner/repo" {
		t.Errorf("got %q, want %q", got, "owner/repo")
	}
}

func TestItemsPanelKeepSelectionFlag(t *testing.T) {
	p := NewItemsPanel(FormatRepoItem, true)
	if !p.KeepSelectionOnBlur {
		t.Fatal("KeepSelectionOnBlur should be true")
	}
}
