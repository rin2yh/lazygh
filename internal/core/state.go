package core

import "fmt"

type Item struct {
	Number int
	Title  string
}

type EnterActionKind int

const (
	EnterNone EnterActionKind = iota
	EnterLoadPRDetail
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

	DetailContent string
	Loading       LoadingKind

	Width  int
	Height int
}

func NewState() *State {
	return &State{
		PRs: []Item{},
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
	if len(prs) == 0 {
		s.DetailContent = "No pull requests"
		return
	}
	s.DetailContent = FormatPRItem(prs[s.PRsSelected])
}

func (s *State) ApplyDetailResult(content string, err error) {
	if err != nil {
		s.showError("Error loading detail", err)
		return
	}
	s.Loading = LoadingNone
	s.DetailContent = sanitizeMultiline(content)
}

func (s *State) NavigateDown() {
	if s.PRsSelected < len(s.PRs)-1 {
		s.PRsSelected++
	}
	s.refreshDetailPreview()
}

func (s *State) NavigateUp() {
	if s.PRsSelected > 0 {
		s.PRsSelected--
	}
	s.refreshDetailPreview()
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
	return EnterAction{Kind: EnterLoadPRDetail, Repo: s.Repo, Number: item.Number}
}

func (s *State) refreshDetailPreview() {
	item, ok := s.selectedPR()
	if !ok {
		return
	}
	s.DetailContent = FormatPRItem(item)
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
