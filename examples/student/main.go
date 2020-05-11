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

func (s *Student) State() stateswitch.State {
	return stateswitch.State(s.Status)
}

func (s *Student) SetGrade(args stateswitch.TransitionArgs) error {
	grade, ok := args.(int)
	if !ok {
		return errors.Errorf("invalid argument type for SetGrade transition")
	}
	s.Grade = grade
	return nil
}

func (s *Student) IsPassed(args stateswitch.TransitionArgs) (bool, error) {
	grade, ok := args.(int)
	if !ok {
		return false, errors.Errorf("invalid arguments for IsPassed condition")
	}
	if grade > 60 {
		return true, nil
	}
	return false, nil
}

func (s *Student) IsFailed(args stateswitch.TransitionArgs) (bool, error) {
	reply, err := s.IsPassed(args)
	return !reply, err
}

func main() {
	student := Student{
		ID:     "123",
		Status: StatePending,
	}

	sm := stateswitch.NewStateMachine(&student)

	sm.AddTransition(
		TransitionTypeSetGrade,
		[]stateswitch.State{StatePending, StateFailed, StatePassed},
		StatePassed,
		student.SetGrade,
		student.IsPassed,
	)

	sm.AddTransition(
		TransitionTypeSetGrade,
		[]stateswitch.State{StatePending, StateFailed, StatePassed},
		StateFailed,
		student.SetGrade,
		student.IsFailed,
	)

	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, 90); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, 50); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, 80); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run("unknown transition", 50); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
	if err := sm.Run(TransitionTypeSetGrade, "invalid args"); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", student)
}
