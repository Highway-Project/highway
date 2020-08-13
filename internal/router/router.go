package router

import (
	"errors"
	"github.com/Highway-Project/highway/pkg/router"
	"github.com/Highway-Project/highway/pkg/router/gorilla"
)

var routerConstructors map[string]func(options router.RouterOptions) (router.Router, error)

func init() {
	routerConstructors = make(map[string]func(options router.RouterOptions) (router.Router, error))
	_ = RegisterRouter("gorilla", gorilla.New)
}

func RegisterRouter(name string, constructor func(options router.RouterOptions) (router.Router, error)) error {
	if _, exists := routerConstructors[name]; exists {
		return errors.New("Router with this name exists: " + name)
	}

	routerConstructors[name] = constructor
	return nil
}

func NewRouter(name string, options router.RouterOptions) (router.Router, error) {
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
