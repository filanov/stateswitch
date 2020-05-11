package stateswitch

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("state_test", func() {
	It("contain", func() {
		states := States{"a", "b", "c"}
		Expect(states.Contain("a")).To(BeTrue())
		Expect(states.Contain("b")).To(BeTrue())
		Expect(states.Contain("c")).To(BeTrue())
	})
	It("not_contain", func() {
		states := States{"a", "b", "c"}
		Expect(states.Contain("d")).To(BeFalse())
	})
	It("empty", func() {
		states := States{}
		Expect(states.Contain("a")).To(BeFalse())
	})
})
