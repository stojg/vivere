package state

type State uint32

const (
	IDLE   State = 0
	DEAD   State = 1
	MOVING State = 2
)

type Stater interface {
	State() State
	SetState(State)
}
