package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
)

func main() {
	parser := NewCliCommandParser(subscribeToSignals(os.Interrupt))
	command := parser.Parse(os.Args)
	code, runErr := command.Run(os.Stderr)
	if runErr != nil {
		fmt.Fprintf(os.Stderr, "gohttp: %s\n", runErr.Error())
	}

	os.Exit(code)
}

func subscribeToSignals(sig os.Signal) <-chan os.Signal {
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, sig)
	return interrupts
}

/* Command parsing */

func NewCliCommandParser(interrupts <-chan os.Signal) *CliCommandParser {
	return &CliCommandParser{
		Interrupts:                interrupts,
		NewCommandToRunHTTPServer: NewCommandToRunHTTPServer,
	}
}

type CliCommandParser struct {
	Interrupts                <-chan os.Signal
	NewCommandToRunHTTPServer MakeCommandToRunHTTPServer
}

func NewCommandToRunHTTPServer(contentRootPath string, host string, port uint16) (CliCommand, chan bool) {
	router := &http.RequestLineRouter{}
	router.AddRoute(fs.NewRoute(contentRootPath))
	server := http.MakeTCPServerWithHandler(
		host,
		port,
		http.NewConnectionHandler(router))
	return NewRunServerCommand(server)
}

type MakeCommandToRunHTTPServer func(contentRootPath string, host string, port uint16) (
	command CliCommand, quit chan bool)

func (parser *CliCommandParser) Parse(args []string) CliCommand {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	path := flagSet.String("d", "", "The root content directory, from which to operate")
	host := "localhost"
	port := flagSet.Uint("p", 0, "The TCP port on which to listen")
	suppressUntimelyOutput(flagSet)

	err := flagSet.Parse(args[1:])
	switch {
	case err == flag.ErrHelp:
		return HelpCommand{FlagSet: flagSet}
	case err != nil:
		return ErrorCommand{Error: err}
	case *path == "":
		return ErrorCommand{Error: fmt.Errorf("missing path")}
	case *port == 0:
		return ErrorCommand{Error: fmt.Errorf("missing port")}
	default:
		command, quit := parser.NewCommandToRunHTTPServer(*path, host, uint16(*port))
		go parser.sendTrueOnFirstInterruption(quit)
		return command
	}
}

func suppressUntimelyOutput(flagSet *flag.FlagSet) {
	flagSet.SetOutput(&bytes.Buffer{})
}

func (parser *CliCommandParser) sendTrueOnFirstInterruption(quit chan<- bool) {
	<-parser.Interrupts
	quit <- true
}

type CliCommand interface {
	Run(stderr io.Writer) (code int, err error)
}

/* ErrorCommand */

type ErrorCommand struct {
	Error error
}

func (command ErrorCommand) Run(stderr io.Writer) (code int, err error) {
	return 1, command.Error
}

/* HelpCommand */

type HelpCommand struct {
	FlagSet *flag.FlagSet
}

func (command HelpCommand) Run(stderr io.Writer) (code int, err error) {
	command.FlagSet.SetOutput(stderr)
	command.FlagSet.Usage()
	return 0, nil
}

/* RunServerCommand */

func NewRunServerCommand(server Server) (command CliCommand, quit chan bool) {
	quit = make(chan bool, 1)
	command = RunServerCommand{Server: server, quit: quit}
	return
}

type RunServerCommand struct {
	Server Server
	quit   <-chan bool
}

type Server interface {
	Address() net.Addr
	Start() error
	Shutdown() error
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
