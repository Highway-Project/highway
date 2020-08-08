package router

import (
	"errors"
	"github.com/Highway-Project/highway/pkg/rules"
	"net/http"
)

var routerConstructors map[string]func(options RouterOptions) (Router, error)

func init() {
	routerConstructors = make(map[string]func(options RouterOptions) (Router, error))
}

type Router interface {
	AddRule(rule rules.Rule) error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type RouterOptions struct {
	Options map[string]string
}

func RegisterRouter(name string, constructor func(options RouterOptions) (Router, error)) error {
	if _, exists := routerConstructors[name]; exists {
		return errors.New("Router with this name exists: " + name)
	}

	routerConstructors[name] = constructor
	return nil
}

func NewRouter(name string, options RouterOptions) (Router, error) {
	constructor, exists := routerConstructors[name]
	if !exists {
		return nil, errors.New("Router with this name does not exists: " + name)
	}
	r, err := constructor(options)
	if err != nil {
		return nil, err
	}
	return r, nil
}
