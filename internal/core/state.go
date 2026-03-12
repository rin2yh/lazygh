package core

import "fmt"

type Item struct {
	Number    int
	Title     string
	Status    string
	Assignees []string
}

type DetailMode int

const (
	DetailModeOverview DetailMode = iota
	DetailModeDiff
)

type EnterActionKind int

const (
	EnterNone EnterActionKind = iota
	EnterLoadPRDetail
	EnterLoadPRDiff
)

type LoadingKind int

const (
	LoadingNone LoadingKind = iota
	LoadingPRs
	LoadingDetail
)

type EnterAction struct {
	Kind   EnterActionKind
	Repo   string
	Number int
}

type State struct {
	Repo string

	PRs         []Item
	PRsLoading  bool
	PRsSelected int

	DetailMode    DetailMode
	DetailContent string
	Loading       LoadingKind

	Width  int
	Height int
}

func NewState() *State {
	return &State{
		PRs:        []Item{},
		DetailMode: DetailModeOverview,
	}
}

func (s *State) SetWindowSize(width int, height int) {
	s.Width = width
	s.Height = height
}

func (s *State) BeginLoadPRs() {
	s.PRsLoading = true
	s.Loading = LoadingPRs
}

func (s *State) ApplyPRsResult(repo string, prs []Item, err error) {
	s.PRsLoading = false
	s.Loading = LoadingNone
	if err != nil {
		s.showError("Error loading PRs", err)
		return
	}

	s.Repo = repo
	s.PRs = prs
	s.PRsSelected = 0
	s.DetailMode = DetailModeOverview
	if len(prs) == 0 {
		s.DetailContent = "No pull requests"
		return
	}
	s.DetailContent = FormatPROverview(prs[s.PRsSelected])
}

func (s *State) ApplyDetailResult(content string, err error) {
	if err != nil {
		s.showError("Error loading detail", err)
		return
	}
	s.Loading = LoadingNone
	s.DetailContent = sanitizeMultiline(content)
}

func (s *State) ApplyDiffResult(content string, err error) {
	if err != nil {
		s.showError("Error loading diff", err)
		return
	}
	s.Loading = LoadingNone
	s.DetailContent = sanitizeMultiline(content)
}

func (s *State) NavigateDown() bool {
	changed := false
	if s.PRsSelected < len(s.PRs)-1 {
		s.PRsSelected++
		changed = true
	}
	if changed && s.DetailMode == DetailModeOverview {
		s.refreshDetailPreview()
	}
	return changed
}

func (s *State) NavigateUp() bool {
	changed := false
	if s.PRsSelected > 0 {
		s.PRsSelected--
		changed = true
	}
	if changed && s.DetailMode == DetailModeOverview {
		s.refreshDetailPreview()
	}
	return changed
}

func (s *State) SwitchToOverview() bool {
	if s.DetailMode == DetailModeOverview {
		return false
	}
	s.DetailMode = DetailModeOverview
	s.Loading = LoadingNone
	s.refreshDetailPreview()
	return true
}

func (s *State) SwitchToDiff() bool {
	if s.DetailMode == DetailModeDiff {
		return false
	}
	s.DetailMode = DetailModeDiff
	s.Loading = LoadingNone
	return true
}

func (s *State) IsDiffMode() bool {
	return s.DetailMode == DetailModeDiff
}

func (s *State) ShouldApplyDetailResult(mode DetailMode, number int) bool {
	if s.DetailMode != mode {
		return false
	}
	item, ok := s.selectedPR()
	if !ok {
		return false
	}
	return item.Number == number
}

func (s *State) PlanEnter(hasClient bool, forcedDetailText string) EnterAction {
	if !hasClient || s.PRsLoading {
		return EnterAction{}
	}
	item, ok := s.selectedPR()
	if !ok {
		return EnterAction{}
	}
	if forcedDetailText != "" {
		s.Loading = LoadingNone
		s.DetailContent = forcedDetailText
		return EnterAction{}
	}
	s.Loading = LoadingDetail
	if s.DetailMode == DetailModeDiff {
		return EnterAction{Kind: EnterLoadPRDiff, Repo: s.Repo, Number: item.Number}
	}
	return EnterAction{Kind: EnterLoadPRDetail, Repo: s.Repo, Number: item.Number}
}

func (s *State) refreshDetailPreview() {
	item, ok := s.selectedPR()
	if !ok {
		return
	}
	s.DetailContent = FormatPROverview(item)
}

func (s *State) selectedPR() (Item, bool) {
	if len(s.PRs) == 0 {
		return Item{}, false
	}
	if s.PRsSelected < 0 || s.PRsSelected >= len(s.PRs) {
		return Item{}, false
	}
	return s.PRs[s.PRsSelected], true
}

func (s *State) showError(msg string, err error) {
	s.Loading = LoadingNone
	s.DetailContent = sanitizeMultiline(fmt.Sprintf("%s: %v", msg, err))
}

func FormatPRItem(item Item) string {
	return fmt.Sprintf("PR #%d %s", item.Number, sanitizeSingleLine(item.Title))
}

func FormatPROverview(item Item) string {
	status := sanitizeSingleLine(item.Status)
	if status == "" {
		status = "OPEN"
	}

	assignee := "unassigned"
	if len(item.Assignees) > 0 {
		list := make([]string, 0, len(item.Assignees))
		for _, name := range item.Assignees {
			n := sanitizeSingleLine(name)
			if n != "" {
				list = append(list, n)
			}
		}
		if len(list) > 0 {
			assignee = list[0]
			if len(list) > 1 {
				assignee = fmt.Sprintf("%s (+%d)", list[0], len(list)-1)
			}
		}
	}

	return fmt.Sprintf(
		"PR #%d %s\nStatus: %s\nAssignee: %s",
		item.Number,
		sanitizeSingleLine(item.Title),
		status,
		assignee,
	)
}
