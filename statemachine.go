package stateswitch

import (
	"github.com/pkg/errors"
)

type StateMachine interface {
	AddTransition(transitionType TransitionType, sourceState States, destinationState State, transition Transition, condition Condition)
	Run(transitionType TransitionType, args TransitionArgs) error
}

func NewStateMachine(stateSwitchObj StateSwitch) *stateMachine {
	return &stateMachine{
		StateSwitchObj:  stateSwitchObj,
		transitionRules: map[TransitionType]TransitionRules{},
	}
}

type stateMachine struct {
	StateSwitchObj  StateSwitch
	transitionRules map[TransitionType]TransitionRules
}

func (sm *stateMachine) Run(transitionType TransitionType, args TransitionArgs) error {
	transByType, ok := sm.transitionRules[transitionType]
	if !ok {
		return errors.Errorf("no match for transition type %s", transitionType)
	}

	objState := sm.StateSwitchObj.State()
	for _, tr := range transByType {
		allow, err := tr.IsAllowedToRun(objState, args)
		if err != nil {
			return err
		}
		if allow {
			if err := tr.Transition(args); err != nil {
				return err
			}
			return sm.StateSwitchObj.SetState(tr.DestinationState)
		}
	}
	return errors.Errorf("no condition passed to run transition %s from state %s",
		transitionType, sm.StateSwitchObj.State())
}

func (sm *stateMachine) AddTransition(transitionType TransitionType, sourceState States, destinationState State, transition Transition, condition Condition) {
	sm.transitionRules[transitionType] = append(sm.transitionRules[transitionType],
		TransitionRule{
			SourceStates:     sourceState,
			Condition:        condition,
			Transition:       transition,
			DestinationState: destinationState,
			TransitionType:   transitionType,
		},
	)
}
