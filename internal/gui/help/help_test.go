package help

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/rin2yh/lazygh/internal/config"
)

func makeBackground(lines, width int) []string {
	bg := make([]string, lines)
	for i := range bg {
		bg[i] = strings.Repeat(" ", width)
	}
	return bg
}

func TestBgLeft(t *testing.T) {
	tests := []struct {
		name string
		bg   string
		x    int
		want string
	}{
		{"truncates when bg longer", "abcde", 3, "abc"},
		{"pads when bg shorter", "ab", 5, "ab   "},
		{"x=0 returns empty", "abc", 0, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bgLeft(tt.bg, tt.x)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBgRight(t *testing.T) {
	tests := []struct {
		name    string
		endX    int
		screenW int
		want    string
	}{
		{"fills remaining space", 10, 15, "     "},
		{"endX equals screenW returns empty", 15, 15, ""},
		{"endX exceeds screenW returns empty", 20, 15, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bgRight(tt.endX, tt.screenW)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOverlayLine(t *testing.T) {
	// bg="abcde"(5), panel="XY"(2), startX=3, panelW=2, screenW=10
	// left  = bgLeft("abcde", 3)   = "abc"
	// middle = PadOrTrim("XY", 2)  = "XY"
	// right  = bgRight(3+2, 10)    = "     " (5 spaces)
	got := overlayLine("abcde", "XY", 3, 2, 10)
	want := "abcXY     "
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBuildPanelLines_WidthMatchesLines(t *testing.T) {
	keys := config.Default().KeyBindings
	lines, w := buildPanelLines(keys, 120)
	for i, line := range lines {
		got := xansi.StringWidth(line)
		if got != w {
			t.Errorf("line[%d] width=%d, want %d: %q", i, got, w, line)
		}
	}
}

func TestBuildPanelLines_ClampsToScreenWidth(t *testing.T) {
	keys := config.Default().KeyBindings
	const screenW = 40
	_, w := buildPanelLines(keys, screenW)
	if w != screenW-2 {
		t.Fatalf("got w=%d, want %d", w, screenW-2)
	}
}

func TestRenderOverlay_LineCountPreserved(t *testing.T) {
	tests := []struct {
		name    string
		bgLines int
		screenW int
	}{
		{"normal", 40, 120},
		{"short bg: 0", 0, 120},
		{"short bg: 1", 1, 120},
		{"short bg: 5", 5, 120},
		{"narrow screen: 10", 40, 10},
		{"narrow screen: 30", 40, 30},
		{"narrow screen: 50", 40, 50},
	}

	keys := config.Default().KeyBindings
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bg := makeBackground(tt.bgLines, tt.screenW)
			got := RenderOverlay(bg, keys, tt.screenW)
			if len(got) != tt.bgLines {
				t.Fatalf("got %d lines, want %d", len(got), tt.bgLines)
			}
		})
	}
}

func TestRenderOverlay_ContainsText(t *testing.T) {
	const screenW = 120
	keys := config.Default().KeyBindings
	bg := makeBackground(40, screenW)
	got := RenderOverlay(bg, keys, screenW)
	joined := strings.Join(got, "\n")

	tests := []struct {
		name string
		want string
	}{
		{"panel title", "Keybindings"},
		{"section Navigation", "Navigation"},
		{"section View", "View"},
		{"section Review", "Review"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(joined, tt.want) {
				t.Errorf("%q not found in output", tt.want)
			}
		})
	}
}

func TestRenderOverlay_UntouchedLinesUnchanged(t *testing.T) {
	const screenW = 120
	keys := config.Default().KeyBindings
	bg := makeBackground(40, screenW)
	got := RenderOverlay(bg, keys, screenW)

	// 40行の背景に対してパネルは中央に配置されるので、先頭・末尾行は変化しないはず
	if got[0] != bg[0] {
		t.Errorf("first line changed: got %q", got[0])
	}
	if got[len(got)-1] != bg[len(bg)-1] {
		t.Errorf("last line changed: got %q", got[len(got)-1])
	}
}
