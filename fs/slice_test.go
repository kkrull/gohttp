package fs_test

import (
	"os"
	"path"

	"github.com/kkrull/gohttp/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseByteRange", func() {
	var (
		slice    fs.FileSlice
		basePath string
		file     string
	)

	BeforeEach(func() {
		basePath = makeEmptyTestDirectory("ParseByteRange", os.ModePerm)
		file = path.Join(basePath, "file.txt")
		Expect(createTextFile(file, "ABCD")).To(Succeed())
	})

	Context("given a file of size n bytes", func() {
		Context("given 'bytes=x-y'", func() {
			It("returns a PartialSlice from x to y (inclusive), when the range is non-decreasing and within the file", func() {
				slice = fs.ParseByteRange("bytes=0-1", file, "text/plain")
				Expect(slice).To(BeEquivalentTo(&fs.PartialSlice{
					ContentType:    "text/plain",
					Path:           file,
					FirstByteIndex: 0,
					LastByteIndex:  1,
				}))
			})

			It("returns an UnsupportedSlice, when the range starts within the file and goes past the end", func() {
				slice = fs.ParseByteRange("bytes=2-5", file, "text/plain")
				Expect(slice).To(BeEquivalentTo(&fs.UnsupportedSlice{
					Path:     file,
					NumBytes: 4,
				}))
			})

			It("returns an UnsupportedSlice, when the entire range is past the end of the file", func() {
				slice = fs.ParseByteRange("bytes=5-6", file, "text/plain")
				Expect(slice).To(BeEquivalentTo(&fs.UnsupportedSlice{
					Path:     file,
					NumBytes: 4,
				}))
			})
		})

		Context("given 'bytes=x-'", func() {
			It("returns a PartialSlice from x to the end of the file, when the first index is within the file", func() {
				slice = fs.ParseByteRange("bytes=2-", file, "text/plain")
				Expect(slice).To(BeEquivalentTo(&fs.PartialSlice{
					ContentType:    "text/plain",
					Path:           file,
					FirstByteIndex: 2,
					LastByteIndex:  3,
				}))
			})

			It("returns an UnsupportedSlice, when the first index is past the end of the file", func() {
				slice = fs.ParseByteRange("bytes=5-", file, "text/plain")
				Expect(slice).To(BeEquivalentTo(&fs.UnsupportedSlice{
					Path:     file,
					NumBytes: 4,
				}))
			})
		})

		Context("given 'bytes=-z'", func() {
			It("returns a PartialSlice of the last z bytes of the file, when z is not bigger than the size of the file", func() {
				slice = fs.ParseByteRange("bytes=-3", file, "text/plain")
				Expect(slice).To(BeEquivalentTo(&fs.PartialSlice{
					ContentType:    "text/plain",
					Path:           file,
					FirstByteIndex: 1,
					LastByteIndex:  3,
				}))
			})

			It("returns a PartialSlice covering the whole file, when z is at least as big as the size of the file", func() {
				slice = fs.ParseByteRange("bytes=-4", file, "text/plain")
				Expect(slice).To(BeEquivalentTo(&fs.PartialSlice{
					ContentType:    "text/plain",
					Path:           file,
					FirstByteIndex: 0,
					LastByteIndex:  3,
				}))
			})
		})

		Context("given 2 or more bytes ranges in the same bytes=... expression", func() {
			It("returns UnsupportedSlice, because that is way too complicated to handle right now", func() {
				slice = fs.ParseByteRange("bytes=0-1,2-3", file, "text/plain")
				Expect(slice).To(BeEquivalentTo(&fs.UnsupportedSlice{
					Path:     file,
					NumBytes: 4,
				}))
			})
		})
	})
})
