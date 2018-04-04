package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"flag"
	"fmt"
	. "github.com/kkrull/gohttp"
	"github.com/kkrull/gohttp/mock"
	"time"
)

var _ = Describe("CliCommandParser", func() {
	Describe("#Build", func() {
		var (
			parser  *CliCommandParser
			command CliCommand
		)

		Context("given --help", func() {
			BeforeEach(func() {
				parser = &CliCommandParser{}
				command = parser.Parse([]string{"/path/to/gohttp", "--help"})
			})

			It("returns HelpCommand", func() {
				Expect(command).To(BeAssignableToTypeOf(HelpCommand{}))
			})
			It("the HelpCommand is configured for the name of the program in the first argument", func() {
				stderr := &bytes.Buffer{}
				command.Run(stderr)
				Expect(stderr.String()).To(HavePrefix("Usage of /path/to/gohttp"))
			})
			It("the HelpCommand shows usage for the gohttp arguments", func() {
				stderr := &bytes.Buffer{}
				command.Run(stderr)
				Expect(stderr.String()).To(ContainSubstring("The root content directory"))
				Expect(stderr.String()).To(ContainSubstring("The TCP port on which to listen"))
			})
		})

		Context("given a complete configuration for the HTTP server", func() {
			It("returns RunServerCommand", func() {
				parser = &CliCommandParser{}
				command = parser.Parse([]string{"gohttp", "-p", "4242", "-d", "/tmp"})
				Expect(command).To(BeAssignableToTypeOf(RunServerCommand{}))
			})
		})

		Context("given unrecognized arguments", func() {
			It("returns ErrorCommand", func() {
				parser = &CliCommandParser{}
				command = parser.Parse([]string{"gohttp", "--bogus"})
				Expect(command).To(BeAssignableToTypeOf(ErrorCommand{}))
			})
		})

		Context("when the path is missing", func() {
			It("returns ErrorCommand", func() {
				parser = &CliCommandParser{}
				command = parser.Parse([]string{"gohttp", "-p", "4242"})
				Expect(command).To(BeAssignableToTypeOf(ErrorCommand{}))
			})
		})

		Context("when the port is missing", func() {
			It("returns ErrorCommand", func() {
				parser = &CliCommandParser{}
				command = parser.Parse([]string{"gohttp", "-d", "/tmp"})
				Expect(command).To(BeAssignableToTypeOf(ErrorCommand{}))
			})
		})
	})
})

var _ = Describe("CliCommands", func() {
	var (
		command CliCommand
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
				command = ErrorCommand{Error: givenError}
				code, err = command.Run(stderr)
				Expect(code).To(Equal(1))
				Expect(err).To(Equal(givenError))
			})
		})
	})

	Describe("HelpCommand", func() {
		Describe("#Run", func() {
			BeforeEach(func() {
				command = HelpCommand{FlagSet: flag.NewFlagSet("widget", flag.ContinueOnError)}
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
					command, quit = MakeRunServerCommand(&server)
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
					command, quit = MakeRunServerCommand(&server)
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
					command, quit = MakeRunServerCommand(&server)
					code, err = command.Run(stderr)
					Expect(code).To(Equal(2))
					Expect(err).To(MatchError("no listening ears"))
				})
			})

			Context("when there is an error shutting down", func() {
				It("returns the error and an exit code indicating failure", func() {
					server = mock.Server{ShutdownFails: "backfire"}
					command, quit = MakeRunServerCommand(&server)
					go scheduleShutdown(quit)
					code, err = command.Run(stderr)
					Expect(code).To(Equal(3))
					Expect(err).To(MatchError("backfire"))
				})
			})
		})
	})
})

func waitForStart() {
	time.Sleep(100 * time.Millisecond)
}

func scheduleShutdown(quit chan bool) {
	waitForStart()
	quit <- true
}
