package fs_test

import (
	"bytes"
	"os"
	"path"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
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

	XIt("returns the new route implementation with a factory")
})

var _ = Describe("NewFileSystemRoute", func() {
	Describe("#Route", func() {
		var (
			route    *fs.NewFileSystemRoute
			factory  *ResourceFactoryMock
			resource *FileSystemResourceMock

			request        http.Request
			responseBuffer = &bytes.Buffer{}
		)

		BeforeEach(func() {
			factory = &ResourceFactoryMock{}
			resource = &FileSystemResourceMock{}
			responseBuffer.Reset()
			route = &fs.NewFileSystemRoute{
				ContentRootPath: makeEmptyTestDirectory("NewFileSystemRoute", os.ModePerm),
				Factory:         factory,
			}
		})

		Context("given a request with a method that is supported by the specified path", func() {
			It("routes to NotFound when the requested path does not exist in the content base directory", func() {
				factory.NonExistingResourceReturns(resource)
				request = route.Route(http.NewGetMessage("/missing.txt"))

				request.Handle(responseBuffer)
				factory.NonExistingResourceShouldHaveReceived("/missing.txt")
				resource.GetShouldHaveReceived("/missing.txt")
			})
			XIt("NotFound implements Get and Head")

			It("routes to ExistingFile when the requested path is a file inside the base directory", func() {
				existingFile := path.Join(route.ContentRootPath, "readable.txt")
				Expect(createTextFile(existingFile, "A")).To(Succeed())

				factory.ExistingFileResourceReturns(resource)
				request = route.Route(http.NewGetMessage("/readable.txt"))

				request.Handle(responseBuffer)
				factory.ExistingFileShouldHaveReceived("/readable.txt", existingFile)
				resource.GetShouldHaveReceived("/readable.txt")
			})
			XIt("ExistingFile implements Get and Head")

			It("routes to DirectoryListing when the requested path is a directory inside the base directory", func() {
				Expect(createTextFile(path.Join(route.ContentRootPath, "A"), "A")).To(Succeed())
				Expect(createTextFile(path.Join(route.ContentRootPath, "B"), "B")).To(Succeed())

				factory.DirectoryListingResourceReturns(resource)
				request = route.Route(http.NewGetMessage("/"))

				request.Handle(responseBuffer)
				factory.DirectoryListingShouldHaveReceived("/", []string{"A", "B"})
				resource.GetShouldHaveReceived("/")
			})
			XIt("DirectoryListing implements Get and Head")
		})
	})
})
