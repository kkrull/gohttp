package httptest

import (
	"bytes"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func ShouldAllowMethods(response *bytes.Buffer, methods ...string) func() {
	return func() {
		responseMessage := ParseResponse(response)
		responseMessage.HeaderShould("Allow", containSubstrings(methods))
	}
}

func containSubstrings(values []string) types.GomegaMatcher {
	valueMatchers := make([]types.GomegaMatcher, len(values))
	for i, value := range values {
		valueMatchers[i] = ContainSubstring(value)
	}

	return SatisfyAll(valueMatchers...)
}

func ShouldHaveNoBody(response *bytes.Buffer, status int, reason string) func() {
	return func() {
		responseMessage := ParseResponse(response)
		responseMessage.ShouldBeWellFormed()
		responseMessage.StatusShouldBe(status, reason)
		responseMessage.HeaderShould("Content-Length", Equal("0"))
		responseMessage.BodyShould(BeEmpty())
	}
}
