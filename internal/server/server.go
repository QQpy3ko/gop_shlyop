package server

import (
	"context"
	"fmt"
	"gop_shlyop/internal/config"
	"gop_shlyop/internal/db"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Server struct {
	cfg    *config.Config
	log    zerolog.Logger
	dbpool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config, log zerolog.Logger) (*Server, error) {
	log.Info().Msg("Connecting to PostgreSQL...")
	dbpool, err := db.NewPostgresDB(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	log.Info().Msg("Successfully connected to PostgreSQL")

	return &Server{
		cfg:    cfg,
		log:    log,
		dbpool: dbpool,
	}, nil
}

// run server
func (s *Server) Run() error {
	s.log.Info().
		Str("host", s.cfg.HTTPServer.Host).
		Str("port", s.cfg.HTTPServer.Port).
		Msg("server is starting")

	// TODO: run web-server in separate goroutine

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // wait for signal

	s.log.Info().Msg("Shutting down server...")

	s.dbpool.Close()

	s.log.Info().Msg("Server stopped gracefully")
	return nil
}
