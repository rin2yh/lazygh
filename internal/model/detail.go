package model

type DetailMode int

const (
	DetailModeOverview DetailMode = iota
	DetailModeDiff
)

type LoadingKind int

const (
	LoadingNone LoadingKind = iota
	LoadingPRs
	LoadingDetail
	LoadingReview
)
