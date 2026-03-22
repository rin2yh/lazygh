package config

import tea "github.com/charmbracelet/bubbletea"

// KeyBinding holds the keys mapped to an action.
type KeyBinding struct {
	Keys []string
}

// KeyBindings maps actions to their key bindings.
type KeyBindings struct {
	bindings map[Action]KeyBinding
}

func newKeyBindings() KeyBindings {
	return KeyBindings{bindings: make(map[Action]KeyBinding, len(actionSpecs))}
}

// Matches reports whether msg matches the keys bound to action.
func (k KeyBindings) Matches(msg tea.KeyMsg, action Action) bool {
	key := msg.String()
	for _, candidate := range k.Binding(action).Keys {
		if candidate == key {
			return true
		}
	}
	return false
}

// ActionFor returns the action matching msg, if any.
func (k KeyBindings) ActionFor(msg tea.KeyMsg) (Action, bool) {
	for _, spec := range actionSpecs {
		if k.Matches(msg, spec.Action) {
			return spec.Action, true
		}
	}
	return 0, false
}

// Binding returns the KeyBinding for action.
func (k KeyBindings) Binding(action Action) KeyBinding {
	if k.bindings == nil {
		return KeyBinding{}
	}
	return k.bindings[action]
}

// SetBinding replaces the key binding for action.
func (k *KeyBindings) SetBinding(action Action, binding KeyBinding) {
	if k.bindings == nil {
		k.bindings = make(map[Action]KeyBinding, len(actionSpecs))
	}
	k.bindings[action] = binding
}
