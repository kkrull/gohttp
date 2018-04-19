package cmd

import (
	"flag"
	"io"
	"net"
)

type ErrorCommand struct {
	Error error
}

func (command ErrorCommand) Run(stderr io.Writer) (code int, err error) {
	return 1, command.Error
}

type HelpCommand struct {
	FlagSet *flag.FlagSet
}

func (command HelpCommand) Run(stderr io.Writer) (code int, err error) {
	command.FlagSet.SetOutput(stderr)
	command.FlagSet.Usage()
	return 0, nil
}

func NewRunServerCommand(server Server) (command RunServerCommand, quit chan bool) {
	quit = make(chan bool, 1)
	command = RunServerCommand{Server: server, quit: quit}
	return
}

type RunServerCommand struct {
	Server Server
	quit   <-chan bool
}

func (command RunServerCommand) Run(stderr io.Writer) (code int, err error) {
	if err := command.Server.Start(); err != nil {
		return 2, err
	}

	command.waitForShutdownRequest()
	if err := command.Server.Shutdown(); err != nil {
		return 3, err
	}

	return 0, nil
}

func (command RunServerCommand) waitForShutdownRequest() {
	<-command.quit
}

type Server interface {
	Address() net.Addr
	Start() error
	Shutdown() error
}
