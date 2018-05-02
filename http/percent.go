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
			return "", fmt.Errorf("%% followed by fewer than 2 characters: %s", field)
		}

		hexCodeCharacters, unencodedRemainder := splitAfterHexCode(hexCodePlusUnencoded)
		outputBuffer.WriteByte(decode(hexCodeCharacters))
		outputBuffer.WriteString(unencodedRemainder)
	}

	return outputBuffer.String(), nil
}

func splitAfterHexCode(hexCodePlusUnencoded string) (hexCode string, unencoded string) {
	return hexCodePlusUnencoded[0:2], hexCodePlusUnencoded[2:]
}

func decode(octetCharacters string) byte {
	asciiCode, _ := strconv.ParseInt(octetCharacters, 16, 8)
	return byte(asciiCode)
}
