package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/main/cmd"
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
	return cmd.NewRunServerCommand(server)
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
		return cmd.HelpCommand{FlagSet: flagSet}
	case err != nil:
		return cmd.ErrorCommand{Error: err}
	case *path == "":
		return cmd.ErrorCommand{Error: fmt.Errorf("missing path")}
	case *port == 0:
		return cmd.ErrorCommand{Error: fmt.Errorf("missing port")}
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
