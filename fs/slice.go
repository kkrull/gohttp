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

func ParseByteRangeSlice(byteRangeSpecifier string, filename string) FileSlice {
	var (
		explicitRangePattern = regexp.MustCompile("^bytes=(\\d+)[-](\\d+)$")
		startingIndexPattern = regexp.MustCompile("^bytes=(\\d+)-$")
		suffixLengthPattern  = regexp.MustCompile("^bytes=-(\\d+)$")
	)

	size, _ := sizeInBytes(filename)
	if matches := explicitRangePattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		lowIndex, _ := strconv.ParseInt(matches[1], base10, bitsInInt64)
		highIndex, _ := strconv.ParseInt(matches[2], base10, bitsInInt64)
		return &PartialSlice{
			Path:           filename,
			FirstByteIndex: lowIndex,
			LastByteIndex:  min(size, highIndex),
		}
	} else if matches := startingIndexPattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		lowIndex, _ := strconv.ParseInt(matches[1], base10, bitsInInt64)
		return &PartialSlice{
			Path:           filename,
			FirstByteIndex: lowIndex,
			LastByteIndex:  size - 1,
		}
	} else if matches := suffixLengthPattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		length, _ := strconv.ParseInt(matches[1], base10, bitsInInt64)
		return &PartialSlice{
			Path:           filename,
			FirstByteIndex: size - length,
			LastByteIndex:  size - 1,
		}
	}

	return nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}

type PartialSlice struct {
	Path           string
	FirstByteIndex int64
	LastByteIndex  int64
}

func (slice *PartialSlice) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, success.PartialContentStatus)
}

func (slice *PartialSlice) WriteContentSizeHeaders(writer io.Writer) {
	msg.WriteHeader(writer, "Content-Length", strconv.FormatInt(slice.len(), base10))
	msg.WriteHeader(writer, "Content-Range", slice.contentRange())
}

func (slice *PartialSlice) WriteBody(writer io.Writer) {
	file, _ := os.Open(slice.Path)
	defer file.Close()

	file.Seek(slice.FirstByteIndex, 0)
	io.CopyN(writer, file, slice.len())
}

func (slice *PartialSlice) contentRange() string {
	size, _ := sizeInBytes(slice.Path)
	return fmt.Sprintf("bytes %d-%d/%d",
		slice.FirstByteIndex,
		slice.LastByteIndex,
		size,
	)
}

func (slice *PartialSlice) len() int64 {
	return slice.LastByteIndex - slice.FirstByteIndex + 1
}

type WholeFile struct {
	Path string
}

func (slice *WholeFile) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, success.OKStatus)
}

func (slice *WholeFile) WriteContentSizeHeaders(writer io.Writer) {
	size, _ := sizeInBytes(slice.Path)
	msg.WriteHeader(writer, "Content-Length", strconv.FormatInt(size, base10))
}

func (slice *WholeFile) WriteBody(writer io.Writer) {
	file, _ := os.Open(slice.Path)
	defer file.Close()
	msg.CopyToBody(writer, file)
}

type FileSlice interface {
	WriteStatus(writer io.Writer)
	WriteContentSizeHeaders(writer io.Writer)
	WriteBody(writer io.Writer)
}

func sizeInBytes(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}
