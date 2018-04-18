package fs

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strconv"

	"github.com/kkrull/gohttp/msg"
)

type DirectoryListing struct {
	Files      []string
	HrefPrefix string
}

func (listing DirectoryListing) WriteTo(client io.Writer) error {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteHeader(client, "Content-Type", "text/html")

	message := listing.messageListingFiles()
	msg.WriteHeader(client, "Content-Length", strconv.Itoa(message.Len()))
	msg.WriteEndOfMessageHeader(client)

	msg.WriteBody(client, message.String())
	return nil
}

func (listing DirectoryListing) WriteHeader(client io.Writer) error {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteHeader(client, "Content-Type", "text/html")

	message := listing.messageListingFiles()
	msg.WriteHeader(client, "Content-Length", strconv.Itoa(message.Len()))
	msg.WriteEndOfMessageHeader(client)
	return nil
}

func (listing DirectoryListing) messageListingFiles() *bytes.Buffer {
	message := &bytes.Buffer{}
	message.WriteString("<html>\n")
	listing.writeHead(message)
	listing.writeBody(message)
	message.WriteString("</html>")
	return message
}

func (listing DirectoryListing) writeHead(message *bytes.Buffer) {
	message.WriteString("<head>\n")
	message.WriteString("<title>gohttp</title>\n")
	message.WriteString("</head>\n")
}

func (listing DirectoryListing) writeBody(message *bytes.Buffer) {
	message.WriteString("<body>\n")
	listing.writeFileListing(message)
	message.WriteString("</body>\n")
}

func (listing DirectoryListing) writeFileListing(message *bytes.Buffer) {
	message.WriteString("<ul>\n")
	for _, file := range listing.Files {
		message.WriteString(makeListItem(listing.makeLink(file)))
		message.WriteString("\n")
	}

	message.WriteString("</ul>\n")
}

func makeListItem(text string) string {
	return fmt.Sprintf("<li>%s</li>", text)
}

func (listing DirectoryListing) makeLink(filename string) string {
	href := path.Join(listing.HrefPrefix, filename)
	return fmt.Sprintf("<a href=\"%s\">%s</a>", href, filename)
}
