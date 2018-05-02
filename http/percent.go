package http

import (
	"bytes"
	"strconv"
	"strings"
)

func PercentDecode(field string) string {
	outputBuffer := &bytes.Buffer{}
	doDecoding2(field, outputBuffer)
	return outputBuffer.String()
}

func doDecoding2(field string, outputBuffer *bytes.Buffer) {
	splits := strings.Split(field, "%")
	inputBeforeFirstTriplet := splits[0]
	outputBuffer.WriteString(inputBeforeFirstTriplet)

	inputContainingTriplets := splits[1:]
	for _, hexCodePlusUnencoded := range inputContainingTriplets {
		hexCode := hexCodePlusUnencoded[0:2]
		asciiCode, _ := strconv.ParseInt(hexCode, 16, 8)
		decodedCharacter := byte(asciiCode)
		outputBuffer.WriteByte(decodedCharacter)

		unencodedInput := hexCodePlusUnencoded[2:]
		outputBuffer.WriteString(unencodedInput)
	}
}
