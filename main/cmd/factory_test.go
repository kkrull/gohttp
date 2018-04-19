package cmd_test

import (
	"flag"

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
