package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	guidiff "github.com/rin2yh/lazygh/internal/gui/diff"
	"github.com/rin2yh/lazygh/internal/model"
	"github.com/rin2yh/lazygh/internal/review"
	appstate "github.com/rin2yh/lazygh/internal/state"
)

type PRClient interface {
	ResolveCurrentRepo() (string, error)
	ListPRs(repo string, state string) ([]gh.PRItem, error)
	ViewPR(repo string, number int) (string, error)
	DiffPR(repo string, number int) (string, error)
}

// ReviewController は gui/ レイヤーが review 機能に要求するインターフェース。
type ReviewController interface {
	ShouldShowDrawer() bool
	IsIndexWithinPendingRange(path string, commentable bool, idx int) bool
	CurrentSummaryValue() string
	CommentInputLines() []string
	SummaryInputLines() []string
	HandleEditorKey(msg tea.KeyMsg) (tea.Cmd, bool)
	HandleSubmit() tea.Cmd
	HandleDiscard() tea.Cmd
	HandleCommentSave() tea.Cmd
	HandleEditCommentSave() tea.Cmd
	HandleDeleteComment() tea.Cmd
	ApplyCommentResult(msg review.CommentSavedMsg)
	ApplyDeleteCommentResult(msg review.CommentDeletedMsg)
	ApplyEditCommentResult(msg review.CommentUpdatedMsg)
	ApplySubmitResult(msg review.SubmittedMsg)
	ApplyDiscardResult(msg review.DiscardedMsg)
	CurrentCommentValue() string
	SetCommentValue(value string)
	StopInput()
	ClearCommentInput()
	CycleReviewEvent()
	BeginEditComment() bool
	IsEditingComment() bool
	SelectNextComment()
	SelectPrevComment()
	ToggleRangeSelection()
	BeginCommentFlow()
	BeginSummaryInput()
	// state accessors
	InputMode() model.ReviewInputMode
	Summary() string
	EventLabel() string
	Notice() string
	RangeStart() *model.ReviewRange
	Comments() []model.ReviewComment
	SelectedCommentIdx() int
	HasRangeStart() bool
	IsInInputMode() bool
	HasPendingReview() bool
	PRNumber() int
	SetNotice(msg string)
	ClearRangeStart()
	Reset()
	SetContext(prNumber int, pullRequestID, commitOID, reviewID string)
	OpenDrawer()
	BeginCommentInput()
}

// DetailViewport は gui/ レイヤーが detail 機能に要求するインターフェース。
type DetailViewport interface {
	Sync(width, height int, body string)
	Height() int
	Update(msg tea.KeyMsg) (bool, tea.Cmd)
	ScrollDown(lines int)
	ScrollUp(lines int)
	View() string
}

type Gui struct {
	config *config.Config
	state  *appstate.State
	client PRClient

	focus    panelFocus
	showHelp bool

	diff   guidiff.Selection
	detail DetailViewport

	review ReviewController
}

func NewGui(cfg *config.Config, prClient PRClient, reviewClient review.PendingReviewClient) (*Gui, error) {
	vp := newViewportState()
	gui := &Gui{
		config: cfg,
		state:  appstate.NewState(),
		client: prClient,
		focus:  panelPRs,
		detail: &vp,
	}
	gui.review = review.NewController(cfg, gui.state, reviewClient, &gui.diff, gui.setReviewFocus)
	return gui, nil
}

func (gui *Gui) Run() error {
	p := tea.NewProgram(&screen{gui: gui}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
