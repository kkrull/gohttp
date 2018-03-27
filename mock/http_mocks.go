package mock

import (
	"fmt"
	. "github.com/onsi/gomega"
)

/* HttpServer */

type HttpServer struct {
	ListenFails  string
	listenCalled bool
}

func (mock *HttpServer) Listen() error {
	mock.listenCalled = true
	if mock.ListenFails != "" {
		return fmt.Errorf(mock.ListenFails)
	}

	return nil
}

func (mock *HttpServer) VerifyListen() {
	Expect(mock.listenCalled).To(BeTrue())
}
