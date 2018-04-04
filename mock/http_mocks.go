package mock

import (
	"fmt"
	. "github.com/onsi/gomega"
	"net"
)

type Server struct {
	StartFails  string
	startCalled bool

	ShutdownFails  string
	shutdownCalled bool
}

func (Server) Address() net.Addr {
	panic("implement me")
}

func (mock *Server) Start() error {
	mock.startCalled = true
	if mock.StartFails != "" {
		return fmt.Errorf(mock.StartFails)
	}

	return nil
}

func (mock Server) VerifyStart() {
	Expect(mock.startCalled).To(BeTrue())
}

func (mock *Server) Shutdown() error {
	mock.shutdownCalled = true
	if mock.ShutdownFails != "" {
		return fmt.Errorf(mock.ShutdownFails)
	}

	return nil
}

func (mock Server) VerifyRunning() {
	Expect(mock.startCalled).To(BeTrue())
	Expect(mock.shutdownCalled).To(BeFalse())
}

func (mock Server) VerifyShutdown() {
	Expect(mock.shutdownCalled).To(BeTrue())
}
