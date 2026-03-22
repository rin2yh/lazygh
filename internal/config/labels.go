package config

import "strings"

// Label returns all key labels for action joined by "/".
func (k KeyBindings) Label(action Action) string {
	return strings.Join(k.labels(action), "/")
}

// PrimaryLabel returns the formatted label for the first key bound to action.
func (k KeyBindings) PrimaryLabel(action Action) string {
	keys := k.Binding(action).Keys
	if len(keys) == 0 {
		return ""
	}
	return formatKeyLabel(keys[0])
}

// MoveLabel returns a combined label for the move-up/down actions.
func (k KeyBindings) MoveLabel() string {
	labels := []string{
		k.primaryNonArrowLabel(ActionMoveDown),
		k.primaryNonArrowLabel(ActionMoveUp),
	}
	if k.hasKey(ActionMoveUp, "up") {
		labels = append(labels, formatKeyLabel("up"))
	}
	if k.hasKey(ActionMoveDown, "down") {
		labels = append(labels, formatKeyLabel("down"))
	}
	return joinUnique(labels...)
}

// PanelLabel returns a combined label for panel-prev/next actions.
func (k KeyBindings) PanelLabel() string {
	return joinUnique(k.PrimaryLabel(ActionPanelPrev), k.PrimaryLabel(ActionPanelNext))
}

// PageLabel returns a combined label for page-down/up actions.
func (k KeyBindings) PageLabel() string {
	return joinUnique(k.pagePrimaryLabel(ActionPageDown), k.pagePrimaryLabel(ActionPageUp))
}

// TopBottomLabel returns a combined label for go-top/bottom actions.
func (k KeyBindings) TopBottomLabel() string {
	return joinUnique(k.primaryNavigationLabel(ActionGoTop), k.primaryNavigationLabel(ActionGoBottom))
}

func (k KeyBindings) labels(action Action) []string {
	keys := k.Binding(action).Keys
	labels := make([]string, 0, len(keys))
	for _, key := range keys {
		labels = append(labels, formatKeyLabel(key))
	}
	return labels
}

func formatKeyLabel(key string) string {
	switch key {
	case "up":
		return "↑"
	case "down":
		return "↓"
	case " ":
		return "space"
	default:
		return key
	}
}

func (k KeyBindings) hasKey(action Action, key string) bool {
	for _, candidate := range k.Binding(action).Keys {
		if candidate == key {
			return true
		}
	}
	return false
}

func (k KeyBindings) primaryNonArrowLabel(action Action) string {
	for _, key := range k.Binding(action).Keys {
		if key == "up" || key == "down" {
			continue
		}
		return formatKeyLabel(key)
	}
	return k.PrimaryLabel(action)
}

// pageKeyPriority ranks page keys for display: prefer space/b > f > pgdown/pgup.
// Space and b are the most recognizable, so rank 0; pgdown/pgup are last-resort.
var pageKeyPriority = map[string]int{
	" ":      0,
	"b":      0,
	"f":      1,
	"pgdown": 2,
	"pgup":   2,
}

func (k KeyBindings) pagePrimaryLabel(action Action) string {
	best := ""
	bestRank := 999
	for _, key := range k.Binding(action).Keys {
		rank, ok := pageKeyPriority[key]
		if !ok {
			return formatKeyLabel(key)
		}
		if rank < bestRank {
			best = key
			bestRank = rank
		}
	}
	return formatKeyLabel(best)
}

func (k KeyBindings) primaryNavigationLabel(action Action) string {
	for _, key := range k.Binding(action).Keys {
		switch key {
		case "home", "end":
			continue
		default:
			return formatKeyLabel(key)
		}
	}
	return k.PrimaryLabel(action)
}

func joinUnique(labels ...string) string {
	unique := make([]string, 0, len(labels))
	seen := make(map[string]struct{}, len(labels))
	for _, label := range labels {
		if label == "" {
			continue
		}
		if _, ok := seen[label]; ok {
			continue
		}
		seen[label] = struct{}{}
		unique = append(unique, label)
	}
	return strings.Join(unique, "/")
}
