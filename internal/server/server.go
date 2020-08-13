package server

import (
	"github.com/Highway-Project/highway/config"
	"github.com/Highway-Project/highway/pkg/router"
	"net/http"
	"time"
)

type Server struct {
	router            router.Router
	httpServer        http.Server
	port              string
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	maxHeaderBytes    int
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func NewServer(global config.GlobalConfig, routerSpec config.RouterSpec, servicesSpec []config.ServiceSpec, rulesSpec []config.RuleSpec, middlewaresSpec []config.MiddlewareSpec) (*Server, error)  {
	return &Server{}, nil
}