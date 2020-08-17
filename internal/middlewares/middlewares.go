package middlewares

import (
	"errors"
	"fmt"
	"github.com/Highway-Project/highway/config"
	"github.com/Highway-Project/highway/logging"
	"github.com/Highway-Project/highway/pkg/middlewares"
	"github.com/Highway-Project/highway/pkg/middlewares/cors"
	"github.com/Highway-Project/highway/pkg/middlewares/nothing"
	"plugin"
)

var middlewareConstructors map[string]func(middlewares.MiddlewareParams) (middlewares.Middleware, error)
var middlewareMap map[string]middlewares.Middleware

func init() {
	middlewareConstructors = make(map[string]func(params middlewares.MiddlewareParams) (middlewares.Middleware, error))
	middlewareMap = make(map[string]middlewares.Middleware)
	_ = RegisterMiddleware("nothing", nothing.New)
	_ = RegisterMiddleware("cors", cors.New)
}

func RegisterMiddleware(name string, constructor func(params middlewares.MiddlewareParams) (middlewares.Middleware, error)) error {
	if _, exists := middlewareConstructors[name]; exists {
		return errors.New("Middleware with this name exists: " + name)
	}

	middlewareConstructors[name] = constructor
	return nil
}

func loadCustomMiddleware(spec config.MiddlewareSpec) error {
	plug, err := plugin.Open(spec.MiddlewarePath)
	if err != nil {
		msg := fmt.Sprintf("could not open custom middleware %s's file", spec.MiddlewareName)
		logging.Logger.WithError(err).Errorf(msg)
		return errors.New(msg)
	}
	constructorSym, err := plug.Lookup("New")
	if err != nil {
		msg := fmt.Sprintf("could not load middleware %s's constructor. New function is not found", spec.MiddlewareName)
		logging.Logger.WithError(err).Errorf(msg)
		return errors.New(msg)
	}
	constructor, ok := constructorSym.(func(map[string]interface{}) (interface{}, error))
	if !ok {
		msg := fmt.Sprintf("New function for middleware %s is not valid", spec.MiddlewareName)
		logging.Logger.WithError(err).Errorf(msg)
		return errors.New(msg)
	}

	var refName string
	if spec.RefName != "" {
		refName = spec.RefName
	} else {
		refName = spec.MiddlewareName
	}

	_, exists := middlewareMap[refName]
	if exists {
		msg := fmt.Sprintf("middleware with name %s already exists", refName)
		logging.Logger.Errorf(msg)
		return errors.New(msg)
	}

	mwInterface, err := constructor(spec.Params)
	if err != nil {
		msg := fmt.Sprintf("could not create middlware %s", refName)
		logging.Logger.Errorf(msg)
		return errors.New(msg)
	}

	mw, ok := mwInterface.(middlewares.Middleware)
	if !ok {
		msg := fmt.Sprintf("output of New function middlware %s is not implementing midlewares.Middleware Interface", refName)
		logging.Logger.Errorf(msg)
		return errors.New(msg)
	}

	middlewareMap[refName] = mw
	return nil
}

func loadBuiltinMiddleware(spec config.MiddlewareSpec) error {
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
	return nil
}

func LoadMiddlewares(specs []config.MiddlewareSpec) error {
	for _, spec := range specs {
		if spec.CustomMiddleware {
			err := loadCustomMiddleware(spec)
			if err != nil {
				return err
			}
		} else {
			err := loadBuiltinMiddleware(spec)
			if err != nil {
				return err
			}
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
