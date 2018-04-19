package cmd_test

import (
	"flag"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/kkrull/gohttp/main/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMainCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "main/cmd")
}

/* Synchronization */

func waitForStart() {
	time.Sleep(100 * time.Millisecond)
}

func scheduleShutdown(quit chan bool) {
	waitForStart()
	quit <- true
}

/* CliCommand */

type CliCommandMock struct{}

func (mock *CliCommandMock) Run(stderr io.Writer) (code int, err error) {
	return -1, nil
}

/* CommandFactoryMock */

type CommandFactoryMock struct {
	ErrorCommandReturns  *CliCommandMock
	errorCommandReceived error

	HelpCommandReturns  *CliCommandMock
	helpCommandReceived *flag.FlagSet

	RunCommandReturns *CliCommandMock

	tcpServerReceivedPath string
	tcpServerReceivedHost string
	tcpServerReceivedPort uint16
}

func (mock *CommandFactoryMock) ErrorCommand(err error) cmd.CliCommand {
	mock.errorCommandReceived = err
	return mock.ErrorCommandReturns
}

func (mock *CommandFactoryMock) ErrorCommandShouldHaveReceived(err error) {
	ExpectWithOffset(1, mock.errorCommandReceived).To(BeEquivalentTo(err))
}

func (mock *CommandFactoryMock) HelpCommand(flagSet *flag.FlagSet) cmd.CliCommand {
	mock.helpCommandReceived = flagSet
	return mock.HelpCommandReturns
}

func (mock *CommandFactoryMock) HelpCommandShouldBeForProgram(name string) {
	ExpectWithOffset(1, mock.helpCommandReceived).NotTo(BeNil())
	ExpectWithOffset(1, mock.helpCommandReceived.Name()).To(Equal(name))
}

func (mock *CommandFactoryMock) HelpCommandShouldHaveFlag(flagName string, usage string) {
	ExpectWithOffset(1, mock.helpCommandReceived).NotTo(BeNil())
	ExpectWithOffset(1, mock.helpCommandReceived.Lookup(flagName).Usage).To(Equal(usage))
}

func (mock *CommandFactoryMock) RunCommand(server cmd.Server) (command cmd.CliCommand, quit chan bool) {
	return mock.RunCommandReturns, nil
}

func (mock *CommandFactoryMock) TCPServer(contentBasePath string, host string, port uint16) cmd.Server {
	mock.tcpServerReceivedPath = contentBasePath
	mock.tcpServerReceivedHost = host
	mock.tcpServerReceivedPort = port
	return nil
}

func (mock *CommandFactoryMock) TCPServerShouldHaveReceived(contentRootPath string, host string, port uint16) {
	ExpectWithOffset(1, mock.tcpServerReceivedPath).To(Equal(contentRootPath))
	ExpectWithOffset(1, mock.tcpServerReceivedHost).To(Equal(host))
	ExpectWithOffset(1, mock.tcpServerReceivedPort).To(Equal(port))
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
