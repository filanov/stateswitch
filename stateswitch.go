package stateswitch

// StateSwitch interface used by state machine
type StateSwitch interface {
	// State return current state
	State() State
	// SetState set a new state
	SetState(state State) error

	RunCondition(fn interface{}, args TransitionArgs) (bool, error)
	RunTransition(fn interface{}, args TransitionArgs) error
}
