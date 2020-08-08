package middlewares

import "net/http"

type Middleware interface {
	Process(handler http.Handler) http.Handler
}

type MiddlewareParams struct {
	Params map[string]string
}
