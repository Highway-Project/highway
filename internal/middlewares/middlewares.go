package middlewares

import (
	"errors"
	"fmt"
	"github.com/Highway-Project/highway/config"
	"github.com/Highway-Project/highway/logging"
	"github.com/Highway-Project/highway/pkg/middlewares"
	"github.com/Highway-Project/highway/pkg/middlewares/nothing"
)

var middlewareConstructors map[string]func(middlewares.MiddlewareParams) (middlewares.Middleware, error)
var middlewareMap map[string]middlewares.Middleware

func init() {
	middlewareConstructors = make(map[string]func(params middlewares.MiddlewareParams) (middlewares.Middleware, error))
	middlewareMap = make(map[string]middlewares.Middleware)
	_ = RegisterMiddleware("nothing", nothing.New)
}

func RegisterMiddleware(name string, constructor func(params middlewares.MiddlewareParams) (middlewares.Middleware, error)) error {
	if _, exists := middlewareConstructors[name]; exists {
		return errors.New("Middleware with this name exists: " + name)
	}

	middlewareConstructors[name] = constructor
	return nil
}

func LoadMiddlewares(specs []config.MiddlewareSpec) error {
	for _, spec := range specs {
		if spec.CustomMiddleware {
			// TODO: load custom middleware
			continue
		} else {
			constructor, exists := middlewareConstructors[spec.MiddlewareName]
			if !exists {
				msg := fmt.Sprintf("could not load middleware %s", spec.MiddlewareName)
				logging.Logger.Errorf(msg)
				return errors.New(msg)
			}

			var refName string
			if spec.RefName != "" {
				refName = spec.RefName
			} else {
				refName = spec.MiddlewareName
			}

			_, exists = middlewareMap[refName]
			if exists {
				msg := fmt.Sprintf("middleware %s already exists", refName)
				logging.Logger.Errorf(msg)
				return errors.New(msg)
			}

			mw, err := constructor(middlewares.MiddlewareParams{Params: spec.Params})
			if err != nil {
				msg := fmt.Sprintf("could not create middlware %s", refName)
				logging.Logger.Errorf(msg)
				return errors.New(msg)
			}

			middlewareMap[refName] = mw
		}
	}

	return nil
}

func GetMiddlewareByName(refName string) (middlewares.Middleware, error) {
	mw, exists := middlewareMap[refName]
	if !exists {
		msg := fmt.Sprintf("middleware %s does not exist", refName)
		logging.Logger.Errorf(msg)
		return nil, errors.New(msg)
	}
	return mw, nil
}
