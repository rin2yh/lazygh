package config

type Theme struct {
	ActiveBorderColor   string
	InactiveBorderColor string
}

type KeyBindings struct {
	Quit         string
	NextPanel    string
	PrevPanel    string
	NavigateDown string
	NavigateUp   string
	Select       string
}

type Config struct {
	Theme       Theme
	KeyBindings KeyBindings
}

func Default() *Config {
	return &Config{
		Theme: Theme{
			ActiveBorderColor:   "green",
			InactiveBorderColor: "white",
		},
		KeyBindings: KeyBindings{
			Quit:         "q",
			NextPanel:    "Tab",
			PrevPanel:    "ShiftTab",
			NavigateDown: "j",
			NavigateUp:   "k",
			Select:       "Enter",
		},
	}
}
