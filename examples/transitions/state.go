package transitions

type StateType int

const (
	StateTypeUnknown StateType = iota
	Created
	Started
	Processing
	Finished
	Failed
)
