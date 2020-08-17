package middlewares

import (
	"errors"
	"fmt"
	"github.com/Highway-Project/highway/logging"
	"net/http"
)

type Middleware interface {
	Process(handler http.HandlerFunc) http.HandlerFunc
}

type MiddlewareParams struct {
	Params map[string]interface{}
}

func (mp *MiddlewareParams) GetList(key string) (res []string, exists bool, err error) {
	v, exists := mp.Params[key]
	if !exists {
		return nil, exists, nil
	}

	_, ok := v.([]interface{})
	if !ok {
		msg := fmt.Sprintf("%s must be of type []string", key)
		logging.Logger.WithError(err).Error(msg)
		return nil, exists, errors.New(msg)
	}

	res = make([]string, len(v.([]interface{})))
	for i, value := range v.([]interface{}) {
		res[i], ok = value.(string)
		if !ok {
			msg := fmt.Sprintf("%s must be of type []string", key)
			logging.Logger.WithError(err).Error(msg)
			return nil, exists, errors.New(msg)
		}
	}

	return res, exists, nil
}

func (mp *MiddlewareParams) GetBool(key string) (res bool, exists bool, err error) {
	v, exists := mp.Params[key]
	if !exists {
		return false, exists, err
	}

	res, ok := v.(bool)
	if !ok {
		msg := fmt.Sprintf("%s must be of type bool", key)
		logging.Logger.WithError(err).Error(msg)
		return false, exists, errors.New(msg)
	}

	return res, exists, nil
}

func (mp *MiddlewareParams) GetInt(key string) (res int, exists bool, err error) {
	v, exists := mp.Params[key]
	if !exists {
		return 0, exists, nil
	}

	res, ok := v.(int)
	if !ok {
		msg := fmt.Sprintf("%s must be of type int", key)
		logging.Logger.WithError(err).Error(msg)
		return 0, exists, errors.New(msg)
	}

	return res, exists, nil
}
