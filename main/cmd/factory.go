package cmd

import (
	"flag"
	"os"

	"github.com/kkrull/gohttp/capability"
	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/log"
	"github.com/kkrull/gohttp/playground"
	"github.com/kkrull/gohttp/teapot"
)

type InterruptFactory struct {
	Interrupts     <-chan os.Signal
	MaxConnections uint
}

func (factory *InterruptFactory) CliCommandParser() *CliCommandParser {
	return &CliCommandParser{
		Factory:    factory,
		Interrupts: factory.Interrupts,
	}
}

func (factory *InterruptFactory) ErrorCommand(err error) CliCommand {
	return &ErrorCommand{Error: err}
}

func (factory *InterruptFactory) HelpCommand(flagSet *flag.FlagSet) CliCommand {
	return &HelpCommand{FlagSet: flagSet}
}

func (factory *InterruptFactory) RunCommand(server Server) (command CliCommand, quit chan bool) {
	quit = make(chan bool, 1)
	command = RunServerCommand{Server: server, Quit: quit}
	return
}

func (factory *InterruptFactory) TCPServer(contentRootPath string, host string, port uint16) Server {
	router := factory.routerWithAllRoutes(contentRootPath)
	handler := http.NewConnectionHandler(router)
	return http.TCPServerBuilder(host).
		ListeningOnPort(port).
		WithConnectionHandler(handler).
		WithMaxConnections(factory.MaxConnections).
		Build()
}

func (factory *InterruptFactory) routerWithAllRoutes(contentRootPath string) http.Router {
	router := http.NewRouter()

	logger := log.NewBufferedRequestLogger()
	router.LogRequests(logger)
	router.AddRoute(log.NewLogRoute("/logs", logger))

	router.AddRoute(capability.NewRoute("*"))
	router.AddRoute(teapot.NewRoute())

	router.AddRoute(playground.NewSingletonRoute("/cat-form"))
	router.AddRoute(playground.NewCookieRoute("/cookie", "/eat_cookie"))
	router.AddRoute(playground.NewNopPostRoute("/form"))
	router.AddRoute(playground.NewReadWriteRoute("/method_options"))
	router.AddRoute(playground.NewReadOnlyRoute("/method_options2"))
	router.AddRoute(playground.NewParameterRoute("/parameters"))
	router.AddRoute(playground.NewRedirectRoute("/redirect"))

	//A couple paths are backed by real files and are meant to support write operations; the rest default to read-only
	router.AddRoute(fs.NewWritableFileRoute("/patch-content.txt", contentRootPath))
	router.AddRoute(fs.NewWritableFileRoute("/put-target", contentRootPath))
	router.AddRoute(fs.NewRoute(contentRootPath))
	return router
}
