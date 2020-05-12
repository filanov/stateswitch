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

type Student struct {
	ID     string
	Grade  int
	Status string
}

func (s *Student) SetState(state stateswitch.State) error {
	s.Status = string(state)
	return nil
}

func (s Student) State() stateswitch.State {
	return stateswitch.State(s.Status)
}

type SetGradeTransitionArgs struct {
	grade   int
	student *Student
}

func SetGrade(args stateswitch.TransitionArgs) error {
	params, ok := args.(*SetGradeTransitionArgs)
	if !ok {
		return errors.Errorf("invalid argument type for SetGrade transition")
	}
	params.student.Grade = params.grade
	return nil
}

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

func IsFailed(args stateswitch.TransitionArgs) (bool, error) {
	reply, err := IsPassed(args)
	return !reply, err
}

func main() {
	student := Student{
		ID:     "123",
		Status: StatePending,
	}

	sm := stateswitch.NewStateMachine()

	sm.AddTransition(stateswitch.TransitionRule{
		SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
		Condition:        IsPassed,
		Transition:       SetGrade,
		TransitionType:   TransitionTypeSetGrade,
		DestinationState: StatePassed,
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeSetGrade,
		SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
		DestinationState: StateFailed,
		Condition:        IsFailed,
		Transition:       SetGrade,
	})

	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, &student, &SetGradeTransitionArgs{
		grade:   90,
		student: &student,
	}); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, &student, &SetGradeTransitionArgs{
		grade:   50,
		student: &student,
	}); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, &student, &SetGradeTransitionArgs{
		grade:   80,
		student: &student,
	}); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run("unknown transition", &student, &SetGradeTransitionArgs{
		grade:   0,
		student: &student,
	}); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, &student, "invalid args"); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
}
