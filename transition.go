package stateswitch

type TransitionType string

type TransitionRule struct {
	SourceStates     States
	Condition        Condition
	Transition       Transition
	TransitionType   TransitionType
	DestinationState State
}

func (tr TransitionRule) IsAllowedToRun(state State, args TransitionArgs) (bool, error) {
	if tr.SourceStates.Contain(state) {
		if tr.Condition == nil {
			return true, nil
		}
		return tr.Condition(args)
	}
	return false, nil
}

type TransitionRules []TransitionRule

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

type Transition func(args TransitionArgs) error

type Condition func(args TransitionArgs) (bool, error)
