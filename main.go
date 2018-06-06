package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/kkrull/gohttp/main/cmd"
)

const maxConnections uint = 4

func main() {
	factory := &cmd.InterruptFactory{
		Interrupts:     subscribeToSignals(os.Interrupt),
		MaxConnections: maxConnections}
	gohttp := &GoHTTP{
		CommandParser: factory.CliCommandParser(),
		Stderr:        os.Stderr}

	code, err := gohttp.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gohttp: %s\n", err.Error())
	}

	os.Exit(code)
}

type GoHTTP struct {
	CommandParser CommandParser
	Stderr        io.Writer
}

func (gohttp *GoHTTP) Run(args []string) (exitCode int, runErr error) {
	command := gohttp.CommandParser.Parse(args)
	return command.Run(gohttp.Stderr)
}

func subscribeToSignals(sig os.Signal) <-chan os.Signal {
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, sig)
	return interrupts
}

type CommandParser interface {
	Parse(args []string) cmd.CliCommand
}
