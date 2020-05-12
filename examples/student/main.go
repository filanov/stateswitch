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

type ConditionFn  func (student *Student, args stateswitch.TransitionArgs) (bool, error)

type TransitionFn func(student *Student, args stateswitch.TransitionArgs) error

func (s *Student) RunCondition(ifn interface{}, args stateswitch.TransitionArgs) (bool, error) {
	fn, ok := ifn.(ConditionFn)
	if !ok {
		return false, fmt.Errorf("Condition function type is not applicable ...")
	}
	return fn(s, args)
}

func (s *Student) RunTransition(ifn interface{}, args stateswitch.TransitionArgs) error {
	fn, ok := ifn.(TransitionFn)
	if !ok {
		return fmt.Errorf("Transition function type is not applicable ...")
	}
	return fn(s, args)
}

// Define arguments for each transition
type SetGradeTransitionArgs struct {
	grade   int
	student *Student
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

// State machine wrapper
type studentMachine struct {
	sm stateswitch.StateMachine
}

func (stm *studentMachine) SetGrade(s *Student, grade int) error {
	return stm.sm.Run(TransitionTypeSetGrade, s, &SetGradeTransitionArgs{
		grade:   grade,
		student: s,
	})
}

func NewStudentMachine() stateswitch.StateMachine {
	sm := stateswitch.NewStateMachine()

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeSetGrade,
		SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
		DestinationState: StatePassed,
		Condition:        ConditionFn(IsPassed),
		Transition:       TransitionFn(SetGradeTransition),
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeSetGrade,
		SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
		DestinationState: StateFailed,
		Condition:        ConditionFn(IsFailed),
		Transition:       TransitionFn(SetGradeTransition),
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
