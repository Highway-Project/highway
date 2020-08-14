package middlewares

import "net/http"

type Middleware interface {
	Process(handler http.HandlerFunc) http.HandlerFunc
}

type MiddlewareParams struct {
	Params map[string]string
}
