package nothing

import (
	"github.com/Highway-Project/highway/pkg/middlewares"
	"net/http"
)

type NothingMiddleware struct {
	msg string
}

func (n NothingMiddleware) Process(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(n.msg, n.msg)
		handler(w, r)
	}
}

func New(params middlewares.MiddlewareParams) (middlewares.Middleware, error) {
	return NothingMiddleware{msg: params.Params["msg"].(string)}, nil
}
