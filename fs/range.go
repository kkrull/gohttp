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

func SingleByteRangeSlice(byteRangeSpecifier string, filename string) []FileSlice {
	info, _ := os.Stat(filename)
	file, _ := os.Open(filename)
	totalSize := info.Size()
	rangePattern, _ := regexp.Compile("^bytes=(\\d+)[-](\\d+)$")
	if matches := rangePattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		lowIndex, _ := strconv.Atoi(matches[1])
		highIndex, _ := strconv.Atoi(matches[2])
		return []FileSlice{&PartialSlice{
			file:            file,
			fileSizeInBytes: totalSize,
			firstByteIndex:  lowIndex,
			lastByteIndex:   highIndex,
		}}
	}

	return nil
}

type PartialSlice struct {
	file            *os.File
	fileSizeInBytes int64
	firstByteIndex  int
	lastByteIndex   int
}

func (slice *PartialSlice) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, success.PartialContentStatus)
}

func (slice *PartialSlice) WriteContentSizeHeaders(writer io.Writer) {
	msg.WriteHeader(writer, "Content-Length", strconv.Itoa(slice.Length()))
	msg.WriteHeader(writer, "Content-Range", slice.ContentRange())
}

func (slice *PartialSlice) WriteBody(writer io.Writer) {
	offset := int64(slice.firstByteIndex)
	copyLength := int64(slice.Length())
	slice.file.Seek(offset, 0)
	io.CopyN(writer, slice.file, copyLength)
}

func (slice *PartialSlice) ContentRange() string {
	return fmt.Sprintf("bytes %d-%d/%d",
		slice.firstByteIndex,
		slice.lastByteIndex,
		slice.fileSizeInBytes,
	)
}

func (slice *PartialSlice) Length() int {
	return slice.lastByteIndex - slice.firstByteIndex + 1
}

type WholeFile struct {
	file        *os.File
	sizeInBytes int64
}

func (slice *WholeFile) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, success.OKStatus)
}

func (slice *WholeFile) WriteContentSizeHeaders(writer io.Writer) {
	msg.WriteHeader(writer, "Content-Length", strconv.FormatInt(slice.sizeInBytes, 10))
}

func (slice *WholeFile) WriteBody(writer io.Writer) {
	msg.CopyToBody(writer, slice.file)
}

type FileSlice interface {
	WriteStatus(writer io.Writer)
	WriteContentSizeHeaders(writer io.Writer)
	WriteBody(writer io.Writer)
}
