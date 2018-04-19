package main_test

import (
	"bytes"
	"fmt"

	. "github.com/kkrull/gohttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoHTTP", func() {
	Describe("#Run", func() {
		var (
			gohttp  *GoHTTP
			parser  *CommandParserMock
			command *CliCommandMock
			stderr  *bytes.Buffer
		)

		BeforeEach(func() {
			command = &CliCommandMock{}
			parser = &CommandParserMock{ParseReturns: command}
			gohttp = &GoHTTP{CommandParser: parser, Stderr: stderr}
		})

		It("parses a command from the given arguments", func() {
			gohttp.Run([]string{"save", "world"})
			parser.ParseShouldHaveReceived([]string{"save", "world"})
		})

		It("runs the command", func() {
			gohttp.Run(nil)
			command.RunShouldHaveReceived(stderr)
		})

		It("returns the exit code and any error from running the command", func() {
			command := &CliCommandMock{RunReturnsCode: 42, RunReturnsError: fmt.Errorf("bang")}
			parser = &CommandParserMock{ParseReturns: command}
			gohttp = &GoHTTP{CommandParser: parser, Stderr: stderr}

			exitCode, returnedRunErr := gohttp.Run(nil)
			Expect(exitCode).To(Equal(42))
			Expect(returnedRunErr).To(BeIdenticalTo(command.RunReturnsError))
		})
	})
})
