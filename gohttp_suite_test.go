package main_test

import (
	"io"
	"testing"

	"github.com/kkrull/gohttp/main/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGohttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "main suite")
}

type CliCommandMock struct {
	RunReturnsCode  int
	RunReturnsError error
	runStderr       io.Writer
}

func (mock *CliCommandMock) Run(stderr io.Writer) (code int, err error) {
	mock.runStderr = stderr
	return mock.RunReturnsCode, mock.RunReturnsError
}

func (mock *CliCommandMock) RunShouldHaveReceived(stderr io.Writer) {
	ExpectWithOffset(1, mock.runStderr).To(BeIdenticalTo(stderr))
}

type CommandParserMock struct {
	ParseReturns cmd.CliCommand
	parseArgs    []string
}

func (mock *CommandParserMock) Parse(args []string) cmd.CliCommand {
	mock.parseArgs = make([]string, len(args))
	copy(mock.parseArgs, args)
	return mock.ParseReturns
}

func (mock *CommandParserMock) ParseShouldHaveReceived(args []string) {
	ExpectWithOffset(1, mock.parseArgs).To(Equal(args))
}
