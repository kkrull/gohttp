package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"flag"
	"fmt"
	. "github.com/kkrull/gohttp"
	"github.com/kkrull/gohttp/mock"
)

var _ = Describe("CliCommandParser", func() {
	Describe("#Build", func() {
		var (
			parser  *CliCommandParser
			command CliCommand
		)

		Context("given --help", func() {
			It("returns HelpCommand", func() {
				parser = &CliCommandParser{}
				command = parser.Parse([]string{"gohttp", "--help"})
				Expect(command).To(BeAssignableToTypeOf(HelpCommand{}))
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
		var server mock.HttpServer

		Describe("#Run", func() {
			Context("when there is no error", func() {
				BeforeEach(func() {
					server = mock.HttpServer{}
					command = RunServerCommand{Server: &server}
					code, err = command.Run(stderr)
				})

				It("starts the server", func() {
					server.VerifyListen()
				})
				It("returns no error and an exit code of 0, to indicate success", func() {
					Expect(code).To(Equal(0))
					Expect(err).To(BeNil())
				})
			})

			Context("when there is an error", func() {
				It("returns the error and an exit code indicating failure", func() {
					server = mock.HttpServer{ListenFails: "no listening ears"}
					command = RunServerCommand{Server: &server}
					code, err = command.Run(stderr)
					Expect(code).To(Equal(2))
					Expect(err).To(MatchError("no listening ears"))
				})
			})
		})
	})
})
