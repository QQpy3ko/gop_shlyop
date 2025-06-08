package server

import (
	"context"
	"fmt"
	"gop_shlyop/internal/config"
	"gop_shlyop/internal/db"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger) // Chi's own logger is fine for now
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/reviews", s.addReview)
		r.Get("/reviews", s.getReviews)
		r.Get("/items/{id}/rating", s.getItemRating)
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.cfg.HTTPServer.Host, s.cfg.HTTPServer.Port),
		Handler: router,
	}

	s.log.Info().
		Str("host", s.cfg.HTTPServer.Host).
		Str("port", s.cfg.HTTPServer.Port).
		Msg("server is starting")

	// run web-server in separate goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // wait for signal

	s.log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		s.log.Error().Err(err).Msg("Server forced to shutdown")
	}

	s.dbpool.Close()

	s.log.Info().Msg("Server stopped gracefully")
	return nil
}

// TODO: Implement handlers
func (s *Server) addReview(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (s *Server) getReviews(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (s *Server) getItemRating(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}
