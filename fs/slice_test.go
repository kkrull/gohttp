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
		Context("given 'bytes=x-y', where y < n and x <= y", func() {
			It("returns PartialSlice from x to y, inclusive", func() {
				slice = fs.ParseByteRangeSlice("bytes=0-1", file)
				Expect(slice).To(BeEquivalentTo(&fs.PartialSlice{
					Path:           file,
					FirstByteIndex: 0,
					LastByteIndex:  1,
				}))
			})
		})

		XIt("x- is a suffix range from x to the end")
		XIt("-y is a suffix range of the last y bytes")

		XIt("a range that starts in the file and runs off the end stops at the end")
		XIt("a range that is completely outside the file is a 416")

		XIt("2 or more ranges is a 416 -- at least for this implementation")
	})
})
