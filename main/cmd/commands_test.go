package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kkrull/gohttp/main/cmd"
)

var _ = Describe("CliCommands", func() {
	var (
		stderr *bytes.Buffer
		code   int
		err    error
	)

	BeforeEach(func() {
		stderr = &bytes.Buffer{}
	})

	Describe("ErrorCommand", func() {
		Describe("#Run", func() {
			It("returns the error and an exit code of 1, to indicate failure", func() {
				givenError := fmt.Errorf("bang")
				command := cmd.ErrorCommand{Error: givenError}
				code, err = command.Run(stderr)
				Expect(code).To(Equal(1))
				Expect(err).To(Equal(givenError))
			})
		})
	})

	Describe("HelpCommand", func() {
		Describe("#Run", func() {
			BeforeEach(func() {
				command := cmd.HelpCommand{FlagSet: flag.NewFlagSet("widget", flag.ContinueOnError)}
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
			command cmd.CliCommand
			factory *cmd.InterruptFactory
			server  ServerMock
			quit    chan bool
		)

		BeforeEach(func() {
			factory = &cmd.InterruptFactory{Interrupts: make(chan os.Signal, 1)}
		})

		Describe("#Run", func() {
			Context("given a workable configuration", func() {
				BeforeEach(func() {
					server = ServerMock{}
					command, quit = factory.RunCommand(&server)
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
					server = ServerMock{}
					command, quit = factory.RunCommand(&server)
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
					server = ServerMock{StartFails: "no listening ears"}
					command, quit = factory.RunCommand(&server)
					code, err = command.Run(stderr)
					Expect(code).To(Equal(2))
					Expect(err).To(MatchError("no listening ears"))
				})
			})

			Context("when there is an error shutting down", func() {
				It("returns the error and an exit code indicating failure", func() {
					server = ServerMock{ShutdownFails: "backfire"}
					command, quit = factory.RunCommand(&server)
					go scheduleShutdown(quit)
					code, err = command.Run(stderr)
					Expect(code).To(Equal(3))
					Expect(err).To(MatchError("backfire"))
				})
			})
		})
	})
})
