package coffee

import "github.com/kkrull/gohttp/http"

func NewRoute() http.Route {
	return &coffeeRoute{}
}

type coffeeRoute struct {
}

func (route *coffeeRoute) Route(requested *http.RequestLine) http.Request {
	if requested.Target == "/coffee" {
		return &GetCoffeeRequest{}
	}

	return nil
}
