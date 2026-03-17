package overview

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rin2yh/lazygh/internal/model"
)

func TestState_ZeroValue(t *testing.T) {
	var s State
	want := State{Mode: model.DetailModeOverview, Content: "", Fetching: model.FetchNone}
	if diff := cmp.Diff(want, s); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func TestState_Fields(t *testing.T) {
	tests := []struct {
		name string
		got  State
		want State
	}{
		{
			name: "overview mode with content",
			got:  State{Mode: model.DetailModeOverview, Content: "PR body", Fetching: model.FetchNone},
			want: State{Mode: model.DetailModeOverview, Content: "PR body", Fetching: model.FetchNone},
		},
		{
			name: "diff mode fetching",
			got:  State{Mode: model.DetailModeDiff, Content: "", Fetching: model.FetchingDetail},
			want: State{Mode: model.DetailModeDiff, Content: "", Fetching: model.FetchingDetail},
		},
		{
			name: "empty content",
			got:  State{Mode: model.DetailModeOverview, Content: "", Fetching: model.FetchNone},
			want: State{Mode: model.DetailModeOverview, Content: "", Fetching: model.FetchNone},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tt.got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
