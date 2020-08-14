package nothing

import (
	"github.com/Highway-Project/highway/pkg/middlewares"
	"net/http"
)

type NothingMiddleware struct{}

func (n NothingMiddleware) Process(handler http.HandlerFunc) http.HandlerFunc {
	return handler

}

func New(params middlewares.MiddlewareParams) middlewares.Middleware {
	return NothingMiddleware{}
}
