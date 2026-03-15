package config

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Theme struct {
	ActiveBorderColor   string
	InactiveBorderColor string
}

type Action int

const (
	ActionQuit Action = iota
	ActionCancel
	ActionFocusNext
	ActionMoveDown
	ActionMoveUp
	ActionPageDown
	ActionPageUp
	ActionGoTop
	ActionGoBottom
	ActionPanelPrev
	ActionPanelNext
	ActionShowOverview
	ActionShowDiff
	ActionOpenSelected
	ActionReviewRange
	ActionReviewComment
	ActionReviewSummary
	ActionReviewSubmit
	ActionReviewDiscard
	ActionReviewSave
	ActionReviewClearComment
)

type KeyBinding struct {
	Keys []string
}

type KeyBindings struct {
	Quit               KeyBinding
	Cancel             KeyBinding
	FocusNext          KeyBinding
	MoveDown           KeyBinding
	MoveUp             KeyBinding
	PageDown           KeyBinding
	PageUp             KeyBinding
	GoTop              KeyBinding
	GoBottom           KeyBinding
	PanelPrev          KeyBinding
	PanelNext          KeyBinding
	ShowOverview       KeyBinding
	ShowDiff           KeyBinding
	OpenSelected       KeyBinding
	ReviewRange        KeyBinding
	ReviewComment      KeyBinding
	ReviewSummary      KeyBinding
	ReviewSubmit       KeyBinding
	ReviewDiscard      KeyBinding
	ReviewSave         KeyBinding
	ReviewClearComment KeyBinding
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
			Quit:               KeyBinding{Keys: []string{"q", "ctrl+c"}},
			Cancel:             KeyBinding{Keys: []string{"esc"}},
			FocusNext:          KeyBinding{Keys: []string{"tab"}},
			MoveDown:           KeyBinding{Keys: []string{"j", "down"}},
			MoveUp:             KeyBinding{Keys: []string{"k", "up"}},
			PageDown:           KeyBinding{Keys: []string{"pgdown", "f", " "}},
			PageUp:             KeyBinding{Keys: []string{"pgup", "b"}},
			GoTop:              KeyBinding{Keys: []string{"home", "g"}},
			GoBottom:           KeyBinding{Keys: []string{"end", "G"}},
			PanelPrev:          KeyBinding{Keys: []string{"h"}},
			PanelNext:          KeyBinding{Keys: []string{"l"}},
			ShowOverview:       KeyBinding{Keys: []string{"o"}},
			ShowDiff:           KeyBinding{Keys: []string{"d"}},
			OpenSelected:       KeyBinding{Keys: []string{"r"}},
			ReviewRange:        KeyBinding{Keys: []string{"v"}},
			ReviewComment:      KeyBinding{Keys: []string{"enter"}},
			ReviewSummary:      KeyBinding{Keys: []string{"R"}},
			ReviewSubmit:       KeyBinding{Keys: []string{"S"}},
			ReviewDiscard:      KeyBinding{Keys: []string{"X"}},
			ReviewSave:         KeyBinding{Keys: []string{"ctrl+s"}},
			ReviewClearComment: KeyBinding{Keys: []string{"x"}},
		},
	}
}

func (k KeyBindings) Matches(msg tea.KeyMsg, action Action) bool {
	key := msg.String()
	for _, candidate := range k.binding(action).Keys {
		if candidate == key {
			return true
		}
	}
	return false
}

func (k KeyBindings) Label(action Action) string {
	return strings.Join(k.labels(action), "/")
}

func (k KeyBindings) PrimaryLabel(action Action) string {
	labels := k.labels(action)
	if len(labels) == 0 {
		return ""
	}
	return labels[0]
}

