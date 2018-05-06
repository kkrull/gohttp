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

	It("returns an error when % is followed by no characters", func() {
		_, err := http.PercentDecode("abc%")
		Expect(err).To(BeEquivalentTo(http.UnfinishedPercentEncoding{EnclosingField: "abc%"}))
	})

	It("returns an error when % is followed by 1 character", func() {
		_, err := http.PercentDecode("abc%1")
		Expect(err).To(BeEquivalentTo(http.UnfinishedPercentEncoding{EnclosingField: "abc%1"}))
	})
})

var _ = Describe("UnfinishedPercentEncoding", func() {
	It("describes the error and the offending text", func() {
		err := &http.UnfinishedPercentEncoding{EnclosingField: "abc%"}
		Expect(err).To(MatchError("% followed by fewer than 2 characters: abc%"))
	})
})
