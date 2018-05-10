package fs

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

func ParseByteRanges(byteRangeSpecifier string, totalSize int64) []*byteRange {
	rangePattern, _ := regexp.Compile("^bytes=(\\d+)[-](\\d+)$")
	if matches := rangePattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		lowIndex, _ := strconv.Atoi(matches[1])
		highIndex, _ := strconv.Atoi(matches[2])
		return []*byteRange{&byteRange{
			firstByteIndex: lowIndex,
			lastByteIndex:  highIndex,
			totalSize:      totalSize,
		}}
	}

	return nil
}

type byteRange struct {
	firstByteIndex int
	lastByteIndex  int
	totalSize      int64
}

func (byteRange *byteRange) Copy(fromFilename string, toWriter io.Writer) {
	file, _ := os.Open(fromFilename)
	defer file.Close()

	offset := int64(byteRange.firstByteIndex)
	copyLength := int64(byteRange.Length())
	file.Seek(offset, 0)
	io.CopyN(toWriter, file, copyLength)
}

func (byteRange *byteRange) ContentRange() string {
	return fmt.Sprintf("bytes %d-%d/%d",
		byteRange.firstByteIndex,
		byteRange.lastByteIndex,
		byteRange.totalSize,
	)
}

func (byteRange *byteRange) Length() int {
	return byteRange.lastByteIndex - byteRange.firstByteIndex + 1
}
