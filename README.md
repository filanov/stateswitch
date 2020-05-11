[![Actions Status](https://github.com/filanov/stateswitch/workflows/make_all/badge.svg)](https://github.com/filanov/stateswitch/actions)

# stateswitch

## Overview

TDB

## Usage

First need to create an entity that implement `StateSwitch` interface

```go
import "github.com/filanov/stateswitch"

const (
	StatePending = "Pending"
	StatePassed  = "Passed"
	StateFailed  = "Failed"
)

type Student struct {
	ID     string
	Grade  int
	Status string // or stateswitch.State if you don't want to handle transtions
}

func (s *Student) SetState(state stateswitch.State) error {
	s.Status = string(state)
	return nil
}

func (s Student) State() stateswitch.State {
	return stateswitch.State(s.Status)
}
```

Second step is to add the logic that you will run on each transition, and the conditions that validate witch of the transitions need to run

```go
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
```

As you can see there is no mention of the states here.

The next step is to create state machine that will receive your student

```go
student := Student{
	ID:     "123",
	Status: StatePending,
}

sm := stateswitch.NewStateMachine(&student)
``` 

Now add transitions to your state machine. 

In this case we have two transitions, both will execute `SetGrade` but according to the condition state machine will set the next step. 
For example if `IsPassed` will return true the next step will be `StatePassed`, 
if this condition will fail then state machine will try the next condition, and if it will pass then the state will change to `StateFailed`.
If none of the conditions will pass then state machone will return an error because there is not valid transitions that support this transition type.

```go
sm.AddTransition(stateswitch.TransitionRule{
	SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
	Condition:        student.IsPassed,
	Transition:       student.SetGrade,
	TransitionType:   TransitionTypeSetGrade,
    DestinationState: StatePassed,
})

sm.AddTransition(stateswitch.TransitionRule{
	TransitionType:   TransitionTypeSetGrade,
	SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
	DestinationState: StateFailed,
	Condition:        student.IsFailed,
	Transition:       student.SetGrade,
})
```

Our state machine is ready to run some logic!
The only thing that state machine is need right now is to get a transition type and arguments, it will automatically select the logic that will be executed.

```go
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
```

example of the output:
```
INFO[0000] {ID:123 Grade:0 Status:Pending}              
INFO[0000] {ID:123 Grade:90 Status:Passed}              
INFO[0000] {ID:123 Grade:50 Status:Failed}              
INFO[0000] {ID:123 Grade:80 Status:Passed}              
ERRO[0000] no match for transition type unknown transition 
INFO[0000] {ID:123 Grade:80 Status:Passed}              
ERRO[0000] invalid arguments for IsPassed condition     
INFO[0000] {ID:123 Grade:80 Status:Passed}   
```

## Examples

Example can be found [here](https://github.com/filanov/stateswitch/tree/master/examples)