package model

type DetailMode int

const (
	DetailModeOverview DetailMode = iota
	DetailModeDiff
)

type FetchKind int

const (
	FetchNone FetchKind = iota
	FetchingPRs
	FetchingDetail
	FetchingReview
)
