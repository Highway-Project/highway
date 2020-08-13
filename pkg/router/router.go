package router

import (
	"github.com/Highway-Project/highway/pkg/rules"
	"net/http"
)

type Router interface {
	AddRule(rule rules.Rule) error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type RouterOptions struct {
	Options map[string]string
}
