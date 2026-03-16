package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/detail"
	guidiff "github.com/rin2yh/lazygh/internal/gui/diff"
	guireview "github.com/rin2yh/lazygh/internal/gui/review"
	appstate "github.com/rin2yh/lazygh/internal/state"
)

type PRClient interface {
	ResolveCurrentRepo() (string, error)
	ListPRs(repo string, state string) ([]gh.PRItem, error)
	ViewPR(repo string, number int) (string, error)
	DiffPR(repo string, number int) (string, error)
}

// ReviewController は gui/ レイヤーが review 機能に要求するインターフェース。
// gui/review/ を internal/review/ へ昇格する際は、tea.Cmd などの
// GUI フレームワーク依存を持たない純粋なドメインインターフェースへ置き換える。
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
	ApplyCommentResult(msg guireview.CommentSavedMsg)
	ApplyDeleteCommentResult(msg guireview.CommentDeletedMsg)
	ApplyEditCommentResult(msg guireview.CommentUpdatedMsg)
	ApplySubmitResult(msg guireview.SubmittedMsg)
	ApplyDiscardResult(msg guireview.DiscardedMsg)
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
}

// DetailViewport は gui/ レイヤーが detail 機能に要求するインターフェース。
// gui/detail/ を internal/detail/ へ昇格する際は、tea.KeyMsg / tea.Cmd などの
// GUI フレームワーク依存を持たない純粋なドメインインターフェースへ置き換える。
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

func NewGui(cfg *config.Config, prClient PRClient, reviewClient guireview.PendingReviewClient) (*Gui, error) {
	d := detail.NewState()
	gui := &Gui{
		config: cfg,
		state:  appstate.NewState(),
		client: prClient,
		focus:  panelPRs,
		detail: &d,
	}
	gui.review = guireview.NewController(cfg, gui.state, reviewClient, &gui.diff, gui.setReviewFocus)
	return gui, nil
}

func (gui *Gui) Run() error {
	p := tea.NewProgram(&screen{gui: gui}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
