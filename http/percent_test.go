package http_test

import (
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PercentDecode", func() {
	It("returns a string unchanged that has no percent triplets in it", func() {
		Expect(http.PercentDecode("abcd")).To(Equal("abcd"))
	})

	It("decodes a % triplet into the ASCII character for its hexadecimal code", func() {
		Expect(http.PercentDecode("%3C")).To(Equal("<"))
	})

	It("retains characters after a percent triplet", func() {
		Expect(http.PercentDecode("%3Cabc")).To(Equal("<abc"))
	})

	It("decodes multiple percent triplets", func() {
		Expect(http.PercentDecode("%3Ca%3Cb")).To(Equal("<a<b"))
	})
})
