package stateswitch

type StateSwitch interface {
	State() State
	SetState(state State) error
}
