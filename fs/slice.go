package fs

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/success"
)

const (
	base10      = 10
	bitsInInt64 = 64
)

func ParseByteRangeSlice(byteRangeSpecifier string, filename string, contentType string) FileSlice {
	var (
		explicitRangePattern = regexp.MustCompile("^bytes=(\\d+)[-](\\d+)$")
		startingIndexPattern = regexp.MustCompile("^bytes=(\\d+)-$")
		suffixLengthPattern  = regexp.MustCompile("^bytes=-(\\d+)$")
	)

	size, _ := sizeInBytes(filename)
	if matches := explicitRangePattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		lowIndex, _ := strconv.ParseInt(matches[1], base10, bitsInInt64)
		highIndex, _ := strconv.ParseInt(matches[2], base10, bitsInInt64)
		if lowIndex >= size {
			return &UnsupportedSlice{
				Path:     filename,
				NumBytes: size,
			}
		} else if highIndex >= size {
			return &UnsupportedSlice{
				Path:     filename,
				NumBytes: size,
			}
		}

		return &PartialSlice{
			Path:           filename,
			ContentType:    contentType,
			FirstByteIndex: lowIndex,
			LastByteIndex:  min(size, highIndex),
		}
	} else if matches := startingIndexPattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		lowIndex, _ := strconv.ParseInt(matches[1], base10, bitsInInt64)
		if lowIndex >= size {
			return &UnsupportedSlice{
				Path:     filename,
				NumBytes: size,
			}
		}

		return &PartialSlice{
			Path:           filename,
			ContentType:    contentType,
			FirstByteIndex: lowIndex,
			LastByteIndex:  size - 1,
		}
	} else if matches := suffixLengthPattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		length, _ := strconv.ParseInt(matches[1], base10, bitsInInt64)
		if length >= size {
			return &WholeFile{
				Path:        filename,
				ContentType: contentType,
			}
		}

		return &PartialSlice{
			Path:           filename,
			ContentType:    contentType,
			FirstByteIndex: size - length,
			LastByteIndex:  size - 1,
		}
	}

	return &UnsupportedSlice{
		Path:     filename,
		NumBytes: size,
	}
}

// A slice of part of a file
type PartialSlice struct {
	Path           string
	ContentType    string
	FirstByteIndex int64
	LastByteIndex  int64
}

func (slice *PartialSlice) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, success.PartialContentStatus)
}

func (slice *PartialSlice) WriteContentHeaders(writer io.Writer) {
	msg.WriteHeader(writer, "Content-Length", strconv.FormatInt(slice.len(), base10))
	msg.WriteHeader(writer, "Content-Range", slice.contentRange())
	msg.WriteContentTypeHeader(writer, slice.ContentType)
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

// An attempt to slice that is not supported
type UnsupportedSlice struct {
	Path     string
	NumBytes int64
}

func (slice *UnsupportedSlice) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, clienterror.RangeNotSatisfiableStatus)
}

func (slice *UnsupportedSlice) WriteContentHeaders(writer io.Writer) {
	msg.WriteContentLengthHeader(writer, 0)
	msg.WriteHeader(writer, "Content-Range", fmt.Sprintf("bytes */%d", slice.NumBytes))
	//msg.WriteContentTypeHeader(writer, "text/plain")
}

func (slice *UnsupportedSlice) WriteBody(writer io.Writer) {
	//msg.WriteBody(writer, "Invalid range")
}

// A slice consisting of the entire file
type WholeFile struct {
	Path        string
	ContentType string
}

func (slice *WholeFile) WriteStatus(writer io.Writer) {
	msg.WriteStatus(writer, success.OKStatus)
}

func (slice *WholeFile) WriteContentHeaders(writer io.Writer) {
	size, _ := sizeInBytes(slice.Path)
	msg.WriteHeader(writer, "Content-Length", strconv.FormatInt(size, base10))
	msg.WriteContentTypeHeader(writer, slice.ContentType)
}

func (slice *WholeFile) WriteBody(writer io.Writer) {
	file, _ := os.Open(slice.Path)
	defer file.Close()
	msg.CopyToBody(writer, file)
}

// A view of all/part of a file
type FileSlice interface {
	WriteStatus(writer io.Writer)
	WriteContentHeaders(writer io.Writer)
	WriteBody(writer io.Writer)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}

func sizeInBytes(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}
