package stateswitch

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("AddTransitionRule", func() {
	It("AddTransitionRule", func() {
		sm := NewStateMachine()

		transitionType := TransitionType("TransitionType")
		srcStates := States{"StateA", "StateB"}
		dstState := State("StateC")
		transitionFunc := func(stateSwitch StateSwitch, args TransitionArgs) error {
			return nil
		}
		condition := func(stateSwitch StateSwitch, args TransitionArgs) (bool, error) {
			return true, nil
		}

		sm.AddTransitionRule(TransitionRule{
			TransitionType:   transitionType,
			SourceStates:     srcStates,
			DestinationState: dstState,
			Transition:       transitionFunc,
			Condition:        condition,
		})
		gomega.Expect(len(sm.transitionRules)).Should(gomega.Equal(1))
		rules, ok := sm.transitionRules[transitionType]
		gomega.Expect(ok).To(gomega.BeTrue())
		gomega.Expect(len(rules)).Should(gomega.Equal(1))
		gomega.Expect(rules[0].Condition).Should(gomega.BeAssignableToTypeOf(condition))
		gomega.Expect(rules[0].Transition).Should(gomega.BeAssignableToTypeOf(transitionFunc))
		gomega.Expect(rules[0].SourceStates).Should(gomega.BeAssignableToTypeOf(srcStates))
		gomega.Expect(rules[0].DestinationState).Should(gomega.BeAssignableToTypeOf(dstState))
	})
})

var _ = Describe("Run", func() {
	var sw *swState
	var swErr *swError
	var sm StateMachine

	BeforeEach(func() {
		sw = &swState{state: stateA}
		swErr = &swError{state: stateA}

		sm = NewStateMachine()
		sm.AddTransitionRule(TransitionRule{
			TransitionType:   ttAToB,
			SourceStates:     []State{stateA},
			DestinationState: stateB,
			Condition:        nil,
			Transition:       nil,
			PostTransition:   func(stateSwitch StateSwitch, args TransitionArgs) error { return nil },
		})
		sm.AddTransitionRule(TransitionRule{
			TransitionType:   ttNotPermittedAToC,
			SourceStates:     []State{stateA},
			DestinationState: stateC,
			Condition:        func(stateSwitch StateSwitch, args TransitionArgs) (bool, error) { return false, nil },
			Transition:       nil,
		})
		sm.AddTransitionRule(TransitionRule{
			TransitionType:   ttConditionError,
			SourceStates:     []State{stateA, stateB},
			DestinationState: stateC,
			Condition:        func(stateSwitch StateSwitch, args TransitionArgs) (bool, error) { return false, errors.Errorf("error") },
			Transition:       nil,
		})
		sm.AddTransitionRule(TransitionRule{
			TransitionType:   ttBToC,
			SourceStates:     []State{stateB},
			DestinationState: stateC,
			Condition:        func(stateSwitch StateSwitch, args TransitionArgs) (bool, error) { return true, nil },
			Transition:       func(stateSwitch StateSwitch, args TransitionArgs) error { return nil },
		})
		sm.AddTransitionRule(TransitionRule{
			TransitionType:   ttBToA,
			SourceStates:     []State{stateB},
			DestinationState: stateA,
			Condition:        func(stateSwitch StateSwitch, args TransitionArgs) (bool, error) { return true, nil },
			Transition:       func(stateSwitch StateSwitch, args TransitionArgs) error { return errors.Errorf("error") },
		})
		sm.AddTransitionRule(TransitionRule{
			TransitionType:   ttAToBPostTransitionErr,
			SourceStates:     []State{stateA},
			DestinationState: stateB,
			Condition:        nil,
			Transition:       nil,
			PostTransition: func(stateSwitch StateSwitch, args TransitionArgs) error {
				return errors.Errorf("post error")
			},
		})
	})

	It("success", func() {
		gomega.Expect(sm.Run(ttAToB, sw, nil)).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(sw.state).Should(gomega.Equal(stateB))
	})
	It("transition type not found", func() {
		gomega.Expect(errors.Is(sm.Run("invalid transition type", sw, nil), NoMatchForTransitionType)).
			Should(gomega.BeTrue())
		gomega.Expect(sw.state).Should(gomega.Equal(stateA))
	})
	It("transition not permitted", func() {
		gomega.Expect(errors.Is(sm.Run(ttNotPermittedAToC, sw, nil), NoConditionPassedToRunTransaction)).
			Should(gomega.BeTrue())
		gomega.Expect(sw.state).Should(gomega.Equal(stateA))
	})
	It("condition error", func() {
		gomega.Expect(sm.Run(ttConditionError, sw, nil)).Should(gomega.HaveOccurred())
		gomega.Expect(sw.state).Should(gomega.Equal(stateA))
	})
	It("run transition", func() {
		gomega.Expect(sw.SetState(stateB)).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(sm.Run(ttBToC, sw, nil)).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(sw.state).Should(gomega.Equal(stateC))
	})
	It("transition error", func() {
		gomega.Expect(sw.SetState(stateB)).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(sm.Run(ttBToA, sw, nil)).Should(gomega.HaveOccurred())
		gomega.Expect(sw.state).Should(gomega.Equal(stateB))
	})
	It("set state error", func() {
		gomega.Expect(sm.Run(ttAToB, swErr, nil)).Should(gomega.HaveOccurred())
		gomega.Expect(swErr.state).Should(gomega.Equal(stateA))
	})
	It("post transition error", func() {
		gomega.Expect(sm.Run(ttAToBPostTransitionErr, sw, nil)).Should(gomega.HaveOccurred())
		gomega.Expect(sw.state).Should(gomega.Equal(stateB))
	})
})

const (
	stateA State = "a"
	stateB State = "b"
	stateC State = "c"
)

const (
	ttAToB                  TransitionType = "ttAToB"
	ttNotPermittedAToC      TransitionType = "notPermittedAToC"
	ttConditionError        TransitionType = "ttConditionError"
	ttBToC                  TransitionType = "ttBToC"
	ttBToA                  TransitionType = "ttBToA"
	ttAToBPostTransitionErr TransitionType = "post transition error"
)

// implement simple state switch object for tests
type swState struct {
	state State
}

func (s *swState) State() State {
	return s.state
}

func (s *swState) SetState(state State) error {
	s.state = state
	return nil
}

// implement simple state switch object for tests
type swError struct {
	state State
}

func (s *swError) State() State {
	return s.state
}

func (s *swError) SetState(state State) error {
	return errors.Errorf("error")
}
