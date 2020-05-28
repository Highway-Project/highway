package nothing

import (
	"github.com/fiust/highway/pkg/middlewares"
	"net/http"
)

type NothingMiddleware struct {}

func (n NothingMiddleware) Process(handler http.Handler) http.Handler {
	return handler
}

func New(params middlewares.MiddlewareParams) middlewares.Middleware {
	return NothingMiddleware{}
}