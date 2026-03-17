package overview

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/model"
)

func TestState_ZeroValue(t *testing.T) {
	var s State
	if s.Mode != model.DetailModeOverview {
		t.Errorf("zero Mode = %v, want DetailModeOverview", s.Mode)
	}
	if s.Content != "" {
		t.Errorf("zero Content = %q, want empty string", s.Content)
	}
	if s.Fetching != model.FetchNone {
		t.Errorf("zero Fetching = %v, want FetchNone", s.Fetching)
	}
}

func TestState_Fields(t *testing.T) {
	tests := []struct {
		name     string
		state    State
		mode     model.DetailMode
		content  string
		fetching model.FetchKind
	}{
		{
			name:     "overview mode with content",
			state:    State{Mode: model.DetailModeOverview, Content: "PR body", Fetching: model.FetchNone},
			mode:     model.DetailModeOverview,
			content:  "PR body",
			fetching: model.FetchNone,
		},
		{
			name:     "diff mode fetching",
			state:    State{Mode: model.DetailModeDiff, Content: "", Fetching: model.FetchingDetail},
			mode:     model.DetailModeDiff,
			content:  "",
			fetching: model.FetchingDetail,
		},
		{
			name:     "empty content",
			state:    State{Mode: model.DetailModeOverview, Content: "", Fetching: model.FetchNone},
			mode:     model.DetailModeOverview,
			content:  "",
			fetching: model.FetchNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.state.Mode != tt.mode {
				t.Errorf("Mode = %v, want %v", tt.state.Mode, tt.mode)
			}
			if tt.state.Content != tt.content {
				t.Errorf("Content = %q, want %q", tt.state.Content, tt.content)
			}
			if tt.state.Fetching != tt.fetching {
				t.Errorf("Fetching = %v, want %v", tt.state.Fetching, tt.fetching)
			}
		})
	}
}
