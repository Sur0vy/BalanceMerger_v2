package Models

type ItemState int

const (
	IsFound ItemState = iota
	IsMissing
	IsCollect
	IsCollectMissing
	IsDifBalance
)
