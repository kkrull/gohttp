package cmd_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/main/cmd"
	"github.com/kkrull/gohttp/teapot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("InterruptFactory", func() {
	var (
		factory    *cmd.InterruptFactory
		interrupts chan os.Signal
	)

	Describe("CliCommand methods", func() {
		var command cmd.CliCommand

		BeforeEach(func() {
			factory = &cmd.InterruptFactory{}
		})

		Describe("#ErrorCommand", func() {
			It("returns an ErrorCommand with the given error", func() {
				err := fmt.Errorf("kaboom")
				command = factory.ErrorCommand(err)
				Expect(command).To(BeEquivalentTo(&cmd.ErrorCommand{Error: err}))
			})
		})

		Describe("#HelpCommand", func() {
			It("returns HelpCommand", func() {
				flagSet := flag.NewFlagSet("program", flag.ContinueOnError)
				command = factory.HelpCommand(flagSet)
				Expect(command).To(BeEquivalentTo(&cmd.HelpCommand{FlagSet: flagSet}))
			})
		})

		Describe("#RunCommand", func() {
			var (
				server *ServerMock
				quit   chan bool
			)

			BeforeEach(func() {
				server = &ServerMock{}
				command, quit = factory.RunCommand(server)
			})

			It("returns RunServerCommand and a channel that signals it to terminate", func() {
				Expect(command).To(BeAssignableToTypeOf(cmd.RunServerCommand{}))
				Expect(command).To(MatchFields(IgnoreExtras, Fields{
					"Server": BeIdenticalTo(server),
				}))
			})

			It("returns a channel that signals it to terminate", func() {
				Expect(quit).To(BeAssignableToTypeOf(make(chan bool, 1)))
			})
		})
	})

	Describe("TCPServer", func() {
		var (
			server      cmd.Server
			typedServer *http.TCPServer
		)

		BeforeEach(func() {
			interrupts = make(chan os.Signal, 1)
			factory = &cmd.InterruptFactory{Interrupts: interrupts}

			server = factory.TCPServer("/public", "localhost", 8421)
			typedServer, _ = server.(*http.TCPServer)
		})

		It("returns an http.TCPServer", func() {
			Expect(server).To(BeAssignableToTypeOf(&http.TCPServer{}))
		})

		It("the teapot route is first", func() {
			firstRoute := typedServer.Routes()[0]
			Expect(firstRoute).To(BeAssignableToTypeOf(teapot.NewRoute()))
		})

		It("the fs route is last", func() {
			firstRoute := typedServer.Routes()[len(typedServer.Routes())-1]
			Expect(firstRoute).To(BeAssignableToTypeOf(fs.NewRoute("/tmp")))
		})
	})
})
