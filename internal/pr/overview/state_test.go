package overview

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestState_ZeroValue(t *testing.T) {
	var s State
	want := State{Mode: DetailModeOverview, Content: "", Fetching: FetchNone}
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
			got:  State{Mode: DetailModeOverview, Content: "PR body", Fetching: FetchNone},
			want: State{Mode: DetailModeOverview, Content: "PR body", Fetching: FetchNone},
		},
		{
			name: "diff mode fetching",
			got:  State{Mode: DetailModeDiff, Content: "", Fetching: FetchingDetail},
			want: State{Mode: DetailModeDiff, Content: "", Fetching: FetchingDetail},
		},
		{
			name: "empty content",
			got:  State{Mode: DetailModeOverview, Content: "", Fetching: FetchNone},
			want: State{Mode: DetailModeOverview, Content: "", Fetching: FetchNone},
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
