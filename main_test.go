package main_test

import (
	"bytes"
	"os"
	"syscall"

	. "github.com/kkrull/gohttp"
	"github.com/kkrull/gohttp/main/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("main", func() {
	XIt("adds a coffee route to the HTTP server")
})

var _ = Describe("CliCommandParser", func() {
	Describe("#Build", func() {
		var (
			parser     *CliCommandParser
			interrupts chan os.Signal

			command CliCommand
			stderr  *bytes.Buffer
		)

		BeforeEach(func() {
			interrupts = make(chan os.Signal, 1)
			parser = NewCliCommandParser(interrupts)
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
				parser = NewCliCommandParser(interrupts)
				command = parser.Parse([]string{"gohttp", "-p", "4242", "-d", "/tmp"})
				Expect(command).To(BeAssignableToTypeOf(cmd.RunServerCommand{}))
			})

			It("the command waits until a signal is sent to the interrupt signal channel", func() {
				quitHttpServer := make(chan bool, 1)
				parser = &CliCommandParser{
					Interrupts: interrupts,
					NewCommandToRunHTTPServer: func(string, string, uint16) (CliCommand, chan bool) {
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
