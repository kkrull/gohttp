package opt_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/opt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StaticCapabilityController", func() {
	var (
		controller     opt.ServerCapabilityController
		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	Describe("#Options", func() {
		BeforeEach(func() {
			responseBuffer = &bytes.Buffer{}
			controller = &opt.StaticCapabilityController{
				AvailableMethods: []string{"CONNECT", "TRACE"},
			}
			controller.Options(responseBuffer)
			response = httptest.ParseResponse(responseBuffer)
		})

		It("writes a well-formed response", func() {
			response.ShouldBeWellFormed()
		})
		It("sets status to 200 OK", func() {
			response.StatusShouldBe(200, "OK")
		})
		It("sets Content-Length to 0", func() {
			response.HeaderShould("Content-Length", Equal("0"))
		})

		Context("given 1 available method", func() {
			BeforeEach(func() {
				responseBuffer = &bytes.Buffer{}
				controller = &opt.StaticCapabilityController{
					AvailableMethods: []string{"OPTIONS"},
				}
				controller.Options(responseBuffer)
				response = httptest.ParseResponse(responseBuffer)
			})

			It("sets Allow to that one method", func() {
				response.HeaderShould("Allow", Equal("OPTIONS"))
			})
		})

		Context("given 2 or more available methods", func() {
			BeforeEach(func() {
				responseBuffer = &bytes.Buffer{}
				controller = &opt.StaticCapabilityController{
					AvailableMethods: []string{"CONNECT", "TRACE"},
				}
				controller.Options(responseBuffer)
				response = httptest.ParseResponse(responseBuffer)
			})

			It("sets Allow to a comma-separated list of the given methods", func() {
				response.HeaderShould("Allow", Equal("CONNECT,TRACE"))
			})
		})
	})
})
