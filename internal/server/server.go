package server

import (
	"context"
	"github.com/fiust/highway/pkg/router"
	"github.com/fiust/highway/pkg/rules"
	"net/http"
	"time"
)

type Server struct {
	Router router.Router
	Rules  []rules.Rule
	srv    http.Server
}

func (s *Server) Run() error {
	srv := http.Server{
		Addr:              ":8080",
		Handler:           s.Router,
		TLSConfig:         nil,
		ReadTimeout:       time.Second * 2,
		ReadHeaderTimeout: 0,
		WriteTimeout:      time.Second * 2,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	s.srv = srv
	return srv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.srv.Shutdown(context.Background())
}

func New(router router.Router, rules []rules.Rule) (*Server, error) {
	s := &Server{
		Router: router,
		Rules:  rules,
	}
	err := s.registerRules()
	if err != nil {
		return nil, nil
	}
	return s, nil
}

func (s *Server) registerRules() error {
	for _, rule := range s.Rules {
		err := s.Router.AddRule(rule)
		if err != nil {
			return err
		}
	}
	return nil
}
