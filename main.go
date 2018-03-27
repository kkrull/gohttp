package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/kkrull/gohttp/http"
	"io"
	"os"
)

func main() {
	parser := CliCommandParser{}
	command := parser.Parse(os.Args)
	code, runErr := command.Run(os.Stderr)
	if runErr != nil {
		fmt.Fprintf(os.Stderr, "gohttp: %s\n", runErr.Error())
	}

	os.Exit(code)
}

/* Command parsing */

type CliCommandParser struct{}

func (parser *CliCommandParser) Parse(args []string) CliCommand {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	path := flagSet.String("d", "", "The root content directory, from which to operate")
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
		return RunServerCommand{Server: http.NewServer(*path, *port)}
	}
}

func suppressUntimelyOutput(flagSet *flag.FlagSet) {
	flagSet.SetOutput(&bytes.Buffer{})
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

type RunServerCommand struct {
	Server http.Server
}

func (command RunServerCommand) Run(stderr io.Writer) (code int, err error) {
	if err := command.Server.Listen(); err != nil {
		return 2, err
	}

	return 0, nil
}
