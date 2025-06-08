package server

import (
	"fmt"
	"gop_shlyop/internal/config"
)

type Server struct {
	cfg *config.Config
	// logger, DB, web-server to be here
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	fmt.Printf("Server is running with config: %+v\n", s.cfg)

	return nil
}
