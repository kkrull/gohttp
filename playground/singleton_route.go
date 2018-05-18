package playground

import "github.com/kkrull/gohttp/http"

func NewSingletonRoute(path string) *SingletonRoute {
	return &SingletonRoute{
		Path:      path,
		Singleton: &SingletonResource{},
	}
}

type SingletonRoute struct {
	Path      string
	Singleton *SingletonResource
}

func (route *SingletonRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Singleton)
}

type SingletonResource struct {
}

func (singleton *SingletonResource) Name() string {
	return "Singleton"
}
