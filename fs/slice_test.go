package fs_test

import (
	"os"
	"path"

	"github.com/kkrull/gohttp/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseByteRangeSlice", func() {
	var (
		slice    fs.FileSlice
		basePath string
		file     string
	)

	BeforeEach(func() {
		basePath = makeEmptyTestDirectory("ParseByteRangeSlice", os.ModePerm)
		file = path.Join(basePath, "file.txt")
		Expect(createTextFile(file, "ABCD")).To(Succeed())
	})

	Context("given a file of size n bytes", func() {
		Context("given 'bytes=x-y', specifying non-decreasing indices within the file", func() {
			It("returns a PartialSlice from x to y, inclusive", func() {
				slice = fs.ParseByteRangeSlice("bytes=0-1", file)
				Expect(slice).To(BeEquivalentTo(&fs.PartialSlice{
					Path:           file,
					FirstByteIndex: 0,
					LastByteIndex:  1,
				}))
			})
		})

		Context("given 'bytes=x-', specifying a first index within the file", func() {
			It("returns a PartialSlice from x to the end of the file", func() {
				slice = fs.ParseByteRangeSlice("bytes=2-", file)
				Expect(slice).To(BeEquivalentTo(&fs.PartialSlice{
					Path:           file,
					FirstByteIndex: 2,
					LastByteIndex:  3,
				}))
			})
		})

		Context("given 'bytes=-z', specifying the last z bytes within the file", func() {
			It("returns a PartialSlice of the last z bytes of the file", func() {
				slice = fs.ParseByteRangeSlice("bytes=-3", file)
				Expect(slice).To(BeEquivalentTo(&fs.PartialSlice{
					Path:           file,
					FirstByteIndex: 1,
					LastByteIndex:  3,
				}))
			})
		})

		Context("given 'bytes=x-y', specifying a first index within the file and a last index beyond the end of the file", func() {
			It("returns a PartialSlice from x to the end of the file", func() {
				slice = fs.ParseByteRangeSlice("bytes=2-5", file)
				Expect(slice).To(BeEquivalentTo(&fs.PartialSlice{
					Path:           file,
					FirstByteIndex: 2,
					LastByteIndex:  4,
				}))
			})
		})

		XIt("a range that is completely outside the file is a 416")

		XIt("2 or more ranges is a 416 -- at least for this implementation")
	})
})
