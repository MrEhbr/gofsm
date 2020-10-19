package transitions

type StateType int

const (
	Created StateType = iota
	Started
	Finished
	Failed
)
