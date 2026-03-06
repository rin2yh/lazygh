package config

import "testing"

func TestDefault(t *testing.T) {
	cfg := Default()
	tests := []struct{ name, got, want string }{
		{"ActiveBorderColor", cfg.Theme.ActiveBorderColor, "green"},
		{"InactiveBorderColor", cfg.Theme.InactiveBorderColor, "white"},
		{"Quit", cfg.KeyBindings.Quit, "q"},
		{"NextPanel", cfg.KeyBindings.NextPanel, "Tab"},
		{"PrevPanel", cfg.KeyBindings.PrevPanel, "ShiftTab"},
		{"NavigateDown", cfg.KeyBindings.NavigateDown, "j"},
		{"NavigateUp", cfg.KeyBindings.NavigateUp, "k"},
		{"Select", cfg.KeyBindings.Select, "Enter"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}
