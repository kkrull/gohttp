package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kkrull/gohttp"
	"github.com/kkrull/gohttp/mock"
	"github.com/kkrull/gohttp/stub"
)

var _ = Describe("Main", func() {
	var (
		main    Main
		builder *stub.ServerBuilder
		server  *mock.HttpServer
	)

	Describe("Run", func() {
		var (
			err            error
			validArguments = []string{"-p", "4242", "-d", "/tmp"}
		)

		Context("given valid command-line arguments", func() {
			BeforeEach(func() {
				server = &mock.HttpServer{}
				builder = &stub.ServerBuilder{BuildReturns: server}
				main = Main{Builder: builder}
				err = main.Run(validArguments)
			})

			It("configures ServerBuilder from the command line", func() {
				builder.VerifyParseCommandLine(validArguments)
			})
			It("builds the HTTP server", func() {
				builder.VerifyBuild()
			})
			It("runs the server", func() {
				server.VerifyListen()
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when there is an error in parsing the command line", func() {
			It("returns the error", func() {
				builder = &stub.ServerBuilder{ParseCommandLineFails: "no parsing for you!"}
				main = Main{Builder: builder}
				err = main.Run(validArguments)
				Expect(err).To(MatchError("no parsing for you!"))
			})
		})

		Context("when there is an error in building the server", func() {
			It("returns the error", func() {
				builder = &stub.ServerBuilder{BuildFails: "bad arguments"}
				main = Main{Builder: builder}
				err = main.Run(validArguments)
				Expect(err).To(MatchError("bad arguments"))
			})
		})

		Context("when there is an error starting the server", func() {
			It("returns the error", func() {
				server = &mock.HttpServer{ListenFails: "where are my listening ears?"}
				builder = &stub.ServerBuilder{BuildReturns: server}
				main = Main{Builder: builder}
				err = main.Run(validArguments)
				Expect(err).To(MatchError("where are my listening ears?"))
			})
		})
	})
})
