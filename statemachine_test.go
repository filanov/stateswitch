package stateswitch

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("AddTransition", func() {
	It("AddTransition", func() {
		sm := NewStateMachine(nil)

		transitionType := TransitionType("TransitionType")
		srcStates := States{"StateA", "StateB"}
		dstState := State("StateC")
		transitionFunc := func(args TransitionArgs) error {
			return nil
		}
		condition := func(args TransitionArgs) (bool, error) {
			return true, nil
		}

		sm.AddTransition(TransitionRule{
			TransitionType:   transitionType,
			SourceStates:     srcStates,
			DestinationState: dstState,
			Transition:       transitionFunc,
			Condition:        condition,
		})

		Expect(len(sm.transitionRules)).Should(Equal(1))
		rules, ok := sm.transitionRules[transitionType]
		Expect(ok).To(BeTrue())
		Expect(len(rules)).Should(Equal(1))
		Expect(rules[0].Condition).Should(BeAssignableToTypeOf(condition))
		Expect(rules[0].Transition).Should(BeAssignableToTypeOf(transitionFunc))
		Expect(rules[0].SourceStates).Should(BeAssignableToTypeOf(srcStates))
		Expect(rules[0].DestinationState).Should(BeAssignableToTypeOf(dstState))
	})
})

var _ = Describe("Run", func() {
	var sw swState
	var sm StateMachine

	BeforeEach(func() {
		sw = swState{state: stateA}
		sm = NewStateMachine(&sw)
		sm.AddTransition(TransitionRule{
			TransitionType:   ttAToB,
			SourceStates:     []State{stateA},
			DestinationState: stateB,
			Condition:        nil,
			Transition:       nil,
		})
		sm.AddTransition(TransitionRule{
			TransitionType:   ttNotPermittedAToC,
			SourceStates:     []State{stateA},
			DestinationState: stateC,
			Condition:        func(args TransitionArgs) (bool, error) { return false, nil },
			Transition:       nil,
		})
		sm.AddTransition(TransitionRule{
			TransitionType:   ttConditionError,
			SourceStates:     []State{stateA, stateB},
			DestinationState: stateC,
			Condition:        func(args TransitionArgs) (bool, error) { return false, errors.Errorf("error") },
			Transition:       nil,
		})
		sm.AddTransition(TransitionRule{
			TransitionType:   ttBToC,
			SourceStates:     []State{stateB},
			DestinationState: stateC,
			Condition:        func(args TransitionArgs) (bool, error) { return true, nil },
			Transition:       func(args TransitionArgs) error { return nil },
		})
		sm.AddTransition(TransitionRule{
			TransitionType:   ttBToA,
			SourceStates:     []State{stateB},
			DestinationState: stateA,
			Condition:        func(args TransitionArgs) (bool, error) { return true, nil },
			Transition:       func(args TransitionArgs) error { return errors.Errorf("error") },
		})
	})

	It("success", func() {
		Expect(sm.Run(ttAToB, nil)).ShouldNot(HaveOccurred())
		Expect(sw.state).Should(Equal(stateB))
	})
	It("transition type not found", func() {
		Expect(sm.Run("invalid transition type", nil)).Should(HaveOccurred())
		Expect(sw.state).Should(Equal(stateA))
	})
	It("transition not permitted", func() {
		Expect(sm.Run(ttNotPermittedAToC, nil)).Should(HaveOccurred())
		Expect(sw.state).Should(Equal(stateA))
	})
	It("condition error", func() {
		Expect(sm.Run(ttConditionError, nil)).Should(HaveOccurred())
		Expect(sw.state).Should(Equal(stateA))
	})
	It("run transition", func() {
		Expect(sw.SetState(stateB)).ShouldNot(HaveOccurred())
		Expect(sm.Run(ttBToC, nil)).ShouldNot(HaveOccurred())
		Expect(sw.state).Should(Equal(stateC))
	})
	It("transition error", func() {
		Expect(sw.SetState(stateB)).ShouldNot(HaveOccurred())
		Expect(sm.Run(ttBToA, nil)).Should(HaveOccurred())
		Expect(sw.state).Should(Equal(stateB))
	})
})

const (
	stateA State = "a"
	stateB State = "b"
	stateC State = "c"
)

const (
	ttAToB             TransitionType = "ttAToB"
	ttNotPermittedAToC TransitionType = "notPermittedAToC"
	ttConditionError   TransitionType = "ttConditionError"
	ttBToC             TransitionType = "ttBToC"
	ttBToA             TransitionType = "ttBToA"
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
