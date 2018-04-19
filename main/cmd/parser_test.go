package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"syscall"

	"github.com/kkrull/gohttp/main/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CliCommandParser", func() {
	Describe("#Build", func() {
		var (
			parser     *cmd.CliCommandParser
			factory    *CommandFactoryMock
			interrupts chan os.Signal

			returned cmd.CliCommand
			stderr   *bytes.Buffer
		)

		BeforeEach(func() {
			interrupts = make(chan os.Signal, 1)
			stderr = &bytes.Buffer{}
		})

		Context("given --help", func() {
			BeforeEach(func() {
				helpCommand := &CliCommandMock{}
				factory = &CommandFactoryMock{HelpCommandReturns: helpCommand}
				parser = &cmd.CliCommandParser{Factory: factory}

				returned = parser.Parse([]string{"/path/to/gohttp", "--help"})
			})

			It("returns a HelpCommand for the program", func() {
				factory.HelpCommandShouldBeForProgram("/path/to/gohttp")
			})
			It("the command has usage for the root directory parameter", func() {
				factory.HelpCommandShouldHaveFlag("d", "The root content directory, from which to operate")
			})
			It("the command has usage for the root directory parameter", func() {
				factory.HelpCommandShouldHaveFlag("p", "The TCP port on which to listen")
			})
		})

		Describe("parsing failures", func() {
			var errorCommand *CliCommandMock

			BeforeEach(func() {
				errorCommand = &CliCommandMock{}
				factory = &CommandFactoryMock{ErrorCommandReturns: errorCommand}
				parser = &cmd.CliCommandParser{Factory: factory}
			})

			Context("given unrecognized arguments", func() {
				It("returns an ErrorCommand", func() {
					returned = parser.Parse([]string{"gohttp", "--bogus"})
					Expect(returned).To(BeIdenticalTo(errorCommand))
				})
			})

			Context("when the path is missing", func() {
				It("returns an ErrorCommand stating that the path is missing", func() {
					returned = parser.Parse([]string{"gohttp", "-p", "4242"})
					factory.ErrorCommandShouldHaveReceived(fmt.Errorf("missing path"))
				})
			})

			Context("when the port is missing", func() {
				It("returns an ErrorCommand stating that the port is missing", func() {
					returned = parser.Parse([]string{"gohttp", "-d", "/tmp"})
					factory.ErrorCommandShouldHaveReceived(fmt.Errorf("missing port"))
				})
			})
		})
	})

	Describe("#Build", func() {
		var (
			parser     *cmd.CliCommandParser
			factory    *cmd.InterruptFactory
			interrupts chan os.Signal

			command cmd.CliCommand
			stderr  *bytes.Buffer
		)

		BeforeEach(func() {
			interrupts = make(chan os.Signal, 1)
			factory = &cmd.InterruptFactory{Interrupts: interrupts}
			parser = &cmd.CliCommandParser{Factory: factory}
			stderr = &bytes.Buffer{}
		})

		Context("given a complete configuration for the HTTP server", func() {
			It("returns a RunServerCommand", func() {
				factory := cmd.InterruptFactory{Interrupts: interrupts}
				parser = factory.NewCliCommandParser()
				command = parser.Parse([]string{"gohttp", "-p", "4242", "-d", "/tmp"})
				Expect(command).To(BeAssignableToTypeOf(cmd.RunServerCommand{}))
			})

			It("the command waits until a signal is sent to the interrupt signal channel", func() {
				quitHttpServer := make(chan bool, 1)
				parser = &cmd.CliCommandParser{
					Interrupts: interrupts,
					NewCommandToRunHTTPServer: func(string, string, uint16) (cmd.CliCommand, chan bool) {
						return nil, quitHttpServer
					}}

				command = parser.Parse([]string{"gohttp", "-p", "4242", "-d", "/tmp"})
				interrupts <- syscall.SIGINT
				Eventually(quitHttpServer).Should(Receive())
			})
		})
	})
})
