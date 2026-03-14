package gui

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

type model struct {
	gui *Gui
}

func (m *model) Init() tea.Cmd {
	if m.gui.client == nil {
		return nil
	}
	m.gui.state.BeginLoadPRs()
	return m.loadPRsCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.gui.state.SetWindowSize(msg.Width, msg.Height)
		return m, nil
	case prsLoadedMsg:
		m.gui.applyPRsResult(msg)
		return m, nil
	case detailLoadedMsg:
		m.gui.applyDetailResult(msg)
		return m, nil
	case reviewCommentSavedMsg:
		m.gui.applyReviewCommentResult(msg)
		return m, nil
	case reviewSubmittedMsg:
		m.gui.applyReviewSubmitResult(msg)
		return m, nil
	case reviewDiscardedMsg:
		m.gui.applyReviewDiscardResult(msg)
		return m, nil
	case tea.KeyMsg:
		if m.gui.state.Review.InputMode != core.ReviewInputNone {
			switch msg.String() {
			case "S":
				return m, m.handleReviewSubmit()
			case "X":
				return m, m.handleReviewDiscard()
			}
			if msg.String() == "ctrl+s" && m.gui.state.Review.InputMode == core.ReviewInputComment {
				return m, m.handleReviewCommentSave()
			}
			if m.gui.handleReviewEditorKey(msg) {
				return m, nil
			}
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.gui.state.Review.InputMode == core.ReviewInputNone && m.gui.state.Review.RangeStart != nil {
				m.gui.state.ClearReviewRangeStart()
				m.gui.state.SetReviewNotice("Range selection cleared.")
				m.gui.focus = panelDiffContent
				return m, nil
			}
			if m.gui.focus == panelReviewDrawer {
				m.gui.stopReviewInput()
				m.gui.focus = panelDiffContent
				return m, nil
			}
			m.gui.focusPRs()
			return m, nil
		case "tab":
			m.gui.cycleFocus()
			return m, nil
		case "j", "down":
			return m, m.handleDownKey()
		case "k", "up":
			return m, m.handleUpKey()
		case "pgdown", "f", " ", "pgup", "b", "home", "g", "end", "G":
			if m.gui.scrollDetailByKey(msg) || m.gui.scrollOverviewByKey(msg) {
				return m, nil
			}
			return m, nil
		case "h":
			return m, m.handleHKey()
		case "l":
			return m, m.handleLKey()
		case "o":
			m.gui.switchToOverview()
			return m, nil
		case "d":
			return m, m.handleDKey()
		case "enter":
			return m, m.handleDetailLoad()
		case "v":
			if !m.gui.state.IsDiffMode() {
				m.gui.state.SetReviewNotice("Review range selection is only available in diff view.")
				return m, nil
			}
			m.gui.toggleReviewRangeSelection()
			return m, nil
		case "c":
			if !m.gui.state.IsDiffMode() {
				m.gui.state.SetReviewNotice("Review comments are only available in diff view.")
				return m, nil
			}
			m.gui.beginReviewCommentFlow()
			return m, nil
		case "R":
			if !m.gui.state.IsDiffMode() {
				m.gui.state.SetReviewNotice("Review summary is only available in diff view.")
				return m, nil
			}
			m.gui.beginReviewSummaryInput()
			return m, nil
		case "S":
			return m, m.handleReviewSubmit()
		case "X":
			return m, m.handleReviewDiscard()
		case "x":
			if m.gui.state.Review.InputMode == core.ReviewInputComment {
				m.gui.commentEditor.SetValue("")
				m.gui.state.SetReviewNotice("Comment input cleared.")
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *model) handleHKey() tea.Cmd {
	m.gui.moveFocus(-1)
	return nil
}

func (m *model) handleLKey() tea.Cmd {
	m.gui.moveFocus(1)
	return nil
}

func (m *model) handleDKey() tea.Cmd {
	if m.gui.switchToDiff() {
		return m.handleDetailLoad()
	}
	return nil
}

func (m *model) handleDownKey() tea.Cmd {
	if m.gui.state.IsDiffMode() {
		switch m.gui.focus {
		case panelPRs:
			changed := m.gui.navigateDown()
			if changed {
				return m.handleDetailLoad()
			}
			return nil
		case panelDiffFiles:
			m.gui.selectNextDiffFile()
			return nil
		case panelDiffContent:
			m.gui.scrollDetailDown()
			return nil
		case panelReviewDrawer:
			return nil
		}
		return nil
	}

	if m.gui.focus == panelPRs {
		m.gui.navigateDown()
	}
	return nil
}

func (m *model) handleUpKey() tea.Cmd {
	if m.gui.state.IsDiffMode() {
		switch m.gui.focus {
		case panelPRs:
			changed := m.gui.navigateUp()
			if changed {
				return m.handleDetailLoad()
			}
			return nil
		case panelDiffFiles:
			m.gui.selectPrevDiffFile()
			return nil
		case panelDiffContent:
			m.gui.scrollDetailUp()
			return nil
		case panelReviewDrawer:
			return nil
		}
		return nil
	}

	if m.gui.focus == panelPRs {
		m.gui.navigateUp()
	}
	return nil
}

func (m *model) View() string {
	return m.gui.render()
}

func toCorePRs(prs []gh.PRItem) []core.Item {
	items := make([]core.Item, 0, len(prs))
	for _, pr := range prs {
		status := pr.State
		if pr.IsDraft {
			status = "DRAFT"
		}
		assignees := make([]string, 0, len(pr.Assignees))
		for _, user := range pr.Assignees {
			name := strings.TrimSpace(user.Login)
			if name != "" {
				assignees = append(assignees, name)
			}
		}
		items = append(items, core.Item{
			Number:    pr.Number,
			Title:     pr.Title,
			Status:    status,
			Assignees: assignees,
		})
	}
	return items
}

func (m *model) loadPRsCmd() tea.Cmd {
	return func() tea.Msg {
		repo, err := m.gui.client.ResolveCurrentRepo()
		if err != nil {
			return prsLoadedMsg{err: err}
		}
		prs, err := m.gui.client.ListPRs(repo)
		if err != nil {
			return prsLoadedMsg{repo: repo, err: err}
		}
		return prsLoadedMsg{repo: repo, prs: toCorePRs(prs)}
	}
}

func (m *model) loadDetailCmd(repo string, number int, mode core.DetailMode) tea.Cmd {
	return func() tea.Msg {
		var (
			content string
			err     error
		)
		switch mode {
		case core.DetailModeDiff:
			content, err = m.gui.client.DiffPR(repo, number)
		default:
			content, err = m.gui.client.ViewPR(repo, number)
		}
		return detailLoadedMsg{mode: mode, number: number, content: content, err: err}
	}
}

func (m *model) handleDetailLoad() tea.Cmd {
	action := m.gui.state.PlanEnter(m.gui.client != nil, os.Getenv("LAZYGH_DEBUG_DETAIL_TEXT"))
	switch action.Kind {
	case core.EnterLoadPRDiff:
		return m.loadDetailCmd(action.Repo, action.Number, core.DetailModeDiff)
	case core.EnterLoadPRDetail:
		return m.loadDetailCmd(action.Repo, action.Number, core.DetailModeOverview)
	default:
		return nil
	}
}

func (m *model) handleReviewCommentSave() tea.Cmd {
	item, ok := m.gui.state.SelectedPR()
	if !ok {
		m.gui.state.SetReviewNotice("No pull request selected.")
		return nil
	}
	comment, err := m.gui.buildReviewCommentDraft(m.gui.commentEditor.Value())
	if err != nil {
		m.gui.state.SetReviewNotice(err.Error())
		return nil
	}
	repo := m.gui.state.Repo
	reviewID := m.gui.state.Review.ReviewID
	ctx := gh.ReviewContext{
		PullRequestID: m.gui.state.Review.PullRequestID,
		CommitOID:     m.gui.state.Review.CommitOID,
	}

	m.gui.state.Loading = core.LoadingReview
	return func() tea.Msg {
		var runErr error
		if reviewID == "" {
			ctx, runErr = m.gui.client.GetReviewContext(repo, item.Number)
			if runErr != nil {
				return reviewCommentSavedMsg{err: runErr}
			}
			reviewID, runErr = m.gui.client.StartPendingReview(repo, item.Number, ctx)
			if runErr != nil {
				return reviewCommentSavedMsg{err: runErr}
			}
		}
		runErr = m.gui.client.AddReviewComment(repo, reviewID, comment)
		return reviewCommentSavedMsg{
			prNumber: item.Number,
			ctx:      ctx,
			reviewID: reviewID,
			comment:  comment,
			err:      runErr,
		}
	}
}

func (m *model) handleReviewSubmit() tea.Cmd {
	if m.gui.state.Review.InputMode == core.ReviewInputSummary {
		m.gui.state.SetReviewSummary(m.gui.summaryEditor.Value())
		m.gui.stopReviewInput()
	}
	if !m.gui.state.HasPendingReview() {
		m.gui.state.SetReviewNotice("No pending review to submit.")
		return nil
	}
	m.gui.state.Loading = core.LoadingReview
	reviewID := m.gui.state.Review.ReviewID
	body := m.gui.state.Review.Summary
	repo := m.gui.state.Repo
	return func() tea.Msg {
		err := m.gui.client.SubmitReview(repo, reviewID, body)
		return reviewSubmittedMsg{reviewID: reviewID, err: err}
	}
}

func (m *model) handleReviewDiscard() tea.Cmd {
	if m.gui.state.Review.InputMode == core.ReviewInputSummary {
		m.gui.stopReviewInput()
	}
	reviewID := m.gui.state.Review.ReviewID
	if reviewID == "" {
		m.gui.state.ResetReviewAfterDiscard("Review draft discarded.")
		return nil
	}
	m.gui.state.Loading = core.LoadingReview
	repo := m.gui.state.Repo
	return func() tea.Msg {
		err := m.gui.client.DeletePendingReview(repo, reviewID)
		return reviewDiscardedMsg{err: err}
	}
}
