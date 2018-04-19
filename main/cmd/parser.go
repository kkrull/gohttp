package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

type CliCommandParser struct {
	Interrupts                <-chan os.Signal
	NewCommandToRunHTTPServer MakeCommandToRunHTTPServer
	Factory                   CommandFactory
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
		return parser.Factory.HelpCommand(flagSet)
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

type CommandFactory interface {
	HelpCommand(flagSet *flag.FlagSet) CliCommand
}

type CliCommand interface {
	Run(stderr io.Writer) (code int, err error)
}