func (k KeyBindings) binding(action Action) KeyBinding {
	switch action {
	case ActionQuit:
		return k.Quit
	case ActionCancel:
		return k.Cancel
	case ActionFocusNext:
		return k.FocusNext
	case ActionMoveDown:
		return k.MoveDown
	case ActionMoveUp:
		return k.MoveUp
	case ActionPageDown:
		return k.PageDown
	case ActionPageUp:
		return k.PageUp
	case ActionGoTop:
		return k.GoTop
	case ActionGoBottom:
		return k.GoBottom
	case ActionPanelPrev:
		return k.PanelPrev
	case ActionPanelNext:
		return k.PanelNext
	case ActionShowOverview:
		return k.ShowOverview
	case ActionShowDiff:
		return k.ShowDiff
	case ActionOpenSelected:
		return k.OpenSelected
	case ActionReviewRange:
		return k.ReviewRange
	case ActionReviewComment:
		return k.ReviewComment
	case ActionReviewSummary:
		return k.ReviewSummary
	case ActionReviewSubmit:
		return k.ReviewSubmit
	case ActionReviewDiscard:
		return k.ReviewDiscard
	case ActionReviewSave:
		return k.ReviewSave
	case ActionReviewClearComment:
		return k.ReviewClearComment
	default:
		return KeyBinding{}
	}
}

func (k KeyBindings) labels(action Action) []string {
	keys := k.binding(action).Keys
	labels := make([]string, 0, len(keys))
	for _, key := range keys {
		labels = append(labels, formatKeyLabel(key))
	}
	return labels
}

func (k KeyBindings) QuitLabel() string {
	return k.PrimaryLabel(ActionQuit)
}

func (k KeyBindings) ReloadLabel() string {
	return k.PrimaryLabel(ActionOpenSelected)
}

func (k KeyBindings) FocusLabel() string {
	return k.PrimaryLabel(ActionFocusNext)
}

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

func (k KeyBindings) PanelLabel() string {
	return joinUnique(k.PrimaryLabel(ActionPanelPrev), k.PrimaryLabel(ActionPanelNext))
}

func (k KeyBindings) PageLabel() string {
	return joinUnique(k.pagePrimaryLabel(ActionPageDown), k.pagePrimaryLabel(ActionPageUp))
}

func (k KeyBindings) TopBottomLabel() string {
	return joinUnique(k.primaryNavigationLabel(ActionGoTop), k.primaryNavigationLabel(ActionGoBottom))
}

func (k KeyBindings) ReviewModeLabel() string {
	return joinUnique(k.PrimaryLabel(ActionReviewComment), k.PrimaryLabel(ActionReviewSummary))
}

func (k KeyBindings) SaveLabel() string {
	return k.PrimaryLabel(ActionReviewSave)
}

func (k KeyBindings) CancelLabel() string {
	return k.PrimaryLabel(ActionCancel)
}

func (k KeyBindings) SubmitLabel() string {
	return k.PrimaryLabel(ActionReviewSubmit)
}

func (k KeyBindings) DiscardLabel() string {
	return k.PrimaryLabel(ActionReviewDiscard)
}

func (k KeyBindings) DiffLabel() string {
	return k.PrimaryLabel(ActionShowDiff)
}

func (k KeyBindings) OverviewLabel() string {
	return k.PrimaryLabel(ActionShowOverview)
}

func (k KeyBindings) RangeLabel() string {
	return k.PrimaryLabel(ActionReviewRange)
}

func (k KeyBindings) CommentLabel() string {
	return k.PrimaryLabel(ActionReviewComment)
}

func (k KeyBindings) SummaryLabel() string {
	return k.PrimaryLabel(ActionReviewSummary)
}

func formatKeyLabel(key string) string {
	switch key {
	case "ctrl+c":
		return "Ctrl+C"
	case "ctrl+s":
		return "Ctrl+S"
	case "esc":
		return "Esc"
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
	for _, candidate := range k.binding(action).Keys {
		if candidate == key {
			return true
		}
	}
	return false
}

func (k KeyBindings) primaryNonArrowLabel(action Action) string {
	for _, key := range k.binding(action).Keys {
		if key == "up" || key == "down" {
			continue
		}
		return formatKeyLabel(key)
	}
	return k.PrimaryLabel(action)
}

func (k KeyBindings) pagePrimaryLabel(action Action) string {
	priority := map[string]int{
		" ":      0,
		"b":      0,
		"f":      1,
		"pgdown": 2,
		"pgup":   2,
	}
	best := ""
	bestRank := 999
	for _, key := range k.binding(action).Keys {
		rank, ok := priority[key]
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
	for _, key := range k.binding(action).Keys {
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
