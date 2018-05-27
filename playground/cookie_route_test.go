package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	setTypePath  = "/set"
	readTypePath = "/read"
)

var _ = Describe("::NewCookieRoute", func() {
	var route *playground.CookieRoute
	BeforeEach(func() {
		route = playground.NewCookieRoute("/set", "/read")
	})

	It("returns a Route at the given paths", func() {
		Expect(route).NotTo(BeNil())
		Expect(route).To(BeAssignableToTypeOf(&playground.CookieRoute{}))
		Expect(route.SetTypePath).To(Equal("/set"))
		Expect(route.ReadTypePath).To(Equal("/read"))
	})

	It("configures a CookieMonster", func() {
		Expect(route.Monster).NotTo(BeNil())
		Expect(route.Monster).To(BeAssignableToTypeOf(&playground.CookieMonster{}))
	})

	It("configures a CookieRegistrar", func() {
		Expect(route.Registrar).NotTo(BeNil())
		Expect(route.Registrar).To(BeAssignableToTypeOf(&playground.CookieRegistrar{}))
	})
})

var _ = Describe("CookieRoute", func() {
	var (
		router    http.Route
		registrar *playground.CookieRegistrar
		monster   *playground.CookieMonster
		response  = &bytes.Buffer{}
	)

	Describe("#Route", func() {
		BeforeEach(func() {
			router = &playground.CookieRoute{
				SetTypePath: setTypePath,
				Registrar:   registrar,

				ReadTypePath: readTypePath,
				Monster:      monster,
			}
			response.Reset()
		})

		Context("when the path is the configured path for setting the type of cookie", func() {
			Context("when the method is OPTIONS", func() {
				BeforeEach(func() {
					requested := http.NewOptionsMessage(setTypePath)
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("allows OPTIONS", httptest.ShouldAllowMethods(response, http.OPTIONS))
				It("allows GET", httptest.ShouldAllowMethods(response, http.GET))
			})

			It("replies Method Not Allowed on any other method", func() {
				requested := http.NewTraceMessage(setTypePath)
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeAssignableToTypeOf(clienterror.MethodNotAllowed()))
			})
		})

		Context("when the path is the configured path for reading the type of cookie", func() {
			Context("when the method is OPTIONS", func() {
				BeforeEach(func() {
					requested := http.NewOptionsMessage(readTypePath)
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("allows OPTIONS", httptest.ShouldAllowMethods(response, http.OPTIONS))
				It("allows GET", httptest.ShouldAllowMethods(response, http.GET))
			})

			It("replies Method Not Allowed on any other method", func() {
				requested := http.NewTraceMessage(readTypePath)
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeAssignableToTypeOf(clienterror.MethodNotAllowed()))
			})
		})

		It("passes on any other path by returning nil", func() {
			requested := http.NewGetMessage("/")
			Expect(router.Route(requested)).To(BeNil())
		})
	})
})

var _ = Describe("CookieMonster", func() {
	var (
		monster  *playground.CookieMonster
		response *httptest.ResponseMessage
	)

	Describe("#Get", func() {
		Context("given 1 Cookie header", func() {
			BeforeEach(func() {
				monster = &playground.CookieMonster{}
				request := &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   readTypePath,
					TargetReturns: readTypePath,
				}
				request.AddHeader("Cookie", "earwax")
				response = invokeResourceMethod(monster.Get, request)
			})

			It("responds 200 OK", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(200, "OK")
			})
			It("sets Content-Type to 'text/plain'", func() {
				response.HeaderShould("Content-Type", Equal("text/plain"))
			})
			It("sets Content-Length to the number of bytes in the message body", func() {
				Expect(response.HeaderAsInt("Content-Length")).To(BeNumerically(">", 0))
			})
			It("expresses satisfaction for the specified type of cookie", func() {
				response.BodyShould(Equal("mmmm earwax"))
			})
		})

		Context("given no Cookie header", func() {
			BeforeEach(func() {
				monster = &playground.CookieMonster{}
				request := &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   readTypePath,
					TargetReturns: readTypePath,
				}
				response = invokeResourceMethod(monster.Get, request)
			})

			It("responds 400 Bad Request", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(400, "Bad Request")
			})
		})

		Context("given 2 or more Cookie headers", func() {
			BeforeEach(func() {
				monster = &playground.CookieMonster{}
				request := &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   readTypePath,
					TargetReturns: readTypePath,
				}
				request.AddHeader("Cookie", "chocolate")
				request.AddHeader("Cookie", "wat")
				response = invokeResourceMethod(monster.Get, request)
			})

			It("responds 400 Bad Request", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(400, "Bad Request")
			})
		})
	})
})

var _ = Describe("CookieRegistrar", func() {
	var (
		registrar *playground.CookieRegistrar
		response  *httptest.ResponseMessage
	)

	BeforeEach(func() {
		registrar = &playground.CookieRegistrar{}
	})

	Describe("#Get", func() {
		Context("given a 'type' query parameter", func() {
			BeforeEach(func() {
				request := &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   setTypePath,
					TargetReturns: setTypePath,
				}
				request.AddQueryParameter("type", "Snickerdoodle")
				response = invokeResourceMethod(registrar.Get, request)
			})

			It("responds 200 OK", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(200, "OK")
			})
			It("Sets a Set-Cookie header", func() {
				response.HeaderShould("Set-Cookie", Not(BeEmpty()))
			})
			It("sets Content-Type to 'text/plain'", func() {
				response.HeaderShould("Content-Type", Equal("text/plain"))
			})
			It("sets Content-Length to the number of bytes in the message body", func() {
				Expect(response.HeaderAsInt("Content-Length")).To(BeNumerically(">", 0))
			})
			It("acknowledges the registration in the message body", func() {
				response.BodyShould(Equal("Eat a Snickerdoodle."))
			})
		})

		Context("given no 'type' parameter", func() {
			BeforeEach(func() {
				request := &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   setTypePath,
					TargetReturns: setTypePath,
				}
				response = invokeResourceMethod(registrar.Get, request)
			})

			It("responds 400 Bad Request", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(400, "Bad Request")
			})
			It("has no body", func() {
				response.BodyShould(BeEmpty())
			})
		})

		Context("given 2 or more 'type' parameters", func() {
			BeforeEach(func() {
				request := &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   setTypePath,
					TargetReturns: setTypePath,
				}
				request.AddQueryParameter("type", "HighlanderCookie")
				request.AddQueryParameter("type", "TheKurganCookie")
				response = invokeResourceMethod(registrar.Get, request)
			})

			It("responds 400 Bad Request", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(400, "Bad Request")
			})
			It("has no body", func() {
				response.BodyShould(BeEmpty())
			})
		})
	})
})
