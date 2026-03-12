package config

type Theme struct {
	ActiveBorderColor   string
	InactiveBorderColor string
}

type KeyBindings struct {
	Quit         string
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
			NavigateDown: "j",
			NavigateUp:   "k",
			Select:       "Enter",
		},
	}
}
