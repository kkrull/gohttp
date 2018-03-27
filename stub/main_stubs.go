package stub

import (
	"fmt"
	. "github.com/onsi/gomega"

	"github.com/kkrull/gohttp/http"
)

type ServerBuilder struct {
	buildCalled  bool
	BuildFails   string
	BuildReturns http.Server

	parseCommandLineArgs  []string
	ParseCommandLineFails string
}

func (stub *ServerBuilder) Build() (server http.Server, err error) {
	stub.buildCalled = true
	if stub.BuildFails == "" {
		server = stub.BuildReturns
	} else {
		err = fmt.Errorf(stub.BuildFails)
	}

	return
}

func (stub *ServerBuilder) VerifyBuild() {
	Expect(stub.buildCalled).To(Equal(true))
}

func (stub *ServerBuilder) ParseCommandLine(args []string) error {
	stub.parseCommandLineArgs = args
	if stub.ParseCommandLineFails == "" {
		return nil
	} else {
		return fmt.Errorf(stub.ParseCommandLineFails)
	}
}

func (stub *ServerBuilder) VerifyParseCommandLine(expected []string) {
	Expect(stub.parseCommandLineArgs).To(Equal(expected))
}
