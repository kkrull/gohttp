package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/kkrull/gohttp/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kkrull/gohttp/main/cmd"
)

var _ = Describe("CliCommandParser", func() {
	Describe("#Build", func() {
		var (
			parser     *cmd.CliCommandParser
			interrupts chan os.Signal

			command cmd.CliCommand
			stderr  *bytes.Buffer
		)

		BeforeEach(func() {
			interrupts = make(chan os.Signal, 1)
			parser = cmd.NewCliCommandParser(interrupts)
			stderr = &bytes.Buffer{}
		})

		Context("given --help", func() {
			BeforeEach(func() {
				command = parser.Parse([]string{"/path/to/gohttp", "--help"})
			})

			It("returns HelpCommand", func() {
				Expect(command).To(BeAssignableToTypeOf(cmd.HelpCommand{}))
			})
			It("the HelpCommand is configured for the name of the program in the first argument", func() {
				command.Run(stderr)
				Expect(stderr.String()).To(HavePrefix("Usage of /path/to/gohttp"))
			})
			It("the HelpCommand shows usage for the gohttp arguments", func() {
				command.Run(stderr)
				Expect(stderr.String()).To(ContainSubstring("The root content directory"))
				Expect(stderr.String()).To(ContainSubstring("The TCP port on which to listen"))
			})
		})

		Context("given a complete configuration for the HTTP server", func() {
			It("returns a RunServerCommand", func() {
				parser = cmd.NewCliCommandParser(interrupts)
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

		Context("given unrecognized arguments", func() {
			It("returns ErrorCommand", func() {
				command = parser.Parse([]string{"gohttp", "--bogus"})
				Expect(command).To(BeAssignableToTypeOf(cmd.ErrorCommand{}))
			})
		})

		Context("when the path is missing", func() {
			It("returns ErrorCommand", func() {
				command = parser.Parse([]string{"gohttp", "-p", "4242"})
				Expect(command).To(BeAssignableToTypeOf(cmd.ErrorCommand{}))
			})
		})

		Context("when the port is missing", func() {
			It("returns ErrorCommand", func() {
				command = parser.Parse([]string{"gohttp", "-d", "/tmp"})
				Expect(command).To(BeAssignableToTypeOf(cmd.ErrorCommand{}))
			})
		})
	})
})

var _ = Describe("CliCommands", func() {
	var (
		command cmd.CliCommand
		stderr  *bytes.Buffer
		code    int
		err     error
	)

	BeforeEach(func() {
		stderr = &bytes.Buffer{}
	})

	Describe("ErrorCommand", func() {
		Describe("#Run", func() {
			It("returns the error and an exit code of 1, to indicate failure", func() {
				givenError := fmt.Errorf("bang")
				command = cmd.ErrorCommand{Error: givenError}
				code, err = command.Run(stderr)
				Expect(code).To(Equal(1))
				Expect(err).To(Equal(givenError))
			})
		})
	})

	Describe("HelpCommand", func() {
		Describe("#Run", func() {
			BeforeEach(func() {
				command = cmd.HelpCommand{FlagSet: flag.NewFlagSet("widget", flag.ContinueOnError)}
				code, err = command.Run(stderr)
			})

			It("prints usage to the given Writer", func() {
				Expect(stderr.String()).To(ContainSubstring("Usage of widget"))
			})
			It("returns no error and an exit code of 0, to indicate success", func() {
				Expect(code).To(Equal(0))
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("RunServerCommand", func() {
		var (
			server mock.Server
			quit   chan bool
		)

		Describe("#Run", func() {
			Context("given a workable configuration", func() {
				BeforeEach(func() {
					server = mock.Server{}
					command, quit = cmd.NewRunServerCommand(&server)
				})

				It("runs the server until the quit channel receives something", func(done Done) {
					go func() {
						defer GinkgoRecover()
						command.Run(stderr)
						server.VerifyShutdown()
						close(done)
					}()

					waitForStart()
					server.VerifyRunning()
					quit <- true
				})
			})

			Context("when everything has run ok", func() {
				BeforeEach(func() {
					server = mock.Server{}
					command, quit = cmd.NewRunServerCommand(&server)
				})

				It("returns 0 and no error", func() {
					go scheduleShutdown(quit)
					code, err = command.Run(stderr)
					Expect(code).To(Equal(0))
					Expect(err).To(BeNil())
				})
			})

			Context("when there is an error starting the server", func() {
				It("returns the error and an exit code indicating failure", func() {
					server = mock.Server{StartFails: "no listening ears"}
					command, quit = cmd.NewRunServerCommand(&server)
					code, err = command.Run(stderr)
					Expect(code).To(Equal(2))
					Expect(err).To(MatchError("no listening ears"))
				})
			})

			Context("when there is an error shutting down", func() {
				It("returns the error and an exit code indicating failure", func() {
					server = mock.Server{ShutdownFails: "backfire"}
					command, quit = cmd.NewRunServerCommand(&server)
					go scheduleShutdown(quit)
					code, err = command.Run(stderr)
					Expect(code).To(Equal(3))
					Expect(err).To(MatchError("backfire"))
				})
			})
		})
	})
})
