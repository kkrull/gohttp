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

const (
	base10      = 10
	bitsInInt64 = 64
)

func ParseByteRangeSlices(byteRangeSpecifier string, filename string) []FileSlice {
	rangePattern, _ := regexp.Compile("^bytes=(\\d+)[-](\\d+)$")
	if matches := rangePattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		lowIndex, _ := strconv.ParseInt(matches[1], base10, bitsInInt64)
		highIndex, _ := strconv.ParseInt(matches[2], base10, bitsInInt64)
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
	firstByteIndex int64
	lastByteIndex  int64
}

func (slice *PartialSlice) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, success.PartialContentStatus)
}

func (slice *PartialSlice) WriteContentSizeHeaders(writer io.Writer) {
	msg.WriteHeader(writer, "Content-Length", strconv.FormatInt(slice.len(), base10))
	msg.WriteHeader(writer, "Content-Range", slice.contentRange())
}

func (slice *PartialSlice) WriteBody(writer io.Writer) {
	file, _ := os.Open(slice.path)
	defer file.Close()

	file.Seek(slice.firstByteIndex, 0)
	io.CopyN(writer, file, slice.len())
}

func (slice *PartialSlice) contentRange() string {
	info, _ := os.Stat(slice.path)
	return fmt.Sprintf("bytes %d-%d/%d",
		slice.firstByteIndex,
		slice.lastByteIndex,
		info.Size(),
	)
}

func (slice *PartialSlice) len() int64 {
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
	msg.WriteHeader(writer, "Content-Length", strconv.FormatInt(info.Size(), base10))
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
