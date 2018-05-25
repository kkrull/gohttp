package cmd

import (
	"bytes"
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
	Interrupts <-chan os.Signal
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
	return http.MakeTCPServerWithHandler(
		host,
		port,
		http.NewConnectionHandler(router))
}

func (factory *InterruptFactory) routerWithAllRoutes(contentRootPath string) http.Router {
	router := http.NewRouter()

	buffer := &bytes.Buffer{}
	logger := &http.TextLogger{Writer: buffer}
	router.LogRequests(logger)

	router.AddRoute(capability.NewRoute())
	router.AddRoute(log.NewLogRoute("/logs", logger))
	router.AddRoute(playground.NewSingletonRoute("/cat-form"))
	router.AddRoute(playground.NewWriteOKRoute("/form"))
	router.AddRoute(playground.NewWriteOKRoute("/put-target"))
	router.AddRoute(playground.NewParameterRoute())
	router.AddRoute(playground.NewReadOnlyRoute())
	router.AddRoute(playground.NewReadWriteRoute())
	router.AddRoute(playground.NewRedirectRoute())
	router.AddRoute(teapot.NewRoute())
	router.AddRoute(fs.NewRoute(contentRootPath))
	return router
}
