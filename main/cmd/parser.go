package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

type CliCommandParser struct {
	Factory    AppFactory
	Interrupts <-chan os.Signal
}

func (parser *CliCommandParser) Parse(args []string) CliCommand {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	path := flagSet.String("d", "", "The root content directory, from which to operate")
	host := "localhost"
	port := flagSet.Uint("p", 0, "The TCP port on which to listen")
	suppressUntimelyOutput(flagSet)

	err := flagSet.Parse(args[1:])
	switch {
	case err == flag.ErrHelp:
		return parser.Factory.HelpCommand(flagSet)
	case err != nil:
		return parser.Factory.ErrorCommand(err)
	case *path == "":
		return parser.Factory.ErrorCommand(fmt.Errorf("missing path"))
	case *port == 0:
		return parser.Factory.ErrorCommand(fmt.Errorf("missing port"))
	default:
		server := parser.Factory.TCPServer(*path, host, uint16(*port))
		command, quit := parser.Factory.RunCommand(server)
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

type AppFactory interface {
	ErrorCommand(err error) CliCommand
	HelpCommand(flagSet *flag.FlagSet) CliCommand
	RunCommand(server Server) (command CliCommand, quit chan bool)
	TCPServer(contentBasePath string, host string, port uint16) Server
}

type CliCommand interface {
	Run(stderr io.Writer) (code int, err error)
}
