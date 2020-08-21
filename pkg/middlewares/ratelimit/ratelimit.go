package ratelimit

import (
	"errors"
	"github.com/Highway-Project/highway/pkg/middlewares"
	"github.com/patrickmn/go-cache"
	"net"
	"net/http"
	"time"
)

type RateLimitMiddleware struct {
	db                *cache.Cache
	rateLimitValue    int
	strategyFunc      func(r *http.Request) string
	rateLimitDuration time.Duration
}

func (r *RateLimitMiddleware) Process(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rateLimitKey := r.strategyFunc(req)
		remainingReqs, found := r.db.Get(rateLimitKey)
		if !found {
			r.db.Set(rateLimitKey, r.rateLimitValue-1, cache.DefaultExpiration)
			handler.ServeHTTP(w, req)
		} else {
			remainingCount, ok := remainingReqs.(int)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`500 Internal Server error`))
				return
			}

			if remainingCount <= 0 {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`429 - Too Many Requests`))
				return
			}

			err := r.db.Decrement(rateLimitKey, 1)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`500 Internal Server error`))
				return
			}

			handler.ServeHTTP(w, req)

		}
	}
}

func New(params middlewares.MiddlewareParams) (middlewares.Middleware, error) {
	strategyObj, ok := params.Params["strategy"]
	if !ok {
		return nil, errors.New("strategy param is required for ratelimit middleware")
	}

	strategy, ok := strategyObj.(string)
	if !ok {
		return nil, errors.New("strategy param is not valid for ratelimit middleware")
	}

	var strategyFunc func(r *http.Request) string

	switch strategy {
	case "ip":
		strategyFunc = func(r *http.Request) string {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				return r.RemoteAddr
			}
			return ip
		}
	default:
		return nil, errors.New("invalid ratelimit strategy")
	}

	limitValueObj, ok := params.Params["limitValue"]
	if !ok {
		return nil, errors.New("limitValue param is required for ratelimit middleware")
	}

	limitValue, ok := limitValueObj.(int)
	if !ok {
		return nil, errors.New("limitValue param is not valid for ratelimit middleware")
	}

	if limitValue <= 0 {
		return nil, errors.New("limitValue param should be positive for ratelimit middleware")
	}

	limitDurationObj, ok := params.Params["limitDuration"]
	if !ok {
		return nil, errors.New("limitUnit param is required for ratelimit middleware")
	}

	limitDuration, ok := limitDurationObj.(string)
	if !ok {
		return nil, errors.New("limitUnit param is not valid for ratelimit middleware")
	}

	duration, err := time.ParseDuration(limitDuration)
	if err != nil {
		return nil, errors.New("limitDuration for limitrate middleware is not valid")
	}

	db := cache.New(duration, duration)

	return &RateLimitMiddleware{
		db:                db,
		strategyFunc:      strategyFunc,
		rateLimitValue:    limitValue,
		rateLimitDuration: duration,
	}, nil
}
