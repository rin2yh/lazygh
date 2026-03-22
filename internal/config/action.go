package config

// Action represents a user-triggerable action in the TUI.
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
	ActionReviewEvent
	ActionReviewDeleteComment
	ActionReviewEditComment
	ActionShowHelp
	ActionFilterPRs
)

// ActionSpec holds an action's canonical name and default key bindings.
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
	{Action: ActionReviewSubmit, Name: "Review Submit", DefaultKeys: []string{"ctrl+r"}},
	{Action: ActionReviewDiscard, Name: "Review Discard", DefaultKeys: []string{"X"}},
	{Action: ActionReviewSave, Name: "Review Save", DefaultKeys: []string{"ctrl+s"}},
	{Action: ActionReviewClearComment, Name: "Review Clear Comment", DefaultKeys: []string{"x"}},
	{Action: ActionReviewEvent, Name: "Review Event", DefaultKeys: []string{"e"}},
	{Action: ActionReviewDeleteComment, Name: "Review Delete Comment", DefaultKeys: []string{"D"}},
	{Action: ActionReviewEditComment, Name: "Review Edit Comment", DefaultKeys: []string{"i"}},
	{Action: ActionShowHelp, Name: "Show Help", DefaultKeys: []string{"?"}},
	{Action: ActionFilterPRs, Name: "Filter PRs", DefaultKeys: []string{"/"}},
}
