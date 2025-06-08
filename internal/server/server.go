package server

import (
	"gop_shlyop/internal/config"

	"github.com/rs/zerolog"
)

type Server struct {
	cfg *config.Config
	log zerolog.Logger
	// logger, DB, web-server to be here
}


func New(cfg *config.Config, log zerolog.Logger) *Server {
	return &Server{
		cfg: cfg,
		log: log,
	}
}

// run server
func (s *Server) Run() error {
	s.log.Info().
		Str("host", s.cfg.HTTPServer.Host).
		Str("port", s.cfg.HTTPServer.Port).
		Msg("server is running")

	return nil
}
