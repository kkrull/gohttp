package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/kkrull/gohttp/http"
)

func main() {
	builder := &serverBuilder{}
	main := Main{Builder: builder}
	if err := main.Run(os.Args[1:]); err != nil {
		fmt.Printf("gohttp: %s\n", err.Error())
		os.Exit(1)
	}
}

type Main struct {
	Builder ServerBuilder
}

func (main Main) Run(args []string) error {
	if err := main.Builder.ParseCommandLine(args); err != nil {
		return err
	}

	server, buildError := main.Builder.Build()
	if buildError != nil {
		return buildError
	}

	return server.Listen()
}

/* ServerBuilder */

type ServerBuilder interface {
	ParseCommandLine(args []string) error
	Build() (http.Server, error)
}

type serverBuilder struct {
	ContentRootDirectory string
	Port                 uint
}

func (builder *serverBuilder) ParseCommandLine(args []string) error {
	port := flag.Uint("p", 0, "The TCP port on which to listen")
	contentRootDirectory := flag.String("d", "", "The root content directory, from which to operate")

	flag.Parse()
	builder.ContentRootDirectory = *contentRootDirectory
	builder.Port = *port
	return nil
}

func (builder *serverBuilder) Build() (server http.Server, err error) {
	if builder.Port == 0 {
		return server, fmt.Errorf("gohttp: missing port")
	} else if builder.ContentRootDirectory == "" {
		return server, fmt.Errorf("gohttp: missing content root directory")
	}

	return http.NewServer(builder.ContentRootDirectory, builder.Port), nil
}
