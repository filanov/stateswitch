package stateswitch

import (
	"fmt"

	"github.com/pkg/errors"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = Describe("test_transition_find", func() {
	transitions := TransitionRules{
		{TransitionType: "a"},
		{TransitionType: "b"},
		{TransitionType: "a"},
		{TransitionType: "c"},
	}

	tests := []struct {
		name            string
		transitionType  TransitionType
		expectedResults int
	}{
		{transitionType: "a", expectedResults: 2},
		{transitionType: "b", expectedResults: 1},
		{transitionType: "c", expectedResults: 1},
		{transitionType: "d", expectedResults: 0},
	}

	for i := range tests {
		t := tests[i]
		It(fmt.Sprintf("find %s expected %d", t.transitionType, t.expectedResults), func() {
			found := transitions.Find(t.transitionType)
			gomega.Expect(len(found)).Should(gomega.Equal(t.expectedResults))
		})
	}
})

var _ = Describe("IsAllowedToRun", func() {
	var srcStateA State = "srcStateA"
	var srcStateB State = "srcStateB"
	transition := TransitionRule{
		SourceStates: []State{srcStateA, srcStateB},
		Condition:    nil,
	}

	tests := []struct {
		name      string
		state     State
		condition Condition
		allow     bool
		fail      bool
	}{
		{
			name:      "no condition",
			state:     srcStateB,
			condition: nil,
			allow:     true,
			fail:      false,
		},
		{
			name:      "invalid source state",
			state:     "some invalid source state",
			condition: nil,
			allow:     false,
			fail:      false,
		},
		{
			name:      "condition allow",
			state:     srcStateA,
			condition: func(stateSwitch StateSwitch, args TransitionArgs) (bool, error) { return true, nil },
			allow:     true,
			fail:      false,
		},
		{
			name:      "condition don't allow",
			state:     srcStateA,
			condition: func(stateSwitch StateSwitch, args TransitionArgs) (bool, error) { return false, nil },
			allow:     false,
			fail:      false,
		},
		{
			name:  "condition error",
			state: srcStateA,
			condition: func(stateSwitch StateSwitch, args TransitionArgs) (bool, error) {
				return false, errors.Errorf("error")
			},
			allow: false,
			fail:  true,
		},
	}

	for i := range tests {
		t := tests[i]
		It(t.name, func() {
			transition.Condition = t.condition
			allow, err := transition.IsAllowedToRun(&swState{state: t.state}, nil)
			gomega.Expect(allow).To(gomega.Equal(t.allow))
			gomega.Expect(err == nil).Should(gomega.Equal(!t.fail))
		})
	}
})
