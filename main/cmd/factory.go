package cmd

import (
	"flag"
	"os"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
)

type InterruptFactory struct {
	Interrupts <-chan os.Signal
}

func (factory *InterruptFactory) ErrorCommand(err error) CliCommand {
	return &ErrorCommand{Error: err}
}

func (factory *InterruptFactory) HelpCommand(flagSet *flag.FlagSet) CliCommand {
	return &HelpCommand{FlagSet: flagSet}
}

func (factory *InterruptFactory) RunCommand(server Server) (command CliCommand, quit chan bool) {
	quit = make(chan bool, 1)
	command = RunServerCommand{Server: server, quit: quit}
	return
}

func (factory *InterruptFactory) TCPServer(contentRootPath string, host string, port uint16) Server {
	router := &http.RequestLineRouter{}
	router.AddRoute(fs.NewRoute(contentRootPath)) //TODO KDK: Add a route here for coffee pots
	return http.MakeTCPServerWithHandler(
		host,
		port,
		http.NewConnectionHandler(router))
}

func (factory *InterruptFactory) NewCliCommandParser() *CliCommandParser {
	return &CliCommandParser{
		Interrupts:                factory.Interrupts,
		NewCommandToRunHTTPServer: factory.NewCommandToRunHTTPServer,
	}
}

func (factory *InterruptFactory) NewCommandToRunHTTPServer(contentRootPath string, host string, port uint16) (CliCommand, chan bool) {
	server := factory.TCPServer(contentRootPath, host, port)
	return factory.RunCommand(server)
}
