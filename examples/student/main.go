package main

import (
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

func (s Student) State() stateswitch.State {
	return stateswitch.State(s.Status)
}

// Define arguments for each transition
type SetGradeTransitionArgs struct {
	grade   int
	student *Student
}

// transition implementation
func SetGradeTransition(args stateswitch.TransitionArgs) error {
	params, ok := args.(*SetGradeTransitionArgs)
	if !ok {
		return errors.Errorf("invalid argument type for SetGrade transition")
	}
	params.student.Grade = params.grade
	return nil
}

// Pass condition
func IsPassed(args stateswitch.TransitionArgs) (bool, error) {
	params, ok := args.(*SetGradeTransitionArgs)
	if !ok {
		return false, errors.Errorf("invalid arguments for IsPassed condition")
	}
	if params.grade > 60 {
		return true, nil
	}
	return false, nil
}

// Failure condition
func IsFailed(args stateswitch.TransitionArgs) (bool, error) {
	reply, err := IsPassed(args)
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

func NewStudentMachine() *studentMachine {
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

	stm := &studentMachine{
		sm: sm,
	}

	return stm
}

func main() {
	student := Student{
		ID:     "123",
		Status: StatePending,
	}

	sm := NewStudentMachine()
	logrus.Infof("%+v", student)
	if err := sm.SetGrade(&student, 90); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.SetGrade(&student, 50); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.SetGrade(&student, 80); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
}
