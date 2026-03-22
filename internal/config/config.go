// Package config provides theme and key binding configuration.
package config

// Theme holds color settings for TUI borders.
type Theme struct {
	ActiveBorderColor   string
	InactiveBorderColor string
}

// Config is the top-level application configuration.
type Config struct {
	Theme       Theme
	KeyBindings KeyBindings
}

// Default returns a Config with built-in defaults.
func Default() *Config {
	keys := newKeyBindings()
	for _, spec := range actionSpecs {
		keys.SetBinding(spec.Action, KeyBinding{Keys: append([]string(nil), spec.DefaultKeys...)})
	}

	return &Config{
		Theme: Theme{
			ActiveBorderColor:   "green",
			InactiveBorderColor: "white",
		},
		KeyBindings: keys,
	}
}
