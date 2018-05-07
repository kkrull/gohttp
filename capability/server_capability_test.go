package capability_test

import (
	"bytes"

	"github.com/kkrull/gohttp/capability"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StaticCapabilityServer", func() {
	var (
		controller     capability.ServerResource
		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	Describe("#Options", func() {
		BeforeEach(func() {
			responseBuffer = &bytes.Buffer{}
			controller = &capability.StaticCapabilityServer{
				AvailableMethods: []string{http.CONNECT, http.TRACE},
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
				controller = &capability.StaticCapabilityServer{
					AvailableMethods: []string{http.OPTIONS},
				}
				controller.Options(responseBuffer)
				response = httptest.ParseResponse(responseBuffer)
			})

			It("sets Allow to that one method", func() {
				response.HeaderShould("Allow", Equal(http.OPTIONS))
			})
		})

		Context("given 2 or more available methods", func() {
			BeforeEach(func() {
				responseBuffer = &bytes.Buffer{}
				controller = &capability.StaticCapabilityServer{
					AvailableMethods: []string{http.CONNECT, http.TRACE},
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
