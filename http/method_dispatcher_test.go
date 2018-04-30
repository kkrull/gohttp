package http_test

import (
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestMessage", func() {
	var message http.RequestMessage

	Context("given a target with no query or fragment", func() {
		BeforeEach(func() {
			message = http.NewGetMessage("/widget")
		})

		It("the path is the full target", func() {
			Expect(message.Path()).To(Equal("/widget"))
		})
		It("there are no query parameters", func() {
			Expect(message.QueryParameters()).To(BeEmpty())
		})
	})

	Context("given a target with a query", func() {
		It("the path is the part before the ?", func() {
			message = http.NewGetMessage("/widget?field=value")
			Expect(message.Path()).To(Equal("/widget"))
		})

		It("parses the part after the ? into query parameters", func() {
			message = http.NewGetMessage("/widget?field=value")
			Expect(message.QueryParameters()).To(
				ContainElement(http.QueryParameter{Name: "field", Value: "value"}))
		})

		It("parses parameters without a value into QueryParameter#Name", func() {
			message = http.NewGetMessage("/widget?flag")
			Expect(message.QueryParameters()).To(
				ContainElement(http.QueryParameter{Name: "flag", Value: ""}))
		})

		It("uses '=' to split a parameter's name and value", func() {
			message = http.NewGetMessage("/widget?field=value")
			Expect(message.QueryParameters()).To(
				ContainElement(http.QueryParameter{Name: "field", Value: "value"}))
		})

		It("uses '&' to split among multiple parameters", func() {
			message = http.NewGetMessage("/widget?one=1&two=2")
			Expect(message.QueryParameters()).To(Equal([]http.QueryParameter{
				{Name: "one", Value: "1"},
				{Name: "two", Value: "2"},
			}))
		})
	})

	Context("given a target with a fragment", func() {
		BeforeEach(func() {
			message = http.NewGetMessage("/widget#section")
		})

		It("the path is the part before the '#'", func() {
			Expect(message.Path()).To(Equal("/widget"))
		})
		It("there are no query parameters", func() {
			Expect(message.QueryParameters()).To(BeEmpty())
		})
	})

	Context("given a target with a query and a fragment", func() {
		BeforeEach(func() {
			message = http.NewGetMessage("/widget?field=value#section")
		})

		It("the path is the part before the ?", func() {
			Expect(message.Path()).To(Equal("/widget"))
		})
		It("query parameters are parsed from the part between the ? and the #", func() {
			Expect(message.QueryParameters()).To(Equal([]http.QueryParameter{
				{Name: "field", Value: "value"},
			}))
		})
	})
})
