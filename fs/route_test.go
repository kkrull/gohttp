package fs_test

import (
	"bytes"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewRoute", func() {
	It("returns a route to files and directories on the local file system", func() {
		route := fs.NewRoute("/public")
		Expect(route).To(BeEquivalentTo(
			&fs.FileSystemRoute{
				ContentRootPath: "/public",
				Resource:        &fs.ReadOnlyFileSystem{BaseDirectory: "/public"},
			}))
	})
})

var _ = Describe("FileSystemRoute", func() {
	var (
		route    *fs.FileSystemRoute
		resource = &FileSystemResourceMock{}
		response = &bytes.Buffer{}
	)

	BeforeEach(func() {
		route = &fs.FileSystemRoute{
			ContentRootPath: "/public",
			Resource:        resource,
		}
		response.Reset()
	})

	Describe("#Route", func() {
		It("routes GET requests to GetRequest", func() {
			requested := &http.RequestLine{Method: "GET", Target: "/foo"}
			routedRequest := route.Route(requested)
			routedRequest.Handle(response)
			resource.GetShouldHaveReceived("/foo")
		})

		It("routes HEAD requests to HeadRequest", func() {
			requested := &http.RequestLine{Method: "HEAD", Target: "/foo"}
			routedRequest := route.Route(requested)
			routedRequest.Handle(response)
			resource.HeadShouldHaveReceived("/foo")
		})

		Context("given any other method", func() {
			BeforeEach(func() {
				requested := &http.RequestLine{Method: "TRACE", Target: "/"}
				routedRequest := route.Route(requested)
				routedRequest.Handle(response)
			})

			It("responds 405 Method Not Allowed", httptest.ShouldHaveNoBody(response, 405, "Method Not Allowed"))
			It("sets Allow to GET and HEAD", httptest.ShouldAllowMethods(response, "GET", "HEAD"))
		})
	})
})
