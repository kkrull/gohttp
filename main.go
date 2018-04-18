package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/kkrull/gohttp/main/cmd"
)

func main() {
	parser := cmd.NewCliCommandParser(subscribeToSignals(os.Interrupt))
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
