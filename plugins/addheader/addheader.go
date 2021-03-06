package main

import (
	"net/http"
)

type AddHeaderMiddleware struct {
	key string
	val string
}

func (a AddHeaderMiddleware) Process(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(a.key, a.val)
		handler(w, r)
	}
}

func New(params map[string]interface{}) (interface{}, error) {
	return AddHeaderMiddleware{
		key: params["key"].(string),
		val: params["val"].(string),
	}, nil
}
