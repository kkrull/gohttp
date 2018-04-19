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

		Context("given a complete configuration for the HTTP server", func() {
			var runCommand *CliCommandMock

			It("creates a TCPServer bound to localhost, for the specified port and content directory", func() {
				runCommand = &CliCommandMock{}
				factory = &CommandFactoryMock{RunCommandReturns: runCommand}
				parser = &cmd.CliCommandParser{
					Factory:    factory,
					Interrupts: interrupts,
				}

				returned = parser.Parse([]string{"gohttp", "-p", "4242", "-d", "/tmp"})
				factory.TCPServerShouldHaveReceived("/tmp", "localhost", 4242)
			})

			It("returns a RunServerCommand", func() {
				runCommand = &CliCommandMock{}
				factory = &CommandFactoryMock{RunCommandReturns: runCommand}
				parser = &cmd.CliCommandParser{
					Factory:    factory,
					Interrupts: interrupts,
				}

				returned = parser.Parse([]string{"gohttp", "-p", "4242", "-d", "/tmp"})
				Expect(returned).To(BeIdenticalTo(runCommand))
			})

			It("wires interrupt signals on the .Interrupts channel to the channel used to terminate the command", func() {
				quitCommand := make(chan bool, 1)
				factory = &CommandFactoryMock{RunCommandReturnsChannel: quitCommand}
				parser = &cmd.CliCommandParser{
					Factory:    factory,
					Interrupts: interrupts,
				}

				parser.Parse([]string{"gohttp", "-p", "4242", "-d", "/tmp"})
				interrupts <- syscall.SIGINT
				Eventually(quitCommand).Should(Receive())
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
})
