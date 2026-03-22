package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/app/layout"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
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

// ReviewReader はレビュー状態の読み取り専用アクセスを提供するインターフェース。
// 描画ロジックなど状態参照のみ必要な箇所はこのインターフェースに依存する。
type ReviewReader interface {
	ShouldShowDrawer() bool
	IsIndexWithinPendingRange(path string, commentable bool, idx int) bool
	SummaryValue() string
	CommentInputLines() []string
	SummaryInputLines() []string
	IsEditingComment() bool
	InputMode() review.InputMode
	Summary() string
	EventLabel() string
	Notice() string
	RangeStart() *review.Range
	Comments() []review.Comment
	SelectedCommentIdx() int
	HasRangeStart() bool
	IsInInputMode() bool
	HasPendingReview() bool
	PRNumber() int
}

// ReviewHandler はユーザー入力によるレビュー操作を処理するインターフェース。
// キー入力ハンドラやフロー開始など、状態変更を伴うアクションを担う。
type ReviewHandler interface {
	EditorKey(msg tea.KeyMsg) (tea.Cmd, bool)
	Submit() tea.Cmd
	Discard() tea.Cmd
	SaveComment() tea.Cmd
	SaveEditComment() tea.Cmd
	DeleteComment() tea.Cmd
	StopInput()
	ClearCommentInput()
	CycleReviewEvent()
	EditComment() bool
	SelectNextComment()
	SelectPrevComment()
	ToggleRangeSelection()
	BeginCommentFlow()
	BeginSummaryInput()
	SetNotice(msg string)
	ClearRangeStart()
}

// ReviewApplier は非同期操作の結果をレビュー状態に適用するインターフェース。
type ReviewApplier interface {
	CommentResult(msg review.CommentSavedMsg)
	DeleteCommentResult(msg review.CommentDeletedMsg)
	EditCommentResult(msg review.CommentUpdatedMsg)
	SubmitResult(msg review.SubmittedMsg)
	DiscardResult(msg review.DiscardedMsg)
}

// ReviewController は app/ レイヤーが review 機能に要求するインターフェース。
// ISP に従い ReviewReader / ReviewHandler / ReviewApplier に分割されており、
// 各呼び出し側は必要な責務のみに依存できる。
type ReviewController interface {
	ReviewReader
	ReviewHandler
	ReviewApplier
}

// DetailViewport は app/ レイヤーが detail 機能に要求するインターフェース。
type DetailViewport interface {
	Sync(width, height int, body string)
	Height() int
	Update(msg tea.KeyMsg) (bool, tea.Cmd)
	ScrollDown(lines int)
	ScrollUp(lines int)
	GotoTop()
	GotoBottom()
	View() string
}

type Gui struct {
	config *config.Config
	coord  *Coordinator
	client PRClient

	focus    layout.Focus
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
		focus:  layout.FocusPRs,
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
