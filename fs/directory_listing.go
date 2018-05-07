package fs

import (
	"bytes"
	"fmt"
	"io"
	"path"

	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

type DirectoryListing struct {
	Files      []string
	HrefPrefix string
	body       *bytes.Buffer
}

func (listing *DirectoryListing) WriteTo(client io.Writer) error {
	listing.WriteHeader(client)
	msg.WriteBody(client, listing.body.String())
	return nil
}

func (listing *DirectoryListing) WriteHeader(client io.Writer) error {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteContentTypeHeader(client, "text/html")

	listing.body = listing.messageListingFiles()
	msg.WriteContentLengthHeader(client, listing.body.Len())
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
