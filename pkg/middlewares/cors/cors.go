package cors

import (
	"github.com/Highway-Project/highway/pkg/middlewares"
	"net/http"
	"strconv"
	"strings"
)

type CORSMiddleware struct {

	// AllowOrigin defines a list of origins that may access the resource.
	// Optional. Default value []string{"*"}.
	AllowOrigins []string

	// AllowMethods defines a list methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	// Optional. Default value DefaultCORSConfig.AllowMethods.
	AllowMethods []string

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This is in response to a preflight request.
	// Optional. Default value []string{}.
	AllowHeaders []string

	// AllowCredentials indicates whether or not the response to the request
	// can be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether or not the
	// actual request can be made using credentials.
	// Optional. Default value false.
	AllowCredentials bool

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	// Optional. Default value []string{}.
	ExposeHeaders []string

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	// Optional. Default value 0.
	MaxAge int
}

const (
	HeaderVary                          = "Vary"
	HeaderOrigin                        = "Origin"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"
)

const (
	AllowOrigins     = "allowOrigins"
	AllowMethods     = "allowMethods"
	AllowHeaders     = "allowHeaders"
	AllowCredentials = "allowCredentials"
	ExposeHeaders    = "exposeHeaders"
	MaxAge           = "maxAge"
)

var (
	// DefaultCORSConfig is the default CORS middleware config.
	DefaultCORSParams = map[string][]string{
		AllowOrigins: {"*"},
		AllowMethods: {http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
)

func New(params middlewares.MiddlewareParams) (middlewares.Middleware, error) {
	allowOrigins, exists, err := params.GetStringList(AllowOrigins)
	if err != nil {
		return CORSMiddleware{}, err
	}
	if !exists {
		allowOrigins = DefaultCORSParams[AllowOrigins]
	}

	allowMethods, exists, err := params.GetStringList(AllowMethods)
	if err != nil {
		return CORSMiddleware{}, err
	}
	if !exists {
		allowMethods = DefaultCORSParams[AllowMethods]
	}

	allowHeaders, _, err := params.GetStringList(AllowHeaders)
	if err != nil {
		return CORSMiddleware{}, err
	}

	allowCredentials, _, err := params.GetBool(AllowCredentials)
	if err != nil {
		return CORSMiddleware{}, err
	}

	exposeHeaders, _, err := params.GetStringList(ExposeHeaders)
	if err != nil {
		return CORSMiddleware{}, err
	}

	maxAge, _, err := params.GetInt(MaxAge)
	if err != nil {
		return CORSMiddleware{}, err
	}

	mw := CORSMiddleware{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: allowCredentials,
		ExposeHeaders:    exposeHeaders,
		MaxAge:           maxAge,
	}

	return mw, nil
}

func (c CORSMiddleware) Process(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get(HeaderOrigin)
		allowOrigin := ""
		allowMethods := strings.Join(c.AllowMethods, ",")
		allowHeaders := strings.Join(c.AllowHeaders, ",")
		exposeHeaders := strings.Join(c.ExposeHeaders, ",")
		maxAge := strconv.Itoa(c.MaxAge)

		for _, o := range c.AllowOrigins {
			if o == "*" && c.AllowCredentials {
				allowOrigin = origin
				break
			}
			if o == "*" || o == origin {
				allowOrigin = o
				break
			}
			if c.matchSubdomain(origin, o) {
				allowOrigin = origin
				break
			}
		}

		// Simple request
		if r.Method != http.MethodOptions {
			w.Header().Add(HeaderVary, HeaderOrigin)
			w.Header().Set(HeaderAccessControlAllowOrigin, allowOrigin)
			if c.AllowCredentials {
				w.Header().Set(HeaderAccessControlAllowCredentials, "true")
			}
			if len(c.ExposeHeaders) > 0 {
				w.Header().Set(HeaderAccessControlExposeHeaders, exposeHeaders)
			}
			handler(w, r)
			return
		}

		// Preflight request
		w.Header().Add(HeaderVary, HeaderOrigin)
		w.Header().Add(HeaderVary, HeaderAccessControlRequestMethod)
		w.Header().Add(HeaderVary, HeaderAccessControlRequestHeaders)
		w.Header().Set(HeaderAccessControlAllowOrigin, allowOrigin)
		w.Header().Set(HeaderAccessControlAllowMethods, allowMethods)
		if c.AllowCredentials {
			w.Header().Set(HeaderAccessControlAllowCredentials, "true")
		}
		if len(c.AllowHeaders) > 0 {
			w.Header().Set(HeaderAccessControlAllowHeaders, allowHeaders)
		} else {
			h := r.Header.Get(HeaderAccessControlRequestHeaders)
			if h != "" {
				w.Header().Set(HeaderAccessControlAllowHeaders, h)
			}
		}
		if c.MaxAge > 0 {
			w.Header().Set(HeaderAccessControlMaxAge, maxAge)
		}
		handler(w, r)
	}
}

func (c CORSMiddleware) matchScheme(domain, pattern string) bool {
	didx := strings.Index(domain, ":")
	pidx := strings.Index(pattern, ":")
	return didx != -1 && pidx != -1 && domain[:didx] == pattern[:pidx]
}

// matchSubdomain compares authority with wildcard
func (c CORSMiddleware) matchSubdomain(domain, pattern string) bool {
	if !c.matchScheme(domain, pattern) {
		return false
	}
	didx := strings.Index(domain, "://")
	pidx := strings.Index(pattern, "://")
	if didx == -1 || pidx == -1 {
		return false
	}
	domAuth := domain[didx+3:]
	// to avoid long loop by invalid long domain
	if len(domAuth) > 253 {
		return false
	}
	patAuth := pattern[pidx+3:]

	domComp := strings.Split(domAuth, ".")
	patComp := strings.Split(patAuth, ".")
	for i := len(domComp)/2 - 1; i >= 0; i-- {
		opp := len(domComp) - 1 - i
		domComp[i], domComp[opp] = domComp[opp], domComp[i]
	}
	for i := len(patComp)/2 - 1; i >= 0; i-- {
		opp := len(patComp) - 1 - i
		patComp[i], patComp[opp] = patComp[opp], patComp[i]
	}

	for i, v := range domComp {
		if len(patComp) <= i {
			return false
		}
		p := patComp[i]
		if p == "*" {
			return true
		}
		if p != v {
			return false
		}
	}
	return false
}
