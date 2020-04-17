package middlewares

import "net/http"

type Middleware interface {
	Process(handler http.Handler)  http.Handler
}