package cmd_test

import (
	"flag"
	"fmt"

	"github.com/kkrull/gohttp/main/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InterruptFactory", func() {
	var factory *cmd.InterruptFactory
	var command cmd.CliCommand

	BeforeEach(func() {
		factory = &cmd.InterruptFactory{}
	})

	Describe("#ErrorCommand", func() {
		It("returns an ErrorCommand with the given error", func() {
			err := fmt.Errorf("kaboom")
			command = factory.ErrorCommand(err)
			Expect(command).To(BeEquivalentTo(cmd.ErrorCommand{Error: err}))
		})
	})

	Describe("#HelpCommand", func() {
		It("returns HelpCommand", func() {
			flagSet := flag.NewFlagSet("program", flag.ContinueOnError)
			command = factory.HelpCommand(flagSet)
			Expect(command).To(BeIdenticalTo(cmd.HelpCommand{FlagSet: flagSet}))
		})
	})

	Describe("NewCommandToRunHTTPServer", func() {
		XIt("configures a coffee route")
	})
})
