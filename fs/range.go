package fs

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

func ParseByteRangeSlices(byteRangeSpecifier string, filename string) []FileSlice {
	rangePattern, _ := regexp.Compile("^bytes=(\\d+)[-](\\d+)$")
	if matches := rangePattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		lowIndex, _ := strconv.Atoi(matches[1])
		highIndex, _ := strconv.Atoi(matches[2])
		return []FileSlice{&PartialSlice{
			path:           filename,
			firstByteIndex: lowIndex,
			lastByteIndex:  highIndex,
		}}
	}

	return nil
}

type PartialSlice struct {
	path           string
	firstByteIndex int
	lastByteIndex  int
}

func (slice *PartialSlice) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, success.PartialContentStatus)
}

func (slice *PartialSlice) WriteContentSizeHeaders(writer io.Writer) {
	msg.WriteHeader(writer, "Content-Length", strconv.Itoa(slice.len()))
	msg.WriteHeader(writer, "Content-Range", slice.contentRange())
}

func (slice *PartialSlice) WriteBody(writer io.Writer) {
	file, _ := os.Open(slice.path)
	defer file.Close()

	offset := int64(slice.firstByteIndex)
	copyLength := int64(slice.len())
	file.Seek(offset, 0)
	io.CopyN(writer, file, copyLength)
}

func (slice *PartialSlice) contentRange() string {
	info, _ := os.Stat(slice.path)
	totalSize := info.Size()
	return fmt.Sprintf("bytes %d-%d/%d",
		slice.firstByteIndex,
		slice.lastByteIndex,
		totalSize,
	)
}

func (slice *PartialSlice) len() int {
	return slice.lastByteIndex - slice.firstByteIndex + 1
}

type WholeFile struct {
	path string
}

func (slice *WholeFile) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, success.OKStatus)
}

func (slice *WholeFile) WriteContentSizeHeaders(writer io.Writer) {
	info, _ := os.Stat(slice.path)
	msg.WriteHeader(writer, "Content-Length", strconv.FormatInt(info.Size(), 10))
}

func (slice *WholeFile) WriteBody(writer io.Writer) {
	file, _ := os.Open(slice.path)
	defer file.Close()
	msg.CopyToBody(writer, file)
}

type FileSlice interface {
	WriteStatus(writer io.Writer)
	WriteContentSizeHeaders(writer io.Writer)
	WriteBody(writer io.Writer)
}
