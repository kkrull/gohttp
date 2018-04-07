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

type CliCommandParser struct {
	Interrupts chan os.Signal
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
		return HelpCommand{FlagSet: flagSet}
	case err != nil:
		return ErrorCommand{Error: err}
	case *path == "":
		return ErrorCommand{Error: fmt.Errorf("missing path")}
	case *port == 0:
		return ErrorCommand{Error: fmt.Errorf("missing port")}
	default:
		command, quit := MakeRunServerCommand(http.MakeTCPServer(*path, host, uint16(*port)))

		go func() {
			for interruptSignal := range parser.Interrupts {
				fmt.Printf("Received signal: %s\n", interruptSignal.String())
				quit <- true
			}
		}()

		return command
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

func MakeRunServerCommand(server http.Server) (command CliCommand, quit chan bool) {
	quit = make(chan bool, 1)
	command = RunServerCommand{Server: server, quit: quit}
	return
}

type RunServerCommand struct {
	Server http.Server
	quit   chan bool
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
