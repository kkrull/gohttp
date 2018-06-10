package cmd_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/kkrull/gohttp/capability"
	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/log"
	"github.com/kkrull/gohttp/main/cmd"
	"github.com/kkrull/gohttp/playground"
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
			factory = &cmd.InterruptFactory{
				Interrupts:     interrupts,
				MaxConnections: 42,
			}

			server = factory.TCPServer("/public", "localhost", 8421)
			typedServer, _ = server.(*http.TCPServer)
		})

		It("returns an http.TCPServer with the specified level of concurrency", func() {
			Expect(server).To(BeAssignableToTypeOf(&http.TCPServer{}))
			Expect(typedServer.MaxConnections).To(Equal(uint(42)))
		})

		Describe("it has built-in routes that are reasonable defaults in many applications", func() {
			It("* - target, not path - is for server-wide capabilities", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeAssignableToTypeOf(capability.NewRoute("*"))))
			})

			It("/coffee and /tea implement the teapot protocol", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeAssignableToTypeOf(teapot.NewRoute())))
			})

			It("/logs shows recently-made requests to administrators", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeEquivalentTo(
					log.NewLogRoute(
						"/logs",
						log.NewBufferedRequestLogger(),
					))))
			})

			It("every other path hits the file system route to read files from the content base path", func() {
				firstRoute := typedServer.Routes()[len(typedServer.Routes())-1]
				Expect(firstRoute).To(BeAssignableToTypeOf(fs.NewRoute("/tmp")))
			})
		})

		Describe("it has playground routes to show off basic capabilities to cob_spec", func() {
			It("/cat-form is good for remembering one writable cat at a time", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeEquivalentTo(playground.NewSingletonRoute("/cat-form"))))
			})

			It("/cookie and /eat_cookie tests round trips for cookies", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeEquivalentTo(playground.NewCookieRoute("/cookie", "/eat_cookie"))))
			})

			It("/form is a black hole that accepts form posts", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeEquivalentTo(playground.NewNopPostRoute("/form"))))
			})

			It("/method_options is an abstract route supporting read and write methods", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeAssignableToTypeOf(playground.NewReadWriteRoute("/method_options"))))
			})

			It("/method_options2 is an abstract route supporting read methods", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeAssignableToTypeOf(playground.NewReadOnlyRoute("/method_options2"))))
			})

			It("/parameters shows you a parsed and decoded query string", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeAssignableToTypeOf(playground.NewParameterRoute("/parameters"))))
			})

			It("/patch-content.txt is a writable file", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeEquivalentTo(fs.NewWritableFileRoute("/patch-content.txt", "/public"))))
			})

			It("/put-target is a writable file", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeEquivalentTo(fs.NewWritableFileRoute("/put-target", "/public"))))
			})

			It("/redirect sends you to your home", func() {
				Expect(typedServer.Routes()).To(ContainElement(BeAssignableToTypeOf(playground.NewRedirectRoute("/redirect"))))
			})
		})
	})
})
