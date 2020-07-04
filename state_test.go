package stateswitch

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = Describe("state_test", func() {
	It("contain", func() {
		states := States{"a", "b", "c"}
		gomega.Expect(states.Contain("a")).To(gomega.BeTrue())
		gomega.Expect(states.Contain("b")).To(gomega.BeTrue())
		gomega.Expect(states.Contain("c")).To(gomega.BeTrue())
	})
	It("not_contain", func() {
		states := States{"a", "b", "c"}
		gomega.Expect(states.Contain("d")).To(gomega.BeFalse())
	})
	It("empty", func() {
		states := States{}
		gomega.Expect(states.Contain("a")).To(gomega.BeFalse())
	})
})
