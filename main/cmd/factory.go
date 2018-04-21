package cmd

import (
	"flag"
	"os"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/opt"
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
	router := &http.RequestLineRouter{}
	router.AddRoute(opt.NewRoute())
	router.AddRoute(teapot.NewRoute())
	router.AddRoute(fs.NewRoute(contentRootPath))
	return router
}
