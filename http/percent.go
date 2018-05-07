package http

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func PercentDecode(field string) (decoded string, malformed error) {
	outputBuffer := &bytes.Buffer{}
	splits := strings.Split(field, "%")
	unencodedPrefix, hexCodePrefixedSubstrings := splits[0], splits[1:]

	outputBuffer.WriteString(unencodedPrefix)
	for _, hexCodePlusUnencoded := range hexCodePrefixedSubstrings {
		if len(hexCodePlusUnencoded) < 2 {
			return "", UnfinishedPercentEncoding{EnclosingField: field}
		}

		hexCodeCharacters, unencodedRemainder := splitAfterHexCode(hexCodePlusUnencoded)
		outputBuffer.WriteByte(decode(hexCodeCharacters))
		outputBuffer.WriteString(unencodedRemainder)
	}

	return outputBuffer.String(), nil
}

func splitAfterHexCode(hexCodePlusUnencoded string) (hexCode string, unencoded string) {
	return hexCodePlusUnencoded[:2], hexCodePlusUnencoded[2:]
}

func decode(octetCharacters string) byte {
	const base16 = 16
	const uintSizeInBits = 8
	asciiCode, _ := strconv.ParseInt(octetCharacters, base16, uintSizeInBits)
	return byte(asciiCode)
}

type UnfinishedPercentEncoding struct {
	EnclosingField string
}

func (invalid UnfinishedPercentEncoding) Error() string {
	return fmt.Sprintf("%% followed by fewer than 2 characters: %s", invalid.EnclosingField)
}
