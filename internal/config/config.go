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
	ActionShowHelp
)

type ActionSpec struct {
	Action      Action
	Name        string
	DefaultKeys []string
}

var actionSpecs = []ActionSpec{
	{Action: ActionQuit, Name: "Quit", DefaultKeys: []string{"q", "ctrl+c"}},
	{Action: ActionCancel, Name: "Cancel", DefaultKeys: []string{"esc"}},
	{Action: ActionFocusNext, Name: "Focus Next", DefaultKeys: []string{"tab"}},
	{Action: ActionMoveDown, Name: "Move Down", DefaultKeys: []string{"j", "down"}},
	{Action: ActionMoveUp, Name: "Move Up", DefaultKeys: []string{"k", "up"}},
	{Action: ActionPageDown, Name: "Page Down", DefaultKeys: []string{"pgdown", "f", " "}},
	{Action: ActionPageUp, Name: "Page Up", DefaultKeys: []string{"pgup", "b"}},
	{Action: ActionGoTop, Name: "Go Top", DefaultKeys: []string{"home", "g"}},
	{Action: ActionGoBottom, Name: "Go Bottom", DefaultKeys: []string{"end", "G"}},
	{Action: ActionPanelPrev, Name: "Panel Prev", DefaultKeys: []string{"h"}},
	{Action: ActionPanelNext, Name: "Panel Next", DefaultKeys: []string{"l"}},
	{Action: ActionShowOverview, Name: "Show Overview", DefaultKeys: []string{"o"}},
	{Action: ActionShowDiff, Name: "Show Diff", DefaultKeys: []string{"d"}},
	{Action: ActionOpenSelected, Name: "Open Selected", DefaultKeys: []string{"r"}},
	{Action: ActionReviewRange, Name: "Review Range", DefaultKeys: []string{"v"}},
	{Action: ActionReviewComment, Name: "Review Comment", DefaultKeys: []string{"enter"}},
	{Action: ActionReviewSummary, Name: "Review Summary", DefaultKeys: []string{"R"}},
	{Action: ActionReviewSubmit, Name: "Review Submit", DefaultKeys: []string{"S"}},
	{Action: ActionReviewDiscard, Name: "Review Discard", DefaultKeys: []string{"X"}},
	{Action: ActionReviewSave, Name: "Review Save", DefaultKeys: []string{"ctrl+s"}},
	{Action: ActionReviewClearComment, Name: "Review Clear Comment", DefaultKeys: []string{"x"}},
	{Action: ActionShowHelp, Name: "Show Help", DefaultKeys: []string{"?"}},
}

type KeyBinding struct {
	Keys []string
}

type KeyBindings struct {
	bindings map[Action]KeyBinding
}

type Config struct {
	Theme       Theme
	KeyBindings KeyBindings
}

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

func newKeyBindings() KeyBindings {
	return KeyBindings{bindings: make(map[Action]KeyBinding, len(actionSpecs))}
}

func (k KeyBindings) Matches(msg tea.KeyMsg, action Action) bool {
	key := msg.String()
	for _, candidate := range k.Binding(action).Keys {
		if candidate == key {
			return true
		}
	}
	return false
}

func (k KeyBindings) ActionFor(msg tea.KeyMsg) (Action, bool) {
	for _, spec := range actionSpecs {
		if k.Matches(msg, spec.Action) {
			return spec.Action, true
		}
	}
	return 0, false
}

func (k KeyBindings) Binding(action Action) KeyBinding {
	if k.bindings == nil {
		return KeyBinding{}
	}
	return k.bindings[action]
}

func (k *KeyBindings) SetBinding(action Action, binding KeyBinding) {
	if k.bindings == nil {
		k.bindings = make(map[Action]KeyBinding, len(actionSpecs))
	}
	k.bindings[action] = binding
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

func (k KeyBindings) HelpLabel() string {
	return k.PrimaryLabel(ActionShowHelp)
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
	for _, key := range k.Binding(action).Keys {
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
