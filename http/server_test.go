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
	server *http.TCPServer
	conn   net.Conn
	err    error
)

var _ = Describe("TCPServer", func() {
	AfterEach(func(done Done) {
		if conn != nil {
			Expect(conn.Close()).To(Succeed())
			conn = nil
			err = nil
		}

		Expect(server.Shutdown()).To(Succeed())
		server = nil
		close(done)
	})

	Describe("#Address", func() {
		Context("when the server is not running", func() {
			BeforeEach(func(done Done) {
				server = http.TCPServerBuilder("localhost").Build()
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
				server = http.TCPServerBuilder("localhost").
					ListeningOnPort(8420).
					Build()
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
				server = http.TCPServerBuilder("localhost").Build()
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
				const invalidHostAddress = "666.666.666.666"
				server = http.TCPServerBuilder(invalidHostAddress).Build()
				Expect(server.Start()).To(MatchError(HaveSuffix("no such host")))
				close(done)
			}, 5)
		})

		Context("when there is an error binding to resolved address", func() {
			It("immediately returns the error", func(done Done) {
				const portOnlyAvailableToRoot uint16 = 1
				server = http.TCPServerBuilder("localhost").
					ListeningOnPort(portOnlyAvailableToRoot).
					Build()
				Expect(server.Start()).To(MatchError(HaveSuffix("bind: permission denied")))
				close(done)
			})
		})

		Context("given an available host address and port", func() {
			It("returns no error as soon as the socket is open", func(done Done) {
				server = http.TCPServerBuilder("localhost").Build()
				Expect(server.Start()).To(Succeed())
				close(done)
			})

			It("accepts connections on the specified address", func(done Done) {
				server = http.TCPServerBuilder("localhost").
					ListeningOnPort(8421).
					Build()
				Expect(server.Start()).To(Succeed())

				conn, err = net.Dial("tcp", "localhost:8421")
				Expect(err).NotTo(HaveOccurred())
				close(done)
			})
		})

		Context("given no port number", func() {
			BeforeEach(func(done Done) {
				server = http.TCPServerBuilder("localhost").Build()
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
				server = http.TCPServerBuilder("localhost").Build()
				Expect(server.Start()).To(Succeed())
			})

			It("responds to HTTP requests", func(done Done) {
				for i := 0; i < 2; i++ {
					conn, err = net.Dial("tcp", server.Address().String())
					Expect(err).NotTo(HaveOccurred())
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
				server = http.TCPServerBuilder("localhost").Build()
				Expect(server.Shutdown()).To(Succeed())
			})
		})

		Context("when the server is running", func() {
			It("stops accepting connections", func(done Done) {
				server = http.TCPServerBuilder("localhost").Build()
				Expect(server.Start()).To(Succeed())
				oldAddress := server.Address().String()

				Expect(server.Shutdown()).To(Succeed())
				conn, err = net.Dial("tcp", oldAddress)
				Expect(err).To(HaveOccurred())
				close(done)
			})
		})

		Context("when the server is stopped", func() {
			BeforeEach(func(done Done) {
				server = http.TCPServerBuilder("localhost").Build()
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
		Context("when it receives a request", func() {
			var handler *HandlerMock

			BeforeEach(func(done Done) {
				handler = &HandlerMock{}
				server = http.TCPServerBuilder("localhost").
					WithConnectionHandler(handler).
					Build()

				Expect(server.Start()).To(Succeed())
				conn, err = net.Dial("tcp", server.Address().String())
				Expect(err).NotTo(HaveOccurred())
				close(done)
			})

			It("passes the request to the configured ConnectionHandler", func(done Done) {
				writeString(conn, "GET / HTTP/1.1\r\n\r\n")
				readString(conn)
				handler.ShouldHandleConnection()
				close(done)
			})
		})

		Context("when the given concurrency limit is 2 or more", func() {
			BeforeEach(func(done Done) {
				server = http.TCPServerBuilder("localhost").
					WithMaxConnections(2).
					Build()
				Expect(server.Start()).To(Succeed())
				close(done)
			})

			FIt("can handle multiple connections at a time", func(done Done) {
				slowConn := dial(server)
				fastConn := dial(server)
				writeString(fastConn, "GET / HTTP/1.1\r\n\r\n")
				readString(fastConn)
				Expect(fastConn.Close()).To(Succeed())

				writeString(slowConn, "GET / HTTP/1.1\r\n\r\n")
				readString(slowConn)
				Expect(slowConn.Close()).To(Succeed())

				close(done)
			})
		})
	})
})

func dial(server *http.TCPServer) net.Conn {
	conn, err = net.Dial("tcp", server.Address().String())
	Expect(err).NotTo(HaveOccurred())
	return conn
}

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
