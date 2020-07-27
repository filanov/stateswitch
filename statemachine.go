package stateswitch

import (
	"github.com/pkg/errors"
)

var (
	NoConditionPassedToRunTransaction = errors.New("no condition found to run transition")
	NoMatchForTransitionType          = errors.New("no match for transition type")
)

type StateMachine interface {
	// AddTransition to state machine
	AddTransition(rule TransitionRule)
	// Run transition by type
	Run(transitionType TransitionType, stateSwitch StateSwitch, args TransitionArgs) error
}

// Create new default state machine
func NewStateMachine() *stateMachine {
	return &stateMachine{
		transitionRules: map[TransitionType]TransitionRules{},
	}
}

type stateMachine struct {
	transitionRules map[TransitionType]TransitionRules
}

// Run transition by type, will search for the first transition that will pass a condition.
func (sm *stateMachine) Run(transitionType TransitionType, stateSwitch StateSwitch, args TransitionArgs) error {
	transByType, ok := sm.transitionRules[transitionType]
	if !ok {
		return NoMatchForTransitionType
	}

	for _, tr := range transByType {
		allow, err := tr.IsAllowedToRun(stateSwitch, args)
		if err != nil {
			return err
		}
		if allow {
			if tr.Transition != nil {
				if err := tr.Transition(stateSwitch, args); err != nil {
					return err
				}
			}
			if err := stateSwitch.SetState(tr.DestinationState); err != nil {
				return err
			}
			if tr.PostTransition != nil {
				return tr.PostTransition(stateSwitch, args)
			}
			return nil
		}
	}
	return NoConditionPassedToRunTransaction
}

// AddTransition to state machine
func (sm *stateMachine) AddTransition(rule TransitionRule) {
	sm.transitionRules[rule.TransitionType] = append(sm.transitionRules[rule.TransitionType], rule)
}
