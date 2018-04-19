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
	return ErrorCommand{Error: err}
}

func (factory *InterruptFactory) HelpCommand(flagSet *flag.FlagSet) CliCommand {
	return HelpCommand{FlagSet: flagSet}
}

func (factory *InterruptFactory) RunCommand(server Server) (CliCommand, chan bool) {
	return NewRunServerCommand(server)
}

func (factory *InterruptFactory) NewCliCommandParser() *CliCommandParser {
	return &CliCommandParser{
		Interrupts:                factory.Interrupts,
		NewCommandToRunHTTPServer: factory.NewCommandToRunHTTPServer,
	}
}

func (factory *InterruptFactory) NewCommandToRunHTTPServer(contentRootPath string, host string, port uint16) (CliCommand, chan bool) {
	router := &http.RequestLineRouter{}
	router.AddRoute(fs.NewRoute(contentRootPath)) //TODO KDK: Add a route here for coffee pots
	server := http.MakeTCPServerWithHandler(
		host,
		port,
		http.NewConnectionHandler(router))
	return NewRunServerCommand(server)
}
