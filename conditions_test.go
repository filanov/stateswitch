package stateswitch

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

type testStateSwitch struct{}

var testError = fmt.Errorf("Blah")

func (t *testStateSwitch) State() State {
	return "blah"
}

func (t *testStateSwitch) SetState(state State) error {
	return nil
}

func True(sw StateSwitch, args TransitionArgs) (bool, error) {
	return true, nil
}

func False(sw StateSwitch, args TransitionArgs) (bool, error) {
	return false, nil
}

func TrueWithError(sw StateSwitch, args TransitionArgs) (bool, error) {
	return true, testError
}

func FalseWithError(sw StateSwitch, args TransitionArgs) (bool, error) {
	return false, testError
}

type resultValue struct {
	b bool
	e error
}

func result(b bool, e error) resultValue {
	return resultValue{
		b: b,
		e: e,
	}
}

func run(c Condition) resultValue {
	b, e := c(&testStateSwitch{}, nil)
	return result(b, e)
}

var _ = Describe("conditions_test", func() {
	It("Basic", func() {
		gomega.Expect(run(True)).To(gomega.Equal(result(true, nil)))
		gomega.Expect(run(False)).To(gomega.Equal(result(false, nil)))
		gomega.Expect(run(TrueWithError)).To(gomega.Equal(result(true, testError)))
		gomega.Expect(run(FalseWithError)).To(gomega.Equal(result(false, testError)))
	})
	It("Not", func() {
		gomega.Expect(run(Not(True))).To(gomega.Equal(result(false, nil)))
		gomega.Expect(run(Not(False))).To(gomega.Equal(result(true, nil)))
		gomega.Expect(run(Not(TrueWithError))).To(gomega.Equal(result(false, testError)))
		gomega.Expect(run(Not(FalseWithError))).To(gomega.Equal(result(false, testError)))
	})
	It("And", func() {
		gomega.Expect(run(And(True, True))).To(gomega.Equal(result(true, nil)))
		gomega.Expect(run(And(True, False))).To(gomega.Equal(result(false, nil)))
		gomega.Expect(run(And(True, TrueWithError))).To(gomega.Equal(result(false, testError)))
		gomega.Expect(run(And(True, FalseWithError))).To(gomega.Equal(result(false, testError)))
		gomega.Expect(run(And(False, True))).To(gomega.Equal(result(false, nil)))
		gomega.Expect(run(And(False, False))).To(gomega.Equal(result(false, nil)))
		gomega.Expect(run(And(False, TrueWithError))).To(gomega.Equal(result(false, nil)))
		gomega.Expect(run(And(False, FalseWithError))).To(gomega.Equal(result(false, nil)))
	})
	It("Or", func() {
		gomega.Expect(run(Or(True, True))).To(gomega.Equal(result(true, nil)))
		gomega.Expect(run(Or(True, False))).To(gomega.Equal(result(true, nil)))
		gomega.Expect(run(Or(True, TrueWithError))).To(gomega.Equal(result(true, nil)))
		gomega.Expect(run(Or(True, FalseWithError))).To(gomega.Equal(result(true, nil)))
		gomega.Expect(run(Or(TrueWithError, False, TrueWithError))).To(gomega.Equal(result(false, testError)))
		gomega.Expect(run(Or(FalseWithError, True, FalseWithError))).To(gomega.Equal(result(false, testError)))
		gomega.Expect(run(Or(False, True))).To(gomega.Equal(result(true, nil)))
		gomega.Expect(run(Or(False, False))).To(gomega.Equal(result(false, nil)))
		gomega.Expect(run(Or(False, TrueWithError))).To(gomega.Equal(result(false, testError)))
		gomega.Expect(run(Or(False, FalseWithError))).To(gomega.Equal(result(false, testError)))
	})
})
