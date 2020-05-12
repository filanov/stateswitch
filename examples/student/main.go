package main

import (
	"fmt"
	"github.com/filanov/stateswitch"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	StatePending = "Pending"
	StatePassed  = "Passed"
	StateFailed  = "Failed"
)

const (
	TransitionTypeSetGrade = "SetGrade"
)

// the model that need state machine will work on
type Student struct {
	ID     string
	Grade  int
	Status string
}

/////////////////////////////////////////////////////////////////////////////
// make sure the model implement stateswich.StateSwitch interface
/////////////////////////////////////////////////////////////////////////////

func (s *Student) SetState(state stateswitch.State) error {
	s.Status = string(state)
	return nil
}

func (s *Student) State() stateswitch.State {
	return stateswitch.State(s.Status)
}

func (s *Student) RunCondition(ifn stateswitch.Condition, args stateswitch.TransitionArgs) (bool, error) {
	fn, ok := ifn.(func (student *Student, args stateswitch.TransitionArgs) (bool, error))
	if !ok {
		return false, fmt.Errorf("Condition function type is not applicable ...")
	}
	return fn(s, args)
}

func (s *Student) RunTransition(ifn stateswitch.Transition, args stateswitch.TransitionArgs) error {
	fn, ok := ifn.(func(student *Student, args stateswitch.TransitionArgs) error)
	if !ok {
		return fmt.Errorf("Transition function type is not applicable ...")
	}
	return fn(s, args)
}

// transition implementation
func SetGradeTransition(s *Student, args stateswitch.TransitionArgs) error {
	grade, ok := args.(int)
	if !ok {
		return errors.Errorf("invalid argument type for SetGrade transition")
	}
	s.Grade = grade
	return nil
}

// Pass condition
func IsPassed(s *Student, args stateswitch.TransitionArgs) (bool, error) {
	grade, ok := args.(int)
	if !ok {
		return false, errors.Errorf("invalid arguments for IsPassed condition")
	}
	if grade > 60 {
		return true, nil
	}
	return false, nil
}

// Failure condition
func IsFailed(s *Student, args stateswitch.TransitionArgs) (bool, error) {
	reply, err := IsPassed(s, args)
	return !reply, err
}

func NewStudentMachine() stateswitch.StateMachine {
	sm := stateswitch.NewStateMachine()

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeSetGrade,
		SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
		DestinationState: StatePassed,
		Condition:        IsPassed,
		Transition:       SetGradeTransition,
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeSetGrade,
		SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
		DestinationState: StateFailed,
		Condition:        IsFailed,
		Transition:       SetGradeTransition,
	})

	return sm
}

func main() {
	student := Student{
		ID:     "123",
		Status: StatePending,
	}

	sm := NewStudentMachine()
	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, &student, 90); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, &student, 50); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, &student, 80); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
}
