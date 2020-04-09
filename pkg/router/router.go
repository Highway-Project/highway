package router

import (
	"github.com/fiust/highway/pkg/rules"
	"net/http"
)

type Router interface {
	AddRoute(rule rules.Rule) error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}