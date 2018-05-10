package fs

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

type ExistingFile struct {
	Filename string
}

func (existingFile *ExistingFile) Name() string {
	return "Existing file"
}

func (existingFile *ExistingFile) Get(client io.Writer, message http.RequestMessage) {
	existingFile.Head(client, message)

	partialRanges := message.HeaderValues("Range")
	if len(partialRanges) == 0 {
		existingFile.writeWholeFile(client)
	} else {
		info, _ := os.Stat(existingFile.Filename)
		contentRange := parseByteRange(partialRanges[0], info.Size())
		contentRange.Copy(existingFile.Filename, client)
	}
}

func (existingFile *ExistingFile) Head(client io.Writer, message http.RequestMessage) {
	partialRanges := message.HeaderValues("Range")
	if len(partialRanges) == 0 {
		msg.WriteStatus(client, success.OKStatus)
		msg.WriteContentTypeHeader(client, contentTypeFromFileExtension(existingFile.Filename))
		info, _ := os.Stat(existingFile.Filename)
		msg.WriteHeader(client, "Content-Length", strconv.FormatInt(info.Size(), 10))
		msg.WriteEndOfMessageHeader(client)
		return
	}

	info, _ := os.Stat(existingFile.Filename)
	msg.WriteStatus(client, success.PartialContentStatus)
	msg.WriteContentTypeHeader(client, contentTypeFromFileExtension(existingFile.Filename))

	contentRange := parseByteRange(partialRanges[0], info.Size())
	msg.WriteHeader(client, "Content-Length", strconv.Itoa(contentRange.Length()))
	msg.WriteHeader(client, "Content-Range", contentRange.ContentRange())
	msg.WriteEndOfMessageHeader(client)
}

func contentTypeFromFileExtension(filename string) string {
	extension := path.Ext(filename)
	if extension == "" {
		return "text/plain"
	}

	return mime.TypeByExtension(extension)
}

func (existingFile *ExistingFile) writeWholeFile(client io.Writer) {
	file, _ := os.Open(existingFile.Filename)
	msg.CopyToBody(client, file)
}

/* byteRange */

func parseByteRange(byteRangeSpecifier string, totalSize int64) *byteRange {
	rangePattern, _ := regexp.Compile("^bytes=(\\d+)[-](\\d+)$")
	if matches := rangePattern.FindStringSubmatch(byteRangeSpecifier); matches != nil {
		lowIndex, _ := strconv.Atoi(matches[1])
		highIndex, _ := strconv.Atoi(matches[2])
		return &byteRange{
			firstByteIndex: lowIndex,
			lastByteIndex:  highIndex,
			totalSize:      totalSize,
		}
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
