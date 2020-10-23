package transitions

type StateType int

const (
	Created StateType = iota
	Started
	Processing
	Finished
	Failed
)
