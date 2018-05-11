package fs_test

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("ParseByteRangeSlices", func() {
	XIt("x-y is a regular range from x to y")
	XIt("x- is a suffix range from x to the end")
	XIt("-y is a suffix range of the last y bytes")

	XIt("a range that starts in the file and runs off the end stops at the end")
	XIt("a range that is completely outside the file is a 416")

	XIt("2 or more ranges is a 416 -- at least for this implementation")
})
