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

// transition implementation
func SetGradeTransition(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
	s, ok := sw.(*Student)
	if !ok {
		return errors.Errorf("StateSwitch object is not of type Student")
	}
	grade, ok := args.(int)
	if !ok {
		return errors.Errorf("invalid argument type for SetGrade transition")
	}
	s.Grade = grade
	return nil
}

// Pass condition
func IsPassed(_ stateswitch.StateSwitch, args stateswitch.TransitionArgs) (bool, error) {
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
func IsFailed(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) (bool, error) {
	reply, err := IsPassed(sw, args)
	return !reply, err
}

// State machine wrapper
type studentMachine struct {
	sm stateswitch.StateMachine
}

func (stm *studentMachine) SetGrade(s *Student, grade int) error {
	return stm.sm.Run(TransitionTypeSetGrade, s, grade)
}

func NewStudentMachine() *studentMachine {
	sm := stateswitch.NewStateMachine()

	sm.AddTransitionRule(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeSetGrade,
		SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
		DestinationState: StatePassed,
		Condition:        IsPassed,
		Transition:       SetGradeTransition,
	})

	sm.AddTransitionRule(stateswitch.TransitionRule{
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
