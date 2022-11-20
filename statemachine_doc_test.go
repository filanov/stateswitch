package stateswitch

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = Describe("State machine documentation", func() {
	Context("A documented state machine", func() {
		sm := NewStateMachine()

		transitionType := TransitionType("TransitionType")
		srcStates := States{"StateA", "StateB"}
		dstState := State("StateC")

		sm.AddTransitionRule(TransitionRule{
			TransitionType:   transitionType,
			SourceStates:     []State{""},
			DestinationState: dstState,
			Condition:        nil,
			Transition:       nil,
			Documentation: TransitionRuleDoc{
				Name:        "Chocobomb Initial",
				Description: "This is documentation!",
			},
		})

		sm.AddTransitionRule(TransitionRule{
			TransitionType:   transitionType,
			SourceStates:     srcStates,
			DestinationState: dstState,
			Condition:        nil,
			Transition:       nil,
			Documentation: TransitionRuleDoc{
				Name:        "Chocobomb",
				Description: "This is documentation?",
			},
		})

		sm.DescribeState("StateA", StateDoc{
			Name:        "State A",
			Description: "State A Documentation",
		})

		sm.DescribeState("StateB", StateDoc{
			Name:        "State B",
			Description: "State B Documentation",
		})

		sm.DescribeState("StateC", StateDoc{
			Name:        "State C",
			Description: "State C Documentation",
		})

		sm.DescribeTransitionType(transitionType, TransitionTypeDoc{
			Name:        "Transition Type",
			Description: "Transition Type Documentation",
		})

		It("should be documented correctly when AsJSON is called", func() {
			docs, err := sm.AsJSON()
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(docs).Should(gomega.MatchJSON(`{
  "transition_rules": [
    {
      "transition_type": "TransitionType",
      "source_states": [
        "initial"
      ],
      "destination_state": "StateC",
      "name": "Chocobomb Initial",
      "description": "This is documentation!"
    },
    {
      "transition_type": "TransitionType",
      "source_states": [
        "StateA",
        "StateB"
      ],
      "destination_state": "StateC",
      "name": "Chocobomb",
      "description": "This is documentation?"
    }
  ],
  "states": {
    "StateA": {
      "name": "State A",
      "description": "State A Documentation"
    },
    "StateB": {
      "name": "State B",
      "description": "State B Documentation"
    },
    "StateC": {
      "name": "State C",
      "description": "State C Documentation"
    },
	"initial": {
	  "name": "Initial",
	  "description": "The initial state of the state machine. This is a synthetic state that is not actually part of the state machine. It appears in documentation when transition rules hold a single source state that is an empty string"
	}
  },
  "transition_types": {
    "TransitionType": {
      "name": "Transition Type",
      "description": "Transition Type Documentation"
    }
  }
}`))

		})
	})
})
