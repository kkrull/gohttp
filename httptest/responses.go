package httptest

import (
	"fmt"
	"strings"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func AllowedMethodsShouldBe(response fmt.Stringer, methods ...string) func() {
	return func() {
		responseMessage := ParseResponse(response)
		responseMessage.HeaderShould("Allow", Equal(strings.Join(methods, ",")))
	}
}

func ShouldAllowMethods(response fmt.Stringer, methods ...string) func() {
	return func() {
		responseMessage := ParseResponse(response)
		responseMessage.HeaderShould("Allow", containSubstrings(methods))
	}
}

func ShouldNotAllowMethod(response fmt.Stringer, method string) func() {
	return func() {
		responseMessage := ParseResponse(response)
		responseMessage.HeaderShouldNot("Allow", ContainSubstring(method))
	}
}

func containSubstrings(values []string) types.GomegaMatcher {
	valueMatchers := make([]types.GomegaMatcher, len(values))
	for i, value := range values {
		valueMatchers[i] = ContainSubstring(value)
	}

	return SatisfyAll(valueMatchers...)
}

func ShouldHaveNoBody(response fmt.Stringer, status int, reason string) func() {
	return func() {
		responseMessage := ParseResponse(response)
		responseMessage.ShouldBeWellFormed()
		responseMessage.StatusShouldBe(status, reason)
		responseMessage.HeaderShould("Content-Length", Equal("0"))
		responseMessage.BodyShould(BeEmpty())
	}
}
