[![Actions Status](https://github.com/filanov/stateswitch/workflows/make_all/badge.svg)](https://github.com/filanov/stateswitch/actions)

# stateswitch

## Overview

TDB

## Usage

First the entity that state machine will work on need to implement `StateSwitch` interface

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
// Define arguments for each transtion
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
```

As you can see there is no mention of the states here.

It's just a suggestion by you should crate a wrapper for your state machine, so you can work with the API and not with the state machine directly
This way you can have DB and other helpers included in your implementation.

```go
// State machine wrapper
type studentMachine struct {
	sm stateswitch.StateMachine
}
``` 

Now add transitions to your state machine. 

In this case we have two transitions, both will execute `SetGrade` but according to the condition state machine will set the next step. 
For example if `IsPassed` will return true the next step will be `StatePassed`, 
if this condition will fail then state machine will try the next condition, and if it will pass then the state will change to `StateFailed`.
If none of the conditions will pass then state machine will return an error because there is not valid transitions that support this transition type.

```go
func NewStudentMachine() *studentMachine {
	sm := stateswitch.NewStateMachine()

	sm.AddTransition(stateswitch.TransitionRule{
		SourceStates:     []stateswitch.State{StatePending, StateFailed, StatePassed},
		Condition:        IsPassed,
		Transition:       SetGradeTransition,
		TransitionType:   TransitionTypeSetGrade,
		DestinationState: StatePassed,
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
```

Our state machine is ready to run some logic!
The only thing that state machine is need right now is to get a transition type and arguments, it will automatically select the logic that will be executed.

```go
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
```

example of the output:
```
INFO[0000] {ID:123 Grade:0 Status:Pending}              
INFO[0000] {ID:123 Grade:90 Status:Passed}              
INFO[0000] {ID:123 Grade:50 Status:Failed}              
INFO[0000] {ID:123 Grade:80 Status:Passed}
```

## Examples

Example can be found [here](https://github.com/filanov/stateswitch/tree/master/examples)