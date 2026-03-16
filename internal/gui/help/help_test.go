package help

import (
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
)

func makeBackground(lines, width int) []string {
	bg := make([]string, lines)
	for i := range bg {
		bg[i] = strings.Repeat(" ", width)
	}
	return bg
}

func TestRenderOverlay_PreservesLineCount(t *testing.T) {
	keys := config.Default().KeyBindings
	bg := makeBackground(40, 120)
	got := RenderOverlay(bg, keys, 120)
	if len(got) != len(bg) {
		t.Fatalf("got %d lines, want %d", len(got), len(bg))
	}
}

func TestRenderOverlay_ShortBackgroundNoPanic(t *testing.T) {
	// background の行数がパネルより少なくてもpanicしない
	keys := config.Default().KeyBindings
	for _, lines := range []int{0, 1, 5} {
		bg := makeBackground(lines, 120)
		got := RenderOverlay(bg, keys, 120)
		if len(got) != lines {
			t.Errorf("lines=%d: got %d, want %d", lines, len(got), lines)
		}
	}
}

func TestRenderOverlay_NarrowScreenNoPanic(t *testing.T) {
	// screenWidth がパネルより小さい場合でもパニックせず行数が保たれる
	for _, screenW := range []int{10, 30, 50} {
		bg := makeBackground(40, screenW)
		keys := config.Default().KeyBindings
		got := RenderOverlay(bg, keys, screenW)
		if len(got) != len(bg) {
			t.Errorf("screenW=%d: got %d lines, want %d", screenW, len(got), len(bg))
		}
	}
}

func TestRenderOverlay_ContainsSections(t *testing.T) {
	const screenW = 120
	keys := config.Default().KeyBindings
	bg := makeBackground(40, screenW)
	got := RenderOverlay(bg, keys, screenW)
	joined := strings.Join(got, "\n")

	for _, want := range []string{"Navigation", "View", "Review"} {
		if !strings.Contains(joined, want) {
			t.Errorf("section %q not found in output", want)
		}
	}
}

func TestRenderOverlay_ContainsPanelTitle(t *testing.T) {
	const screenW = 120
	keys := config.Default().KeyBindings
	bg := makeBackground(40, screenW)
	got := RenderOverlay(bg, keys, screenW)

	found := false
	for _, line := range got {
		if strings.Contains(line, "Keybindings") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("no line contains 'Keybindings' title")
	}
}

func TestRenderOverlay_UntouchedLinesUnchanged(t *testing.T) {
	const screenW = 120
	keys := config.Default().KeyBindings
	bg := makeBackground(40, screenW)
	got := RenderOverlay(bg, keys, screenW)

	// panelH は実際の描画を見て確認するため、先頭行と末尾行がそのままかを検証する
	// 40行の背景に対してパネルは中央に配置されるので、先頭数行と末尾数行は変化しないはず
	first := got[0]
	last := got[len(got)-1]
	if first != bg[0] {
		t.Errorf("first line changed: got %q", first)
	}
	if last != bg[len(bg)-1] {
		t.Errorf("last line changed: got %q", last)
	}
}
