package http_test

import (
	"bytes"
	"io/ioutil"
	"net"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	server       *http.TCPServer
	conn         net.Conn
	connectError error
)

var _ = Describe("TCPServer", func() {
	AfterEach(func(done Done) {
		if conn != nil {
			Expect(conn.Close()).To(Succeed())
			conn = nil
			connectError = nil
		}

		Expect(server.Shutdown()).To(Succeed())
		server = nil
		close(done)
	})

	Describe("#Address", func() {
		Context("when the server is not running", func() {
			BeforeEach(func(done Done) {
				server = http.MakeTCPServerOnAvailablePort("localhost")
				Expect(server.Start()).To(Succeed())
				Expect(server.Shutdown()).To(Succeed())
				close(done)
			})

			It("returns nil", func() {
				Expect(server.Address()).To(BeNil())
			})
		})

		Context("when the server is running", func() {
			BeforeEach(func(done Done) {
				server = http.MakeTCPServer("localhost", 8420)
				Expect(server.Start()).To(Succeed())
				close(done)
			})

			It("reports the bound TCP address", func() {
				var netAddress = server.Address()
				var tcpAddress = netAddress.(*net.TCPAddr)
				Expect(tcpAddress.IP.String()).NotTo(Equal(""))
				Expect(tcpAddress.Port).To(Equal(8420))
			})
		})
	})

	Describe("#Start", func() {
		Context("when the server is already running", func() {
			BeforeEach(func(done Done) {
				server = http.MakeTCPServerOnAvailablePort("localhost")
				Expect(server.Start()).To(Succeed())
				close(done)
			})

			It("returns an error", func(done Done) {
				Expect(server.Start()).To(MatchError("TCPServer: already running"))
				close(done)
			})
		})

		Context("when there is an error resolving the given host and port to an address", func() {
			It("immediately returns an error", func(done Done) {
				invalidHostAddress := "666.666.666.666"
				server = http.MakeTCPServerOnAvailablePort(invalidHostAddress)
				Expect(server.Start()).To(MatchError(HaveSuffix("no such host")))
				close(done)
			}, 5)
		})

		Context("when there is an error binding to resolved address", func() {
			It("immediately returns the error", func(done Done) {
				var portOnlyAvailableToRoot uint16 = 1
				server = http.MakeTCPServer("localhost", portOnlyAvailableToRoot)
				Expect(server.Start()).To(MatchError(HaveSuffix("bind: permission denied")))
				close(done)
			})
		})

		Context("given an available host address and port", func() {
			It("returns no error as soon as the socket is open", func(done Done) {
				server = http.MakeTCPServerOnAvailablePort("localhost")
				Expect(server.Start()).To(Succeed())
				close(done)
			})

			It("accepts connections on the specified address", func(done Done) {
				server = http.MakeTCPServer("localhost", 8421)
				Expect(server.Start()).To(Succeed())

				conn, connectError = net.Dial("tcp", "localhost:8421")
				Expect(connectError).NotTo(HaveOccurred())
				close(done)
			})
		})

		Context("given no port number", func() {
			BeforeEach(func(done Done) {
				server = http.MakeTCPServerOnAvailablePort("localhost")
				Expect(server.Start()).To(Succeed())
				close(done)
			})

			It("listens on a random, available port", func() {
				var netAddress = server.Address()
				var tcpAddress = netAddress.(*net.TCPAddr)
				Expect(tcpAddress.Port).To(BeNumerically(">", 0))
				Expect(tcpAddress.Port).To(BeNumerically("<=", 1<<16))
			})
		})

		Context("when the server is running", func() {
			BeforeEach(func() {
				server = http.MakeTCPServerOnAvailablePort("localhost")
				Expect(server.Start()).To(Succeed())
			})

			It("responds to HTTP requests", func(done Done) {
				for i := 0; i < 2; i++ {
					conn, connectError = net.Dial("tcp", server.Address().String())
					Expect(connectError).NotTo(HaveOccurred())
					writeString(conn, "GET / HTTP/1.1\r\n\r\n")
					expectHttpResponse(conn)
				}

				close(done)
			})
		})
	})

	Describe("#Shutdown", func() {
		Context("when the server has not been started", func() {
			It("returns no error", func() {
				server = http.MakeTCPServerOnAvailablePort("localhost")
				Expect(server.Shutdown()).To(Succeed())
			})
		})

		Context("when the server is running", func() {
			It("stops accepting connections", func(done Done) {
				server = http.MakeTCPServerOnAvailablePort("localhost")
				Expect(server.Start()).To(Succeed())
				oldAddress := server.Address().String()

				Expect(server.Shutdown()).To(Succeed())
				conn, connectError = net.Dial("tcp", oldAddress)
				Expect(connectError).To(HaveOccurred())
				close(done)
			})
		})

		Context("when the server is stopped", func() {
			BeforeEach(func(done Done) {
				server = http.MakeTCPServerOnAvailablePort("localhost")
				Expect(server.Start()).To(Succeed())
				Expect(server.Shutdown()).To(Succeed())
				close(done)
			})

			It("returns no error", func(done Done) {
				Expect(server.Shutdown()).To(Succeed())
				close(done)
			})
		})
	})

	Describe("when running", func() {
		var handler *HandlerMock

		Context("when it receives a request", func() {
			BeforeEach(func(done Done) {
				handler = &HandlerMock{}
				server = &http.TCPServer{
					Handler: handler,
					Host:    "localhost"}

				Expect(server.Start()).To(Succeed())
				conn, connectError = net.Dial("tcp", server.Address().String())
				Expect(connectError).NotTo(HaveOccurred())
				close(done)
			})

			It("passes the request to the configured ConnectionHandler", func(done Done) {
				writeString(conn, "GET / HTTP/1.1\r\n\r\n")
				readString(conn)
				handler.ShouldHandleConnection()
				close(done)
			})
		})
	})
})

func expectHttpResponse(conn net.Conn) {
	rfc7230StatusLinePattern := "^HTTP/1[.]1 \\d{3} [\\w\\s]+[\r][\n]"
	Expect(readString(conn)).To(MatchRegexp(rfc7230StatusLinePattern))
}

func readString(conn net.Conn) (string, error) {
	readBytes, err := ioutil.ReadAll(conn)
	return string(readBytes), err
}

func writeString(conn net.Conn, s string) {
	buffer := bytes.NewBufferString(s)
	Expect(conn.Write(buffer.Bytes())).To(Equal(len(s)))
}
