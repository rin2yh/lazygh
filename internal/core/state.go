package core

import "fmt"

type PanelType int

const (
	PanelRepos PanelType = iota
	PanelIssues
	PanelPRs
	PanelDetail
	panelCount
)

type Item struct {
	Number int
	Title  string
}

type EnterActionKind int

const (
	EnterNone EnterActionKind = iota
	EnterLoadItems
	EnterLoadIssueDetail
	EnterLoadPRDetail
)

type EnterAction struct {
	Kind   EnterActionKind
	Repo   string
	Number int
}

type State struct {
	ActivePanel PanelType

	Repos         []Item
	ReposLoading  bool
	ReposSelected int
	ReposLoaded   bool

	Issues         []Item
	IssuesLoading  bool
	IssuesSelected int

	PRs         []Item
	PRsLoading  bool
	PRsSelected int

	DetailContent string

	Width  int
	Height int
}

func NewState() *State {
	return &State{
		ActivePanel: PanelRepos,
		Repos:       []Item{},
		Issues:      []Item{},
		PRs:         []Item{},
	}
}

func (s *State) SetWindowSize(width int, height int) {
	s.Width = width
	s.Height = height
}

func (s *State) BeginLoadRepos() {
	s.ReposLoading = true
}

func (s *State) ApplyReposResult(repos []string, err error) {
	s.ReposLoading = false
	if err != nil {
		s.showError("Error loading repos", err)
		return
	}
	s.Repos = toRepoItems(repos)
	s.ReposSelected = 0
	s.ReposLoaded = true
}

func (s *State) ApplyItemsResult(repo string, issues []Item, prs []Item, err error) {
	s.IssuesLoading = false
	s.PRsLoading = false
	if err != nil {
		s.showError("Error loading items", err)
		return
	}
	currentRepo, ok := s.SelectedRepo()
	if !ok || currentRepo != repo {
		return
	}

	s.Issues = issues
	s.PRs = prs
	s.IssuesSelected = 0
	s.PRsSelected = 0
	s.DetailContent = ""
}

func (s *State) ApplyDetailResult(content string, err error) {
	if err != nil {
		s.showError("Error loading detail", err)
		return
	}
	s.DetailContent = sanitizeMultiline(content)
}

func (s *State) NextPanel() {
	s.ActivePanel = (s.ActivePanel + 1) % panelCount
}

func (s *State) PrevPanel() {
	if s.ActivePanel == PanelRepos {
		s.ActivePanel = PanelDetail
		return
	}
	s.ActivePanel--
}

func (s *State) NavigateDown() {
	switch s.ActivePanel {
	case PanelRepos:
		if s.ReposSelected < len(s.Repos)-1 {
			s.ReposSelected++
		}
	case PanelIssues:
		if s.IssuesSelected < len(s.Issues)-1 {
			s.IssuesSelected++
		}
		s.refreshDetailPreview()
	case PanelPRs:
		if s.PRsSelected < len(s.PRs)-1 {
			s.PRsSelected++
		}
		s.refreshDetailPreview()
	}
}

func (s *State) NavigateUp() {
	switch s.ActivePanel {
	case PanelRepos:
		if s.ReposSelected > 0 {
			s.ReposSelected--
		}
	case PanelIssues:
		if s.IssuesSelected > 0 {
			s.IssuesSelected--
		}
		s.refreshDetailPreview()
	case PanelPRs:
		if s.PRsSelected > 0 {
			s.PRsSelected--
		}
		s.refreshDetailPreview()
	}
}

func (s *State) PlanEnter(hasClient bool, forcedDetailText string) EnterAction {
	switch s.ActivePanel {
	case PanelRepos:
		repo, ok := s.SelectedRepo()
		if !ok || !hasClient {
			return EnterAction{}
		}
		s.IssuesLoading = true
		s.PRsLoading = true
		s.DetailContent = "Loading items..."
		return EnterAction{Kind: EnterLoadItems, Repo: repo}
	case PanelIssues:
		if !hasClient {
			return EnterAction{}
		}
		repo, ok := s.SelectedRepo()
		if !ok {
			return EnterAction{}
		}
		item, ok := s.selectedIssue()
		if !ok {
			return EnterAction{}
		}
		if forcedDetailText != "" {
			s.DetailContent = forcedDetailText
			return EnterAction{}
		}
		s.DetailContent = "Loading detail..."
		return EnterAction{Kind: EnterLoadIssueDetail, Repo: repo, Number: item.Number}
	case PanelPRs:
		if !hasClient {
			return EnterAction{}
		}
		repo, ok := s.SelectedRepo()
		if !ok {
			return EnterAction{}
		}
		item, ok := s.selectedPR()
		if !ok {
			return EnterAction{}
		}
		if forcedDetailText != "" {
			s.DetailContent = forcedDetailText
			return EnterAction{}
		}
		s.DetailContent = "Loading detail..."
		return EnterAction{Kind: EnterLoadPRDetail, Repo: repo, Number: item.Number}
	default:
		return EnterAction{}
	}
}

func (s *State) SelectedRepo() (string, bool) {
	if len(s.Repos) == 0 {
		return "", false
	}
	if s.ReposSelected < 0 || s.ReposSelected >= len(s.Repos) {
		return "", false
	}
	return FormatRepoItem(s.Repos[s.ReposSelected]), true
}

func (s *State) refreshDetailPreview() {
	switch s.ActivePanel {
	case PanelIssues:
		item, ok := s.selectedIssue()
		if !ok {
			return
		}
		s.DetailContent = FormatIssueItem(item)
	case PanelPRs:
		item, ok := s.selectedPR()
		if !ok {
			return
		}
		s.DetailContent = FormatPRItem(item)
	}
}

func (s *State) selectedIssue() (Item, bool) {
	if len(s.Issues) == 0 {
		return Item{}, false
	}
	if s.IssuesSelected < 0 || s.IssuesSelected >= len(s.Issues) {
		return Item{}, false
	}
	return s.Issues[s.IssuesSelected], true
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
	s.DetailContent = sanitizeMultiline(fmt.Sprintf("%s: %v", msg, err))
}

func toRepoItems(repos []string) []Item {
	items := make([]Item, 0, len(repos))
	for _, repo := range repos {
		items = append(items, Item{Title: sanitizeSingleLine(repo)})
	}
	return items
}

func FormatRepoItem(item Item) string {
	return sanitizeSingleLine(item.Title)
}

func FormatIssueItem(item Item) string {
	return fmt.Sprintf("Issue #%d %s", item.Number, sanitizeSingleLine(item.Title))
}

func FormatPRItem(item Item) string {
	return fmt.Sprintf("PR #%d %s", item.Number, sanitizeSingleLine(item.Title))
}
