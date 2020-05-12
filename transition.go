package stateswitch

type TransitionType string

type TransitionRule struct {
	TransitionType   TransitionType
	SourceStates     States
	DestinationState State
	Condition        Condition
	Transition       Transition
	PostTransition   Transition
}

// IsAllowedToRun validate if current state supported, after then check the condition,
// if it pass then transition is a allowed. Nil condition is automatic approval.
func (tr TransitionRule) IsAllowedToRun(stateSwitch StateSwitch, args TransitionArgs) (bool, error) {
	state := stateSwitch.State()
	if tr.SourceStates.Contain(state) {
		if tr.Condition == nil {
			return true, nil
		}
		return stateSwitch.RunCondition(tr.Condition, args)
	}
	return false, nil
}

type TransitionRules []TransitionRule

// Find search for all matching transitions by transition type
func (tr TransitionRules) Find(transitionType TransitionType) TransitionRules {
	match := TransitionRules{}
	for i := range tr {
		if tr[i].TransitionType == transitionType {
			match = append(match, tr[i])
		}
	}
	return match
}

type TransitionArgs interface{}

// Transition is users business logic, should not set the state or return next state
// If condition return true this function will be executed
type Transition interface{}

// Condition for the transition, transition will be executed only if this function return true
// Can be nil, in this case it's considered as return true, nil
type Condition interface{}
