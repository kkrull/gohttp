package fs_test

import (
	"bytes"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
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
			requested := http.NewGetMessage("/foo")
			routedRequest := route.Route(requested)
			routedRequest.Handle(response)
			resource.GetShouldHaveReceived("/foo")
		})

		It("routes HEAD requests to HeadRequest", func() {
			requested := http.NewHeadMessage("/foo")
			routedRequest := route.Route(requested)
			routedRequest.Handle(response)
			resource.HeadShouldHaveReceived("/foo")
		})

		It("routes any other method to MethodNotAllowed", func() {
			requested := http.NewTraceMessage("/")
			routedRequest := route.Route(requested)
			Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed(http.GET, http.HEAD, http.OPTIONS)))
		})
	})
})
