package stateswitch

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestStateSwitch(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "transition_test")
}
