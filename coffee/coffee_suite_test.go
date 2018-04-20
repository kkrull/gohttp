package coffee_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCoffee(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "coffee")
}

type ControllerMock struct {

}
