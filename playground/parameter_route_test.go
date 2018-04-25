package playground_test

import (
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewParameterRoute", func() {
	It("returns a ParameterRoute", func() {
		Expect(playground.NewParameterRoute()).NotTo(BeNil())
		Expect(playground.NewParameterRoute()).To(BeAssignableToTypeOf(&playground.ParameterRoute{}))
	})
})
