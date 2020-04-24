package server

import (
	"github.com/fiust/highway/pkg/router"
	"github.com/fiust/highway/pkg/rules"
)

type Server struct {
	Router      router.Router
	Rules       []rules.Rule
}

func (s *Server) Run() error {
	return nil
}

func (s *Server) Stop() error {
	return nil
}

func New(router router.Router, rules []rules.Rule) (*Server, error) {
	return &Server{
		Router: router,
		Rules:  rules,
	}, nil
}

