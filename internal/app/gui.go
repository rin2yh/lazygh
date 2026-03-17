package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/model"
	"github.com/rin2yh/lazygh/internal/pr/diff"
	"github.com/rin2yh/lazygh/internal/review"
	"github.com/rin2yh/lazygh/pkg/gui/viewport"
)

type PRClient interface {
	ResolveCurrentRepo() (string, error)
	ListPRs(repo string, state string) ([]gh.PRItem, error)
	ViewPR(repo string, number int) (string, error)
	DiffPR(repo string, number int) (string, error)
}

// ReviewController は app/ レイヤーが review 機能に要求するインターフェース。
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

// DetailViewport は app/ レイヤーが detail 機能に要求するインターフェース。
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
	coord  *Coordinator
	client PRClient

	focus    panelFocus
	showHelp bool

	diff   diff.Selection
	detail DetailViewport

	review ReviewController
}

func NewGui(cfg *config.Config, coord *Coordinator, prClient PRClient, reviewClient review.PendingReviewClient) (*Gui, error) {
	vp := viewport.New()
	g := &Gui{
		config: cfg,
		coord:  coord,
		client: prClient,
		focus:  panelPRs,
		detail: &vp,
	}
	revCtrl := review.NewController(cfg, coord, reviewClient, &g.diff, g.setReviewFocus)
	g.review = revCtrl
	coord.SetReviewHook(revCtrl)
	return g, nil
}

func (gui *Gui) Run() error {
	p := tea.NewProgram(&screen{gui: gui}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
