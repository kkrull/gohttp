package http

import (
	"bytes"
	"fmt"
	"strconv"
)

func PercentDecode(field string) string {
	inputBuffer := bytes.NewBufferString(field)
	inputBytes := make([]byte, len(field))
	inputSize, _ := inputBuffer.Read(inputBytes)

	outputBuffer := &bytes.Buffer{}
	for i := 0; i < inputSize; {
		input := inputBytes[i]
		if input != '%' {
			fmt.Printf("- Writing %x %s\n", input, string(input))
			outputBuffer.WriteByte(input)
			i += 1
			continue
		}

		var asciiCode int = 0
		asciiCode += 16 * base10Value(inputBytes[i+1])
		asciiCode += base10Value(inputBytes[i+2])
		fmt.Printf("- Writing %x %s\n", asciiCode, string(asciiCode))
		outputBuffer.WriteByte(byte(asciiCode))
		i += 3
	}

	fmt.Printf("field=<%s>, inputBytes_10[%d]=<%v>, inputBytes_16[%d]=<0x%x>, output[%d]=<%s>\n",
		field,
		inputSize, inputBytes,
		inputSize, inputBytes,
		outputBuffer.Len(), outputBuffer.String(),
	)
	return outputBuffer.String()
}

func base10Value(character byte) int {
	value, _ := strconv.ParseInt(string(character), 16, 8)
	return int(value)
}
