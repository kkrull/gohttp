package cmd_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cmd Suite")
}

/* Synchronization */

func waitForStart() {
	time.Sleep(100 * time.Millisecond)
}

func scheduleShutdown(quit chan bool) {
	waitForStart()
	quit <- true
}

/* ServerMock */

type ServerMock struct {
	StartFails  string
	startCalled bool

	ShutdownFails  string
	shutdownCalled bool
}

func (ServerMock) Address() net.Addr {
	panic("implement me")
}

func (mock *ServerMock) Start() error {
	mock.startCalled = true
	if mock.StartFails != "" {
		return fmt.Errorf(mock.StartFails)
	}

	return nil
}

func (mock ServerMock) VerifyStart() {
	Expect(mock.startCalled).To(BeTrue())
}

func (mock *ServerMock) Shutdown() error {
	mock.shutdownCalled = true
	if mock.ShutdownFails != "" {
		return fmt.Errorf(mock.ShutdownFails)
	}

	return nil
}

func (mock ServerMock) VerifyRunning() {
	Expect(mock.startCalled).To(BeTrue())
	Expect(mock.shutdownCalled).To(BeFalse())
}

func (mock ServerMock) VerifyShutdown() {
	Expect(mock.shutdownCalled).To(BeTrue())
}
